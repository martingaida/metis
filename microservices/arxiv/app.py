from handlers.arxiv_client import ArXivClient
import json


def lambda_handler(event, context):
    try:
        scraper = ArXivClient()
        papers = scraper.get_random_papers(num_categories=3, papers_per_category=1)
        
        # Ensure we're returning an array of papers
        if not isinstance(papers, list):
            papers = []
            
        print(f'Returning {len(papers)} papers from cache')
        
        return {
            'statusCode': 200,
            'headers': {
                'Content-Type': 'application/json',
                'Access-Control-Allow-Origin': '*'
            },
            'body': json.dumps(papers)
        }
        
    except Exception as e:
        print(f"Error in lambda_handler: {str(e)}")
        return {
            'statusCode': 500,
            'headers': {
                'Content-Type': 'application/json',
                'Access-Control-Allow-Origin': '*'
            },
            'body': json.dumps([])
        }