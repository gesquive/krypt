package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"

	"github.com/pkg/errors"

	"golang.org/x/crypto/pbkdf2"
)

const aes256Name = "AES256"

// AES256Cipher encrypts using AES256-GCM
type AES256Cipher struct {
	description string
	name        string
	cipherType  CipherType
}

// NewAES256Cipher constructor
func NewAES256Cipher() *AES256Cipher {
	r := &AES256Cipher{}
	r.description = "AES256-GCM cipher"
	r.name = aes256Name
	r.cipherType = AES256

	return r
}

// GetDescription returns description string
func (c *AES256Cipher) GetDescription() string {
	return c.description
}

// GetName returns name string
func (c *AES256Cipher) GetName() string {
	return c.name
}

// GetType returns CryptType
func (c *AES256Cipher) GetType() CipherType {
	return c.cipherType
}

// Encrypt data using AES256-GCM cipher. This both hides the content of
// the data and provides a check that it hasn't been altered. Output takes the
// form nonce|ciphertext|tag|salt where '|' indicates concatenation.
func (c *AES256Cipher) Encrypt(data []byte, password []byte) ([]byte, error) {
	// we need the salt as random as possible
	salt := make([]byte, defaultSaltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, errors.Wrapf(err, "randomizing salt")
	}
	// derive a key from password using HMAC-SHA-256 based PBKDF2 key derivation function
	key := pbkdf2.Key(password, salt, 4096, 32, sha256.New)

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrapf(err, "creating block cipher")
	}

	modeCipher, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, errors.Wrapf(err, "creating block mode cipher")
	}

	nonce := make([]byte, modeCipher.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, errors.Wrapf(err, "randomizing nonce")
	}

	ciphertext := modeCipher.Seal(nonce, nonce, data, nil)
	ciphertext = append(ciphertext, salt...)

	return ciphertext, nil
}

// Decrypt data using AES256-GCM cipher. This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag|salt where '|' indicates concatenation.
func (c *AES256Cipher) Decrypt(data []byte, password []byte) ([]byte, error) {
	salt := data[len(data)-defaultSaltSize:]

	key := pbkdf2.Key(password, salt, 4096, 32, sha256.New)

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrapf(err, "creating block cipher")
	}

	modeCipher, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, errors.Wrapf(err, "creating block mode cipher")
	}

	if len(data) < modeCipher.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	plaintext, err := modeCipher.Open(nil, data[:modeCipher.NonceSize()], data[modeCipher.NonceSize():len(data)-defaultSaltSize], nil)
	if err != nil {
		return nil, errors.Wrapf(err, "decrypting data")
	}

	return plaintext, nil
}
