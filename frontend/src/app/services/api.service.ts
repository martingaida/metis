import { Injectable, Inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Layer {
  layer_name: string;
  explanation: string;
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
  explanations: Topic[];
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
    return this.http.post<ExplanationResponse>(`${this.apiUrl}/api/explain`, { text });
  }
}