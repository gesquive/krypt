package crypto

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTwofishCipher(t *testing.T) {
	cipher := NewTwofishCipher()

	data := []byte("This is the test data to compare")
	pass := []byte("geronimo")

	encryptedData, err := cipher.Encrypt(data, pass)
	if err != nil {
		t.Fatal()
	}

	decryptedData, err := cipher.Decrypt(encryptedData, pass)
	if err != nil {
		t.Fatal()
	}

	if bytes.Compare(data, decryptedData) != 0 {
		t.Fail()
	}
}

func TestTwofishEncrypt(t *testing.T) {
	cipher := NewTwofishCipher()

	data := []byte(`There is a theory which states that if ever anyone discovers
exactly what the Universe is for and why it is here, it will
instantly disappear and be replaced by something even more
bizarre and inexplicable. There is another theory which states
that this has already happened.`)
	pass := []byte("password")

	encryptedData, err := cipher.Encrypt(data, pass)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, encryptedData, "no encrypted data")
	assert.False(t, bytes.Compare(data, encryptedData) == 0, "encrypted data should be different")
}
