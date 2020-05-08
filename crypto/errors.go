package crypto

// DataIsEncryptedError when trying to encrypt already encrypted data
type DataIsEncryptedError struct {
	msg string // description of error
}

func (e *DataIsEncryptedError) Error() string { return e.msg }

// NewDataIsEcryptedError returns a new error
func NewDataIsEcryptedError() *DataIsEncryptedError {
	return &DataIsEncryptedError{"data is already encrypted"}
}

// DataIsNotEncryptedError when trying to decrypt data that is not encrypted
type DataIsNotEncryptedError struct {
	msg string // description of error
}

func (e *DataIsNotEncryptedError) Error() string { return e.msg }

// NewDataIsNotEncryptedError returns a new error
func NewDataIsNotEncryptedError() *DataIsNotEncryptedError {
	return &DataIsNotEncryptedError{"data is not encrypted"}
}

// UnknownCipherNameError when cipher name is not known
type UnknownCipherNameError struct {
	msg string // description of error
}

func (e *UnknownCipherNameError) Error() string { return e.msg }

// NewUnknownCipherNameError returns a new error
func NewUnknownCipherNameError() *UnknownCipherNameError {
	return &UnknownCipherNameError{"cipher name not recognized"}
}

// UnknownCipherTypeError when cipher name is not known
type UnknownCipherTypeError struct {
	msg string // description of error
}

func (e *UnknownCipherTypeError) Error() string { return e.msg }

// NewUnknownCipherTypeError returns a new error
func NewUnknownCipherTypeError() *UnknownCipherTypeError {
	return &UnknownCipherTypeError{"cipher type not recognized"}
}
