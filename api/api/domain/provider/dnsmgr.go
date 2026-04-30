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
	baseURL    string
	uid        string
	apiKey     string
	httpClient *http.Client
}

type DnsmgrPageResult struct {
	List  []map[string]interface{} `json:"list"`
	Total int64                    `json:"total"`
	Raw   interface{}              `json:"raw,omitempty"`
}

type DnsmgrRecordRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	Line     string `json:"line"`
	TTL      string `json:"ttl"`
	MX       string `json:"mx"`
	Priority string `json:"priority"`
	Remark   string `json:"remark"`
	Status   string `json:"status"`
	RecordID string `json:"recordId"`
}

type DnsmgrBatchRequest struct {
	Action     string                   `json:"action"`
	RecordInfo []map[string]interface{} `json:"recordInfo"`
	Remark     string                   `json:"remark"`
	Status     string                   `json:"status"`
}

func NewDnsmgrClient(system config.ExternalSystem) *DnsmgrClient {
	return &DnsmgrClient{
		baseURL:    strings.TrimRight(system.BaseURL, "/"),
		uid:        system.Metadata["uid"],
		apiKey:     system.Metadata["apiKey"],
		httpClient: &http.Client{Timeout: 20 * time.Second},
	}
}

func (c *DnsmgrClient) HealthCheck(ctx context.Context) error {
	if _, err := c.ListDomains(ctx, url.Values{"offset": []string{"0"}, "limit": []string{"1"}}); err != nil {
		return err
	}
	return nil
}

func (c *DnsmgrClient) ListDomains(ctx context.Context, query url.Values) (*DnsmgrPageResult, error) {
	return c.postPage(ctx, "/api/domain", query)
}

func (c *DnsmgrClient) DomainDetail(ctx context.Context, domainID string, loginURL bool) (map[string]interface{}, error) {
	form := url.Values{}
	if loginURL {
		form.Set("loginurl", "1")
	}
	payload, err := c.post(ctx, "/api/domain/"+url.PathEscape(domainID), form)
	if err != nil {
		return nil, err
	}
	return normalizeObject(payload), nil
}

func (c *DnsmgrClient) ListRecords(ctx context.Context, domainID string, query url.Values) (*DnsmgrPageResult, error) {
	return c.postPage(ctx, "/api/record/data/"+url.PathEscape(domainID), query)
}

func (c *DnsmgrClient) AddRecord(ctx context.Context, domainID string, record DnsmgrRecordRequest) (interface{}, error) {
	return c.post(ctx, "/api/record/add/"+url.PathEscape(domainID), record.toValues())
}

func (c *DnsmgrClient) UpdateRecord(ctx context.Context, domainID, recordID string, record DnsmgrRecordRequest) (interface{}, error) {
	form := record.toValues()
	form.Set("recordid", recordID)
	return c.post(ctx, "/api/record/update/"+url.PathEscape(domainID), form)
}

func (c *DnsmgrClient) DeleteRecord(ctx context.Context, domainID, recordID string) (interface{}, error) {
	return c.post(ctx, "/api/record/delete/"+url.PathEscape(domainID), url.Values{"recordid": []string{recordID}})
}

func (c *DnsmgrClient) SetRecordStatus(ctx context.Context, domainID, recordID, status string) (interface{}, error) {
	return c.post(ctx, "/api/record/status/"+url.PathEscape(domainID), url.Values{"recordid": []string{recordID}, "status": []string{status}})
}

func (c *DnsmgrClient) SetRecordRemark(ctx context.Context, domainID, recordID, remark string) (interface{}, error) {
	return c.post(ctx, "/api/record/remark/"+url.PathEscape(domainID), url.Values{"recordid": []string{recordID}, "remark": []string{remark}})
}

func (c *DnsmgrClient) BatchRecords(ctx context.Context, domainID string, batch DnsmgrBatchRequest) (interface{}, error) {
	recordInfo, err := json.Marshal(batch.RecordInfo)
	if err != nil {
		return nil, fmt.Errorf("encode recordinfo failed: %w", err)
	}
	form := url.Values{
		"action":     []string{batch.Action},
		"recordinfo": []string{string(recordInfo)},
	}
	if batch.Remark != "" {
		form.Set("remark", batch.Remark)
	}
	if batch.Status != "" {
		form.Set("status", batch.Status)
	}
	return c.post(ctx, "/api/record/batch/"+url.PathEscape(domainID), form)
}

func (r DnsmgrRecordRequest) toValues() url.Values {
	form := url.Values{}
	setIfNotEmpty(form, "name", r.Name)
	setIfNotEmpty(form, "type", strings.ToUpper(r.Type))
	setIfNotEmpty(form, "value", r.Value)
	setIfNotEmpty(form, "line", r.Line)
	setIfNotEmpty(form, "ttl", r.TTL)
	setIfNotEmpty(form, "mx", firstNonEmpty(r.MX, r.Priority))
	setIfNotEmpty(form, "remark", r.Remark)
	setIfNotEmpty(form, "status", r.Status)
	setIfNotEmpty(form, "recordid", r.RecordID)
	return form
}

func (c *DnsmgrClient) postPage(ctx context.Context, path string, form url.Values) (*DnsmgrPageResult, error) {
	payload, err := c.post(ctx, path, form)
	if err != nil {
		return nil, err
	}
	list := findObjectList(payload)
	return &DnsmgrPageResult{
		List:  list,
		Total: pickTotal(payload, len(list)),
		Raw:   payload,
	}, nil
}

func (c *DnsmgrClient) post(ctx context.Context, path string, form url.Values) (interface{}, error) {
	if c.baseURL == "" {
		return nil, fmt.Errorf("dnsmgr baseUrl is not configured")
	}
	if form == nil {
		form = url.Values{}
	}
	c.addAuthParams(form)

	endpoint, err := url.JoinPath(c.baseURL, strings.TrimPrefix(path, "/"))
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader([]byte(form.Encode())))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
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
	if message := pickErrorMessage(payload); message != "" {
		return nil, fmt.Errorf("dnsmgr error: %s", message)
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
	return []map[string]interface{}{}
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

func normalizeObject(payload interface{}) map[string]interface{} {
	if object, ok := payload.(map[string]interface{}); ok {
		if data, ok := object["data"].(map[string]interface{}); ok {
			return data
		}
		return object
	}
	return map[string]interface{}{"raw": payload}
}

func pickTotal(payload interface{}, fallback int) int64 {
	if object, ok := payload.(map[string]interface{}); ok {
		for _, key := range []string{"total", "count", "recordsTotal"} {
			if value, ok := object[key]; ok {
				if total, ok := numberToInt64(value); ok {
					return total
				}
			}
		}
		if data, ok := object["data"].(map[string]interface{}); ok {
			for _, key := range []string{"total", "count", "recordsTotal"} {
				if value, ok := data[key]; ok {
					if total, ok := numberToInt64(value); ok {
						return total
					}
				}
			}
		}
	}
	return int64(fallback)
}

func numberToInt64(value interface{}) (int64, bool) {
	switch typed := value.(type) {
	case float64:
		return int64(typed), true
	case int64:
		return typed, true
	case string:
		parsed, err := strconv.ParseInt(typed, 10, 64)
		return parsed, err == nil
	default:
		return 0, false
	}
}

func pickErrorMessage(payload interface{}) string {
	object, ok := payload.(map[string]interface{})
	if !ok {
		return ""
	}
	code := strings.TrimSpace(fmt.Sprintf("%v", object["code"]))
	status := strings.TrimSpace(fmt.Sprintf("%v", object["status"]))
	success := strings.TrimSpace(fmt.Sprintf("%v", object["success"]))
	if code == "0" || code == "200" || status == "success" || status == "1" || success == "true" {
		return ""
	}
	if code != "" && code != "<nil>" {
		return firstNonEmpty(fmt.Sprintf("%v", object["msg"]), fmt.Sprintf("%v", object["message"]))
	}
	return ""
}

func setIfNotEmpty(form url.Values, key, value string) {
	if strings.TrimSpace(value) != "" {
		form.Set(key, strings.TrimSpace(value))
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
