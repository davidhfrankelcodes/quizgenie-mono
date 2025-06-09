// src/app/components/nav-bar/nav-bar.ts
import { Component }    from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink }   from '@angular/router';
import { Observable }   from 'rxjs';
import { AuthService }  from '../../services/auth.service';

@Component({
  selector: 'app-nav-bar',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './nav-bar.html',
  styleUrls: ['./nav-bar.css']
})
export class NavBar {
  // will hold login state
  isLoggedIn$!: Observable<boolean>;

  constructor(private auth: AuthService) {
    // initialize after auth is available
    this.isLoggedIn$ = this.auth.isLoggedIn();
  }

  logout() {
    this.auth.logout();
  }
}
