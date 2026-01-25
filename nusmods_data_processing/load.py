import psycopg2
from psycopg2.extras import execute_batch
import os, json, dotenv

dotenv.load_dotenv(os.path.join(os.path.dirname(__file__), '..', '.env'))
# DB_NAME = os.getenv("POSTGRES_DB")
# USER = os.getenv("POSTGRES_USER")
# PASSWORD = os.getenv("POSTGRES_PASSWORD")
DATABASE_URL = os.getenv("DATABASE_URL")

with open("filtered_modules.json") as f:
    modules = json.load(f)

conn = psycopg2.connect(DATABASE_URL, sslmode='require')

# conn = psycopg2.connect(
#     host="localhost",
#     port=5432,
#     dbname= DB_NAME,
#     user=USER,
#     password=PASSWORD,
# )

insert_sql = """
INSERT INTO modules (code, title, description, faculty)
VALUES (%(module_code)s, %(title)s, %(description)s, %(faculty)s)
ON CONFLICT (code) DO UPDATE SET
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    faculty = EXCLUDED.faculty;
"""

with conn:
    with conn.cursor() as cur:
        execute_batch(cur, insert_sql, modules, page_size=500)

conn.close()
