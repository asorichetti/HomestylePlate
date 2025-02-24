package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "modernc.org/sqlite"
)

// Database path
const dbPath = "../Meals.db"

// getMeals retrieves meals and ratings from a table (handles both HF_Meal and P3_Meal)
func getMeals(db *sql.DB, table string) (map[string]int, error) {
	var query string
	if table == "HF_Meal" {
		query = "SELECT HFMealID, HFMealRating, HFMealName FROM HF_Meal"
	} else if table == "P3_Meal" {
		query = "SELECT P3MealID, P3MealRating, P3MealName FROM P3_Meal"
	} else {
		return nil, fmt.Errorf("invalid table name: %s", table)
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	meals := make(map[string]int)
	for rows.Next() {
		var id, rating int
		var name string
		if err := rows.Scan(&id, &rating, &name); err != nil {
			return nil, err
		}
		meals[name] = rating
	}
	return meals, nil
}

// Function which randomly selects meals based on the weight of the entered meals
func weightedSelection(meals map[string]int, numMeals int) []string {
	rand.Seed(time.Now().UnixNano())
	selected := make(map[string]bool)
	result := make([]string, 0, numMeals)

	mealList := []string{}
	for name, weight := range meals {
		for i := 0; i < weight; i++ {
			mealList = append(mealList, name)
		}
	}

	for len(result) < numMeals && len(mealList) > 0 {
		meal := mealList[rand.Intn(len(mealList))]
		if !selected[meal] {
			result = append(result, meal)
			selected[meal] = true
		}
	}

	return result
}

func main() {
	// Open database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()
	fmt.Println("Connected to the DB Successfully")

	// Get user input for desired amounts of recipes
	var hfAmount, p3Amount int
	fmt.Print("Please enter number of Hello Fresh Recipes: ")
	fmt.Scan(&hfAmount)
	fmt.Print("Please enter number of Paprika3 Recipes: ")
	fmt.Scan(&p3Amount)

	// Fetch meals from database in the correct areas
	hfMeals, err := getMeals(db, "HF_Meal")
	if err != nil {
		log.Fatal("Error fetching Hello Fresh meals:", err)
	}
	p3Meals, err := getMeals(db, "P3_Meal")
	if err != nil {
		log.Fatal("Error fetching Paprika3 meals:", err)
	}

	// Perform weighted selections
	var selectedHF, selectedP3 []string
	if hfAmount > 0 {
		selectedHF = weightedSelection(hfMeals, hfAmount)
	}
	if p3Amount > 0 {
		selectedP3 = weightedSelection(p3Meals, p3Amount)
	}

	// Print results for each area of meal options
	fmt.Println("\nYour selected meals for the week are the following:")
	if len(selectedHF) > 0 {
		fmt.Println("🍽️ Hello Fresh Meals:")
		for _, meal := range selectedHF {
			fmt.Println("-", meal)
		}
	} else {
		//Error message for if no meals are selected
		fmt.Println("⚠️ No Hello Fresh meals selected.")
	}

	if len(selectedP3) > 0 {
		fmt.Println("\n🥗 Paprika3 Meals:")
		for _, meal := range selectedP3 {
			fmt.Println("-", meal)
		}
	} else {
		//Error message for if no meals are selected
		fmt.Println("⚠️ No Paprika3 meals selected.")
	}
}
