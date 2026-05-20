#!/usr/bin/env python3
"""Convert the PostgreSQL dumps into bundled JSON for the offline iOS app.

Reads the COPY blocks from Backend/initdb/*.sql and writes Frontend/data/*.json.
The JSON shapes intentionally match what the Go APIs already return, so the
frontend can swap fetch() calls for local data with minimal changes.

Run from the repository root:  python3 tools/convert_data.py
"""

import json
import os
import sys

REPO_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
MEALS_SQL = os.path.join(REPO_ROOT, "Backend", "initdb", "meals_dump.sql")
DESSERTS_SQL = os.path.join(REPO_ROOT, "Backend", "initdb", "desserts_dump.sql")
OUT_DIR = os.path.join(REPO_ROOT, "Frontend", "data")


def extract_copy_rows(sql_path, table_marker):
    """Return the tab-separated data rows for one COPY block.

    table_marker is the substring identifying the COPY line, e.g.
    'COPY public."HF_Meal"'. Rows run until a line that is just '\\.'.
    Postgres escapes tab/newline/backslash inside fields; this dump's data
    contains none of those, so a plain tab split is correct here.
    """
    with open(sql_path, encoding="utf-8") as fh:
        lines = fh.read().splitlines()

    rows = []
    in_block = False
    for line in lines:
        if not in_block:
            if line.startswith(table_marker):
                in_block = True
            continue
        if line == r"\.":
            break
        rows.append(line.split("\t"))
    return rows


def convert_meals():
    """HF_Meal + P3_Meal -> list of meal dicts matching the /meals API shape."""
    meals = []
    for table, source in (("HF_Meal", "HF_Meal"), ("P3_Meal", "P3_Meal")):
        marker = f'COPY public."{table}"'
        for row in extract_copy_rows(MEALS_SQL, marker):
            # Columns: ID, Rating, Name, RecipeLink, PhotoLink
            _id, rating, name, recipe_link, photo_link = row
            meals.append({
                "name": name,
                "rating": int(rating),
                "image_url": photo_link,
                "recipe_link": recipe_link,
                "source": source,
            })
    return meals


def convert_desserts():
    """desserts -> list of dessert dicts matching the /api/desserts shape."""
    desserts = []
    for row in extract_copy_rows(DESSERTS_SQL, "COPY public.desserts"):
        # Columns: id, name, image_url, ingredients, bake_time, recipe_link
        _id, name, image_url, ingredients, bake_time, recipe_link = row
        desserts.append({
            "id": int(_id),
            "name": name,
            "image_url": image_url,
            "ingredients": ingredients,
            "bake_time": bake_time,
            "recipe_link": recipe_link,
        })
    return desserts


def main():
    os.makedirs(OUT_DIR, exist_ok=True)

    meals = convert_meals()
    desserts = convert_desserts()

    meals_path = os.path.join(OUT_DIR, "meals.json")
    desserts_path = os.path.join(OUT_DIR, "desserts.json")

    with open(meals_path, "w", encoding="utf-8") as fh:
        json.dump(meals, fh, indent=2, ensure_ascii=False)
        fh.write("\n")
    with open(desserts_path, "w", encoding="utf-8") as fh:
        json.dump(desserts, fh, indent=2, ensure_ascii=False)
        fh.write("\n")

    hf = sum(1 for m in meals if m["source"] == "HF_Meal")
    p3 = sum(1 for m in meals if m["source"] == "P3_Meal")
    print(f"meals.json:    {len(meals)} meals ({hf} HF_Meal, {p3} P3_Meal)")
    print(f"desserts.json: {len(desserts)} desserts")

    # Basic sanity checks so a malformed dump fails loudly.
    if not meals or not desserts:
        print("ERROR: produced an empty dataset", file=sys.stderr)
        return 1
    for m in meals:
        if not m["name"] or m["rating"] < 1:
            print(f"ERROR: bad meal record: {m}", file=sys.stderr)
            return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
