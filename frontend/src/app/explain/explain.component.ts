import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApiService, Topic, Concept, Layer } from '../services/api.service';

@Component({
  selector: 'app-explain',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="container mt-5">
      <div class="mb-3">
        <label for="inputText" class="form-label">Enter text to explain:</label>
        <textarea class="form-control" id="inputText" rows="3" [(ngModel)]="inputText"></textarea>
      </div>
      <button class="btn btn-primary" (click)="explainText()">Explain</button>
      <div *ngIf="explanations.length > 0" class="mt-3">
        <h3>Explanations:</h3>
        <div *ngFor="let topic of explanations">
          <h4>{{ topic.topic }}</h4>
          <div *ngFor="let concept of topic.concepts">
            <h5>{{ concept.concept }}</h5>
            <div *ngFor="let layer of concept.layers; let i = index">
              <h6>Layer {{ i + 1 }}</h6>
              <p><strong>What:</strong> {{ layer.What }}</p>
              <p><strong>Why:</strong> {{ layer.Why }}</p>
              <p><strong>How:</strong> {{ layer.How }}</p>
            </div>
          </div>
        </div>
      </div>
      <div *ngIf="mostSignificantTakeaway" class="mt-3">
        <h3>Most Significant Takeaway:</h3>
        <p>{{ mostSignificantTakeaway }}</p>
      </div>
      <div *ngIf="totalTime" class="mt-3">
        <p>Total processing time: {{ totalTime }} seconds</p>
      </div>
    </div>
  `,
  styles: []
})
export class ExplainComponent {
  inputText = '';
  explanations: Topic[] = [];
  mostSignificantTakeaway: string = '';
  totalTime: number = 0;

  constructor(private apiService: ApiService) {}

  explainText() {
    console.log('Explaining text:', this.inputText);
    this.apiService.explainText(this.inputText).subscribe(
      (response) => {
        console.log('Received response:', response);
        this.explanations = response.explanations.topics;
        this.mostSignificantTakeaway = response.explanations.most_significant_takeaway;
      },
      (error) => {
        console.error('Error:', error);
        // Handle error appropriately
      }
    );
  }
}