package crypto

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAES256Crypt(t *testing.T) {
	crypt := NewAES256Cipher()

	data := []byte("This is the test data to compare")
	pass := []byte("geronimo")

	encryptedData, err := crypt.Encrypt(data, pass)
	if err != nil {
		t.Fatal()
	}

	decryptedData, err := crypt.Decrypt(encryptedData, pass)
	if err != nil {
		t.Fatal()
	}

	if bytes.Compare(data, decryptedData) != 0 {
		t.Fail()
	}
}

func TestAES256Encrypt(t *testing.T) {
	cipher := NewAES256Cipher()

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
