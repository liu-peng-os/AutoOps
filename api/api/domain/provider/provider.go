package provider

import "context"

type Zone struct {
	Name        string
	DisplayName string
	Provider    string
	Status      string
	ExternalID  string
	RawData     string
}

type Record struct {
	ZoneExternalID string
	ZoneName       string
	Name           string
	Type           string
	Value          string
	Line           string
	TTL            int
	Priority       int
	Status         string
	ExternalID     string
	RawData        string
}

type DomainProvider interface {
	HealthCheck(ctx context.Context) error
	ListZones(ctx context.Context) ([]Zone, error)
	ListRecords(ctx context.Context, zone Zone) ([]Record, error)
}
