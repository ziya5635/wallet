package main

import (
	"database/sql"
	"errors"
	"fmt"
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

func InsertRow(wallet Wallet) (int64, error) {
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

func QueryWallet(username string) (*Wallet, error) {
    var wallet Wallet
    query:= `SELECT username, password FROM wallet WHERE username = ?`
    err := Db.QueryRow(query, username).Scan(&wallet.username, &wallet.password)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("username '%s' not found", username)
        }
        return nil, err
    }
    decrypted, err := Decrypt(wallet.password, os.Getenv("WALLET_SECRET"))
    if err != nil {
        return nil, err
    }
    wallet.password = decrypted
    return &wallet, nil
}
