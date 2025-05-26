package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Dessert struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ImageURL    string `json:"image_url"`
	Ingredients string `json:"ingredients"`
	BakeTime    string `json:"bake_time"`
	RecipeLink  string `json:"recipe_link"`
}

var db *sql.DB

func initDB() {
	var err error
	dbConn := os.Getenv("DESSERTS_DB_CONN")
	if dbConn == "" {
		log.Fatal("DESSERTS_DB_CONN environment variable is not set")
	}

	// Retry connecting to DB up to 5 times with 2 second intervals
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", dbConn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("Connected to database successfully")
				break
			}
		}
		log.Println("Waiting for DB to be ready... retrying in 2 seconds")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Database connection failed after retries:", err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS desserts (
		id SERIAL PRIMARY KEY,
		name TEXT,
		image_url TEXT,
		ingredients TEXT,
		bake_time TEXT,
		recipe_link TEXT
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
}

func getDessertsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rows, err := db.Query("SELECT id, name, image_url, ingredients, bake_time, recipe_link FROM desserts")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var desserts []Dessert

	for rows.Next() {
		var d Dessert
		err := rows.Scan(&d.ID, &d.Name, &d.ImageURL, &d.Ingredients, &d.BakeTime, &d.RecipeLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		desserts = append(desserts, d)
	}
	json.NewEncoder(w).Encode(desserts)
}

func main() {
	initDB()

	http.HandleFunc("/api/desserts", getDessertsHandler)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
