import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterOutlet } from '@angular/router';
import { ExplainComponent } from './explain/explain.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterOutlet, ExplainComponent],
  template: `
    <div class="app-container">
      <img src="./assets/images/plex-logo.svg" alt="PLEX Logo" class="logo">
      <app-explain></app-explain>
    </div>
  `,
  styles: [`
    .app-container {
      display: flex;
      flex-direction: column;
      align-items: center;
      width: 100%;
      max-width: 800px;
      margin: 0 auto;
      padding: 1rem;
    }
    .logo {
      max-width: 200px;
      width: 100%;
      height: auto;
      margin-top: 2rem;
      margin-bottom: 2rem;
    }
  `],
})
export class AppComponent {
  title = 'PLEX';
}