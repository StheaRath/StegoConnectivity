package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"math/big"

	"golang.org/x/crypto/pbkdf2"
)

const (
	SaltSize   = 16
	NonceSize  = 12
	Iterations = 100000
)

func deriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, Iterations, 32, sha256.New)
}

func Encrypt(plaintext []byte, password string) ([]byte, []byte, []byte, error) {
	salt := make([]byte, SaltSize)
	io.ReadFull(rand.Reader, salt)
	nonce := make([]byte, NonceSize)
	io.ReadFull(rand.Reader, nonce)
	key := deriveKey(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, nil, err
	}
	return gcm.Seal(nil, nonce, plaintext, nil), salt, nonce, nil
}

func Decrypt(ciphertext, salt, nonce []byte, password string) ([]byte, error) {
	key := deriveKey(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func GenerateRSA(bits int) (string, string, error) {
	if bits != 2048 && bits != 3072 && bits != 4096 {
		return "", "", errors.New("invalid RSA key size")
	}
	key, _ := rsa.GenerateKey(rand.Reader, bits)
	priv := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	pubBytes, _ := x509.MarshalPKIXPublicKey(&key.PublicKey)
	pub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
	return string(priv), string(pub), nil
}

func GenerateECC(curveName string) (string, string, error) {
	var c elliptic.Curve
	switch curveName {
	case "P-256":
		c = elliptic.P256()
	case "P-384":
		c = elliptic.P384()
	case "P-521":
		c = elliptic.P521()
	default:
		return "", "", errors.New("invalid curve")
	}
	key, _ := ecdsa.GenerateKey(c, rand.Reader)
	b, _ := x509.MarshalECPrivateKey(key)
	priv := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: b})
	pubBytes, _ := x509.MarshalPKIXPublicKey(&key.PublicKey)
	pub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
	return string(priv), string(pub), nil
}

func GenerateDHGroup(groupName string) (string, string, error) {
	var bitSize int
	switch groupName {
	case "Group 14":
		bitSize = 2048
	case "Group 15":
		bitSize = 3072
	case "Group 16":
		bitSize = 4096
	default:
		return "", "", errors.New("invalid DH group")
	}
	// Simulate DH Key pair generation
	privInt, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), uint(bitSize)))
	priv := pem.EncodeToMemory(&pem.Block{Type: "DH PRIVATE KEY", Bytes: privInt.Bytes()})
	pubHash := sha256.Sum256(privInt.Bytes())
	pub := pem.EncodeToMemory(&pem.Block{Type: "DH PUBLIC KEY", Bytes: pubHash[:]})
	return string(priv), string(pub), nil
}
