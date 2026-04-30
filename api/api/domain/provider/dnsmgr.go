package provider

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"dodevops-api/common/config"
)

type DnsmgrClient struct {
	baseURL     string
	uid         string
	apiKey      string
	token       string
	authHeader  string
	authScheme  string
	healthPath  string
	zonesPath   string
	recordsPath string
	httpClient  *http.Client
}

func NewDnsmgrClient(system config.ExternalSystem) *DnsmgrClient {
	metadata := system.Metadata
	return &DnsmgrClient{
		baseURL:     strings.TrimRight(system.BaseURL, "/"),
		uid:         metadata["uid"],
		apiKey:      metadata["apiKey"],
		token:       metadata["token"],
		authHeader:  firstNonEmpty(metadata["authHeader"], "Authorization"),
		authScheme:  metadata["authScheme"],
		healthPath:  firstNonEmpty(metadata["healthPath"], "/"),
		zonesPath:   firstNonEmpty(metadata["zonesPath"], "/api/domain"),
		recordsPath: firstNonEmpty(metadata["recordsPath"], "/api/record/data/{zone_id}"),
		httpClient:  &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *DnsmgrClient) HealthCheck(ctx context.Context) error {
	if c.baseURL == "" {
		return fmt.Errorf("dnsmgr baseUrl is not configured")
	}

	req, err := c.newRequest(ctx, http.MethodGet, c.healthPath, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("dnsmgr health check returned HTTP %d", resp.StatusCode)
	}
	return nil
}

func (c *DnsmgrClient) ListZones(ctx context.Context) ([]Zone, error) {
	body, err := c.postJSON(ctx, c.zonesPath, url.Values{
		"offset": []string{"0"},
		"limit":  []string{"1000"},
	})
	if err != nil {
		return nil, err
	}

	items := findObjectList(body)
	zones := make([]Zone, 0, len(items))
	for i, item := range items {
		name := pickString(item, "domain", "name", "zone", "domain_name")
		if name == "" {
			continue
		}
		externalID := pickString(item, "id", "domain_id", "zone_id")
		if externalID == "" {
			externalID = name
		}
		raw, _ := json.Marshal(item)
		zones = append(zones, Zone{
			Name:        name,
			DisplayName: firstNonEmpty(pickString(item, "title", "remark", "note"), name),
			Provider:    firstNonEmpty(pickString(item, "typename"), pickString(item, "type", "provider", "dns", "dns_provider")),
			Status:      normalizeStatus(pickString(item, "status", "state", "enabled", "checkstatus")),
			ExternalID:  externalID,
			RawData:     string(raw),
		})
		_ = i
	}
	return zones, nil
}

func (c *DnsmgrClient) ListRecords(ctx context.Context, zone Zone) ([]Record, error) {
	path := strings.ReplaceAll(c.recordsPath, "{zone_id}", url.QueryEscape(zone.ExternalID))
	path = strings.ReplaceAll(path, "{domain}", url.QueryEscape(zone.Name))

	body, err := c.postJSON(ctx, path, url.Values{
		"offset": []string{"0"},
		"limit":  []string{"5000"},
	})
	if err != nil {
		return nil, err
	}

	items := findObjectList(body)
	records := make([]Record, 0, len(items))
	for _, item := range items {
		recordType := strings.ToUpper(pickString(item, "type", "record_type"))
		name := firstNonEmpty(pickString(item, "name", "host", "sub_domain", "rr"), "@")
		value := pickString(item, "value", "record", "content", "target")
		if recordType == "" || value == "" {
			continue
		}
		externalID := pickString(item, "id", "record_id", "recordid")
		if externalID == "" {
			externalID = fmt.Sprintf("%s:%s:%s:%s", zone.ExternalID, name, recordType, value)
		}
		raw, _ := json.Marshal(item)
		records = append(records, Record{
			ZoneExternalID: zone.ExternalID,
			ZoneName:       zone.Name,
			Name:           name,
			Type:           recordType,
			Value:          value,
			Line:           firstNonEmpty(pickString(item, "linename"), pickString(item, "line", "view", "route")),
			TTL:            pickInt(item, "ttl"),
			Priority:       pickInt(item, "priority", "mx"),
			Status:         normalizeStatus(pickString(item, "status", "state", "enabled")),
			ExternalID:     externalID,
			RawData:        string(raw),
		})
	}
	return records, nil
}

func (c *DnsmgrClient) postJSON(ctx context.Context, path string, form url.Values) (interface{}, error) {
	if form == nil {
		form = url.Values{}
	}
	c.addAuthParams(form)

	req, err := c.newRequest(ctx, http.MethodPost, path, []byte(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("dnsmgr returned HTTP %d: %s", resp.StatusCode, string(body))
	}

	var payload interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("parse dnsmgr response failed: %w", err)
	}
	return payload, nil
}

func (c *DnsmgrClient) addAuthParams(form url.Values) {
	if c.uid == "" || c.apiKey == "" {
		return
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	form.Set("uid", c.uid)
	form.Set("timestamp", timestamp)
	form.Set("sign", fmt.Sprintf("%x", md5.Sum([]byte(c.uid+timestamp+c.apiKey))))
}

func (c *DnsmgrClient) newRequest(ctx context.Context, method, path string, body []byte) (*http.Request, error) {
	if c.baseURL == "" {
		return nil, fmt.Errorf("dnsmgr baseUrl is not configured")
	}
	endpoint, err := url.JoinPath(c.baseURL, strings.TrimPrefix(path, "/"))
	if err != nil {
		return nil, err
	}
	if strings.Contains(path, "?") {
		endpoint = c.baseURL + "/" + strings.TrimPrefix(path, "/")
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if c.token != "" {
		value := c.token
		if c.authScheme != "" {
			value = c.authScheme + " " + c.token
		}
		req.Header.Set(c.authHeader, value)
	}
	return req, nil
}

func findObjectList(payload interface{}) []map[string]interface{} {
	switch value := payload.(type) {
	case []interface{}:
		return toObjectList(value)
	case map[string]interface{}:
		for _, key := range []string{"data", "list", "rows", "items", "records", "domains"} {
			if child, ok := value[key]; ok {
				if list := findObjectList(child); len(list) > 0 {
					return list
				}
			}
		}
	}
	return nil
}

func toObjectList(items []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		if object, ok := item.(map[string]interface{}); ok {
			result = append(result, object)
		}
	}
	return result
}

func pickString(item map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		value, ok := lookupValue(item, key)
		if ok && value != nil {
			switch typed := value.(type) {
			case string:
				return strings.TrimSpace(typed)
			case float64:
				return strconv.FormatInt(int64(typed), 10)
			case bool:
				if typed {
					return "true"
				}
				return "false"
			default:
				return strings.TrimSpace(fmt.Sprintf("%v", typed))
			}
		}
	}
	return ""
}

func pickInt(item map[string]interface{}, keys ...string) int {
	for _, key := range keys {
		value, ok := lookupValue(item, key)
		if ok && value != nil {
			switch typed := value.(type) {
			case float64:
				return int(typed)
			case string:
				parsed, _ := strconv.Atoi(typed)
				return parsed
			}
		}
	}
	return 0
}

func lookupValue(item map[string]interface{}, key string) (interface{}, bool) {
	if value, ok := item[key]; ok {
		return value, true
	}
	normalizedKey := strings.ToLower(strings.ReplaceAll(key, "_", ""))
	for itemKey, value := range item {
		normalizedItemKey := strings.ToLower(strings.ReplaceAll(itemKey, "_", ""))
		if normalizedItemKey == normalizedKey {
			return value, true
		}
	}
	return nil, false
}

func normalizeStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "1", "true", "enabled", "enable", "active", "ok", "normal":
		return "enabled"
	case "0", "false", "disabled", "disable", "inactive", "paused":
		return "disabled"
	case "":
		return "unknown"
	default:
		return status
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
