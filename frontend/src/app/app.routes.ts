import { Routes } from '@angular/router';
import { LoginComponent } from './components/login/login.component';
import { SignupComponent } from './components/signup/signup.component';
import { HomeComponent } from './home/home.component';
import { AuthGuard } from './guards/auth.guard';

export const routes: Routes = [
  // guarded home
  { path: '',     component: HomeComponent, canActivate: [AuthGuard] },
  // public
  { path: 'login',  component: LoginComponent },
  { path: 'signup', component: SignupComponent },
  // catch-all
  { path: '**',    redirectTo: '' }
];