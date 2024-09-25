import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../services/api.service';
import { Topic } from '../services/api.service';

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
      <div *ngIf="explanation.length > 0" class="mt-3">
        <h2>Explanation:</h2>
        <div *ngFor="let topic of explanation">
          <h3>{{ topic.topic }}</h3>
          <div *ngFor="let concept of topic.concepts">
            <h4>{{ concept.concept }}</h4>
            <div *ngFor="let layer of concept.layers">
              <h5>{{ layer.layer_name }}</h5>
              <p>{{ layer.explanation }}</p>
            </div>
          </div>
        </div>
      </div>
      <div *ngIf="mostSignificantTakeaway" class="mt-3">
        <h3>Most Significant Takeaway:</h3>
        <p>{{ mostSignificantTakeaway }}</p>
      </div>
    </div>
  `,
  styles: []
})
export class HomeComponent {
  inputText = '';
  explanation: Topic[] = [];
  mostSignificantTakeaway: string = '';

  constructor(private apiService: ApiService) {}

  explainText() {
    this.apiService.explainText(this.inputText).subscribe(
      (response) => {
        this.explanation = response.explanations.topics;
        this.mostSignificantTakeaway = response.explanations.most_significant_takeaway;
      },
      (error) => {
        console.error('Error:', error);
        this.explanation = [];
        this.mostSignificantTakeaway = '';
      }
    );
  }
}