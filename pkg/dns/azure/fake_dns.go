package azure

import (
	"github.com/openshift/cluster-ingress-operator/pkg/dns"
	"github.com/openshift/cluster-ingress-operator/pkg/dns/azure/client"
)

func NewFakeManager(config Config, client client.DNSClient, operatorReleaseVersion string) (dns.Manager, error) {
	return &manager{config: config, client: client}, nil
}
