from handlers.arxiv_scraper import ArXivScraperClient
import json


def lambda_handler(event, context):
    scraper = ArXivScraperClient()
    papers = scraper.get_random_papers(num_categories=4, papers_per_category=1)

    return {
        'statusCode': 200,
        'body': json.dumps(papers, indent=2)
    }