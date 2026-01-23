import json

with open("raw_modules.json") as f:
    modules = json.load(f)

filtered = []

for m in modules:
    try:
        filtered.append({
            "module_code": m["moduleCode"],
            "title": m.get("title"),
            "description": m.get("description"),
            "faculty": m.get("faculty"),
        })
    except KeyError:
        # Skip malformed entries
        continue

with open("filtered_modules.json", "w") as f:
    json.dump(filtered, f, indent=2)