import json
import boto3
import os
from handlers import explain


def lambda_handler(event, context):
    if 'body' in event:
        # This is an HTTP request
        try:
            body = json.loads(event['body'])
            text = body.get('text', '')

            # Process the text synchronously
            result = explain.process_text(text)

            return {
                'statusCode': 200,
                'body': json.dumps(result)
            }
        except Exception as e:
            return {
                'statusCode': 500,
                'body': json.dumps({'error': str(e)})
            }
    else:
        # This is not an HTTP request (should not happen with Function URL)
        return {
            'statusCode': 400,
            'body': json.dumps({'error': 'Invalid request'})
        }