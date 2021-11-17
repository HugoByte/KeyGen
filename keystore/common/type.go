package common

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"

	"golang.org/x/crypto/scrypt"
)

type KeyStoreData struct {
	Address  []byte     `json:"address"`
	ID       string     `json:"id"`
	Version  int        `json:"version"`
	CoinType string     `json:"coinType"`
	Crypto   CryptoData `json:"crypto"`
}
type CryptoData struct {
	Cipher       string          `json:"cipher"`
	CipherParams json.RawMessage `json:"cipherparams"`
	CipherText   RawHexBytes     `json:"ciphertext"`
	KDF          string          `json:"kdf"`
	KDFParams    json.RawMessage `json:"kdfparams"`
	MAC          RawHexBytes     `json:"mac"`
}

type ScryptParams struct {
	DKLen int         `json:"dklen"`
	N     int         `json:"n"`
	R     int         `json:"r"`
	P     int         `json:"p"`
	Salt  RawHexBytes `json:"salt"`
}
type AES128CTRParams struct {
	IV RawHexBytes `json:"iv"`
}

type RawHexBytes []byte

func (rh RawHexBytes) MarshalJSON() ([]byte, error) {
	if rh == nil {
		return []byte("null"), nil
	}
	s := hex.EncodeToString(rh)
	return json.Marshal(s)
}

func (rh *RawHexBytes) UnmarshalJSON(b []byte) error {
	var os *string
	if err := json.Unmarshal(b, &os); err != nil {
		return err
	}
	if os == nil {
		*rh = nil
		return nil
	}
	s := *os
	if bin, err := hex.DecodeString(s); err != nil {
		return err
	} else {
		*rh = bin
		return nil
	}
}

func (rh RawHexBytes) Bytes() []byte {
	if rh == nil {
		return nil
	}
	return rh[:]
}

func (rh RawHexBytes) String() string {
	if rh == nil {
		return "null"
	}
	return hex.EncodeToString(rh)
}

func (p *ScryptParams) Init(pKey []byte) error {
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return err
	}
	p.DKLen = len(pKey)
	p.P = 1
	p.R = 16
	p.N = 131072
	p.Salt = salt
	return nil
}

func (p *ScryptParams) Key(pw []byte) ([]byte, error) {
	return scrypt.Key(pw, p.Salt.Bytes(), p.N, p.R, p.P, p.DKLen)
}
