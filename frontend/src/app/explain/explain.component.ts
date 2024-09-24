import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../services/api.service';

@Component({
  selector: 'app-explain',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="container mt-5">
      <h2>Explain Text</h2>
      <div class="mb-3">
        <label for="inputText" class="form-label">Enter text to explain:</label>
        <textarea class="form-control" id="inputText" rows="3" [(ngModel)]="inputText"></textarea>
      </div>
      <button class="btn btn-primary" (click)="explainText()">Explain</button>
      <div *ngIf="explanation" class="mt-3">
        <h3>Explanation:</h3>
        <p>{{ explanation }}</p>
      </div>
    </div>
  `,
  styles: []
})

export class ExplainComponent {
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