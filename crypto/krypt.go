package crypto

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"

	"github.com/pkg/errors"
)

const name = "krypt"
const libVersion = uint8(1)

const defaultSaltSize = 12

// CipherType is the cipher type
type CipherType uint8

// cipher types
const (
	Unknown CipherType = iota
	AES256
	TWOFISH
	SERPENT
)

// TODO: add secretbox https://godoc.org/golang.org/x/crypto/nacl/secretbox

// TODO: find golang implementations of MARS & RC6

// For CDC mode ops
// TODO: add blowfish https://godoc.org/golang.org/x/crypto/blowfish
// TODO: add tripledes https://golang.org/pkg/crypto/des/
// 						https://gist.github.com/cuixin/10612934

// Cipher interface represents a en/decrypting module
type Cipher interface {
	Encrypt(data []byte, password []byte) ([]byte, error)
	Decrypt(data []byte, password []byte) ([]byte, error)
	GetDescription() string
	GetName() string
	GetType() CipherType
}

// getCipher gets the cipher object type
func getCipher(cipherType CipherType) (cipher Cipher, err error) {
	switch cipherType {
	case AES256:
		cipher = NewAES256Cipher()
	case TWOFISH:
		cipher = NewTwofishCipher()
	case SERPENT:
		cipher = NewSerpentCipher()
	default:
		err = NewUnknownCipherTypeError()
	}
	return
}

// getCipherByName gets the cipher by name
func getCipherByName(name string) (cipher Cipher, err error) {
	switch name {
	case aes256Name:
		cipher = NewAES256Cipher()
	case twofishName:
		cipher = NewTwofishCipher()
	case serpentName:
		cipher = NewSerpentCipher()
	default:
		err = NewUnknownCipherNameError()
	}
	return
}

// GetCipherTypeByName gets the CipherType by name
func GetCipherTypeByName(name string) (c CipherType, err error) {
	//TODO: make this not care about whitespace or letter case
	switch strings.TrimSpace(name) {
	case aes256Name:
		c = AES256
	case twofishName:
		c = TWOFISH
	case serpentName:
		c = SERPENT
	default:
		err = NewUnknownCipherNameError()
	}
	return
}

// GetCipherList returns a list of strings
func GetCipherList() []Cipher {
	cipherList := []Cipher{}
	cipherList = append(cipherList, NewAES256Cipher())
	cipherList = append(cipherList, NewTwofishCipher())
	cipherList = append(cipherList, NewSerpentCipher())
	return cipherList
}

// Encrypt data with password in the given CryptType format
func Encrypt(cipherType CipherType, password []byte, data []byte) ([]byte, error) {
	cipher, cerr := getCipher(cipherType)
	if cerr != nil {
		return nil, cerr
	}

	cipherText, err := cipher.Encrypt(data, password)
	if err != nil {
		return nil, err
	}

	kryptData := writeKrypt(cipherType, cipherText)

	return kryptData, nil
}

// Decrypt data block with the given password, encryption type
// 	is derived from data block metadata
func Decrypt(password []byte, data []byte) ([]byte, error) {
	kryptVersion, cipherType, kerr := getKryptInfo(data)
	if kerr != nil {
		return nil, errors.Wrap(kerr, "reading krypt")
	}

	cipher, cerr := getCipher(cipherType)
	if cerr != nil {
		return nil, cerr
	}

	if kryptVersion != 1 {
		return nil, errors.New("unknown krypt version")
	}

	payloadLen := len(data) - 2
	if payloadLen == 0 {
		return nil, errors.New("no payload found")
	}
	reader := bytes.NewReader(data[2:])
	payload := make([]byte, payloadLen)
	if _, err := io.ReadFull(reader, payload); err != nil {
		return nil, errors.Wrap(err, "reading payload")
	}

	plainText, err := cipher.Decrypt(payload, password)
	if err != nil {
		return nil, errors.Wrapf(err, "decrpyting payload")
	}

	return plainText, nil
}

// Write the data to a byte stream
func writeKrypt(cipherType CipherType, data []byte) []byte {
	// var buffer bytes.Buffer
	buffer := new(bytes.Buffer)

	binary.Write(buffer, binary.LittleEndian, libVersion)
	binary.Write(buffer, binary.LittleEndian, cipherType)

	buffer.Write(data)

	return buffer.Bytes()
}

func getKryptInfo(data []byte) (uint8, CipherType, error) {
	dataLen := len(data)
	if dataLen <= 2 {
		return 0, 0, errors.New("krypt info missing")
	}

	reader := bytes.NewReader(data)

	var foundVersion uint8
	if err := binary.Read(reader, binary.LittleEndian, &foundVersion); err != nil {
		return 0, 0, errors.Wrap(err, "reading krypt version")
	}
	if foundVersion != libVersion {
		return 0, 0, errors.New("cannot read krypt info")
	}

	var foundCipherType uint8
	if err := binary.Read(reader, binary.LittleEndian, &foundCipherType); err != nil {
		return 0, 0, errors.Wrap(err, "reading krypt cipher")
	}
	if _, err := getCipher(CipherType(foundCipherType)); err != nil {
		return 0, 0, errors.Wrap(err, "cannot determine cipher used")
	}
	return foundVersion, CipherType(foundCipherType), nil
}
