import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Bucket {
  id:   number;
  name: string;
}

@Injectable({ providedIn: 'root' })
export class BucketService {
  constructor(private http: HttpClient) {}

  list(): Observable<Bucket[]> {
    return this.http.get<Bucket[]>('/buckets');
  }

  create(): Observable<Bucket> {
    return this.http.post<Bucket>('/buckets', {});
  }
}