import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AuthService } from '../services/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="home-container">
      <h1>Welcome to Quizgenie!</h1>
      <button (click)="logout()">Log out</button>
    </div>
  `,
  styles: [`
    .home-container {
      max-width: 400px;
      margin: 100px auto;
      text-align: center;
    }
    button {
      margin-top: 20px;
      padding: 8px 16px;
    }
  `]
})
export class HomeComponent {
  constructor(private auth: AuthService, private router: Router) {}

  logout() {
    this.auth.logout();
  }
}