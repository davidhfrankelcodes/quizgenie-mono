// src/app/components/login/login.component.ts
import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [
    CommonModule,   // for *ngIf
    FormsModule     // for [(ngModel)]
  ],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent {
  username = '';
  password = '';
  errorMessage = '';

  constructor(private router: Router) {}

  login() {
    // ← Your actual login logic (e.g. call AuthService) goes here.
    // For now you can just navigate to “/” once the user logs in successfully:
    if (this.username && this.password) {
      // e.g. AuthService.login(...).subscribe(...),
      // on success:
      this.router.navigate(['/']);
    } else {
      this.errorMessage = 'Username and password are required.';
    }
  }

  clearError() {
    this.errorMessage = '';
  }
}