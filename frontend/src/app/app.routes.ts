// src/app/app.routes.ts
import { Routes } from '@angular/router';

import { LoginComponent }   from './components/login/login.component';
import { SignupComponent }  from './components/signup/signup.component';
import { BucketList }       from './components/bucket-list/bucket-list';
import { BucketDetail }     from './components/bucket-detail/bucket-detail';
import { AuthGuard }        from './guards/auth.guard';

export const routes: Routes = [
  // public
  { path: 'login',  component: LoginComponent },
  { path: 'signup', component: SignupComponent },

  // authenticated views
  { path: '',                component: BucketList,   canActivate: [AuthGuard] },
  { path: 'buckets/:id',     component: BucketDetail, canActivate: [AuthGuard] },

  // fallback
  { path: '**', redirectTo: '' }
];
