package main

//add auth to the whole db (one-way hash)?
//add changing password feature
import (
	"errors"
	"fmt"
	"log"
	"os"
)

func main()  {
	err := InitDatabase()
	reportErrorAndExit("Failed to initialize database,", err)

	defer CloseDb()

	_ ,err = CreateTable()
	reportErrorAndExit("Failed to initialize database,", err)

	for {
		userInput, err := outputServices()
		reportErrorAndExit("Invalid user input given,", err)

		err = runService(userInput)
		reportErrorAndExit("oops, an error thrown,", err)
	}
}

func outputServices() (int, error)  {
	var userInput int
	fmt.Println("1.Generate a new random password")
	fmt.Println("2.revive an old password")
	fmt.Println("3.Change a password")
	fmt.Println("4.quit")
	fmt.Print("your choice:")
	n, err := fmt.Scanln(&userInput)
	if err != nil {
		return 0, err
	}
	if n!=1 {
		return 0, errors.New("only one argument required")
	}
	return userInput, nil
}

func runService(pickedService int) error  {
	switch pickedService {
	case 1:
		var wallet Wallet
		userInput, err:=getUserInput("username:")
		if err!=nil {
			return err
		}
		password, err := GenRandomPassword()
		if err != nil {
			return err
		}
		wallet = Wallet{username: userInput, password: password}
		id, err := wallet.Save()
		if err != nil {
			return err
		}
		fmt.Printf("Credential saved with ID: %d\n", id)
		username, password := wallet.ToString()
		log.Println("Your username:", username)
		log.Println("Your password:", password)

		return nil
	case 2:
		userInput, err:=getUserInput("username to query:")
		if err!=nil {
			return err
		}
		wallet, err:=QueryWallet(userInput)
		if err!=nil {
			return err
		}
		username, password := wallet.ToString()
		log.Println("Your username:", username)
		log.Println("Your password:", password)
		return nil
	case 3:
		log.Print("Changing a password")
		return nil
	case 4:
		return fmt.Errorf("exit requested")
	default:
		return fmt.Errorf("invalid user input given: ('%d')", pickedService)
	}
}

func reportErrorAndExit(message string, err error) {
	if err != nil {
		if err.Error() == "exit requested"{
			log.Print("Goodbye!")
			CloseDb() // ← Manual cleanup before exiting
        	os.Exit(0)
    	}
	    log.Print(message, err)
		CloseDb() // ← Manual cleanup before exiting
    	os.Exit(1)	
	}
}

func getUserInput(text string) (string, error)  {
	var userInput string
	fmt.Print(text)
	n, err := fmt.Scanln(&userInput)
	if n!=1 {
		return "", errors.New("too many args given")
	}
	if err != nil {
		return "", err
	}
	return  userInput,nil
}