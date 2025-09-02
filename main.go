package main

//To build for raspberry5:
//GOOS=linux GOARCH=arm64 go build -o wallet.bin -ldflags="-s -w" .

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func main()  {
	err := InitDatabase()
	reportError("Failed to initialize database,", err)

	defer CloseDb()

	_ ,err = CreateWalletTable()
	reportError("Failed to initialize wallet table,", err)

	_, err = CreateUserTable()
	reportError("Failed to initialize user table", err)

	doesAnyUserExist, err := CheckAnyUserExists()
	if err != nil {
		CloseDb() // ← Manual cleanup before exiting
    	os.Exit(1)
	}
	if !doesAnyUserExist {
		fmt.Println("====================================================================================")
		log.Println("initial setup, you won't be asked to go through this step again")
		fmt.Println("====================================================================================")
		name, err := getUserInput("admin user name:")
		if err != nil {
			reportError("Unable to get first admin user", err)
		}
		fmt.Print("admin password:")
		password, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			reportError("Unable to get first admin password", err)
		}
		user, err := NewUser(name, string(password))
		if err != nil {
			reportError("Unable to create new user", err)
		}
		err = user.Save()
		if err != nil {
			reportError("Unable to save new user into db", err)
		}

	}else {
		for {
		name, err := getUserInput("admin user name:")
		if err != nil {
			log.Println("Unable to get user's name")
			continue
		}
		fmt.Print("admin password:")
		password, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Println("Unable to get user's password")
			continue
		}
		isValid, err := LoginUser(name, string(password))
		if err != nil {
			log.Println("Unable to login, Invalid credentials given")
			continue
		}
		if isValid {
			break
		}
		log.Println("Unable to login, Invalid credentials given")
	}
	}
	for {
		userInput, err := outputServices()
		if err != nil {
			reportError("oops, an error thrown,", err)
			continue
		}

		err = runService(userInput)
		reportError("oops, an error thrown,", err)
	}
}

func outputServices() (string, error)  {
	var userInput string
	fmt.Println("1.Create a Wallet")
	fmt.Println("2.Restore a Wallet")
	fmt.Println("3.Update a Wallet")
	fmt.Println("4.Remove a wallet")
	fmt.Println("5.Print all wallets")
	fmt.Println("6.Exit")
	fmt.Print("Your choice:")
	_, err := fmt.Scanln(&userInput)
	if err != nil {
		return "0", err
	}
	return userInput, nil
}

func runService(userInput string) error  {
	switch userInput {
	case "1":
		userInput, err:=getUserInput("username:")
		if err != nil {
			return err
		}
		wallet, err := NewWallet(userInput)
		if err != nil {
			return err
		}
		 err = wallet.Save()
		if err != nil {
			return err
		}
		fmt.Printf("Credential saved \n")
		walletText, err := wallet.ToString()
		if err != nil {
			return err
		}
		log.Println(walletText)
		return nil
	case "2":
		userInput, err := getUserInput("username to query:")
		if err != nil {
			return err
		}
		walletText, err := QueryWallet(userInput)
		if err != nil {
			return err
		}
		log.Println(walletText)
		return nil
	case "3":
		username, err := getUserInput("username:")
		if err != nil {
			return err
		}
		oldWallet, err := QueryWallet(username)
		if err!=nil {
			return err
		}
		log.Printf(`Old wallet info: %s`, oldWallet)
		wallet ,err := UpdatePassword(username)
		if err != nil {
			return err
		}
		if wallet == nil {
			return nil
		}
		walletText, err := wallet.ToString()
		if err != nil {
			return err
		}
		fmt.Println("========================================================================================")
		log.Printf(`New wallet info: %s`, walletText)
		fmt.Println("========================================================================================")
		return nil
	case "4":
		var wallet Wallet
		username, err := getUserInput("username to remove:")
		wallet.username = username
		if err != nil {
			return err
		}
		err = wallet.Remove()
		if err != nil {
			return err
		}
		log.Printf("%s deleted successfully.", username)
		return nil
	case "5":
		wallets, err := GetAllWallets()
		if err != nil {
			return err
		}
    for i, wallet := range wallets {
		fmt.Println("======================")
        fmt.Printf("%d. Username: %s\n", 
            i+1, wallet.username)
		fmt.Println("======================")
    }
	return nil
	case "6":
		return fmt.Errorf("exit requested")
	default:
		return fmt.Errorf("invalid user input given: %s", userInput)
	}
}

func reportError(message string, err error) {
	if err != nil {
		if err.Error() == "exit requested"{
			log.Print("Goodbye!")
			CloseDb() // ← Manual cleanup before exiting
        	os.Exit(0)
    	}
	    log.Print(message, err)
		// CloseDb() // ← Manual cleanup before exiting
    	// os.Exit(1)
	}
}

func getUserInput(text string) (string, error)  {
	var userInput string
	fmt.Print(text)
	_, err := fmt.Scanln(&userInput)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(userInput) == "" {
		return "", errors.New("invalid input given")
	}
	return  userInput,nil
}