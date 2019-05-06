package azure

import (
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/cluster-ingress-operator/pkg/dns"
	"github.com/openshift/cluster-ingress-operator/pkg/dns/azure/client"
	logf "github.com/openshift/cluster-ingress-operator/pkg/log"
)

var (
	_   dns.Manager = &manager{}
	log             = logf.Logger.WithName("dns")
)

// Config is the necessary input to configure the manager for azure.
type Config struct {
	// Environment is the azure cloud environment.
	Environment string
	// ClientID is an azure service principal appID.
	ClientID string
	// ClientSecret is an azure service principal's credential.
	ClientSecret string
	// TenantID is the azure identity's tenant ID.
	TenantID string
	// SubscriptionID is the azure identity's subscription ID.
	SubscriptionID string
	// DNS is public and private DNS zone configuration for the cluster.
	DNS *configv1.DNS
}

type manager struct {
	config       Config
	client       client.DNSClient
	clientConfig client.Config
}

func NewManager(config Config, operatorReleaseVersion string) (dns.Manager, error) {
	c, err := client.New(client.Config{
		Environment:    config.Environment,
		SubscriptionID: config.SubscriptionID,
		ClientID:       config.ClientID,
		ClientSecret:   config.ClientSecret,
		TenantID:       config.TenantID,
	})
	if err != nil {
		return nil, err
	}
	return &manager{config: config, client: c}, nil
}

func (m *manager) Ensure(record *dns.Record) error {
	return nil
}

func (m *manager) Delete(record *dns.Record) error {
	panic("not implemented")
}
