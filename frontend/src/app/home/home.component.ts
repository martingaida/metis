import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../services/api.service';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="container mt-5">
      <h1>Text Explainer</h1>
      <div class="mb-3">
        <label for="textInput" class="form-label">Enter text to explain:</label>
        <textarea class="form-control" id="textInput" rows="3" [(ngModel)]="inputText"></textarea>
      </div>
      <button class="btn btn-primary" (click)="explainText()">Explain</button>
      <div *ngIf="explanation" class="mt-3">
        <h2>Explanation:</h2>
        <p>{{ explanation }}</p>
      </div>
    </div>
  `,
  styles: []
})
export class HomeComponent {
  inputText = '';
  explanation = '';

  constructor(private apiService: ApiService) {}

  explainText() {
    this.apiService.explainText(this.inputText).subscribe(
      (response) => {
        this.explanation = response.explanation;
      },
      (error) => {
        console.error('Error:', error);
        this.explanation = 'An error occurred while explaining the text.';
      }
    );
  }
}