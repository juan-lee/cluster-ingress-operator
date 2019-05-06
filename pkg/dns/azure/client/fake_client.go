package client

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2017-10-01/dns"
	"github.com/Azure/go-autorest/autorest/to"
)

type FakeDNSClient struct {
	fakeARM map[string]dns.RecordSet
}

func NewFake(config Config) (*FakeDNSClient, error) {
	return &FakeDNSClient{}, nil
}

func (c *FakeDNSClient) Put(ctx context.Context, arec ARecord) error {
	rs := dns.RecordSet{RecordSetProperties: &dns.RecordSetProperties{
		TTL: to.Int64Ptr(arec.TTL),
		ARecords: &[]dns.ARecord{
			{Ipv4Address: &arec.Address},
		},
	}}
	c.fakeARM[arec.ResourceGroup+arec.Zone+arec.RelativeRecordSet] = rs
	return nil
}

func (c *FakeDNSClient) Delete(ctx context.Context, arec ARecord) error {
	delete(c.fakeARM, arec.ResourceGroup+arec.Zone+arec.RelativeRecordSet)
	return nil
}

func (c *FakeDNSClient) Get(ctx context.Context, rg, zone, rel string) (dns.RecordSet, bool) {
	rs, ok := c.fakeARM[rg+zone+rel]
	return rs, ok
}
