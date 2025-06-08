import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule }  from '@angular/forms';
import { Router, RouterLink } from '@angular/router';

import { EnvService }  from '../../services/env.service';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [ CommonModule, FormsModule, RouterLink ],
  templateUrl: './login.component.html',
  styleUrls:   ['./login.component.css']
})
export class LoginComponent {
  username = '';
  password = '';
  errorMessage = '';
  allowSignup: boolean;

  constructor(
    private router: Router,
    private auth: AuthService,
    env: EnvService,
  ) {
    this.allowSignup = env.allowSignup;
  }

  login() {
    if (!this.username || !this.password) {
      this.errorMessage = 'Username and password are required.';
      return;
    }
    this.auth.login(this.username, this.password).subscribe({
      next: () => this.router.navigate(['']),    // ← go home
      error: () => this.errorMessage = 'Invalid credentials.'
    });
  }

  clearError() {
    this.errorMessage = '';
  }
}