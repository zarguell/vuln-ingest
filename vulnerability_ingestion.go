package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	RequireToken bool   `json:"require_token"`
	Port         string `json:"port"`
}

type Vulnerability struct {
	Title       string `json:"title"`
	Cve         string `json:"cve"`
	Cwe         string `json:"cwe"`
	Evidence    string `json:"evidence"`
	Found       string `json:"found"`
	Description string `json:"description"`
}

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "config.json", "Path to config file")
	flag.Parse()

	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config Config
	err = json.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	db, err := sql.Open("sqlite3", "vulnerabilities.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	tokendb, err := sql.Open("sqlite3", "tokens.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer tokendb.Close()

	http.HandleFunc("/ingest", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if config.RequireToken {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Authorization token required", http.StatusUnauthorized)
				return
			}

			var hash []byte
			err := tokendb.QueryRow("SELECT hash FROM tokens WHERE token = ?", token).Scan(&hash)
			if err != nil {
				http.Error(w, "Token not found", http.StatusUnauthorized)
				return
			}

			err = bcrypt.CompareHashAndPassword(hash, []byte(token))
						if err != nil {
				http.Error(w, "Token not valid", http.StatusUnauthorized)
				return
			}
		}

		var vulnerability Vulnerability
		err = json.NewDecoder(r.Body).Decode(&vulnerability)
		if err != nil {
			http.Error(w, "Failed to parse JSON body", http.StatusBadRequest)
			return
		}

		stmt, err := db.Prepare("INSERT INTO vulnerabilities(title, cve, cwe, evidence, found, description, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			http.Error(w, "Failed to prepare SQL statement", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(vulnerability.Title, vulnerability.Cve, vulnerability.Cwe, vulnerability.Evidence, vulnerability.Found, vulnerability.Description, time.Now().UTC().String())
		if err != nil {
			http.Error(w, "Failed to execute SQL statement", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	log.Printf("Starting server on port %s", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Port), nil))
}

				
