// frontend/src/app/app.component.ts

import { Component } from '@angular/core';
import { CommonModule, NgIf } from '@angular/common';
import { RouterOutlet } from '@angular/router';

import { NavBar }       from './components/nav-bar/nav-bar';
import { BucketList }   from './components/bucket-list/bucket-list';
import { HomeComponent }from './home/home.component';

import { AuthService }  from './services/auth.service';
import { DrawerService }from './services/drawer.service';
import { Observable }   from 'rxjs';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    CommonModule,
    NgIf,
    RouterOutlet,
    NavBar,
    BucketList,    // for the drawer
    HomeComponent, // for the default route
  ],
  templateUrl: './app.html',
  styleUrls: ['./app.css']
})
export class AppComponent {
  isLoggedIn$: Observable<boolean>;
  drawerOpen$: Observable<boolean>;
  currentYear = new Date().getFullYear();

  constructor(
    private auth: AuthService,
    public  drawer: DrawerService
  ) {
    this.isLoggedIn$ = this.auth.isLoggedIn();
    this.drawerOpen$ = this.drawer.open$;
  }
}
