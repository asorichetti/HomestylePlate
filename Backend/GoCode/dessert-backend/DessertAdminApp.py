import sqlite3
import tkinter as tk
from tkinter import ttk, messagebox

# Connect to Dessert Database #
DB_FILE = 'Backend/GoCode/dessert-backend/dessert.db'
def get_connection():
    return sqlite3.connect(DB_FILE)

# Tooltip Helper Class #
class ToolTip:
    def __init__(self, widget, text):
        self.widget = widget
        self.text = text
        self.tooltip = None
        self.widget.bind("<Enter>", self.show_tooltip)
        self.widget.bind("<Leave>", self.hide_tooltip)

    def show_tooltip(self, event):
        x, y, _, _ = self.widget.bbox("insert")
        self.tooltip = tk.Toplevel(self.widget)
        self.tooltip.wm_overrideredirect(True)
        self.tooltip.wm_geometry(f"+{x+25}+{y+25}")

        text_widget = tk.Text(
            self.tooltip,
            height=5,  # Adjust height as needed
            width=40,  # Adjust width as needed
            wrap="word",  # Word wrapping
            background="lightyellow",
            relief="solid",
            borderwidth=1,
            padx=5,  # Padding inside the Text widget
            pady=5   # Padding inside the Text widget
            )

        text_widget.insert("1.0", text)  # Insert the ingredients text
        text_widget.config(state=tk.DISABLED)  # Make it read-only so it doesn't allow editing
        text_widget.pack()

    def hide_tooltip(self, event):
        if self.tooltip:
            self.tooltip.destroy()

# Create GUI Elements #
class DessertAdminApp:
    def __init__(self, root):
        self.root = root
        self.root.title("Dessert Admin Panel")
        self.setup_widgets()
        self.load_data()

    def setup_widgets(self):
        # Form to Add/Edit Desserts #
        self.form_frame = tk.Frame(self.root)
        self.form_frame.pack(pady=10)

        tk.Label(self.form_frame, text="Dessert Name:").grid(row=0, column=0, sticky="e", padx=5)
        self.name_entry = tk.Entry(self.form_frame)
        self.name_entry.grid(row=0, column=1, padx=5)

        tk.Label(self.form_frame, text="Image URL:").grid(row=1, column=0, sticky="e", padx=5)
        self.image_entry = tk.Entry(self.form_frame)
        self.image_entry.grid(row=1, column=1, padx=5)

        tk.Label(self.form_frame, text="Ingredients:").grid(row=2, column=0, sticky="e", padx=5)
        self.ingredients_entry = tk.Entry(self.form_frame)
        self.ingredients_entry.grid(row=2, column=1, padx=5)

        tk.Label(self.form_frame, text="Bake Time:").grid(row=3, column=0, sticky="e", padx=5)
        self.bake_time_entry = tk.Entry(self.form_frame)
        self.bake_time_entry.grid(row=3, column=1, padx=5)

        tk.Label(self.form_frame, text="Recipe Link:").grid(row=4, column=0, sticky="e", padx=5)
        self.recipe_link_entry = tk.Entry(self.form_frame)
        self.recipe_link_entry.grid(row=4, column=1, padx=5)

        # Buttons #
        self.add_button = tk.Button(self.root, text="Add Dessert", command=self.add_record)
        self.add_button.pack(pady=5)

        self.edit_button = tk.Button(self.root, text="Edit Dessert", command=self.edit_record)
        self.edit_button.pack(pady=5)

        self.delete_button = tk.Button(self.root, text="Delete Dessert", command=self.delete_record)
        self.delete_button.pack(pady=5)

        # Table Display - Treeview #
        self.tree = ttk.Treeview(
            self.root,
            columns=("ID", "Name", "Image", "Ingredients", "BakeTime", "RecipeLink"),
            show="headings"
        )
        self.tree.heading("ID", text="ID")
        self.tree.column("ID", width=30)
        self.tree.heading("Name", text="Name")
        self.tree.column("Name", width=200)
        self.tree.heading("Image", text="Image URL")
        self.tree.column("Image", width=200)
        self.tree.heading("Ingredients", text="Ingredients")
        self.tree.column("Ingredients", width=250)
        self.tree.heading("BakeTime", text="Bake Time")
        self.tree.column("BakeTime", width=100)
        self.tree.heading("RecipeLink", text="Recipe Link")
        self.tree.column("RecipeLink", width=200)
        self.tree.pack(pady=10)
        self.tree.bind("<ButtonRelease-1>", self.on_record_select)

        # Tooltip behavior setup #
        self.setup_tooltip()

    def setup_tooltip(self):
        self.tooltip = None
        self.tree.bind("<Motion>", self.show_tooltip)
        self.tree.bind("<Leave>", self.hide_tooltip)

    def show_tooltip(self, event):
        region = self.tree.identify("region", event.x, event.y)
        if region != "cell":
            self.hide_tooltip(event)
            return

        row_id = self.tree.identify_row(event.y)
        col = self.tree.identify_column(event.x)

        if not row_id or col not in ('#2', '#4'):
            self.hide_tooltip(event)
            return

        item = self.tree.item(row_id)
        if col == '#2':
            text = item["values"][1]  # Index 1 = Name
        else:
            text = item["values"][3]

        if self.tooltip and self.tooltip.winfo_exists():
            self.tooltip.destroy()

        self.tooltip = tk.Toplevel(self.root)
        self.tooltip.wm_overrideredirect(True)
        self.tooltip.wm_geometry(f"+{event.x_root + 10}+{event.y_root + 10}")

        text_widget = tk.Text(
            self.tooltip,
            height=10,  # Adjust height as needed
            width=40,  # Adjust width as needed
            wrap="word",  # Word wrapping
            background="lightyellow",
            relief="solid",
            borderwidth=1,
            padx=5,  # Padding inside the Text widget
            pady=5   # Padding inside the Text widget
        )

        text_widget.insert("1.0", text)  # Insert the ingredients text
        text_widget.config(state=tk.DISABLED)  # Make it read-only so it doesn't allow editing
        text_widget.pack()

    def hide_tooltip(self, event):
        if self.tooltip and self.tooltip.winfo_exists():
            self.tooltip.destroy()
            self.tooltip = None

    def load_data(self):
        conn = get_connection()
        cursor = conn.cursor()

        # Clear Current Treeview #
        for row in self.tree.get_children():
            self.tree.delete(row)

        # Query All Desserts #
        cursor.execute("SELECT * FROM Desserts")
        rows = cursor.fetchall()

        for row in rows:
            self.tree.insert("", "end", values=row)

        conn.close()

    def on_record_select(self, event):
        selected_item = self.tree.selection()
        if not selected_item:
            return

        record = self.tree.item(selected_item)["values"]
        self.name_entry.delete(0, tk.END)
        self.name_entry.insert(0, record[1])

        self.image_entry.delete(0, tk.END)
        self.image_entry.insert(0, record[2])

        self.ingredients_entry.delete(0, tk.END)
        self.ingredients_entry.insert(0, record[3])

        self.bake_time_entry.delete(0, tk.END)
        self.bake_time_entry.insert(0, record[4])

        self.recipe_link_entry.delete(0, tk.END)
        self.recipe_link_entry.insert(0, record[5])

    def add_record(self):
        name = self.name_entry.get()
        image = self.image_entry.get()
        ingredients = self.ingredients_entry.get()
        bake_time = self.bake_time_entry.get()
        recipe = self.recipe_link_entry.get()

        if not name or not image or not ingredients or not bake_time or not recipe:
            messagebox.showerror("Input Error", "All fields are required for a valid entry")
            return

        conn = get_connection()
        cursor = conn.cursor()
        cursor.execute(
            "INSERT INTO Desserts (name, image_url, ingredients, bake_time, recipe_link) VALUES (?, ?, ?, ?, ?)",
            (name, image, ingredients, bake_time, recipe)
        )
        conn.commit()
        conn.close()
        self.load_data()
        self.clear_form()

    def edit_record(self):
        selected_item = self.tree.selection()
        if not selected_item:
            messagebox.showerror("Selection Error", "Please Select a Dessert to Edit")
            return

        record_id = self.tree.item(selected_item)["values"][0]
        name = self.name_entry.get()
        image = self.image_entry.get()
        ingredients = self.ingredients_entry.get()
        bake_time = self.bake_time_entry.get()
        recipe = self.recipe_link_entry.get()

        if not name or not image or not ingredients or not bake_time or not recipe:
            messagebox.showerror("Input Error", "All Fields are Required")
            return

        conn = get_connection()
        cursor = conn.cursor()
        cursor.execute(
            "UPDATE Desserts SET name=?, image_url=?, ingredients=?, bake_time=?, recipe_link=? WHERE id=?",
            (name, image, ingredients, bake_time, recipe, record_id)
        )
        conn.commit()
        conn.close()
        self.load_data()
        self.clear_form()

    def delete_record(self):
        selected_item = self.tree.selection()
        if not selected_item:
            messagebox.showerror("Selection Error", "Please Select a Dessert to Delete")
            return

        record_id = self.tree.item(selected_item)["values"][0]
        conn = get_connection()
        cursor = conn.cursor()
        cursor.execute("DELETE FROM Desserts WHERE id=?", (record_id,))
        conn.commit()
        conn.close()
        self.load_data()

    def clear_form(self):
        self.name_entry.delete(0, tk.END)
        self.image_entry.delete(0, tk.END)
        self.ingredients_entry.delete(0, tk.END)
        self.bake_time_entry.delete(0, tk.END)
        self.recipe_link_entry.delete(0, tk.END)

# Start App #
root = tk.Tk()
app = DessertAdminApp(root)
root.mainloop()
