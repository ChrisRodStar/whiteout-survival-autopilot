package domain

// GiftCodes — root of the YAML file db/giftCodes.yaml
type GiftCodes struct {
	Codes []GiftCode `yaml:"codes"`
}

// GiftCode — a single promo code
type GiftCode struct {
	Name    string            `yaml:"name"`
	Expires string            `yaml:"expires,omitempty"` // RFC 3339 UTC
	UserFor map[string]string `yaml:"userFor"`           // uid → status
}
