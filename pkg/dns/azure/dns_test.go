package azure_test

import (
	"context"
	"testing"

	v1 "github.com/openshift/api/config/v1"
	"github.com/openshift/cluster-ingress-operator/pkg/dns"
	"github.com/openshift/cluster-ingress-operator/pkg/dns/azure"
	"github.com/openshift/cluster-ingress-operator/pkg/dns/azure/client"
)

func TestEnsureDNS(t *testing.T) {
	c := client.Config{}
	fc, err := client.NewFake(c)
	if err != nil {
		t.Error("failed to create manager")
		return
	}

	cfg := azure.Config{}
	mgr, err := azure.NewFakeManager(cfg, fc, "")
	if err != nil {
		t.Error("failed to create manager")
		return
	}

	rg := "test-rg"
	zone := "test-rg.dnszone.io"
	record := dns.Record{
		Zone: v1.DNSZone{
			ID: "/subscriptions/E540B02D-5CCE-4D47-A13B-EB05A19D696E/resourceGroups/test-rg/providers/Microsoft.Network/dnszones/test-rg.dnszone.io",
		},
		Type: dns.ARecordType,
		ARecord: &dns.ARecord{
			Domain:  "dnszone.io",
			Address: "55.11.22.33",
		},
	}
	err = mgr.Ensure(&record)
	if err != nil {
		t.Error("failed to ensure dns")
		return
	}

	result := fc.Get(context.TODO(), rg, zone)
}
