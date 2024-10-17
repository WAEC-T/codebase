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
psql -h waectdatabase-waectresultstoragedbe1ccf3e0-tmed4lgituj0.c10ia6ywc903.eu-central-1.rds.amazonaws.com -p 5432 -U waect -d waectResultStorage
````

4. **Enter Password**:
When prompted, enter the password retrieved from the secret.

## Security Considerations

- The database is configured to allow connections from any IP address on port 5432. Ensure your network security measures are appropriate for your use case.
- The database instance is deployed in a public subnet. Consider using a bastion host or VPN for added security in production environments.
- Regularly rotate the database password by updating the `AWS_DATABASE_PASSWORD` environment variable and redeploying the stack.

## Troubleshooting

- If you cannot connect, verify that your IP is allowed in the security group rules.
- Ensure the `AWS_DATABASE_PASSWORD` environment variable was set correctly during stack deployment.
- Check that the VPC and security group configurations allow incoming traffic on port 5432.

Remember to handle the database credentials securely and avoid exposing them in your application code or version control systems.

