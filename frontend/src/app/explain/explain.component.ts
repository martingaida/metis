import { ApiService, Topic, Concept, Layer, ArXivPaper } from '../services/api.service';
import { trigger, transition, style, animate } from '@angular/animations';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { Component, ViewChild, ElementRef, OnInit } from '@angular/core';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatTabsModule } from '@angular/material/tabs';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { MatButtonModule } from '@angular/material/button';

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
    MatIconModule,
    MatButtonToggleModule,
    MatButtonModule
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

export class ExplainComponent implements OnInit {
  inputText = '';
  explanations: Topic[] = [];
  mainTakeaway: string = '';
  totalTime: number = 0;
  isLoading = false;
  isExplanationVisible = false;
  mode: 'arXiv' | 'Custom' = 'arXiv';
  arXivPapers: ArXivPaper[] = [];
  isLoadingArXiv = false; // New property to track arXiv loading state

  constructor(private apiService: ApiService) {}

  @ViewChild('explanationContainer') private explanationContainer!: ElementRef;

  ngOnInit() {
    this.loadArXivPapers();
  }

  loadArXivPapers() {
    this.isLoadingArXiv = true; // Set to true when starting to load
    this.apiService.getArXivPapers().subscribe(
      (papers) => {
        this.arXivPapers = papers;
        this.isLoadingArXiv = false; // Set to false when loading is complete
      },
      (error) => {
        console.error('Error fetching arXiv papers:', error);
        this.isLoadingArXiv = false; // Set to false if there's an error
      }
    );
  }

  get isExplainDisabled(): boolean {
    return this.isLoading || this.isLoadingArXiv || (!this.inputText.trim() && this.mode === 'Custom');
  }

  explainText() {
    if (this.isExplainDisabled) return;

    this.isLoading = true;
    this.isExplanationVisible = false;
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

  explainArXiv(paper: ArXivPaper) {
    this.inputText = `${paper.title}\n\n${paper.abstract}`;
    this.explainText();
  }
}