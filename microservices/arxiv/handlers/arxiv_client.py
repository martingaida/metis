import random
import arxiv
import json
import time
import os


class ArXivClient:
    def __init__(self):
        # Get the directory where this file is located
        current_dir = os.path.dirname(os.path.abspath(__file__))
        self.cache_file = os.path.join(current_dir, 'arxiv_papers.json')
        self.categories = [
            'astro-ph', 'cond-mat', 'gr-qc', 'hep-ex', 'hep-lat', 
            'hep-ph', 'hep-th', 'math-ph', 'nlin', 'nucl-ex', 
            'nucl-th', 'physics', 'quant-ph', 'math', 'cs'
        ]

    def get_random_papers(self, num_categories=3, papers_per_category=1):
        """Load papers directly from cache file"""
        try:
            with open(self.cache_file, 'r') as f:
                cached_papers = json.load(f)
                
            # Convert the cached data into the expected format
            papers = []
            for paper in cached_papers:
                formatted_paper = {
                    'id': paper.get('id', ''),
                    'title': paper.get('title', ''),
                    'abstract': paper.get('abstract', ''),
                    'category': paper.get('category', ''),
                    'authors': paper.get('authors', ''),
                    'published': paper.get('published', ''),
                    'abstract_url': paper.get('abstract_url', ''),
                    'pdf_url': paper.get('pdf_url', '')
                }
                papers.append(formatted_paper)
            
            # Randomly select papers
            if papers:
                selected_papers = random.sample(papers, min(num_categories * papers_per_category, len(papers)))
                return selected_papers
            else:
                return []
                
        except Exception as e:
            print(f"Error loading cached papers: {str(e)}")
            return []

    def get_random_papers_cache(self, cache_path, num_papers=3):
        with open(cache_path, 'r') as f:
            all_papers = json.load(f)

        # Ensure we don't try to sample more papers than are available
        num_papers = min(num_papers, len(all_papers))

        # Randomly sample the specified number of papers
        selected_papers = random.sample(all_papers, num_papers)

        return selected_papers

    def get_random_papers_from_arxiv(self, num_categories=4, papers_per_category=1):
        if os.path.exists(self.cache_file):
            print('Loading papers from cache...')
            time.sleep(1)
            return self.get_random_papers_cache(self.cache_file, 3)
        
        print('Cache not found. Fetching from arXiv...')
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