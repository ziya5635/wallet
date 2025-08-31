package main

//add auth to the whole db (one-way hash)?
//add admin table to enable admin user to authenticate before running services
import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func main()  {
	err := InitDatabase()
	reportError("Failed to initialize database,", err)

	defer CloseDb()

	_ ,err = CreateTable()
	reportError("Failed to initialize database,", err)

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
	fmt.Println("5.Show all wallets")
	fmt.Println("6.quit")
	fmt.Print("your choice:")
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
		wallet, err := New(userInput)
		if err != nil {
			return err
		}
		id, err := wallet.Save()
		if err != nil {
			return err
		}
		fmt.Printf("Credential saved with ID: %d\n", id)
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
		log.Print(walletText)
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
        fmt.Printf("%d. Username: %s\n", 
            i, wallet.username)
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