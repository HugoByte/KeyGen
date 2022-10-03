package near

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/btcsuite/btcutil/base58"
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

func NewKeyPairFromPrivateKey(privateKey string) (*Keypair, error) {

	key := strings.Trim(privateKey, "ed25519:")

	keyToByte := base58.Decode(key)

	keyPair := &Keypair{}
	keyPair.privateKey = ed25519.PrivateKey(keyToByte)
	keyPair.publicKey = keyPair.privateKey.Public().(ed25519.PublicKey)

	return keyPair, nil
}

//hub bacon history scissors joy thing surprise friend spin dust winter rural
