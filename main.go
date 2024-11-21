package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Konfigurasi database
const (
	dsn  = "ux4zlfmbja6zvv0b:siOXUPKGQVLLZ86meGtT@tcp(bf4vkxzzfwc0oridaiqg-mysql.services.clever-cloud.com:3306)/bf4vkxzzfwc0oridaiqg"
	port = ":8380"
)

var db *sql.DB

func main() {
	var err error
	// Membuka koneksi ke database
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Gagal membuka koneksi ke database: %v\n", err)
	}
	defer db.Close()

	// Memastikan koneksi berhasil
	err = db.Ping()
	if err != nil {
		log.Fatalf("Gagal menghubungkan ke database: %v\n", err)
	}
	fmt.Println("Berhasil terhubung ke database!")

	// Route handler
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/submit", submitHandler)

	// Menjalankan server
	fmt.Printf("Server berjalan di http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// Menampilkan form
func formHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/form.html")
	if err != nil {
		http.Error(w, "Gagal memuat template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Menangani pengiriman form
func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Hanya metode POST yang diizinkan", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")

	if name == "" || email == "" {
		http.Error(w, "Semua field harus diisi", http.StatusBadRequest)
		return
	}

	// Menyimpan data ke database
	query := "INSERT INTO users (name, email) VALUES (?, ?)"
	_, err := db.Exec(query, name, email)
	if err != nil {
		http.Error(w, "Gagal menyimpan data", http.StatusInternalServerError)
		log.Printf("Error: %v\n", err)
		return
	}

	// Redirect ke halaman sukses atau form
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
