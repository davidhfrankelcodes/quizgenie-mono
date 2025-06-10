// frontend/src/app/app.routes.ts

import { Routes } from '@angular/router';

import { LoginComponent }   from './components/login/login.component';
import { SignupComponent }  from './components/signup/signup.component';
import { HomeComponent }    from './home/home.component';
import { BucketList }       from './components/bucket-list/bucket-list';
import { BucketDetail }     from './components/bucket-detail/bucket-detail';
import { AuthGuard }        from './guards/auth.guard';

export const routes: Routes = [
  // public
  { path: 'login',  component: LoginComponent },
  { path: 'signup', component: SignupComponent },

  // authenticated views
  { path: '',        component: HomeComponent,  canActivate: [AuthGuard] },
  { path: 'buckets', component: BucketList,     canActivate: [AuthGuard] },
  { path: 'buckets/:id', component: BucketDetail, canActivate: [AuthGuard] },

  // fallback
  { path: '**', redirectTo: '' }
];
