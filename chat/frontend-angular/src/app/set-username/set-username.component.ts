import { Component, Inject, PLATFORM_ID } from '@angular/core';
import { Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { isPlatformBrowser } from '@angular/common';

@Component({
  selector: 'app-set-username',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './set-username.component.html',
  styleUrls: ['./set-username.component.css']
})
export class SetUsernameComponent {
  username: string = '';
  private isBrowser: boolean;

  constructor(
    private router: Router,
    @Inject(PLATFORM_ID) platformId: Object
  ) {
    this.isBrowser = isPlatformBrowser(platformId);
  }

  setUsername() {
    if (this.username.trim()) {
      if (this.isBrowser) {
        localStorage.setItem('username', this.username);
        this.router.navigate(['/chat']);
      }
    } else {
      if (this.isBrowser) {
        alert('Te rog introdu un nume de utilizator!');
      }
    }
  }
}