package azure_test

import (
	"errors"
	"testing"

	"github.com/openshift/cluster-ingress-operator/pkg/dns/azure"
)

func TestFunction(t *testing.T) {
	var zoneTests = []struct {
		name  string
		id    string
		err   error
		subID string
		rg    string
		zone  string
	}{
		{"TestValidZoneID", "/subscriptions/E540B02D-5CCE-4D47-A13B-EB05A19D696E/resourceGroups/test-rg/providers/Microsoft.Network/dnszones/test-rg.dnszone.io",
			nil, "E540B02D-5CCE-4D47-A13B-EB05A19D696E", "test-rg", "test-rg.dnszone.io"},
		{"TestInvalidZoneID", "/subscriptions/E540B02D-5CCE-4D47-A13B-EB05A19D696E/resourceGroups/test-rg/providers/Microsoft.Network/dnszones",
			errors.New("invalid azure dns zone id"), "E540B02D-5CCE-4D47-A13B-EB05A19D696E", "test-rg", "test-rg.dnszone.io"},
		{"TestEmptyZoneID", "", errors.New("invalid azure dns zone id"),
			"E540B02D-5CCE-4D47-A13B-EB05A19D696E", "test-rg", "test-rg.dnszone.io"},
	}
	for _, tt := range zoneTests {
		t.Run(tt.name, func(t *testing.T) {
			zone, err := azure.NewZone(tt.id)
			switch tt.err {
			case nil:
				if err != nil {
					t.Errorf("%s: expected [%s] actual [%s]", tt.name, tt.err, err)
					return
				}
			default:
				if tt.err.Error() != err.Error() {
					t.Errorf("%s: expected [%s] actual [%s]", tt.name, tt.err, err)
				}
				return
			}
			if zone.SubscriptionID != tt.subID {
				t.Errorf("expected [%s] actual [%s]", tt.subID, zone.SubscriptionID)
			}
			if zone.ResourceGroup != tt.rg {
				t.Errorf("expected [%s] actual [%s]", tt.rg, zone.ResourceGroup)
			}
			if zone.Zone != tt.zone {
				t.Errorf("expected [%s] actual [%s]", tt.zone, zone.Zone)
			}
		})
	}
}
