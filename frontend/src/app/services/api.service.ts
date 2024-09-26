import { Injectable, Inject } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Layer {
  what: string;
  why: string;
  how: string;
}

export interface Concept {
  concept: string;
  layers: Layer[];
}

export interface Topic {
  topic: string;
  concepts: Concept[];
}

export interface ExplanationResponse {
  explanations: {
    topics: Topic[];
    main_takeaway: string;
  };
}

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  constructor(
    private http: HttpClient,
    @Inject('API_URL') private apiUrl: string
  ) {}

  explainText(text: string): Observable<ExplanationResponse> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json'
    });

    return this.http.post<ExplanationResponse>(
      `${this.apiUrl}`,
      { text },
      { headers }
    );
  }
}