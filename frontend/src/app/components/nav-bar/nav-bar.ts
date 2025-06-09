import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { Observable } from 'rxjs';
import { AuthService } from '../../services/auth.service';
import { DrawerService } from '../../services/drawer.service';

@Component({
  selector: 'app-nav-bar',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './nav-bar.html',
  styleUrls: ['./nav-bar.css']
})
export class NavBar {
  isLoggedIn$: Observable<boolean>;
  constructor(
    private auth: AuthService,
    public drawer: DrawerService
  ) {
    this.isLoggedIn$ = this.auth.isLoggedIn();
  }

  logout() {
    this.auth.logout();
    this.drawer.close();
  }
}