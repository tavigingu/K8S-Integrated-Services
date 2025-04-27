import { Component } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-navigation',
  templateUrl: './navigation.component.html',
  styleUrls: ['./navigation.component.css'],
  standalone: false
  // Remove standalone: true property
})
export class NavigationComponent {
  constructor(private router: Router) {}
}