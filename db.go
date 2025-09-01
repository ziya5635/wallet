package main

import (
	"database/sql"
	"fmt"

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

func CreateWalletTable() (sql.Result, error){
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

func CreateUserTable() (sql.Result, error){
    query := `
    CREATE TABLE IF NOT EXISTS user (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`
    
    res, err := Db.Exec(query)
    if err != nil {
        return nil, err
    }
    return res, nil
}


func InsertUser(user *User) (int64, error)  {
    query := `INSERT INTO user (name, password) VALUES (?, ?)`
    result, err := Db.Exec(query, user.name, user.password)
    if err != nil {
        return 0, err
    }
    return result.LastInsertId()
}

func QueryUser(name string) (*User, error)  {
    var user User
    query:= `SELECT name, password FROM user WHERE name = ?`
    err := Db.QueryRow(query, name).Scan(&user.name, &user.password)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("username '%s' not found", name)
        }
        return nil, err
    }
    return &user, nil
}

func CheckAnyUserExists() (bool, error) {
    var count int
    err := Db.QueryRow("SELECT COUNT(*) FROM user").Scan(&count)
    if err != nil {
        return false, err
    }
    if count > 0 {
        return true, nil
    } else {
        fmt.Println("No users yet")
        return false, nil
    }
}

func InsertWallet(wallet *Wallet) (int64, error) {
    query := `INSERT INTO wallet (username, password) VALUES (?, ?)`
    result, err := Db.Exec(query, wallet.username, wallet.password)
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
            return "", fmt.Errorf("username '%s' not found", username)
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
        return nil, fmt.Errorf("username '%s' not found", username)
    }
    wallet, err := NewWallet(username)
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
    var exists bool
    query:= `SELECT Exists (SELECT 1 FROM wallet WHERE username = ?)`
    err := Db.QueryRow(query, w.username).Scan(&exists)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("username '%s' not found", w.username)
    }
    _,err = Db.Exec("DELETE FROM wallet WHERE username = ?", w.username )
    return err
}

func GetAllWallets() ([]Wallet, error) {
    rows, err := Db.Query("SELECT username, password FROM wallet")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var wallets []Wallet
    for rows.Next() {
        var wallet Wallet
        err := rows.Scan(&wallet.username, &wallet.password)
        if err != nil {
            return nil, err
        }
        wallets = append(wallets, wallet)
    }
    return wallets, rows.Err()
}