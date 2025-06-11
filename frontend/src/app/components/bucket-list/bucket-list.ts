import { Component, OnInit } from '@angular/core';
import { CommonModule }      from '@angular/common';
import { RouterLink, Router } from '@angular/router';

import { BucketService, Bucket } from '../../services/bucket.service';

@Component({
  selector: 'app-bucket-list',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './bucket-list.html',
  styleUrls: ['./bucket-list.css']
})
export class BucketList implements OnInit {
  buckets: Bucket[] = [];

  constructor(
    private bucketSvc: BucketService,
    private router:    Router,
  ) {}

  ngOnInit() {
    this.refresh();
  }

  refresh() {
    this.bucketSvc.list().subscribe(bs => this.buckets = bs);
  }

  createBucket() {
    this.bucketSvc.create().subscribe(b => {
      this.router.navigate(['/buckets', b.id]);
    });
  }
}