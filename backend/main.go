package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

func updateBarang(w http.ResponseWriter, r *http.Request) {
	var barang Barang
	if err := json.NewDecoder(r.Body).Decode(&barang); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if barang.ID == 0 {
		http.Error(w, "Missing id in request", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("UPDATE barang SET nama = $1, harga = $2, stok = $3 WHERE id = $4", barang.Nama, barang.Harga, barang.Stok, barang.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Barang not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(barang)
}

func deleteBarang(w http.ResponseWriter, r *http.Request) {
	// Mengambil id
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("DELETE FROM barang WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Barang not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	initDB()
	http.HandleFunc("/barang", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getBarang(w, r)
		case http.MethodPost:
			createBarang(w, r)
		case http.MethodPut:
			updateBarang(w, r)
		case http.MethodDelete:
			deleteBarang(w, r)
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
