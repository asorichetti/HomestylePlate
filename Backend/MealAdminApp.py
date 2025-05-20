import psycopg2
import tkinter as tk
from tkinter import ttk, messagebox

# Connect to Database #
def get_connection():
    return psycopg2.connect(
        dbname='meals',
        user='alexsorichetti',
        host='localhost'
    )


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
        self.tooltip.wm_overrideredirect(True)  # Removes window decorations
        self.tooltip.wm_geometry(f"+{x+25}+{y+25}")

        label = tk.Label(self.tooltip, text=self.text, background="yellow", relief="solid", borderwidth=1)
        label.pack()

    def hide_tooltip(self, event):
        if self.tooltip:
            self.tooltip.destroy()

# Create GUI Elements #
class MealAdminApp:
    def __init__(self, root):
        self.root=root
        self.root.title("Meal Admin Panel")
        self.table_name = tk.StringVar(value="HF_Meal")
        self.setup_widgets()
        self.load_data()
        # Reset sequence once on startup
        table = self.table_name.get()
        id_column = "HFMealID" if table == "HF_Meal" else "P3MealID"
        self.reset_sequence()

    # Reset PostgreSQL sequence #
    def reset_sequence(self):
        table = self.table_name.get()
        id_column = "HFMealID" if table == "HF_Meal" else "P3MealID"
        conn = get_connection()
        cursor = conn.cursor()
        query = f"""
            SELECT setval(
                pg_get_serial_sequence('"{table}"', '{id_column}'),
                (SELECT MAX("{id_column}") FROM "{table}")
            )
            """
        cursor.execute(query)
        conn.commit()
        conn.close()

    # Widget Creation #
    def setup_widgets(self):
        # Table Selection #
        table_frame = tk.Frame(self.root)
        table_frame.pack(pady=10)

        tk.Label(table_frame, text="Select Table:").pack(side=tk.LEFT)
        table_dropdown = ttk.Combobox(
            table_frame,
            textvariable=self.table_name,
            values=["HF_Meal", "P3_Meal"],
            state="readonly"
        )
        table_dropdown.pack(side=tk.LEFT)
        table_dropdown.bind("<<ComboboxSelected>>", lambda e: self.load_data())

        # Form to Add/Edit Meals #
        self.form_frame = tk.Frame(self.root)
        self.form_frame.pack(pady=10)

        tk.Label(self.form_frame, text="Meal Name:").grid(row=0, column=0, sticky="e", padx=5)
        self.name_entry = tk.Entry(self.form_frame)
        self.name_entry.grid(row=0, column=1, padx=5)

        tk.Label(self.form_frame, text="Meal Rating:").grid(row=1, column=0, sticky="e", padx=5)
        self.rating_entry = tk.Entry(self.form_frame)
        self.rating_entry.grid(row=1, column=1, padx=5)

        tk.Label(self.form_frame, text="Recipe Link:").grid(row=2, column=0, sticky="e", padx=5)
        self.recipe_entry = tk.Entry(self.form_frame)
        self.recipe_entry.grid(row=2, column=1, padx=5)

        tk.Label(self.form_frame, text="Photo Link:").grid(row=3, column=0, sticky="e", padx=5)
        self.photo_entry = tk.Entry(self.form_frame)
        self.photo_entry.grid(row=3, column=1, padx=5)

        # Buttons #
        self.add_button = tk.Button(self.root, text="Add Record", command=self.add_record)
        self.add_button.pack(pady=5)

        self.edit_button = tk.Button(self.root, text="Edit Record", command=self.edit_record)
        self.edit_button.pack(pady=5)

        self.delete_button = tk.Button(self.root, text="Delete Record", command=self.delete_record)
        self.delete_button.pack(pady=5)

        # Table Display - Treeview #
        self.tree = ttk.Treeview(
            self.root,
            columns=("ID", "Rating", "Name", "RecipeLink", "PhotoLink"),
            show="headings"
        )
        self.tree.heading("ID", text="ID")
        self.tree.column("ID", width=30)
        self.tree.heading("Rating", text="Rating")
        self.tree.column("Rating", width=10)
        self.tree.heading("Name", text="Name")
        self.tree.column("Name", width=300)
        self.tree.heading("RecipeLink", text="Recipe Link")
        self.tree.column("RecipeLink", width=200)
        self.tree.heading("PhotoLink", text="Photo Link")
        self.tree.column("PhotoLink", width=200)
        self.tree.pack(pady=10)
        self.tree.bind("<ButtonRelease-1>", self.on_record_select)

        # Tooltip behavior setup #
        self.setup_tooltip()

    # Tooltip Logic #
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

        if not row_id or col != '#3':  # '#3' is the Name column
            self.hide_tooltip(event)
            return

        item = self.tree.item(row_id)
        name = item["values"][2]  # Index 2 = Name

        if self.tooltip and self.tooltip.winfo_exists():
            self.tooltip.destroy()

        self.tooltip = tk.Toplevel(self.root)
        self.tooltip.wm_overrideredirect(True)
        self.tooltip.wm_geometry(f"+{event.x_root + 10}+{event.y_root + 10}")

        label = tk.Label(self.tooltip, text=name, background="lightyellow", relief="solid", borderwidth=1)
        label.pack()

    def hide_tooltip(self, event):
        if self.tooltip and self.tooltip.winfo_exists():
            self.tooltip.destroy()
            self.tooltip = None

    # Load Data #
    def load_data(self):
        table = self.table_name.get()
        conn = get_connection()
        cursor = conn.cursor()

        # Clear Current Treeview #
        for row in self.tree.get_children():
            self.tree.delete(row)

        # Query the Selected Table #
        id_column = "HFMealID" if table == "HF_Meal" else "P3MealID"
        query = f'SELECT * FROM "{table}" ORDER BY "{id_column}" ASC'
        cursor.execute(query)
        rows = cursor.fetchall()

        # Insert Data into Treeview #
        for row in rows:
            self.tree.insert("", "end", values=row) 
        conn.close()

    # Populate Form Field with Record Date #    
    def on_record_select(self, event):
        selected_item = self.tree.selection()
        if not selected_item:
            return

        record = self.tree.item(selected_item)["values"]
        self.name_entry.delete(0, tk.END)
        self.name_entry.insert(0, record[2])

        self.rating_entry.delete(0, tk.END)
        self.rating_entry.insert(0, record[1])

        self.recipe_entry.delete(0, tk.END)
        self.recipe_entry.insert(0, record[3])

        self.photo_entry.delete(0, tk.END)
        self.photo_entry.insert(0, record[4])

    def add_record(self):
        # Add New Record to DB #
        name = self.name_entry.get()
        rating = self.rating_entry.get()
        recipe_link = self.recipe_entry.get()
        photo_link = self.photo_entry.get()

        if not name or not rating or not recipe_link or not photo_link:
            messagebox.showerror("Input Error", "All fields are required for a valid entry")
            return

        table= self.table_name.get()
        conn = get_connection()
        cursor = conn.cursor()

        # Insert New Record into the Selected Table #
        query = f'INSERT INTO "{table}" ("HFMealName", "HFMealRating", "RecipeLink", "PhotoLink") VALUES (%s, %s, %s, %s)' if table == "HF_Meal" else f'INSERT INTO "{table}" ("P3MealName", "P3MealRating", "RecipeLink", "PhotoLink") VALUES (%s, %s, %s, %s)'
        cursor.execute(query, (name, rating, recipe_link, photo_link))
        conn.commit()

        conn.close()
        self.load_data()
        self.clear_form()

    def edit_record(self):
        # Edit an Existing Record #
        selected_item = self.tree.selection()
        if not selected_item:
            messagebox.showerror("Selection Error", "Please Select a Record to Edit")
            return

        record_id = self.tree.item(selected_item)["values"][0]
        name = self.name_entry.get()
        rating = self.rating_entry.get()
        recipe_link = self.recipe_entry.get()
        photo_link = self.photo_entry.get()

        if not name or not rating or not recipe_link or not photo_link:
            messagebox.showerror ("Input Error", "All Fields are Required to Create a Valid Entry")
            return

        table = self.table_name.get()
        conn = get_connection()
        cursor = conn.cursor()

        # Update Selected Record #
        query = f'UPDATE "{table}" SET "HFMealName"=%s, "HFMealRating"=%s, "RecipeLink"=%s, "PhotoLink"=%s WHERE "HFMealID"=%s' if table == "HF_Meal" else f'UPDATE "{table}" SET "P3MealName"=%s, "P3MealRating"=%s, "RecipeLink"=%s, "PhotoLink"=%s WHERE "P3MealID"=%s'
        cursor.execute (query, (name, rating, recipe_link, photo_link, record_id))
        conn.commit()
        conn.close()
        self.load_data()
        self.clear_form()

    def delete_record(self):
        # Delete a Selected Record #
        selected_item = self.tree.selection()
        if not selected_item: 
            messagebox.showerror("Selection Error", "Please select a record to delete")
            return

        record_id = self.tree.item(selected_item)["values"][0]
        table = self.table_name.get()

        conn = get_connection()
        cursor = conn.cursor()

        # Delete the Record #
        query = f'DELETE FROM "{table}" WHERE "HFMealID"=%s' if table == "HF_Meal" else f'DELETE FROM "{table}" WHERE "P3MealID"=%s'
        cursor.execute(query, (record_id,))
        conn.commit()
        conn.close()
        self.load_data()
        id_column = "HFMealID" if table == "HF_Meal" else "P3MealID"
        self.reset_sequence()


    def clear_form(self):
        # Clear the Form Fields After Adding or Editing #
        self.name_entry.delete(0, tk.END)
        self.rating_entry.delete(0, tk.END)
        self.recipe_entry.delete(0, tk.END)
        self.photo_entry.delete(0, tk.END)

# Create Tkinter Window and Start Application
root = tk.Tk()
app = MealAdminApp(root)
root.mainloop()
