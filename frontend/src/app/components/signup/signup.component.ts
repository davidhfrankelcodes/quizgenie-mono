import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';

import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-signup',
  standalone: true,
  imports: [ CommonModule, FormsModule ],
  templateUrl: './signup.component.html',
  styleUrls:   ['./signup.component.css']
})
export class SignupComponent {
  username = '';
  email    = '';
  password = '';
  error    = '';

  constructor(
    private auth: AuthService,
    private router: Router,
  ) {}

  signup() {
    if (!this.username || !this.email || !this.password) {
      this.error = 'All fields are required.';
      return;
    }
    this.auth.signup(this.username, this.password, this.email).subscribe({
      next: () => this.router.navigate(['']),  // ← go home
      error: err => this.error = err.status === 409
        ? 'Username already taken'
        : 'Signup failed'
    });
  }
}