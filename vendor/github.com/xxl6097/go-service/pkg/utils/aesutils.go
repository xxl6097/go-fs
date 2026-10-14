package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
)

var (
	ErrEntityInvalid      = errors.New(`common.ENTITY_INVALID`)
	ErrFailedVerification = errors.New(`common.ENTITY_CHECK_FAILED`)
)
var KEY = []byte("0123456789abcdef0123456789abcdef")

func EncAES(data []byte, key []byte) ([]byte, error) {
	hash, _ := GetMD5(data)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(block, hash)
	encBuffer := make([]byte, len(data))
	stream.XORKeyStream(encBuffer, data)
	return append(hash, encBuffer...), nil
}

func DecAES(data []byte, key []byte) ([]byte, error) {
	// MD5[16 bytes] + Data[n bytes]
	dataLen := len(data)
	if dataLen <= 16 {
		return nil, ErrEntityInvalid
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(block, data[:16])
	decBuffer := make([]byte, dataLen-16)
	stream.XORKeyStream(decBuffer, data[16:])
	hash, _ := GetMD5(decBuffer)
	if !bytes.Equal(hash, data[:16]) {
		return nil, ErrFailedVerification
	}
	return decBuffer[:dataLen-16], nil
}
func GenRandByte(n int) []byte {
	secBuffer := make([]byte, n)
	rand.Reader.Read(secBuffer)
	return secBuffer
}
func GetUUID() []byte {
	return GenRandByte(16)
}

func GetMD5(data []byte) ([]byte, string) {
	hash := md5.New()
	hash.Write(data)
	result := hash.Sum(nil)
	hash.Reset()
	return result, hex.EncodeToString(result)
}

// AES‑256‑GCM 加密：返回 nonce(12) + cipher+tag
func Aes256GcmEncrypt(plain []byte, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("key must be 32 bytes for aes‑256")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plain, nil), nil
}

func Aes256GcmDecrypt(cipherData []byte, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("key must be 32 bytes for aes‑256")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(cipherData) < nonceSize+16 {
		return nil, fmt.Errorf("cipher data too short")
	}
	nonce, ct := cipherData[:nonceSize], cipherData[nonceSize:]
	return gcm.Open(nil, nonce, ct, nil)
}
