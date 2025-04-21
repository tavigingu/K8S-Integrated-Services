import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { FormsModule } from '@angular/forms'; // Pentru ngModel

@Component({
  selector: 'app-set-username',
  standalone: true,
  imports: [FormsModule], // Importăm FormsModule direct
  templateUrl: './set-username.component.html',
  styleUrls: ['./set-username.component.css']
})
export class SetUsernameComponent {
    username: string = '';

  constructor(private router: Router) {}

  setUsername() {
    if (this.username.trim()) {
      localStorage.setItem('username', this.username);
      this.router.navigate(['/chat']);
    } else {
      alert('Te rog introdu un nume de utilizator!');
    }
  }
}