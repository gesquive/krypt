package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKrypt(t *testing.T) {
	data := []byte("This is the test data to compare")
	pass := []byte("geronimo")

	encryptedData, err := Encrypt(AES256, pass, data)
	if err != nil {
		t.Fail()
	}

	decryptedData, err := Decrypt(pass, encryptedData)
	if err != nil {
		t.Fatal("error decrypting: ", err)
	}

	assert.Equal(t, data, decryptedData,
		"decrypted data does not match original data")
}

func TestValidKryptHeader(t *testing.T) {
	mockPayload := []byte("unimportant data")
	mockData := mockKrypt(libVersion, AES256, mockPayload)

	kryptVersion, cipherType, err := getKryptInfo(mockData)
	if err != nil {
		t.Fatal("error: ", err)
	}

	assert.Equal(t, libVersion, kryptVersion, "krypt version mismatch")
	assert.Equal(t, AES256, cipherType, "cipher type mismatch")
}

func TestInvalidKryptVersion(t *testing.T) {
	mockPayload := []byte("completely random data")
	mockData := mockKrypt(255, AES256, mockPayload)

	kryptVersion, cipherType, err := getKryptInfo(mockData)
	if err == nil {
		t.Fatal("no error returned")
	}

	assert.Equal(t, uint8(0), kryptVersion, "krypt version mismatch")
	assert.Equal(t, CipherType(0), cipherType, "cipher type mismatch")
}

func TestInvalidCipherVersion(t *testing.T) {
	mockPayload := []byte("completely random data")
	mockData := mockKrypt(libVersion, 255, mockPayload)

	kryptVersion, cipherType, err := getKryptInfo(mockData)
	if err == nil {
		t.Fatal("no error returned")
	}

	assert.Equal(t, uint8(0), kryptVersion, "krypt version mismatch")
	assert.Equal(t, CipherType(0), cipherType, "cipher type mismatch")
}

func TestNoCipherVersion(t *testing.T) {
	mockData := []byte{byte(libVersion)}

	kryptVersion, cipherType, err := getKryptInfo(mockData)
	if err == nil {
		t.Fatal("no error returned")
	}

	assert.Equal(t, uint8(0), kryptVersion, "krypt version mismatch")
	assert.Equal(t, CipherType(0), cipherType, "cipher type mismatch")
}
func TestNoPayload(t *testing.T) {
	mockData := mockKrypt(libVersion, AES256, []byte{})

	kryptVersion, cipherType, err := getKryptInfo(mockData)
	if err == nil {
		t.Fatal("no error returned")
	}

	assert.Equal(t, uint8(0), kryptVersion, "krypt version mismatch")
	assert.Equal(t, CipherType(0), cipherType, "cipher type mismatch")
}

func TestEmptyKrypt(t *testing.T) {
	mockData := []byte("")

	kryptVersion, cipherType, err := getKryptInfo(mockData)
	if err == nil {
		t.Fatal("no error returned")
	}

	assert.Equal(t, uint8(0), kryptVersion, "krypt version mismatch")
	assert.Equal(t, CipherType(0), cipherType, "cipher type mismatch")
}

func TestDataNotEncrypted(t *testing.T) {
	mockData := []byte("completely random data")
	kryptVersion, cipherType, err := getKryptInfo(mockData)
	if err == nil {
		t.Fatal("no error returned")
	}

	assert.Equal(t, uint8(0), kryptVersion, "krypt version mismatch")
	assert.Equal(t, CipherType(0), cipherType, "cipher type mismatch")
}

func TestGetCipherTypeByName(t *testing.T) {
	cryptType := GetCipherTypeByName("no-crypt")
	assert.Equal(t, Unknown, cryptType, "crypt types mismatch")

	cryptType = GetCipherTypeByName("AES256")
	assert.Equal(t, AES256, cryptType, "crypt types mismatch")
}

func mockKrypt(kryptVersion uint8, cipherType CipherType, data []byte) []byte {
	var buf []byte
	buf = append(buf, byte(kryptVersion), byte(cipherType))
	buf = append(buf, data...)
	return buf
}
