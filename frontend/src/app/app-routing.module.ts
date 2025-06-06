import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LoginComponent } from './components/login/login.component';
// Placeholder for other components

import { AuthGuard } from './guards/auth.guard';

const routes: Routes = [
  { path: 'login', component: LoginComponent },
  // All other routes will be protected by AuthGuard
  { path: '', redirectTo: '/login', pathMatch: 'full' },
  // { path: 'buckets', component: BucketListComponent, canActivate: [AuthGuard] },
  // { path: 'buckets/:id', component: BucketDetailComponent, canActivate: [AuthGuard] },
  // { path: 'quizzes/:quizId/status', component: QuizStatusComponent, canActivate: [AuthGuard] },
  // { path: 'quizzes/:quizId/take', component: QuizTakingComponent, canActivate: [AuthGuard] },
  // { path: 'attempts/:attemptId', component: QuizReportComponent, canActivate: [AuthGuard] },
  // { path: 'buckets/:id/history', component: ReportHistoryComponent, canActivate: [AuthGuard] },
  { path: '**', redirectTo: '/login' }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }