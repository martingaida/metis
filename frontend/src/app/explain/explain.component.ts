import { ApiService, Topic, Concept, Layer } from '../services/api.service';
import { trigger, transition, style, animate } from '@angular/animations';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { Component, ViewChild, ElementRef } from '@angular/core';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatTabsModule } from '@angular/material/tabs';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-explain',
  standalone: true,
  imports: [
    CommonModule, 
    FormsModule, 
    MatProgressBarModule, 
    MatCardModule, 
    MatExpansionModule, 
    MatTabsModule, 
    MatTooltipModule, 
    MatIconModule
  ],
  templateUrl: './explain.component.html',
  styleUrls: ['./explain.component.scss'],
  animations: [
    trigger('fadeIn', [
      transition(':enter', [
        style({ opacity: 0, transform: 'translateY(10px)' }),
        animate('300ms ease-out', style({ opacity: 1, transform: 'translateY(0)' }))
      ])
    ])
  ]
})
export class ExplainComponent {
  inputText = '';
  explanations: Topic[] = [];
  mainTakeaway: string = '';
  totalTime: number = 0;
  isLoading = false;
  isExplanationVisible = false;

  constructor(private apiService: ApiService) {}

  @ViewChild('explanationContainer') private explanationContainer!: ElementRef;

  explainText() {
    this.isLoading = true;
    this.isExplanationVisible = false; // Hide explanation while loading
    console.log('Explaining text:', this.inputText);
    this.apiService.explainText(this.inputText).subscribe(
      (response) => {
        console.log('Received response:', response);
        this.explanations = response.explanations.topics;
        this.mainTakeaway = response.explanations.main_takeaway;
        this.isLoading = false;
        setTimeout(() => {
          this.isExplanationVisible = true;
          this.scrollToExplanation();
        }, 100);
      },
      (error) => {
        console.error('Error:', error);
        this.isLoading = false;
      }
    );
  }

  private scrollToExplanation() {
    if (this.explanationContainer) {
      this.explanationContainer.nativeElement.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
  }
}