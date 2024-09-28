from handlers.arxiv_scraper import ArXivScraperClient
import json


def lambda_handler(event, context):
    scraper = ArXivScraperClient()
    papers = scraper.get_random_papers()
    print(papers)
    return {
        'statusCode': 200,
        'body': json.dumps(papers, indent=2)
    }