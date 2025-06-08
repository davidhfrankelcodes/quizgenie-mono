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
  username        = '';
  email           = '';
  password        = '';
  confirmPassword = '';
  error           = '';

  constructor(
    private auth: AuthService,
    private router: Router,
  ) {}

  signup() {
    // all fields required
    if (!this.username || !this.email || !this.password || !this.confirmPassword) {
      this.error = 'All fields are required.';
      return;
    }

    // passwords must match
    if (this.password !== this.confirmPassword) {
      this.error = 'Passwords do not match.';
      return;
    }

    // proceed
    this.auth.signup(this.username, this.password, this.email).subscribe({
      next: () => this.router.navigate(['']),
      error: err => this.error = err.status === 409
        ? 'Username already taken'
        : 'Signup failed'
    });
  }
}