package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// ====================== DATA MODELS ======================

type Meal struct {
	Name       string `json:"name"`
	Rating     int    `json:"rating"`
	PhotoLink  string `json:"image_url"`
	RecipeLink string `json:"recipe_link"`
	Source     string `json:"source"`
}

type Dessert struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ImageURL    string `json:"image_url"`
	Ingredients string `json:"ingredients"`
	BakeTime    string `json:"bake_time"`
	RecipeLink  string `json:"recipe_link"`
}

// ====================== DATABASE ======================

var db *sql.DB

func initDB() error {
	dbPath := filepath.Join(getExecutableDir(), "meals.db")
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for better performance
	_, err = db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		return fmt.Errorf("failed to set WAL mode: %w", err)
	}

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS HF_Meal (
			HFMealID INTEGER PRIMARY KEY AUTOINCREMENT,
			HFMealName TEXT NOT NULL,
			HFMealRating INTEGER NOT NULL DEFAULT 1,
			RecipeLink TEXT NOT NULL,
			PhotoLink TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS P3_Meal (
			P3MealID INTEGER PRIMARY KEY AUTOINCREMENT,
			P3MealName TEXT NOT NULL,
			P3MealRating INTEGER NOT NULL DEFAULT 1,
			RecipeLink TEXT NOT NULL,
			PhotoLink TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS desserts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			image_url TEXT,
			ingredients TEXT,
			bake_time TEXT,
			recipe_link TEXT
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	// Seed data if tables are empty
	seedData()

	return nil
}

func getExecutableDir() string {
	exePath, err := os.Executable()
	if err == nil {
		return filepath.Dir(exePath)
	}
	return "."
}

// findInPath searches for a file in multiple possible locations
func findInPath(filename string) string {
	candidates := []string{
		// Direct copy next to executable
		filepath.Join(getExecutableDir(), filename),
		// Relative to executable (macOS app bundle structure)
		filepath.Join(getExecutableDir(), "..", "..", "Frontend", "data", filename),
		filepath.Join(getExecutableDir(), "..", filename),
		// Relative to current directory
		filename,
		"Frontend/data/" + filename,
		"data/" + filename,
		"Frontend/data/" + filename,
	}

	// Also check parent of executable dir (Contents/MacOS -> Contents -> root)
	if exePath, err := os.Executable(); err == nil {
		appDir := filepath.Dir(exePath)
		parentDir := filepath.Dir(appDir)
		greatGrandparentDir := filepath.Dir(parentDir)
		
		candidates = append(candidates,
			filepath.Join(appDir, "Frontend", "data", filename),
			filepath.Join(appDir, filename),
			filepath.Join(parentDir, "Frontend", "data", filename),
			filepath.Join(parentDir, filename),
			filepath.Join(greatGrandparentDir, "Frontend", "data", filename),
			filepath.Join(greatGrandparentDir, filename),
			filepath.Join(appDir, "data", filename),
			filepath.Join(parentDir, "data", filename),
		)
	}

	for _, path := range candidates {
		// Clean the path
		cleanPath := filepath.Clean(path)
		if _, err := os.Stat(cleanPath); err == nil {
			return cleanPath
		}
	}

	return ""
}

// ====================== DATA SEEDING ======================

func seedData() {
	// Count existing meals
	var hfCount, p3Count, dessertCount int
	db.QueryRow("SELECT COUNT(*) FROM HF_Meal").Scan(&hfCount)
	db.QueryRow("SELECT COUNT(*) FROM P3_Meal").Scan(&p3Count)
	db.QueryRow("SELECT COUNT(*) FROM desserts").Scan(&dessertCount)

	if hfCount == 0 && p3Count == 0 {
		log.Println("Seeding meals data from JSON file...")
		seedMeals()
	}
	if dessertCount == 0 {
		log.Println("Seeding desserts data from JSON file...")
		seedDesserts()
	}
}

func seedMeals() {
	mealsPath := findInPath("meals.json")
	if mealsPath == "" {
		log.Println("Warning: Could not find meals.json - using empty meal database")
		return
	}

	data, err := os.ReadFile(mealsPath)
	if err != nil {
		log.Printf("Warning: Could not read meals JSON: %v", err)
		return
	}

	var meals []Meal
	if err := json.Unmarshal(data, &meals); err != nil {
		log.Printf("Warning: Could not parse meals JSON: %v", err)
		return
	}

	tx, _ := db.Begin()
	hfStmt, _ := tx.Prepare(`INSERT INTO HF_Meal (HFMealName, HFMealRating, RecipeLink, PhotoLink) VALUES (?, ?, ?, ?)`)
	p3Stmt, _ := tx.Prepare(`INSERT INTO P3_Meal (P3MealName, P3MealRating, RecipeLink, PhotoLink) VALUES (?, ?, ?, ?)`)
	defer hfStmt.Close()
	defer p3Stmt.Close()

	hfCount := 0
	p3Count := 0
	for _, m := range meals {
		if m.Source == "HF_Meal" {
			hfStmt.Exec(m.Name, m.Rating, m.RecipeLink, m.PhotoLink)
			hfCount++
		} else if m.Source == "P3_Meal" {
			p3Stmt.Exec(m.Name, m.Rating, m.RecipeLink, m.PhotoLink)
			p3Count++
		}
	}
	tx.Commit()
	log.Printf("Seeded %d HelloFresh meals and %d Paprika3 meals", hfCount, p3Count)
}

func seedDesserts() {
	dessertsPath := findInPath("desserts.json")
	if dessertsPath == "" {
		log.Println("Warning: Could not find desserts.json - desserts will be empty")
		return
	}

	data, err := os.ReadFile(dessertsPath)
	if err != nil {
		log.Printf("Warning: Could not read desserts JSON: %v", err)
		return
	}

	var desserts []Dessert
	if err := json.Unmarshal(data, &desserts); err != nil {
		log.Printf("Warning: Could not parse desserts JSON: %v", err)
		return
	}

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(`INSERT OR REPLACE INTO desserts (id, name, image_url, ingredients, bake_time, recipe_link) VALUES (?, ?, ?, ?, ?, ?)`)
	defer stmt.Close()

	for _, d := range desserts {
		stmt.Exec(d.ID, d.Name, d.ImageURL, d.Ingredients, d.BakeTime, d.RecipeLink)
	}
	tx.Commit()
	log.Printf("Seeded %d desserts", len(desserts))
}

// ====================== MEAL API ======================

func getMeals(table string) ([]Meal, error) {
	var query string
	switch table {
	case "HF_Meal":
		query = `SELECT HFMealName, HFMealRating, PhotoLink, RecipeLink FROM HF_Meal`
	case "P3_Meal":
		query = `SELECT P3MealName, P3MealRating, PhotoLink, RecipeLink FROM P3_Meal`
	default:
		return nil, fmt.Errorf("invalid table name: %s", table)
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meals []Meal
	for rows.Next() {
		var meal Meal
		if err := rows.Scan(&meal.Name, &meal.Rating, &meal.PhotoLink, &meal.RecipeLink); err != nil {
			return nil, err
		}
		meal.Source = table
		meals = append(meals, meal)
	}
	return meals, nil
}

func weightedSelection(meals []Meal, numMeals int) []Meal {
	if len(meals) == 0 {
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	selected := make(map[string]bool)
	var result []Meal

	// Weight by rating
	weighted := []Meal{}
	for _, meal := range meals {
		for i := 0; i < meal.Rating; i++ {
			weighted = append(weighted, meal)
		}
	}

	// Shuffle
	for i := len(weighted) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		weighted[i], weighted[j] = weighted[j], weighted[i]
	}

	for _, meal := range weighted {
		if len(result) >= numMeals {
			break
		}
		if !selected[meal.Name] {
			result = append(result, meal)
			selected[meal.Name] = true
		}
	}

	return result
}

func parseTypeCounts(s string) (map[string]int, error) {
	result := make(map[string]int)
	pairs := strings.Split(s, ",")
	for _, pair := range pairs {
		parts := strings.Split(strings.TrimSpace(pair), ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid type format: %s", pair)
		}
		table := strings.TrimSpace(parts[0])
		countStr := strings.TrimSpace(parts[1])
		count, err := strconv.Atoi(countStr)
		if err != nil || count < 1 {
			return nil, fmt.Errorf("invalid count for table %s: %s", table, countStr)
		}
		result[table] = count
	}
	return result, nil
}

func mealsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	rawTypes := r.URL.Query().Get("type")
	if rawTypes == "" {
		http.Error(w, "Missing 'type' parameter (e.g., HF_Meal:4,P3_Meal:4)", http.StatusBadRequest)
		return
	}

	typeCountMap, err := parseTypeCounts(rawTypes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var finalSelected []Meal
	for table, count := range typeCountMap {
		meals, err := getMeals(table)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get meals from table: %s", table), http.StatusInternalServerError)
			return
		}
		if len(meals) == 0 {
			continue
		}
		selected := weightedSelection(meals, count)
		finalSelected = append(finalSelected, selected...)
	}

	if len(finalSelected) == 0 {
		http.Error(w, "No meals found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(finalSelected)
}

// ====================== DESSERT API ======================

func dessertsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rows, err := db.Query("SELECT id, name, image_url, ingredients, bake_time, recipe_link FROM desserts ORDER BY id")
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

// ====================== ADMIN API ======================

// AddMealHandler handles adding a new meal
func addMealHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var meal struct {
		Table      string `json:"table"`
		Name       string `json:"name"`
		Rating     int    `json:"rating"`
		RecipeLink string `json:"recipe_link"`
		PhotoLink  string `json:"photo_link"`
	}

	if err := json.NewDecoder(r.Body).Decode(&meal); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var query string
	switch meal.Table {
	case "HF_Meal":
		query = `INSERT INTO HF_Meal (HFMealName, HFMealRating, RecipeLink, PhotoLink) VALUES (?, ?, ?, ?)`
	case "P3_Meal":
		query = `INSERT INTO P3_Meal (P3MealName, P3MealRating, RecipeLink, PhotoLink) VALUES (?, ?, ?, ?)`
	default:
		http.Error(w, "Invalid table name", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(query, meal.Name, meal.Rating, meal.RecipeLink, meal.PhotoLink)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    id,
		"msg":   "Meal added successfully",
		"table": meal.Table,
	})
}

// DeleteMealHandler handles deleting a meal
func deleteMealHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Table string `json:"table"`
		ID    int64  `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var query string
	switch req.Table {
	case "HF_Meal":
		query = `DELETE FROM HF_Meal WHERE HFMealID = ?`
	case "P3_Meal":
		query = `DELETE FROM P3_Meal WHERE P3MealID = ?`
	default:
		http.Error(w, "Invalid table name", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(query, req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"msg": "Meal deleted successfully"})
}

// AdminMeal represents a meal for admin operations (includes ID)
type AdminMeal struct {
	ID         int64  `json:"id"`
	IDCol      string `json:"id_column"`
	Name       string `json:"name"`
	Rating     int    `json:"rating"`
	PhotoLink  string `json:"image_url"`
	RecipeLink string `json:"recipe_link"`
	Source     string `json:"source"`
}

// getAdminMeals fetches meals with their IDs
func getAdminMeals(table string) ([]AdminMeal, error) {
	var query string
	var idCol string
	switch table {
	case "HF_Meal":
		query = `SELECT HFMealID, HFMealName, HFMealRating, PhotoLink, RecipeLink FROM HF_Meal`
		idCol = "HFMealID"
	case "P3_Meal":
		query = `SELECT P3MealID, P3MealName, P3MealRating, PhotoLink, RecipeLink FROM P3_Meal`
		idCol = "P3MealID"
	default:
		return nil, fmt.Errorf("invalid table name: %s", table)
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meals []AdminMeal
	for rows.Next() {
		var meal AdminMeal
		if err := rows.Scan(&meal.ID, &meal.Name, &meal.Rating, &meal.PhotoLink, &meal.RecipeLink); err != nil {
			return nil, err
		}
		meal.IDCol = idCol
		meal.Source = table
		meals = append(meals, meal)
	}
	return meals, nil
}

// GetAllMealsHandler returns all meals from both tables (with IDs)
func getAllMealsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	hfMeals, _ := getAdminMeals("HF_Meal")
	p3Meals, _ := getAdminMeals("P3_Meal")

	allMeals := append(hfMeals, p3Meals...)
	json.NewEncoder(w).Encode(allMeals)
}

// ====================== MAIN ======================

func main() {
	// Initialize database
	if err := initDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Determine the directory for static files
	var frontendDir string
	searchDirs := []string{
		filepath.Join(getExecutableDir(), "Frontend"),
		filepath.Join(getExecutableDir(), "..", "Frontend"),
		"Frontend",
	}

	// Try each directory
	for _, dir := range searchDirs {
		if _, err := os.Stat(dir); err == nil {
			frontendDir = dir
			break
		}
	}
	if frontendDir == "" {
		log.Fatal("Could not find Frontend directory - please ensure Frontend/ is next to the executable")
	}

	// Serve static frontend files
	// Use a custom handler for proper SPA-like routing
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Don't interfere with API routes or static assets
		if strings.HasPrefix(r.URL.Path, "/meals") ||
			strings.HasPrefix(r.URL.Path, "/api/") ||
			strings.HasPrefix(r.URL.Path, "/health") {
			http.NotFound(w, r)
			return
		}

		// Try to serve the exact path first
		reqPath := r.URL.Path
		if reqPath == "/" {
			reqPath = "/LandingPage/index.html"
		} else if !strings.Contains(reqPath, ".") {
			// It's a page path, try to find the file
			reqPath = reqPath + "/index.html"
		}

		// Build the file path
		filePath := filepath.Join(frontendDir, reqPath)
		
		// Security: ensure we're not escaping the frontend directory
		cleanPath, err := filepath.Abs(filePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		cleanFrontend, err := filepath.Abs(frontendDir)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if !strings.HasPrefix(cleanPath, cleanFrontend) {
			http.NotFound(w, r)
			return
		}

		// Serve the file
		file, err := os.Open(cleanPath)
		if err != nil {
			// Try serving index.html as fallback
			indexFile, indexErr := os.Open(filepath.Join(frontendDir, "LandingPage/index.html"))
			if indexErr != nil {
				http.NotFound(w, r)
				return
			}
			defer indexFile.Close()
			info, _ := indexFile.Stat()
			http.ServeContent(w, r, "index.html", info.ModTime(), indexFile)
			return
		}
		defer file.Close()
		
		info, _ := file.Stat()
		http.ServeContent(w, r, r.URL.Path, info.ModTime(), file)
	})

	// API routes
	http.HandleFunc("/meals", mealsHandler)
	http.HandleFunc("/api/desserts", dessertsHandler)
	http.HandleFunc("/api/admin/meals", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			action := r.URL.Query().Get("action")
			switch action {
			case "add":
				addMealHandler(w, r)
			case "delete":
				deleteMealHandler(w, r)
			case "all":
				getAllMealsHandler(w, r)
			default:
				http.Error(w, "Unknown action", http.StatusBadRequest)
			}
		default:
			getAllMealsHandler(w, r)
		}
	})

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	exePath, _ := os.Executable()
	log.Printf("🍽️  Homestyle Plate starting on http://localhost:%s", port)
	log.Printf("📁 Serving from: %s", frontendDir)
	log.Printf("📦 Database: %s", getExecutableDir())
	log.Printf("📦 Executable: %s", exePath)

	log.Printf("🚀 Starting server...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
