#This program will provide a desired amount of reciped for a given week based on the desired number of recipes given a selection
# Hello Fresh (HF) meals or Recipes from our Paprika3 Library
import random

#This function handles selection of a number of HF meals with weights provided based on how much the given meal was enjoyed. 
def weighted_selection_hf(ui):
    meal = ["Almond Crusted Chicken", "Bacon, Apple and Cheddar Melts", "Baked Chicken Parmesan", "BBQ Beef Meatballs", "BBQ Pork and Apple Meatballs", "BBQ-Sauced Chicken", "BBQ Steak Sandwiches", "Beef and Roasted Pepper Ragu", "Beef Taquitos", "Bodega-Inspired Beef Burgers", "Braised Chicken and Mushrooms", "Breaded Baked Chicken", "Caesar Turkey Burgers", "Cal Smart Sweet Soy Turkey Patties", "Cheesy Beef and Pork Hash", "Cheesy Beef and Rigatoni Bake", "Cheesy Pork Enchiladas", "Cheesy Pork Quesadillas", "Cheesy Stuffed Chicken and Sweet Potato Mash", "Cheesy Tex-Mex Orzo Skillet", "Chicken Chow Mein-Style Noodles", "Chinese Cashew Pork Stir-Fry", "Chicken Schnitzel", "Chicken Tikka Masala", "Chopped Cheese Quesadillas", "Fajita-Inspired Chicken Flatbreads", "Fajita Style Beef Bowls", "Farmhouse Chicken", "Fig Maple Pork Tenderloin", "French Dip Burgers", "Galactic Pork Souvlaki-Inspired Burgers", "Golden Breaded Tilapia", "Grilled Sausage Flatbreads", "Grilled Tuscan Chicken Penne", "Hearty Beef and Black Bean Chili", "Hearty Beef and Pork Ragu", "Hearty Meatball and Mushroom Stew", "Honey-Garlic Chicken Wraps", "Honey Mustard Baked Salmon", "Italian Inspired Beef Burgers", "Lemony Beef and Orzo Bowls", "Maple Balsamic Chicken", "Mediterranean Tortellini", "Mexican-Inspired Pork Quesadillas", "Miso-Honey Glazed Salmon", "Moroccan-Inspired Chicken", "One Pan Chicken and Orzo", "One-Pot Italian Sausage Soup", "One-Pot Soutwest-Style Beef and Cavatappi", "Pan-Seared Pork Chops", "Panko-Crusted Chicken", "Parmesan-Crusted Pork Chops", "Pecan-Crusted Roasted Salmon", "Pesto Turkey Bowls", "Pork and Apple Burgers", "Pork and Sweet Pepper Tacos", "Pork Chops and Mushroom-Sour Cream Sauce", "Pork Spring Roll-Inspired Bowls", "Retro Burgers", "Salisbury Steaks in Onion Gravy", "Salisbury Steak-Style Pork Patties and Mash", "Sausage Gnocchi Skillet", "Scandi-Style Salmon", "Sheet Pan Beef Koftas", "Silky Peanut-Ginger Beef and Pork Noodle Bowls", "Smart Cruchy Clementine Chicken Salad", "Smart Mexican-Inspired Chicken Patties", "Smash Burgers", "SuperQuick Southwest Turkey Salad", "Speedy BBQ Chicken Sandwiches", "Speedy Chicken Chili", "St.Baptiste Cheesy Blueberry Burgers", "Steak with Rosemary Sous Vide Potatoes", "Sticky Glazed Chicken Sammies", "Super Quick Open-Faced Speedy Sausage Sandos", "Super Quick Turkey Tacos", "Sweet Ginger Pork Stir-Fry", "Sweet n Savory Beef and Pork Noodles", "Tangy Sesame Chicken", "Super Quick Teriyaki Beef Ramen Bowls", "Tex-Mex Beef and Rice Skillet", "Tex-Mex Style Beef and Pork Skillet Rice", "Turkey Burrito Bowls", "Turkey Sloppy Joes", "Turkey Taco Pizzas", "Turkey and Fig Burgers", "Turkey and Zucchini Rigatoni", "Umami Steak and Noodle Stir-Fry", "Zesty Beef Bowls"]
    weight = [1, 2, 2, 5, 2, 3, 4, 3, 7, 1, 2, 1, 3, 2, 5, 2, 7, 3, 3, 4, 4, 2, 4, 1, 2, 4, 2, 3, 3, 2, 3, 2, 1, 2, 2, 5, 3, 2, 2, 4, 2, 1, 1, 3, 4, 4, 2, 1, 2, 1, 1, 2, 1, 4, 2, 3, 7, 6, 5, 5, 2, 1, 6, 5, 1, 2, 1, 1, 3, 1, 5, 5, 5, 3, 2, 2, 3, 4, 2, 5, 3, 2, 2, 1, 3, 3, 3, 3, 2]
    x = False
    selected_hf = []
    for i in range (0, ui, 1):
        selected_hf.append(random.choices(meal, weights=weight, k=1)[0])
        for h in range (0, i, 1):
            if selected_hf[h] == selected_hf[i]:
                selected_hf[i] = random.choices(meal, weights=weight, k=1)[0]
                h -= 1
    return selected_hf

hf_amount = int(input("Welcome to the menu selector. Please input how many Hello Fresh Recipes you'd like selected for this week: "))
if (hf_amount > 0):
    selected_hf = weighted_selection_hf(hf_amount)

print("Your selected meals for the week are the following:")
for l in range (0, hf_amount, 1):
    print(selected_hf[l])