package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)
var Db *sql.DB

func InitDatabase() error {
    var err error
    Db, err = sql.Open("sqlite3", "./wallet.db")
    if err != nil {
        return err
    }
    // Verify the connection works
    err = Db.Ping()
    if err != nil {
		return err
    }
    return nil
}

func CloseDb()  {
    if Db != nil{
        Db.Close()
    }
}

func CreateTable() (sql.Result, error){
    query := `
    CREATE TABLE IF NOT EXISTS wallet (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`
    
    res, err := Db.Exec(query)
    if err != nil {
        return nil, err
    }
    return res, nil
}

func InsertWallet(wallet Wallet) (int64, error) {
    query := `INSERT INTO wallet (username, password) VALUES (?, ?)`
    secret := os.Getenv("WALLET_SECRET")
    if secret == "" {
        return 0, errors.New("WALLET_SECRET env variable not set")
    }
    encryptedPassword, err := Encrypt(wallet.password, secret)
    if err != nil {
        return 0, err
    }
    result, err := Db.Exec(query, wallet.username, encryptedPassword)
    if err != nil {
        return 0, err
    }
    return result.LastInsertId()
}

func QueryWallet(username string) (string, error) {
    var wallet Wallet
    query:= `SELECT username, password FROM wallet WHERE username = ?`
    err := Db.QueryRow(query, username).Scan(&wallet.username, &wallet.password)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Printf("username '%s' not found!", username)
            return "", nil
        }
        return "", err
    }
    text, err := wallet.ToString()
    if err != nil {
        return "", err
    }
    return text, nil
}

func UpdateWalletPassword(username string) (*Wallet, error) {
    var exists bool
    query:= `SELECT Exists (SELECT 1 FROM wallet WHERE username = ?)`
    err := Db.QueryRow(query, username).Scan(&exists)
    if err != nil {
        return nil, err
    }
    if !exists {
        log.Printf("username '%s' not found", username)
        return nil, nil
    }
    wallet, err := New(username)
    if err != nil {
        return nil,err
    }
    _,err = Db.Exec("UPDATE wallet SET username = ?, password = ? WHERE username = ?", wallet.username, wallet.password, wallet.username )
        if err != nil {
        return nil,err
    }
    return wallet, nil
}

func RemoveWallet(w *Wallet) error {
    _,err := Db.Exec("DELETE FROM wallet WHERE username = ?", w.username )
    return err
}
