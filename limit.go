package nostr

// NIP-11 - https://github.com/nostr-protocol/nips/blob/master/11.md

type RelayInformation struct {
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	PubKey        string           `json:"pubkey"`
	Contact       string           `json:"contact"`
	SupportedNIPs []int            `json:"supported_nips"`
	Version       string           `json:"version"`
	Limitation    *RelayLimitation `json:"limitation,omitempty"`
}

type RelayLimitation struct {
	MaxFilters      int  `json:"max_filters,omitempty"`
	MaxLimit        int  `json:"max_limit,omitempty"`
	AuthRequired    bool `json:"auth_required"`
	PaymentRequired bool `json:"payment_required"`
}
