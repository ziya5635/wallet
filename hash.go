package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

func Hash(plainText , secret string) string {
    h := hmac.New(sha256.New, []byte(secret))
    h.Write([]byte(plainText))
    return hex.EncodeToString(h.Sum(nil))
}

func parseKey(key string) ([]byte, error) {
    keyBytes, err := base64.StdEncoding.DecodeString(key)
    if err != nil {
        return nil, fmt.Errorf("failed to base64 decode key: %w", err)
    }
    if len(keyBytes) != 32 {
        return nil, fmt.Errorf("key must be 32 bytes after base64 decode, got %d", len(keyBytes))
    }
    return keyBytes, nil
}

func Encrypt(plainText, key string) (string, error) {
    keyBytes, err := parseKey(key)
    if err != nil {
        return "", err
    }

    block, err := aes.NewCipher(keyBytes)
    if err != nil {
        return "", err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }

    cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
    return hex.EncodeToString(cipherText), nil
}

func Decrypt(encrypted, key string) (string, error) {
    keyBytes, err := parseKey(key)
    if err != nil {
        return "", err
    }

    data, err := hex.DecodeString(encrypted)
    if err != nil {
        return "", err
    }
    block, err := aes.NewCipher(keyBytes)
    if err != nil {
        return "", err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonceSize := gcm.NonceSize()
    if len(data) < nonceSize {
        return "", fmt.Errorf("ciphertext too short")
    }

    nonce, cipherText := data[:nonceSize], data[nonceSize:]
    plainText, err := gcm.Open(nil, nonce, cipherText, nil)
    if err != nil {
        return "", err
    }
    return string(plainText), nil
}