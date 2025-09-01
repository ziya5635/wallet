package main

import (
	"fmt"
	"os"
)

type Wallet struct{
	username string
	password string
}

func NewWallet(username string) (*Wallet, error)  {
	var w Wallet
	key := os.Getenv("WALLET_SECRET")
	password, err := GenRandomPassword()
		if err!=nil {
		return nil, err
	}
	encryptedPassword, err := Encrypt(password, key)
	if err!=nil {
		return nil, err
	}
	w.username = username
	w.password = encryptedPassword
	return &w, nil
}

func (w *Wallet) ToString()(string, error){
	decrypted, err := Decrypt(w.password, os.Getenv("WALLET_SECRET"))
    if err != nil {
        return "", err
    }
	return fmt.Sprintf("username:%s | password:%s", w.username, decrypted), nil
}

func (w *Wallet) Save() error {
	_, err:=InsertWallet(w)
	    if err != nil {
		return err
    }
	return nil
}

func UpdatePassword(username string) (*Wallet, error) {
	wallet, err := UpdateWalletPassword(username)
	return wallet,err
}

func (w *Wallet) Remove() error {
	err := RemoveWallet(w)
	return err
}