package azure

import (
	"errors"
	"strings"
)

type Zone struct {
	SubscriptionID string
	ResourceGroup  string
	Zone           string
}

func NewZone(id string) (*Zone, error) {
	s := strings.Split(id, "/")
	if len(s) < 9 {
		return nil, errors.New("invalid azure dns zone id")
	}
	return &Zone{SubscriptionID: s[2], ResourceGroup: s[4], Zone: s[8]}, nil
}
