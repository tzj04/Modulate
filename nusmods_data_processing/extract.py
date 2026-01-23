import requests
import json

URL = "https://api.nusmods.com/v2/2025-2026/moduleInfo.json"

response = requests.get(URL)
response.raise_for_status()

with open("raw_modules.json", "w") as f:
    json.dump(response.json(), f)

