package crypto

// DataIsEncryptedError when trying to encrypt already encrypted data
type DataIsEncryptedError struct {
	msg string // description of error
}

func (e *DataIsEncryptedError) Error() string { return e.msg }

// NewDataIsEcryptedError returns a new error
func NewDataIsEcryptedError() *DataIsEncryptedError {
	return &DataIsEncryptedError{"Data is already encrypted"}
}

// DataIsNotEncryptedError when trying to decrypt data that is not encrypted
type DataIsNotEncryptedError struct {
	msg string // description of error
}

func (e *DataIsNotEncryptedError) Error() string { return e.msg }

// NewDataIsNotEncryptedError returns a new error
func NewDataIsNotEncryptedError() *DataIsNotEncryptedError {
	return &DataIsNotEncryptedError{"Data is not encrypted"}
}
