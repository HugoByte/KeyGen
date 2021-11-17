package near

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
)

type Keypair struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

func NewKeyPair() (*Keypair, error) {

	keypair := &Keypair{}

	public, private, err := ed25519.GenerateKey(rand.Reader)

	if err != nil {

		return nil, fmt.Errorf("can't genetrate Ed25519 Keys :%s", err)
	}

	keypair.privateKey = private
	keypair.publicKey = public

	return keypair, nil
}
