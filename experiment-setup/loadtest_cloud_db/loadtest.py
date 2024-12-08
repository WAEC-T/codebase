import asyncio
import asyncpg
import time
from pathlib import Path
import pandas as pd
from execution.orchestration import clean_database

DATABASE_URL = ""

async def load_test(user_id: int):
    """Load test operation to insert a user into the database and fetch the inserted user.
    """
    try:
        # Establish an asynchronous connection to the database
        conn = await asyncpg.connect(DATABASE_URL)
        try:
            # Generate random data for the user
            username = f"User_{user_id}"
            email = f"{username}@example.com"
            pw_hash = "pbkdf2:sha256:50000$example_hash$example_salt"

            # Perform an INSERT operation
            insert_query = """
                INSERT INTO users VALUES ($1, $2, $3, $4)
                RETURNING user_id;
            """
            inserted_user_id = await conn.fetchval(insert_query, user_id, username, email, pw_hash)

            # Perform a SELECT operation to fetch the inserted user
            select_query = "SELECT * FROM users WHERE user_id = $1;"
            result = await conn.fetchrow(select_query, inserted_user_id)

            return True, dict(result)  # Return success and fetched result as a dictionary
        finally:
            await conn.close()
    except Exception as e:
        print(f"Load test operation failed: {e}")
        return False, None

def save_data(df, out_path, recording_name):
    """Save as CSV file."""
    service_path = Path(out_path)
    target_folder = service_path
    target_folder.mkdir(parents=True, exist_ok=True)
    csv_path = target_folder / f"{recording_name}.csv"
    df.to_csv(csv_path, index=False)


async def main():
    """Main function to execute load test and capture results.
    """
    results = []

    print("Started load test!...")
    for i in range(1, 500):
        start_time = time.time()
        start = time.strftime("%Y-%m-%d %H:%M:%S", time.localtime(start_time))

        # Execute load test operations
        success, result = await load_test(user_id=i)

        end_time = time.time()
        end = time.strftime("%Y-%m-%d %H:%M:%S", time.localtime(end_time))

        delta = end_time - start_time

        # Append results
        results.append({
            "run": i,
            "start": start,
            "end": end,
            "time_taken": delta,
            "success": success,
        })

    # Save
    df = pd.DataFrame(results)
    save_data(df, Path("data/out"), "load_test")

    clean_database(with_dump=False)
    print("Finished load test!")

if __name__ == "__main__":
    asyncio.run(main())