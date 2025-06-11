import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class DrawerService {
  private openSubject = new BehaviorSubject<boolean>(false);
  open$: Observable<boolean> = this.openSubject.asObservable();

  toggle() {
    this.openSubject.next(!this.openSubject.value);
  }

  close() {
    this.openSubject.next(false);
  }
}