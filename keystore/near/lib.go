package near

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io"

	"github.com/gofrs/uuid"
	"github.com/hugobyte/keygen/keystore/common"
	"github.com/pkg/errors"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/sha3"
)

const (
	coin            = "near"
	kdfScrypt       = "scrypt-sha256"
	cipherAES128CTR = "aes-128-ctr"
)

func GenerateNewKeystore(file string, pw string) error {

	///Generate New KeyPair

	keypair, err := NewKeyPair()
	if err != nil {
		return err
	}

	/// Genreate KeyStore from the Private Key obtained from Keypair

	err = EncryptKey(keypair.privateKey, pw, file)
	if err != nil {
		return err
	}

	return nil
}

func sha3sumkeccak256(data ...[]byte) []byte {
	s := sha3.NewLegacyKeccak256()
	for _, d := range data {
		s.Write(d)
	}
	return s.Sum([]byte{})
}

func EncryptKey(s ed25519.PrivateKey, pw string, file string) error {
	var ks common.KeyStoreData
	var k common.ScryptParams
	var c common.AES128CTRParams

	///password to bytes of size 64
	passBytes := sha512.Sum512([]byte(pw))

	if err := k.Init(s); err != nil {
		return err
	}

	key, err := k.Key(passBytes[:])
	if err != nil {
		return err
	}

	ks.Crypto.KDF = kdfScrypt
	ks.Crypto.KDFParams, err = json.Marshal(&k)
	if err != nil {
		return err
	}

	cipherBlock, err := aes.NewCipher(key[0:16])
	if err != nil {
		return err
	}

	//Initialization Vector
	c.IV = make([]byte, cipherBlock.BlockSize())
	_, err = io.ReadFull(rand.Reader, c.IV)
	if err != nil {
		return err
	}

	//creating cipher text

	cipherText := make([]byte, len(s))
	//create cipher stream
	cipherStrem := cipher.NewCTR(cipherBlock, c.IV)

	//creating encrypted key

	cipherStrem.XORKeyStream(cipherText, s)

	ks.Crypto.Cipher = cipherAES128CTR
	ks.Crypto.CipherParams, err = json.Marshal(&c)
	if err != nil {
		return err
	}

	ks.Crypto.CipherText = cipherText

	ks.Crypto.MAC = sha3sumkeccak256(key[16:32], cipherText)
	ks.CoinType = coin
	ks.Version = 3
	ks.ID = uuid.Must(uuid.NewV4()).String()

	address := hex.EncodeToString(s.Public().(ed25519.PublicKey))

	ks.Address = []byte(address)

	output, err := json.Marshal(&ks)

	if err != nil {
		return err
	}

	return common.WriteFile(file, output, 0600)
}

func DecryptKey(ksData *common.KeyStoreData, password string) (ed25519.PrivateKey, error) {

	pwb := sha512.Sum512([]byte(password))
	if ksData.Crypto.Cipher != "aes-128-ctr" {
		return nil, errors.Errorf("UnsupportedCipher(cipher=%s)",
			ksData.Crypto.Cipher)
	}
	var cipherParams common.AES128CTRParams
	if err := json.Unmarshal(ksData.Crypto.CipherParams, &cipherParams); err != nil {
		return nil, err
	}

	if ksData.Crypto.KDF != "scrypt-sha256" {
		return nil, errors.Errorf("UnsupportedKDF(kdf=%s)", ksData.Crypto.KDF)
	}
	var kdfParams common.ScryptParams
	if err := json.Unmarshal(ksData.Crypto.KDFParams, &kdfParams); err != nil {
		return nil, err
	}

	key, err := scrypt.Key(pwb[:], kdfParams.Salt, int(kdfParams.N), int(kdfParams.R), int(kdfParams.P), kdfParams.DKLen)
	if err != nil {
		return nil, err
	}

	cipheredBytes := ksData.Crypto.CipherText.Bytes()

	s := sha3.NewLegacyKeccak256()
	s.Write(key[16:32])
	s.Write(cipheredBytes)
	mac := s.Sum([]byte{})

	if !bytes.Equal(mac, ksData.Crypto.MAC.Bytes()) {
		return nil, errors.Errorf("InvalidPassword")
	}

	block, err := aes.NewCipher(key[0:16])
	if err != nil {
		return nil, err
	}

	secretBytes := make([]byte, len(cipheredBytes))

	stream := cipher.NewCTR(block, cipherParams.IV.Bytes())
	stream.XORKeyStream(secretBytes, cipheredBytes)

	secret := secretBytes
	if err != nil {
		return nil, err
	}

	private := ed25519.PrivateKey(secret)

	return private, nil
}
