import psycopg2
from dotenv import load_dotenv
import os

load_dotenv("../../../.env.local")
DATABASE_URL = os.environ.get("DATABASE_URL")

def clean_database():
    try:
        with psycopg2.connect(DATABASE_URL) as conn:
            with conn.cursor() as cur:
                cur.execute("TRUNCATE TABLE users CASCADE;")
                cur.execute("TRUNCATE TABLE messages CASCADE;")
                cur.execute("TRUNCATE TABLE followers CASCADE;")
                conn.commit()
        return True
    except Exception as e:
        print(f"Database cleaning failed: {e}")
        return False 