import { ApiService, ExplanationResponse, ArXivPaper, Topic } from '../services/api.service';
import { trigger, transition, style, animate } from '@angular/animations';
import { Component, ViewChild, ElementRef, OnInit, OnChanges, SimpleChanges, Input, OnDestroy } from '@angular/core';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatButtonModule } from '@angular/material/button';
import { MatTabsModule } from '@angular/material/tabs';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Subject } from 'rxjs';
import { takeUntil, catchError, finalize } from 'rxjs/operators';
import { of } from 'rxjs';

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
    MatButtonModule,
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

export class ExplainComponent implements OnInit, OnChanges, OnDestroy {
  @Input() selectedLevel: string = 'Basic';
  explainedPapers: { [id: string]: { [level: string]: ExplanationResponse } } = {};
  customExplanations: { [text: string]: { [level: string]: ExplanationResponse } } = {};
  inputText = '';
  explanations: Topic[] = [];
  mainTakeaway: string = '';
  totalTime: number = 0;
  isLoading = false;
  isExplanationVisible = false;
  mode: 'arXiv' | 'Custom' = 'arXiv';
  arXivPapers: ArXivPaper[] = [];
  isLoadingArXiv = false;
  maxRetries = 3;
  currentPaperTitle: string | null = null;
  currentPaperPdfUrl: string | null = null;
  private destroy$ = new Subject<void>();
  errorMessage: string | null = null;

  levels: { value: string, viewValue: string }[] = [
    { value: 'Basic', viewValue: 'K3' },
    { value: 'Elementary', viewValue: 'K6' },
    { value: 'Middle School', viewValue: 'K9' },
    { value: 'High School', viewValue: 'K12' },
    { value: 'College', viewValue: 'College' },
    { value: 'Graduate', viewValue: 'Graduate' }
  ];

  constructor(private apiService: ApiService) {}

  @ViewChild('explanationContainer') private explanationContainer!: ElementRef;

  ngOnInit() {
    this.loadArXivPapers();
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes['selectedLevel'] && !changes['selectedLevel'].firstChange) {
      this.resetExplanationState();
      if (this.mode === 'arXiv') {
        this.updateArXivExplanation();
      } else if (this.mode === 'Custom' && this.inputText.trim()) {
        this.updateCustomExplanation();
      }
    }
  }

  ngOnDestroy() {
    // Unsubscribe from all observables when the component is destroyed
    this.destroy$.next();
    this.destroy$.complete();
  }

  loadArXivPapers() {
    if (this.arXivPapers.length > 0) return;

    this.isLoadingArXiv = true;
    this.fetchArXivPapers(this.maxRetries);
  }

  fetchArXivPapers(retriesLeft: number) {
    this.apiService.getArXivPapers().pipe(takeUntil(this.destroy$)).subscribe(
      (papers) => {
        console.log('Received arXiv papers:', papers);
        this.arXivPapers = papers;
        this.isLoadingArXiv = false;
      },
      (error) => {
        console.error('Error fetching arXiv papers:', error);
        if (retriesLeft > 0) {
          console.log(`Retrying... ${retriesLeft} attempts left`);
          setTimeout(() => this.fetchArXivPapers(retriesLeft - 1), 2000);
        } else {
          this.isLoadingArXiv = false;
          this.errorMessage = 'Failed to load arXiv papers. Please try again later.';
        }
      }
    );
  }

  get isExplainDisabled(): boolean {
    return this.isLoading || this.isLoadingArXiv || (!this.inputText.trim() && this.mode === 'Custom');
  }

  explainText() {
    this.currentPaperTitle = null;
    if (this.isExplainDisabled) return;

    const trimmedText = this.inputText.trim();

    this.checkOrFetchExplanation(
      trimmedText,
      () => this.apiService.explainText(trimmedText, this.selectedLevel),
      this.customExplanations
    );
  }

  explainArXiv(paper: ArXivPaper) {
    this.currentPaperTitle = paper.title;
    this.currentPaperPdfUrl = paper.pdf_url;

    this.checkOrFetchExplanation(
      paper.id,
      () => this.apiService.explainText(`${paper.title}\n\n${paper.abstract}`, this.selectedLevel),
      this.explainedPapers
    );
  }

  private checkOrFetchExplanation(key: string, fetchCallback: () => any, cache: { [key: string]: { [level: string]: ExplanationResponse } }) {
    if (cache[key] && cache[key][this.selectedLevel]) {
      const savedExplanation = cache[key][this.selectedLevel];
      this.updateExplanationDisplay(savedExplanation);
    } else {
      this.isLoading = true;
      this.isExplanationVisible = false;
      this.errorMessage = null;
      fetchCallback().pipe(
        catchError(error => {
          this.errorMessage = 'An error occurred while fetching the explanation. Please try again.';
          console.error('Error fetching explanation:', error);
          return of(null);
        }),
        finalize(() => this.isLoading = false),
        takeUntil(this.destroy$)
      ).subscribe(
        (response: ExplanationResponse | null) => {
          if (response) {
            if (!cache[key]) {
              cache[key] = {};
            }
            cache[key][this.selectedLevel] = response;
            this.updateExplanationDisplay(response);
          }
        }
      );
    }
  }

  isPaperExplained(paperId: string): boolean {
    return !!this.explainedPapers[paperId] && !!this.explainedPapers[paperId][this.selectedLevel];
  }

  private updateArXivExplanation() {
    if (this.currentPaperTitle) {
      const currentPaper = this.arXivPapers.find(paper => paper.title === this.currentPaperTitle);
      if (currentPaper) {
        this.explainArXiv(currentPaper);
      }
    }
  }

  private updateCustomExplanation() {
    const trimmedText = this.inputText.trim();
    if (this.customExplanations[trimmedText] && this.customExplanations[trimmedText][this.selectedLevel]) {
      const savedExplanation = this.customExplanations[trimmedText][this.selectedLevel];
      this.updateExplanationDisplay(savedExplanation);
    }
  }

  private updateExplanationDisplay(explanation: ExplanationResponse) {
    this.explanations = explanation.explanations;
    this.mainTakeaway = explanation.main_takeaway;
    this.isLoading = false;
    setTimeout(() => {
      this.isExplanationVisible = true;
      this.scrollToExplanation();
      this.loadConceptImages();
    }, 100);
  }

  private scrollToExplanation() {
    if (this.explanationContainer) {
      this.explanationContainer.nativeElement.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
  }

  private resetExplanationState() {
    this.isExplanationVisible = false;
    this.explanations = [];
    this.mainTakeaway = '';
  }

  private loadConceptImages() {
    // Implement lazy loading of images here
  }
}