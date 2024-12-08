import boto3
import os
from botocore.exceptions import NoCredentialsError, PartialCredentialsError

def upload_folder_to_s3(bucket_name, folder_path, s3_prefix=""):
    """
    Upload an entire folder to an S3 bucket.

    :param bucket_name: Name of the target S3 bucket.
    :param folder_path: Path to the local folder to upload.
    :param s3_prefix: Prefix for the S3 object keys (optional).
    :return: True if all files were uploaded successfully, False otherwise.
    """
    s3_client = boto3.client('s3')

    try:
        # Walk through the folder, and upload each file
        for root, _, files in os.walk(folder_path):
            for file_name in files:
                # Full path to the local file
                local_path = os.path.join(root, file_name)

                # Relative path within the folder
                relative_path = os.path.relpath(local_path, folder_path)

                # Construct the S3 object key
                s3_key = os.path.join(s3_prefix, relative_path).replace("\\", "/")  # Use '/' for S3 key format

                # Upload the file
                s3_client.upload_file(local_path, bucket_name, s3_key)
                print(f"Uploaded {local_path} to s3://{bucket_name}/{s3_key}")

        print("All files uploaded successfully.")
        return True
    except FileNotFoundError as e:
        print(f"Error: {e}")
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


if __name__ == "__main__":
    import argparse

    # Set up command-line arguments
    parser = argparse.ArgumentParser(description="Upload an entire folder to an S3 bucket.")
    parser.add_argument("bucket", help="Name of the S3 bucket")
    parser.add_argument("folder", help="Path to the local folder to upload")
    parser.add_argument("--prefix", help="Prefix for the S3 keys (optional)", default="")

    args = parser.parse_args()

    # Upload the folder
    success = upload_folder_to_s3(args.bucket, args.folder, args.prefix)
    if success:
        print("Upload completed successfully.")
    else:
        print("Upload failed.")
