import boto3
from botocore.exceptions import NoCredentialsError, PartialCredentialsError

def upload_csv_to_s3(bucket_name, file_path, object_name=None):
    """
    Upload a CSV file to an S3 bucket.

    :param bucket_name: Name of the target S3 bucket.
    :param file_path: Path to the CSV file to upload.
    :param object_name: Key name for the file in S3 (defaults to the file name).
    :return: True if the file was uploaded successfully, False otherwise.
    """
    # Initialize S3 client
    s3_client = boto3.client('s3')
    
    # Default the object name to the file name if not provided
    if object_name is None:
        object_name = file_path.split('/')[-1]
    
    try:
        # Upload the file to the bucket
        s3_client.upload_file(file_path, bucket_name, object_name)
        print(f"File {file_path} uploaded to {bucket_name}/{object_name}")
        return True
    except FileNotFoundError:
        print(f"Error: File {file_path} not found.")
        return False
    except NoCredentialsError:
        print("Error: AWS credentials not found.")
        return False
    except PartialCredentialsError:
        print("Error: Incomplete AWS credentials.")
        return False
    except Exception as e:
        print(f"An unexpected error occurred: {e}")
        return False

# Example of how this function can be used dynamically
if __name__ == "__main__":
    import argparse
    
    # Set up command-line arguments
    parser = argparse.ArgumentParser(description="Upload a CSV file to an S3 bucket.")
    parser.add_argument("bucket", help="Name of the S3 bucket")
    parser.add_argument("file", help="Path to the CSV file to upload")
    parser.add_argument("--key", help="Key name for the file in S3 (optional)", default=None)
    
    args = parser.parse_args()
    
    # Upload the file
    success = upload_csv_to_s3(args.bucket, args.file, args.key)
    if success:
        print("Upload successful.")
    else:
        print("Upload failed.")
