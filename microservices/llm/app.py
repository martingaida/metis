import json
from handlers import explain


def lambda_handler(event, context):
    try:
        # Extract text input from the POST request body
        body = json.loads(event['body'])
        text = body.get('text', '')

        # Call Explanation handler
        result = explain.generate_response(text)

        return {
            'statusCode': 200,
            'body': json.dumps({'explanation': result})
        }
    except Exception as e:
        return {
            'statusCode': 500,
            'body': json.dumps({'error': str(e)})
        }