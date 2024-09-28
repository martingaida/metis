import random
import arxiv


class ArXivScraperClient:
    def __init__(self):
        self.categories = [
            'astro-ph', 'cond-mat', 'gr-qc', 'hep-ex', 'hep-lat', 
            'hep-ph', 'hep-th', 'math-ph', 'nlin', 'nucl-ex', 
            'nucl-th', 'physics', 'quant-ph', 'math', 'cs'
        ]

    def get_random_papers(self, num_categories=4, papers_per_category=1):
        selected_categories = random.sample(self.categories, num_categories)
        papers = []

        for category in selected_categories:
            search = arxiv.Search(
                query=f"cat:{category}",
                max_results=papers_per_category,
                sort_by=arxiv.SortCriterion.SubmittedDate
            )
            for result in search.results():
                papers.append({
                    'id': result.entry_id.split('/')[-1],
                    'title': result.title,
                    'abstract': self.extract_abstract(result.summary),
                    'category': category,
                    'authors': ', '.join(author.name for author in result.authors),
                    'published': result.published.strftime("%Y-%m-%d"),
                    'abstract_url': result.entry_id,
                    'pdf_url': result.pdf_url
                })

        return papers

    @staticmethod
    def extract_abstract(text):
        return text.replace('\n', ' ').strip()