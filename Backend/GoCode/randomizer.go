package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// Database path
const dbConn = "user=alexsorichetti dbname=meals sslmode=disable"

// Meal represents a single meal record returned to the frontend
type Meal struct {
	Name       string `json:"name"`
	Rating     int    `json:"rating"`
	PhotoLink  string `json:"image_url"`
	RecipeLink string `json:"recipe_link"`
}

// getMeals fetches all meals from a specified table
func getMeals(db *sql.DB, table string) ([]Meal, error) {
	var query string
	if table == "HF_Meal" {
		query = "SELECT HFMealName, HFMealRating, PhotoLink, RecipeLink FROM HF_Meal"
	} else if table == "P3_Meal" {
		query = "SELECT P3MealName, P3MealRating, PhotoLink, RecipeLink FROM P3_Meal"
	} else {
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
		meals = append(meals, meal)
	}
	return meals, nil
}

// weightedSelection chooses meals using rating as a weight
func weightedSelection(meals []Meal, numMeals int) []Meal {
	rand.Seed(time.Now().UnixNano())
	selected := make(map[string]bool)
	var result []Meal

	weighted := []Meal{}
	for _, meal := range meals {
		for i := 0; i < meal.Rating; i++ {
			weighted = append(weighted, meal)
		}
	}

	for len(result) < numMeals && len(weighted) > 0 {
		meal := weighted[rand.Intn(len(weighted))]
		if !selected[meal.Name] {
			result = append(result, meal)
			selected[meal.Name] = true
		}
	}

	return result
}

// parseTypeCounts parses the type query param, e.g., "HF_Meal:4,P3_Meal:4"
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

// mealServer handles the /meals API route
func mealServer(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rawTypes := r.URL.Query().Get("type")
	if rawTypes == "" {
		http.Error(w, "Missing 'type' parameter", http.StatusBadRequest)
		return
	}

	typeCountMap, err := parseTypeCounts(rawTypes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var finalSelected []Meal
	for table, count := range typeCountMap {
		meals, err := getMeals(db, table)
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

func main() {
	http.HandleFunc("/meals", mealServer)
	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
