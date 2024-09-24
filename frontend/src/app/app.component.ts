import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterOutlet } from '@angular/router';
import { ExplainTextComponent } from './explain/explain.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterOutlet, ExplainTextComponent],
  template: `
    <main class="container">
      <h1>PLEX/h1>
      <h4>Please Explain</h4>
      <app-explain></app-explain>
    </main>
  `,
  styles: [],
})
export class AppComponent {
  title = 'PLEX';
}