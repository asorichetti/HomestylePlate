import psycopg2
import tkinter as tk
from tkinter import ttk, messagebox



def get_connection():
    return psycopg2.connect(
        dbname='desserts',
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
        self.tooltip.wm_overrideredirect(True)
        self.tooltip.wm_geometry(f"+{x+25}+{y+25}")

        text_widget = tk.Text(
            self.tooltip,
            height=5,
            width=40,
            wrap="word",
            background="lightyellow",
            relief="solid",
            borderwidth=1,
            padx=5,
            pady=5
        )
        text_widget.insert("1.0", self.text)
        text_widget.config(state=tk.DISABLED)
        text_widget.pack()

    def hide_tooltip(self, event):
        if self.tooltip:
            self.tooltip.destroy()

# Create GUI Elements #
class DessertAdminApp:
    def __init__(self, root):
        self.root = root
        self.root.title("Dessert Admin Panel")
        self.tooltip = None  # Initialize tooltip as None
        self.setup_widgets()
        self.load_data()

    def setup_widgets(self):
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

        self.add_button = tk.Button(self.root, text="Add Dessert", command=self.add_record)
        self.add_button.pack(pady=5)

        self.edit_button = tk.Button(self.root, text="Edit Dessert", command=self.edit_record)
        self.edit_button.pack(pady=5)

        self.delete_button = tk.Button(self.root, text="Delete Dessert", command=self.delete_record)
        self.delete_button.pack(pady=5)

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

        self.tree.bind("<Motion>", self.show_tooltip)
        self.tree.bind("<Leave>", self.hide_tooltip)
        self.tree.bind("<ButtonRelease-1>", self.on_tree_click)

    def show_tooltip(self, event):
        # Identify the region and check if it's a valid cell in the 'Name' column
        region = self.tree.identify("region", event.x, event.y)
        if region != "cell":
            self.hide_tooltip(event)
            return

        row_id = self.tree.identify_row(event.y)
        col = self.tree.identify_column(event.x)

        # Only show tooltip for the 'Name' column (col == '#2')
        if not row_id or col != '#2':
            self.hide_tooltip(event)
            return

        item = self.tree.item(row_id)
        text = item["values"][1]  # The "Name" column (index 1)

        # Destroy existing tooltip if it exists
        if self.tooltip and self.tooltip.winfo_exists():
            self.tooltip.destroy()

        # Create the tooltip
        self.tooltip = tk.Toplevel(self.root)
        self.tooltip.wm_overrideredirect(True)
        text_widget = tk.Text(
            self.tooltip,
            wrap="word",
            width=40,
            height=2,
            background="lightyellow",
            relief="solid",
            borderwidth=1,
            padx=5,
            pady=5
        )
        text_widget.insert(tk.END, text)
        text_widget.config(state=tk.DISABLED)
        text_widget.pack()

        # Position the tooltip near the mouse pointer
        self.tooltip.wm_geometry(f"+{event.x_root + 10}+{event.y_root + 10}")

    def hide_tooltip(self, event):
        # Ensure that the tooltip exists before trying to destroy it
        if self.tooltip and self.tooltip.winfo_exists():
            self.tooltip.destroy()
            self.tooltip = None

    def on_tree_click(self, event):
        region = self.tree.identify("region", event.x, event.y)
        if region != "cell":
            return

        row_id = self.tree.identify_row(event.y)
        col = self.tree.identify_column(event.x)

        if col == '#4':  # Ingredients column
            item = self.tree.item(row_id)
            ingredients = item["values"][3]

            # Show a new popup window with full ingredients
            popup = tk.Toplevel(self.root)
            popup.title("Ingredients")

            text = tk.Text(popup, wrap="word", width=60, height=15)
            text.insert("1.0", ingredients)
            text.config(state=tk.DISABLED)
            text.pack(padx=10, pady=10)

            tk.Button(popup, text="Close", command=popup.destroy).pack(pady=5)

    def load_data(self):
        conn = get_connection()
        cursor = conn.cursor()
        for row in self.tree.get_children():
            self.tree.delete(row)
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
            "INSERT INTO Desserts (name, image_url, ingredients, bake_time, recipe_link) VALUES (%s, %s, %s, %s, %s)",
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
            "UPDATE Desserts SET name=%s, image_url=%s, ingredients=%s, bake_time=%s, recipe_link=%s WHERE id=%s",
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
        cursor.execute("DELETE FROM Desserts WHERE id=%s", (record_id,))
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
