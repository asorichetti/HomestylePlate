# This program will provide a desired amount of reciped for a given week based on the desired number of recipes given a selection
# Hello Fresh (HF) meals or Recipes from our Paprika3 (P3) Library
# All ratings are based on my families enjoyment of the meals listed. 
import random
import sqlite3 as db
#This function handles selection of a number of HF meals with weights provided based on how much the given meal was enjoyed. 
def weighted_selection_hf(ui):
    crsr.execute("SELECT * FROM HF_Meal")
    hf_meals = crsr.fetchall()
    meals = []
    rating = []
    for i in range(0, len(hf_meals), 1):
        meals.append(hf_meals[i][2])
        rating.append(hf_meals[i][1])
    selected_hf = []
    for i in range (0, ui, 1):
        selected_hf.append(random.choices(meals, weights=rating, k=1)[0])
        for h in range (0, i, 1):
            if selected_hf[h] == selected_hf[i]:
                selected_hf[i] = random.choices(meals, weights=rating, k=1)[0]
                h -= 1
    return selected_hf

# This function is to handle selection of user inputted number of Paprika3 meals with weights provided based on the star rating in the P3 app
def weighted_selection_P3(ui):
    crsr.execute("SELECT * FROM P3_Meal")
    p3_meals = crsr.fetchall()
    meals = []
    rating = []
    for i in range(0, len(p3_meals), 1):
        meals.append(p3_meals[i][2])
        rating.append(p3_meals[i][1])
    selected_p3 = []
    for i in range (0, ui, 1):
        selected_p3.append(random.choices(meals, weights=rating, k=1)[0])
        for h in range (0, i, 1):
            if selected_p3[h] == selected_p3[i]:
                selected_p3[i] = random.choices(meals, weights=rating, k=1)[0]
                h -= 1
    return selected_p3


conn = db.connect('./Meals.db')
crsr = conn.cursor()
print ("Connected to the DB")

hf_amount = int(input("Welcome to the menu selector. Please input how many Hello Fresh Recipes you'd like selected for this week: "))
if (hf_amount > 0):
    selected_hf = weighted_selection_hf(hf_amount)
P3_amount = int(input("Now please enter the number of Paprika3 Recipes you would like for the week: "))
if (P3_amount > 0):
    selected_p3 = weighted_selection_P3(P3_amount)

print("Your selected meals for the week are the following:")
for l in range (0, hf_amount, 1):
    print(selected_hf[l])
for p in range (0, P3_amount, 1):
    print(selected_p3[p])

conn.close()