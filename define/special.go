package define

import "golang.org/x/crypto/ed25519"

type SpecialOP struct {
	Type   string            // only support ed25519
	PubKey ed25519.PublicKey // ed25519 is 32 bytes
	Power  uint32            // 0-10
}
