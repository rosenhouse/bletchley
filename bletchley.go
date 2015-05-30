// Package bletchly is a simple, high-level library for asymmetric encryption and decryption.
package bletchley

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"hash"
	"io"
)

var (
	_hash        hash.Hash = sha256.New()
	randomReader io.Reader = rand.Reader
)

func randomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

// EncryptedMessage is an encrypted (but not authenticated) representation of a plaintext message.
// The consumer of this package should not need to understand or manipulate the fields except for serialization.
// Decryption requires possession of the private key.
type EncryptedMessage struct {
	Ciphertext   []byte `json:"ciphertext"`
	EncryptedKey []byte `json:"encrypted_key"`
}

// Encrypt encrypts a given plaintext using the provided public key.
// The encryption process uses random data. Therefore, this function is not deterministic.
func Encrypt(publicKey *rsa.PublicKey, plaintext []byte) (EncryptedMessage, error) {
	aesKey, err := randomBytes(symmetricKeyLength)
	if err != nil {
		return EncryptedMessage{}, err
	}

	ciphertext, err := symmetricEncrypt(aesKey, plaintext)
	if err != nil {
		return EncryptedMessage{}, err
	}

	encryptedKey, err := rsa.EncryptOAEP(_hash, randomReader, publicKey, aesKey, nil)
	if err != nil {
		return EncryptedMessage{}, err
	}

	return EncryptedMessage{
		Ciphertext:   ciphertext,
		EncryptedKey: encryptedKey,
	}, nil
}

// Decrypt decrypts a given EncryptedMessage using the provided private key.
// If the provided key is invalid then Decrypt will return an empty slice and an error.
// Decrypt does not validate the authenticity of the encrypted message.
func Decrypt(privateKey *rsa.PrivateKey, msg EncryptedMessage) ([]byte, error) {
	aesKey, err := rsa.DecryptOAEP(_hash, randomReader, privateKey, msg.EncryptedKey, nil)
	if err != nil {
		return []byte{}, err
	}

	return symmetricDecrypt(aesKey, msg.Ciphertext)
}
