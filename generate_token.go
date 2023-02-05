package main

import (
	"crypto/rand"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/mattn/go-sqlite3"
)

func generateToken() ([]byte, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

func hashToken(token []byte) ([]byte, error) {
	hashed, err := bcrypt.GenerateFromPassword(token, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashed, nil
}

func main() {
	token, err := generateToken()
	if err != nil {
		log.Fatalf("Error generating token: %v", err)
	}
	hashed, err := hashToken(token)
	if err != nil {
		log.Fatalf("Error hashing token: %v", err)
	}

	db, err := sql.Open("sqlite3", "tokens.db")
	if err != nil {
		log.Fatalf("Error opening tokens database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS tokens (id INTEGER PRIMARY KEY AUTOINCREMENT, hash BLOB NOT NULL)")
	if err != nil {
		log.Fatalf("Error creating tokens table: %v", err)
	}

	_, err = db.Exec("INSERT INTO tokens (hash) VALUES (?)", hashed)
	if err != nil {
		log.Fatalf("Error inserting token into database: %v", err)
	}

	fmt.Printf("Generated API Token: %s\n", string(token))
}
