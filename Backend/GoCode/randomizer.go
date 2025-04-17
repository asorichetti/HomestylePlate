package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

// Database path
const dbPath = "./Backend/Meals.db"

// JSON response structure
type Meal struct {
	Name   string `json:"name"`
	Rating int    `json:"rating"`
}

// getMeals retrieves meals and ratings from a table (handles both HF_Meal and P3_Meal)
func getMeals(db *sql.DB, table string) ([]Meal, error) {
	var query string
	if table == "HF_Meal" {
		query = "SELECT HFMealName, HFMealRating FROM HF_Meal"
	} else if table == "P3_Meal" {
		query = "SELECT P3MealName, P3MealRating FROM P3_Meal"
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
		if err := rows.Scan(&meal.Name, &meal.Rating); err != nil {
			return nil, err
		}
		meals = append(meals, meal)
	}
	return meals, nil
}

// Function which randomly selects meals based on the weight of the entered meals
func weightedSelection(meals []Meal, numMeals int) []Meal {
	rand.Seed(time.Now().UnixNano())
	selected := make(map[string]bool)
	var result []Meal

	mealList := []Meal{}
	for _, meal := range meals {
		for i := 0; i < meal.Rating; i++ {
			mealList = append(mealList, meal)
		}
	}

	for len(result) < numMeals && len(mealList) > 0 {
		meal := mealList[rand.Intn(len(mealList))]
		if !selected[meal.Name] {
			result = append(result, meal)
			selected[meal.Name] = true
		}
	}

	return result
}

// API handler for Meal Selection
func mealServer(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	mealType := r.URL.Query().Get("type")
	mealCount := r.URL.Query().Get("count")

	count, err := strconv.Atoi(mealCount)
	if err != nil || count < 1 {
		http.Error(w, "Invalid meal count", http.StatusBadRequest)
		return
	}

	meals, err := getMeals(db, mealType)
	if err != nil || len(meals) == 0 {
		http.Error(w, "No meals found", http.StatusNotFound)
		return
	}

	selectedMeals := weightedSelection(meals, count)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(selectedMeals)
}

func main() {
	http.HandleFunc("/meals", mealServer)
	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
