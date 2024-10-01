import { Injectable, Inject } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { catchError } from 'rxjs/operators';

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
  topic: string;
  concepts: Concept[];
}

export interface ExplanationResponse {
  explanations: Topic[];
  main_takeaway: string;
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
    
    console.log('Requesting URL:', this.apiUrl, 'with body:', body);
    
    return this.http.post<ArXivPaper[]>(
      this.apiUrl,
      body,
      { 
        headers: this.getHeaders(),
      }
    ).pipe(
      catchError(this.handleError<ArXivPaper[]>('getArXivPapers', []))
    );
  }

  explainText(text: string, level: string): Observable<ExplanationResponse> {
    const body = { action: 'explain', text, level };

    console.log('Requesting URL:', this.apiUrl, 'with body:', body);
    
    return this.http.post<ExplanationResponse>(
      this.apiUrl,
      body,
      { 
        headers: this.getHeaders(),
      }
    ).pipe(
      catchError(this.handleError<ExplanationResponse>('explainText'))
    );
  }

  private handleError<T>(operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {
      console.error(`${operation} failed: ${error.message}`);
      return of(result as T);
    };
  }
}