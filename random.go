package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

func RandomInt(max int) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, errors.New("unable to make random Int")
	}
	return int(n.Int64()), nil
}

func GenRandomPassword() (string, error) {
	var chars = "abcdefghijklmnopqrstuvwxyz!@#$%^&*()+1234567890_ABCSDEFGHJKQLZWERTYUIOP{XCVNM}"
	var password string
	for i := 0; i < 5; i++ {
		for j := 0; j < 4; j++ {
			length:=strings.Count(chars, "")
			index, err := RandomInt(length-1)
			if err != nil {
				fmt.Println(err)
				return "", errors.New("unable to create a new password")
			}
			password += string(chars[index])
		}
		if i != 4 {
			password += "-"
		}
	}
	return password, nil
}