import boto3
import time
import logging

# configure logging.
logger = logging.getLogger()
logger.setLevel(logging.INFO)

def lambda_handler(event, context):
    
    logger.info(f"Received event: {event}")

    # extract the S3 bucket name and object key from the event.
    bucket_name = event['Records'][0]['s3']['bucket']['name']
    object_key = event['Records'][0]['s3']['object']['key']
    log_group_name = 'REDACTED'
    log_stream_name = 'REDACTED'

    # create S3 and CloudWatch Logs clients.
    s3 = boto3.client('s3')
    logs = boto3.client('logs')

    try:
        # get the object content.
        s3_object = s3.get_object(Bucket=bucket_name, Key=object_key)
        file_content = s3_object['Body'].read().decode('utf-8')

        # send the file content to CloudWatch Logs.
        logs.put_log_events(
            logGroupName=log_group_name,
            logStreamName=log_stream_name,
            logEvents=[
                {
                    'timestamp': int(time.time() * 1000),
                    'message': file_content
                }
            ]
        )
        logger.info(f"Processed file: {object_key}")

    except Exception as e:
        logger.error(f"Error processing file '{object_key}' from bucket '{bucket_name}': {str(e)}")
