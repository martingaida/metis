import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../services/api.service';
import { Topic } from '../services/api.service';

interface Layer {
  what: string;
  why: string;
  how: string;
}

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
  mainTakeaway: string = '';
  layers: Layer[] = [
    {
      what: "Some description",
      why: "Some reason",
      how: "Some method",
    },
  ];

  constructor(private apiService: ApiService) {}

  explainText() {
    this.apiService.explainText(this.inputText).subscribe(
      (response) => {
        this.explanation = response.explanations.topics;
        this.mainTakeaway = response.explanations.main_takeaway;
      },
      (error) => {
        console.error('Error:', error);
        this.explanation = [];
        this.mainTakeaway = '';
      }
    );
  }
}