import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { ApiService, Topic, Concept, Layer } from '../services/api.service';

@Component({
  selector: 'app-explain',
  standalone: true,
  imports: [CommonModule, FormsModule, MatProgressBarModule],
  templateUrl: './explain.component.html',
  styleUrls: ['./explain.component.scss']
})
export class ExplainComponent {
  inputText = '';
  explanations: Topic[] = [];
  mainTakeaway: string = '';
  totalTime: number = 0;
  isLoading = false;

  constructor(private apiService: ApiService) {}

  explainText() {
    this.isLoading = true;
    console.log('Explaining text:', this.inputText);
    this.apiService.explainText(this.inputText).subscribe(
      (response) => {
        console.log('Received response:', response);
        this.explanations = response.explanations.topics;
        this.mainTakeaway = response.explanations.main_takeaway;
        this.isLoading = false;
      },
      (error) => {
        console.error('Error:', error);
        // Handle error appropriately
        this.isLoading = false;
      }
    );
  }
}