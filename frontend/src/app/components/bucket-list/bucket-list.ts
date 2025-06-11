import { Component, OnInit } from '@angular/core';
import { CommonModule }      from '@angular/common';
import { Router, RouterLink }from '@angular/router';
import { BucketService, Bucket } from '../../services/bucket.service';
import { FileUploadModal }   from '../file-upload-modal/file-upload-modal';
import { Observable }        from 'rxjs';

@Component({
  selector: 'app-bucket-list',
  standalone: true,
  imports: [ CommonModule, RouterLink, FileUploadModal ],
  templateUrl: './bucket-list.html',
  styleUrls: ['./bucket-list.css']
})
export class BucketList implements OnInit {
  buckets: Bucket[] = [];
  showUpload = false;
  newBucketId!: number;

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

  // instead of navigate directly, create & then open modal
  createBucket() {
    this.bucketSvc.create().subscribe(b => {
      this.newBucketId = b.id;
      this.showUpload = true;
    });
  }

  closeModal() {
    this.showUpload = false;
    this.refresh();
  }
}