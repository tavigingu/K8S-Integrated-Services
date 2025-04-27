import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  // Remove the standalone property
  standalone: false,
  template: `
    <app-navigation></app-navigation>
    <main class="container">
      <router-outlet></router-outlet>
    </main>
  `,
  styles: [`
    .container {
      max-width: 1200px;
      margin: 0 auto;
      padding: 20px;
    }
  `]
})
export class AppComponent {
  title = 'AI Entity Extraction';
}