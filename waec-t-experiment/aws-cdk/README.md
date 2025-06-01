# Connecting to the PostgreSQL Database

This guide explains how to connect to the PostgreSQL database instance created using the provided AWS CDK stack.

## Database Details

- **Engine**: PostgreSQL 16.4
- **Instance Type**: t3.micro
- **Port**: 5432
- **Database Name**: waectResultStorage

## Prerequisites

1. Ensure you have access to the AWS account where this stack is deployed.
2. Install a PostgreSQL client (e.g., psql) on your local machine or the instance from which you'll connect.

## Connection Steps

1. **Retrieve Connection Information**:
   The database connection details are stored in AWS Secrets Manager. Retrieve the secret named 'waectResultStorageSECRET' using the AWS Management Console or AWS CLI.

2. **Extract Connection Details**:
   The secret contains the following information:
   - Username
   - Password
   - Hostname
   - Port
   - Database name

3. **Connect using psql**:
   Use the following command structure to connect:
```
psql -h {AWS_POSTGRES_HOST}  -p 5432 -U waect -d waect
````

4. **Enter Password**:
When prompted, enter the password retrieved from the secret.


