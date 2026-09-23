package out

type PublicKey struct {
	ID        string
	KeyType   string
	Curve     string
	Algorithm string
	X         string
}

type SigningKeys interface {
	PublicKeys() []PublicKey
}
