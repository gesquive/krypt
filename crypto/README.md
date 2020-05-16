# crypto

This library converts bytes into an encrypted binary format and back.

crypto uses the (GCM mode of operation)[https://en.wikipedia.org/wiki/Galois/Counter_Mode] with the specified block cipher to create cipher text that is then packaged in a binary file format.

Keys are derived from the given password using HMAC-SHA-256 based PBKDF2 key derivation function.

Supported Ciphers
 - AES256 (default)
 - Twofish
 - Serpent

The binary format is meant to be as efficient as possible, and thus minimally invasive
