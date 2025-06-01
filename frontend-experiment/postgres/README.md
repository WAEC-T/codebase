# Postgres Container

The postgres database is a modified standard postgres image edited with the [schema.sql](./schema.sql) for the MiniTwit implementations, as well as a default set of credentials for connecting to the database

See the [dockerfile](./dockerfile) for more information

## Default connection string

- Username: user
- Password: password
- Database name: minitwit

`postgresql://postgres:1234@postgres-implementations:5432/postgres`

## How to edit the connection string

### Modify credentials

To modify the credentials used, add an environment variable modifying a specific credential, e.g. changing the user to `updated`:

`docker run --detach --env POSTGRES_USER=updated --name postgres ghcr.io/3-1-research-project/postgres:latest`

### Modify the schema

By modifying the schema.sql, you can modify the schema applied to the database
