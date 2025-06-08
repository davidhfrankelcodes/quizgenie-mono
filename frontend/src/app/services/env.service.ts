import { Injectable } from '@angular/core';

declare global {
  interface Window { __env: any; }
}

@Injectable({ providedIn: 'root' })
export class EnvService {
  // will be the string "true" or "false"
  readonly allowSignup = window.__env?.ALLOW_SIGNUP === 'true';
}
