# How to Run the Script to Upload a Folder to S3
This guide explains the prerequisites and steps to run the provided Python script to upload an entire folder to an Amazon S3 bucket.

## Prerequisites
### 1. Python Environment
Ensure you have Python 3.7 or later installed on your system.
You can check the Python version by running:
```bash
python --version
```
### 2. Install Required Libraries
Install the required libraries using pip. Run the following command:
```bash
pip install -r requirements.txt
```
### 3. AWS Credentials
Configure AWS credentials to authenticate with the S3 service. You can set this up in one of the following ways:
- Use the AWS CLI to configure credentials:

```bash
aws configure
```
Provide your Access Key ID, Secret Access Key, region, and default output format.

- Alternatively, save the credentials in ~/.aws/credentials:
```bash
aws_access_key_id=YOUR_ACCESS_KEY
aws_secret_access_key=YOUR_SECRET_KEY
```
### 4. S3 Bucket
Ensure that the target S3 bucket exists and that your AWS user has permissions to upload files to it.

## Script Usage
### Command-Line Arguments
The script accepts the following arguments:

1. bucket (required): The name of the S3 bucket where files will be uploaded.
2. folder (required): The path to the local folder to upload.
3. --prefix (optional): The prefix (path in S3) where the files will be stored. Defaults to the root of the bucket.

### Example Commands
#### Basic Usage:
Upload a folder /my/local/folder to the S3 bucket my-bucket:

```bash
python3 upload_folder_to_s3.py my-bucket /my/local/folder
````

Upload the same folder to the data/ prefix in the bucket:

```bash
python3 upload_folder_to_s3.py my-bucket /my/local/folder --prefix data
```
This will upload files to S3 in the following structure:

```bash
s3://my-bucket/data/<file_structure_from_local_folder>
```
## Output
The script prints the progress of each file uploaded. If successful, you’ll see:

```bash
Uploaded /my/local/folder/file1.txt to s3://my-bucket/backup/file1.txt
Uploaded /my/local/folder/subfolder/file2.txt to s3://my-bucket/backup/subfolder/file2.txt
All files uploaded successfully.
```
If the upload fails, the error message will provide details, such as missing credentials or the inability to find the specified folder.

