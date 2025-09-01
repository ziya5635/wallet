package main

import (
	"errors"
	"os"
)

type User struct{
	name string
	password string
}

func NewUser(name, password string) (*User,error) {
	secret := os.Getenv("WALLET_SECRET")
    if secret == "" {
        return nil, errors.New("WALLET_SECRET env variable not set")
    }
	encrypted, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	return &User{name: name, password: encrypted }, nil
}

func (u *User) Save() error {
	_, err := InsertUser(u)
	if err != nil {
		return err
	}
	return nil
}

func LoginUser(name,userPassword string) (bool, error) {
		user, err := QueryUser(name)
		if err != nil {
			return false, err
		}
	isValid := CheckPasswordHash(user.password, userPassword)
	return isValid, nil
}