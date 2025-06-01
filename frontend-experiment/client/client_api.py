from fastapi import (
    FastAPI,
    Request,
    status,
    Path,
    Body,
)
from fastapi.responses import JSONResponse
from experiment.scenario import run

# from frontend_scenarios.runner import run
import os
import logging
import psycopg2

# NAME_DESCRIPTION = "The name of the schema file"
# NAME_EXAMPLE = "my-schema"
# SCHEMA_DESCRIPTION = "The schema content"
NAME = "schema"
SCHEMA_EXAMPLE = {}
URL_DESCRIPTION = "The url of the MiniTwit application"
URL_EXAMPLE = "http://0.0.0.0:5000/"
LOGGER = logging.getLogger("uvicorn")

DATABASE_URL = os.environ.get(
    "DATABASE_URL"
)  # "postgresql://user:pass@localhost:5432/waect"


def create_folder_if_not_exists(folder_path):
    if not os.path.exists(folder_path):
        os.makedirs(folder_path)


create_folder_if_not_exists("logs")
# create_folder_if_not_exists("schemas")

app = FastAPI()


@app.exception_handler(Exception)
async def exception_handler(request: Request, exception: Exception):
    return JSONResponse(
        status_code=500,
        content={"args": str(exception.args), "message": exception},
    )


# def get_schema_file_path_from_name(name: str):
# return os.path.join("schemas", f"{name}.json")


@app.post("/start", status_code=status.HTTP_200_OK)
async def start(
    url: str = Body(..., description=URL_DESCRIPTION, example=URL_EXAMPLE),
):
    await run(url)
    return {"status": "Tests done"}


@app.post("/cleardb", status_code=status.HTTP_200_OK)
async def clear_db(
    db_url: str = Body(..., description="The url of the database", nullable=True),
):
    with psycopg2.connect(db_url or DATABASE_URL) as conn:
        with conn.cursor() as cur:
            tables_to_truncate = [
                "users",
                "messages",
                "followers",
                "AspNetRoles",
                "AspNetRoleClaims",
                "AspNetUserClaims",
                "AspNetUserLogins",
                "AspNetUserRoles",
                "AspNetUserTokens",
            ]  # AspNet tables for Identity

            for table in tables_to_truncate:
                cur.execute(
                    """
                    SELECT EXISTS (
                        SELECT FROM information_schema.tables
                        WHERE table_schema = 'public'
                        AND table_name = %s
                    );
                """,
                    (table,),
                )

                exists = cur.fetchone()[0]

                if exists:
                    cur.execute(f'TRUNCATE TABLE "{table}" CASCADE;')
                    print(f"Table {table} truncated")

            conn.commit()
