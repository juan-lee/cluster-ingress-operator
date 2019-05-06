package main

import (
	"context"
	"fmt"
	"os"

	"github.com/openshift/cluster-ingress-operator/pkg/dns"
	awsdns "github.com/openshift/cluster-ingress-operator/pkg/dns/aws"
	azuredns "github.com/openshift/cluster-ingress-operator/pkg/dns/azure"
	logf "github.com/openshift/cluster-ingress-operator/pkg/log"
	"github.com/openshift/cluster-ingress-operator/pkg/operator"
	operatorclient "github.com/openshift/cluster-ingress-operator/pkg/operator/client"
	operatorconfig "github.com/openshift/cluster-ingress-operator/pkg/operator/config"
	"github.com/openshift/cluster-ingress-operator/pkg/operator/controller"

	configv1 "github.com/openshift/api/config/v1"

	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/types"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

const (
	// cloudCredentialsSecretName is the name of the secret in the
	// operator's namespace that will hold the credentials that the operator
	// will use to authenticate with the cloud API.
	cloudCredentialsSecretName = "cloud-credentials"
)

var log = logf.Logger.WithName("entrypoint")

func main() {
	metrics.DefaultBindAddress = ":60000"

	// Get a kube client.
	kubeConfig, err := config.GetConfig()
	if err != nil {
		log.Error(err, "failed to get kube config")
		os.Exit(1)
	}
	kubeClient, err := operatorclient.NewClient(kubeConfig)
	if err != nil {
		log.Error(err, "failed to create kube client")
		os.Exit(1)
	}

	// Collect operator configuration.
	operatorNamespace := os.Getenv("WATCH_NAMESPACE")
	if len(operatorNamespace) == 0 {
		log.Error(fmt.Errorf("missing environment variable"), "'WATCH_NAMESPACE' environment variable must be set")
		os.Exit(1)
	}
	routerImage := os.Getenv("IMAGE")
	if len(routerImage) == 0 {
		log.Error(fmt.Errorf("missing environment variable"), "'IMAGE' environment variable must be set")
		os.Exit(1)
	}
	releaseVersion := os.Getenv("RELEASE_VERSION")
	if len(releaseVersion) == 0 {
		releaseVersion = controller.UnknownReleaseVersionName
		log.Info("RELEASE_VERSION environment variable missing", "release version", controller.UnknownReleaseVersionName)
	}

	// Retrieve the cluster infrastructure config.
	infraConfig := &configv1.Infrastructure{}
	err = kubeClient.Get(context.TODO(), types.NamespacedName{Name: "cluster"}, infraConfig)
	if err != nil {
		log.Error(err, "failed to get infrastructure 'config'")
		os.Exit(1)
	}

	dnsConfig := &configv1.DNS{}
	err = kubeClient.Get(context.TODO(), types.NamespacedName{Name: "cluster"}, dnsConfig)
	if err != nil {
		log.Error(err, "failed to get dns 'cluster'")
		os.Exit(1)
	}

	operatorConfig := operatorconfig.Config{
		OperatorReleaseVersion: releaseVersion,
		Namespace:              operatorNamespace,
		RouterImage:            routerImage,
	}

	// Set up the DNS manager.
	dnsManager, err := createDNSManager(kubeClient, operatorConfig, infraConfig, dnsConfig)
	if err != nil {
		log.Error(err, "failed to create DNS manager")
		os.Exit(1)
	}

	// Set up and start the operator.
	op, err := operator.New(operatorConfig, dnsManager, kubeConfig)
	if err != nil {
		log.Error(err, "failed to create operator")
		os.Exit(1)
	}
	if err := op.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "failed to start operator")
		os.Exit(1)
	}
}

// createDNSManager creates a DNS manager compatible with the given cluster
// configuration.
func createDNSManager(cl client.Client, operatorConfig operatorconfig.Config, infraConfig *configv1.Infrastructure, dnsConfig *configv1.DNS) (dns.Manager, error) {
	var dnsManager dns.Manager
	switch infraConfig.Status.Platform {
	case configv1.AWSPlatformType:
		awsCreds := &corev1.Secret{}
		err := cl.Get(context.TODO(), types.NamespacedName{Namespace: operatorConfig.Namespace, Name: cloudCredentialsSecretName}, awsCreds)
		if err != nil {
			return nil, fmt.Errorf("failed to get aws creds from secret %s/%s: %v", awsCreds.Namespace, awsCreds.Name, err)
		}
		log.Info("using aws creds from secret", "namespace", awsCreds.Namespace, "name", awsCreds.Name)
		manager, err := awsdns.NewManager(awsdns.Config{
			AccessID:  string(awsCreds.Data["aws_access_key_id"]),
			AccessKey: string(awsCreds.Data["aws_secret_access_key"]),
			DNS:       dnsConfig,
		}, operatorConfig.OperatorReleaseVersion)
		if err != nil {
			return nil, fmt.Errorf("failed to create AWS DNS manager: %v", err)
		}
		dnsManager = manager
	case configv1.AzurePlatformType:
		azureCreds := &corev1.Secret{}
		err := cl.Get(context.TODO(), types.NamespacedName{Namespace: operatorConfig.Namespace, Name: cloudCredentialsSecretName}, azureCreds)
		if err != nil {
			return nil, fmt.Errorf("failed to get azure creds from secret %s/%s: %v", azureCreds.Namespace, azureCreds.Name, err)
		}
		log.Info("using azure creds from secret", "namespace", azureCreds.Namespace, "name", azureCreds.Name)
		manager, err := azuredns.NewManager(azuredns.Config{
			Environment:    "AzurePublicCloud",
			ClientID:       string(azureCreds.Data["azure_client_id"]),
			ClientSecret:   string(azureCreds.Data["azure_client_secret"]),
			TenantID:       string(azureCreds.Data["azure_tenant_id"]),
			SubscriptionID: string(azureCreds.Data["azure_subscription_id"]),
			DNS:            dnsConfig,
		}, operatorConfig.OperatorReleaseVersion)
		if err != nil {
			return nil, fmt.Errorf("failed to create Azure DNS manager: %v", err)
		}
		dnsManager = manager
	default:
		dnsManager = &dns.NoopManager{}
	}
	return dnsManager, nil
}
