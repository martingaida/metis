import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApiService, Topic, Concept, Layer } from '../services/api.service';

@Component({
  selector: 'app-explain',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './explain.component.html',
  styleUrls: ['./explain.component.scss']
})
export class ExplainComponent {
  inputText = '';
  explanations: Topic[] = [];
  mainTakeaway: string = '';
  totalTime: number = 0;

  constructor(private apiService: ApiService) {}

  explainText() {
    console.log('Explaining text:', this.inputText);
    this.apiService.explainText(this.inputText).subscribe(
      (response) => {
        console.log('Received response:', response);
        this.explanations = response.explanations.topics;
        this.mainTakeaway = response.explanations.mainTakeaway;
      },
      (error) => {
        console.error('Error:', error);
        // Handle error appropriately
      }
    );
  }
}