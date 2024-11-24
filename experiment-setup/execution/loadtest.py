import psycopg2
import time
import pandas as pd

DATABASE_URL = ""

# Define the function

def clean_database(with_dump: bool):
    try:
        with psycopg2.connect(DATABASE_URL) as conn:
            with conn.cursor() as cur:
                cur.execute("TRUNCATE TABLE users CASCADE;")
                cur.execute("TRUNCATE TABLE messages CASCADE;")
                cur.execute("TRUNCATE TABLE followers CASCADE;")
                if with_dump:
                    with open("dump.sql", "r") as f:
                        sql_commands = f.read()
                        cur.execute(sql_commands)
                    conn.commit()
        return True
    except Exception as e:
        print(f"Database cleaning failed: {e}")
        return False

# Prepare for load testing
results = []

# Execute the function 100 times and capture timing
print("Started load test!...")
for i in range(1, 101):
    start_time = time.time()
    start = time.strftime("%Y-%m-%d %H:%M:%S", time.localtime(start_time))

    # Call the function
    success = clean_database(with_dump=True)

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

# Save results to a DataFrame
df = pd.DataFrame(results)

# Save to CSV for further analysis
df.to_csv("loadtest_results.csv", index=False)

# Display DataFrame
df.head()

print("Finished load test!...")