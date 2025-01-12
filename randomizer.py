import random

def weighted_selection(meal, weight):
    return random.choices(meal, weights=weight, k=1)[0]

meal = ["Almond Crusted Chicken", "Bacon, Apple and Cheddar Melts", "Baked Chicken Parmesan", "BBQ Beef Meatballs", "BBQ Pork and Apple Meatballs", "BBQ-Sauced Chicken", "BBQ Steak Sandwiches", "Beef and Roasted Pepper Ragu", "Beef Taquitos", "Bodega-Inspired Beef Burgers", "Braised Chicken and Mushrooms", "Breaded Baked Chicken", "Caesar Turkey Burgers", "Cal Smart Sweet Soy Turkey Patties", "Cheesy Beef and Pork Hash", "Cheesy Beef and Rigatoni Bake", "Cheesy Pork Enchiladas", "Cheesy Pork Quesadillas", "Cheesy Stuffed Chicken and Sweet Potato Mash", "Cheesy Tex-Mex Orzo Skillet", "Chicken Chow Mein-Style Noodles", "Chinese Cashew Pork Stir-Fry", "Chicken Schnitzel", "Chicken Tikka Masala", "Chopped Cheese Quesadillas", "Fajita-Inspired Chicken Flatbreads", "Fajita Style Beef Bowls", "Farmhouse Chicken", "Fig Maple Pork Tenderloin", "French Dip Burgers", "Galactic Pork Souvlaki-Inspired Burgers", "Golden Breaded Tilapia", "Grilled Sausage Flatbreads", "Grilled Tuscan Chicken Penne", "Hearty Beef and Black Bean Chili", "Hearty Beef and Pork Ragu", "Hearty Meatball and Mushroom Stew", "Honey-Garlic Chicken Wraps", "Honey Mustard Baked Salmon", "Italian Inspired Beef Burgers", "Lemony Beef and Orzo Bowls", "Maple Balsamic Chicken", "Mediterranean Tortellini", "Mexican-Inspired Pork Quesadillas", "Miso-Honey Glazed Salmon", "Moroccan-Inspired Chicken", "One Pan Chicken and Orzo", "One-Pot Italian Sausage Soup", "One-Pot Soutwest-Style Beef and Cavatappi", "Pan-Seared Pork Chops", "Panko-Crusted Chicken", "Parmesan-Crusted Pork Chops", "Pecan-Crusted Roasted Salmon", "Pesto Turkey Bowls", "Pork and Apple Burgers", "Pork and Sweet Pepper Tacos", "Pork Chops and Mushroom-Sour Cream Sauce", "Pork Spring Roll-Inspired Bowls", "Retro Burgers", "Salisbury Steaks in Onion Gravy", "Salisbury Steak-Style Pork Patties and Mash", "Sausage Gnocchi Skillet", "Scandi-Style Salmon", "Sheet Pan Beef Koftas", "Silky Peanut-Ginger Beef and Pork Noodle Bowls", "Smart Cruchy Clementine Chicken Salad", "Smart Mexican-Inspired Chicken Patties", "Smash Burgers", "SuperQuick Southwest Turkey Salad", "Speedy BBQ Chicken Sandwiches", "Speedy Chicken Chili", "St.Baptiste Cheesy Blueberry Burgers", "Steak with Rosemary Sous Vide Potatoes", "Sticky Glazed Chicken Sammies", "Super Quick Open-Faced Speedy Sausage Sandos", "Super Quick Turkey Tacos", "Sweet Ginger Pork Stir-Fry", "Sweet n Savory Beef and Pork Noodles", "Tangy Sesame Chicken", "Super Quick Teriyaki Beef Ramen Bowls", "Tex-Mex Beef and Rice Skillet", "Tex-Mex Style Beef and Pork Skillet Rice", "Turkey Burrito Bowls", "Turkey Sloppy Joes", "Turkey Taco Pizzas", "Turkey and Fig Burgers", "Turkey and Zucchini Rigatoni", "Umami Steak and Noodle Stir-Fry", "Zesty Beef Bowls"]
weight = [1, 2, 2, 5, 2, 3, 4, 3, 7, 1, 2, 1, 3, 2, 5, 2, 7, 3, 3, 4, 4, 2, 4, 1, 2, 4, 2, 3, 3, 2, 3, 2, 1, 2, 2, 5, 3, 2, 2, 4, 2, 1, 1, 3, 4, 4, 2, 1, 2, 1, 1, 2, 1, 4, 2, 3, 7, 6, 5, 5, 2, 1, 6, 5, 1, 2, 1, 1, 3, 1, 5, 5, 5, 3, 2, 2, 3, 4, 2, 5, 3, 2, 2, 1, 3, 3, 3, 3, 2]
x = False
selected = []
selected.append(weighted_selection(meal, weight))
while (x != True):
    selected.append(weighted_selection(meal, weight))
    if selected[1] != selected[0]:
        x = True
x = False
while (x != True):
    selected.append(weighted_selection(meal, weight))
    if selected[2] != selected[0] and selected[2] != selected[1]:
        x = True

print(f"Your selected meals for the week are the following: {selected[0]}, {selected[1]}, and {selected[2]}")