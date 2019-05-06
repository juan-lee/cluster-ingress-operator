package client

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2017-10-01/dns"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/pkg/errors"
)

type DNSClient interface {
	Put(ctx context.Context, arec ARecord) error
	Delete(ctx context.Context, arec ARecord) error
}

type Config struct {
	Environment    string
	SubscriptionID string
	ClientID       string
	ClientSecret   string
	TenantID       string
}

type ARecord struct {
	ResourceGroup     string
	Zone              string
	RelativeRecordSet string
	Address           string
	TTL               int64
}

type dnsClient struct {
	zones      dns.ZonesClient
	recordSets dns.RecordSetsClient
	config     Config
}

func New(config Config) (DNSClient, error) {
	authorizer, err := getAuthorizerForResource(config)
	if err != nil {
		return nil, err
	}
	zc := dns.NewZonesClient(config.SubscriptionID)
	zc.Authorizer = authorizer
	rc := dns.NewRecordSetsClient(config.SubscriptionID)
	rc.Authorizer = authorizer
	return &dnsClient{zones: zc, recordSets: rc, config: config}, nil
}

func (c *dnsClient) Put(ctx context.Context, arec ARecord) error {
	rs := dns.RecordSet{RecordSetProperties: &dns.RecordSetProperties{
		TTL: to.Int64Ptr(arec.TTL),
		ARecords: &[]dns.ARecord{
			{Ipv4Address: &arec.Address},
		},
	}}
	_, err := c.recordSets.CreateOrUpdate(ctx, arec.ResourceGroup, arec.Zone, arec.RelativeRecordSet, dns.A, rs, "", "")
	if err != nil {
		return errors.Wrap(err, "failed to update dns a record")
	}
	return nil
}

func (c *dnsClient) Delete(ctx context.Context, arec ARecord) error {
	rs := dns.RecordSet{RecordSetProperties: &dns.RecordSetProperties{
		TTL: to.Int64Ptr(arec.TTL),
		ARecords: &[]dns.ARecord{
			{Ipv4Address: &arec.Address},
		},
	}}
	_, err := c.recordSets.CreateOrUpdate(ctx, arec.ResourceGroup, arec.Zone, arec.RelativeRecordSet, dns.A, rs, "", "")
	if err != nil {
		return errors.Wrap(err, "failed to update dns a record")
	}
	return nil
}
