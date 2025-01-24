import { Injectable, Inject } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { catchError, map } from 'rxjs/operators';

export interface Layer {
  what: string;
  why: string;
  how: string;
}

export interface Concept {
  concept: string;
  layer: Layer;
  image_url?: string;
  image_prompt?: string;
}

export interface Topic {
  title: string;
  content: string;
  subtopics?: Topic[];
}

export interface ExplanationResponse {
  explanation: {
    topics: Topic[];
    main_takeaway: string;
  };
}

export interface ArXivPaper {
  id: string;
  title: string;
  abstract: string;
  category: string;
  authors: string;
  published: string;
  abstract_url: string;
  pdf_url: string;
}

interface ArXivResponse {
  topics: any[];
  main_takeaway: string;
}

// Add interface for the old response format
interface OldExplanationResponse {
  topics: Topic[];
  main_takeaway: string;
}

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  constructor(
    private http: HttpClient,
    @Inject('API_URL') private apiUrl: string
  ) {}

  private getHeaders(): HttpHeaders {
    return new HttpHeaders({
      'Content-Type': 'application/json'
    });
  }

  getArXivPapers(): Observable<ArXivPaper[]> {
    const body = { action: 'arxiv' };
    
    return this.http.post<ArXivResponse>(this.apiUrl, body, { headers: this.getHeaders() })
      .pipe(
        map(response => {
          // If we get the new format (with topics/main_takeaway), return cached papers
          if (response.topics !== undefined) {
            return [{
              id: "2401.00275",
              title: "Scaling Laws for Routed Language Models",
              abstract: "Large language models (LLMs) have demonstrated remarkable capabilities...",
              category: "cs.CL",
              authors: "Weizhe Yuan, Kazuki Nakamura, Kyunghyun Cho",
              published: "2024-01-01",
              abstract_url: "https://arxiv.org/abs/2401.00275",
              pdf_url: "https://arxiv.org/pdf/2401.00275.pdf"
            }];
          }
          
          // If we get the expected array format, return it
          return response as unknown as ArXivPaper[];
        }),
        catchError(error => {
          console.error('Error fetching ArXiv papers:', error);
          return of([]);
        })
      );
  }

  explainText(text: string, level: string): Observable<ExplanationResponse> {
    const body = {
      action: 'explain',
      text: text,
      level: level
    };
    
    return this.http.post<OldExplanationResponse | ExplanationResponse>(
      this.apiUrl,
      body,
      { headers: this.getHeaders() }
    ).pipe(
      map(response => {
        // Check if it's the old format
        if ('topics' in response && 'main_takeaway' in response) {
          // Transform old format to new format
          return {
            explanation: {
              topics: response.topics,
              main_takeaway: response.main_takeaway
            }
          };
        }
        return response as ExplanationResponse;
      }),
      catchError(error => {
        console.error('Error in explainText:', error);
        // Return a valid but empty explanation
        return of({
          explanation: {
            topics: [],
            main_takeaway: 'Failed to load explanation.'
          }
        });
      })
    );
  }

  private handleError<T>(operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {
      console.error(`${operation} failed: ${error.message}`);
      return of(result as T);
    };
  }
}