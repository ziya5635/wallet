package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

func Hash(plainText , secret string) string {
    h := hmac.New(sha256.New, []byte(secret))
    h.Write([]byte(plainText))
    return hex.EncodeToString(h.Sum(nil))
}

func Encrypt(plainText, key string)(string, error)  {
    block,err:=aes.NewCipher([]byte(key))
    if err != nil {
        return "", err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    nonce := make([]byte, gcm.NonceSize())
    io.ReadFull(rand.Reader, nonce)
    cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
    return hex.EncodeToString(cipherText), nil
 }

 func Decrypt(encrypted, key string)(string, error)  {
    data, err := hex.DecodeString(encrypted)
    if err != nil{
        return "", err
    }
    block, err := aes.NewCipher([]byte(key))
    if err != nil {
        return "", err
    }
    gcm, err := cipher.NewGCM(block)
        if err != nil {
        return "", err
    }
    nonceSize := gcm.NonceSize()
    nonce, cipherText := data[:nonceSize], data[nonceSize:]
    plainText, err := gcm.Open(nil, nonce, cipherText, nil)
    return string(plainText), err
 }