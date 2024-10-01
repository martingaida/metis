import json
from handlers import explain

def lambda_handler(event, context):
    if 'body' in event:
        try:
            body = json.loads(event['body'])
            text = body.get('text', '')
            level = body.get('level', 'Basic')
            generate_images = body.get('generate_images', False)

            result = explain.generate_response(text, level, generate_images)

            # Convert Layer objects to dictionaries
            for topic in result['topics']:
                for concept in topic['concepts']:
                    concept['layer'] = concept['layer'].__dict__
                    
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
        return {
            'statusCode': 400,
            'body': json.dumps({'error': 'Invalid request'})
        }