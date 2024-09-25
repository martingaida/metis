import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../services/api.service';
import { Topic } from '../services/api.service';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './home.component.html',
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