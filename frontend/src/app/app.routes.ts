// src/app/app.routes.ts
import { Routes } from '@angular/router';
import { LoginComponent } from './components/login/login.component';

export const routes: Routes = [
  // Route for login page
  { path: 'login', component: LoginComponent },

  // (You can add additional routes here, e.g. a home/dashboard component once the user is logged in.)
  // Example:
  // { path: '', component: HomeComponent, canActivate: [AuthGuard] },

  // Redirect “empty” to login (or wherever you want)
  { path: '', redirectTo: 'login', pathMatch: 'full' },

  // Wildcard catch-all (optional)
  { path: '**', redirectTo: 'login' }
];