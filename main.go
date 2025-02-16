package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Barang struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

var db *sql.DB

func initDB() {
	var err error
	dsn := "host=localhost user=postgres password=superuser dbname=tokogw port=5432 sslmode=disable"
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	log.Println("Connected to database successfully")
}

func getBarang(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, nama, harga, stok FROM barang")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var barangList []Barang
	for rows.Next() {
		var barang Barang
		if err := rows.Scan(&barang.ID, &barang.Nama, &barang.Harga, &barang.Stok); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		barangList = append(barangList, barang)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(barangList)
}

func createBarang(w http.ResponseWriter, r *http.Request) {
	var barang Barang
	if err := json.NewDecoder(r.Body).Decode(&barang); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err := db.Exec("INSERT into barang (nama, harga, stok) VALUES ($1, $2, $3)", barang.Nama, barang.Harga, barang.Stok)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(barang)
}

func main() {
	initDB()
	http.HandleFunc("/barang", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getBarang(w, r)
		case http.MethodPost:
			createBarang(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server running on port 5000")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal("Server failed to start", err)
	}
}
