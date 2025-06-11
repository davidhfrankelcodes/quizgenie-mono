import { Component, EventEmitter, Input, Output } from '@angular/core';
import { CommonModule }                      from '@angular/common';
import { HttpClient, HttpEventType }         from '@angular/common/http';
import { FormsModule }                       from '@angular/forms';

@Component({
  selector: 'app-file-upload-modal',
  standalone: true,
  imports: [ CommonModule, FormsModule ],
  templateUrl: './file-upload-modal.html',
  styleUrls:   ['./file-upload-modal.css']
})
export class FileUploadModal {
  @Input() bucketId!: number;
  @Output() closed = new EventEmitter<void>();

  files: File[] = [];
  progress: Record<string, number> = {};
  uploading = false;
  error = '';

  // drag & drop handlers
  onDrop(event: DragEvent) {
    event.preventDefault();
    if (event.dataTransfer?.files) {
      this.addFiles(Array.from(event.dataTransfer.files));
    }
  }
  onDragOver(event: DragEvent) { event.preventDefault(); }

  // file dialog
  onFileSelect(ev: Event) {
    const input = ev.target as HTMLInputElement;
    if (input.files) {
      this.addFiles(Array.from(input.files));
    }
  }

  private addFiles(newFiles: File[]) {
    newFiles.forEach(f => {
      if (!this.files.find(x => x.name === f.name)) {
        this.files.push(f);
        this.progress[f.name] = 0;
      }
    });
  }

  // kick off uploads
  uploadAll() {
    this.uploading = true;
    this.error = '';
    this.files.forEach(file => {
      const form = new FormData();
      form.append('file', file);
      // POST to your backend endpoint
      this.http.post<any>(
        `/buckets/${this.bucketId}/files`,
        form,
        { reportProgress: true, observe: 'events' }
      ).subscribe({
        next: ev => {
          if (ev.type === HttpEventType.UploadProgress && ev.total) {
            this.progress[file.name] = Math.round(100 * ev.loaded / ev.total);
          }
        },
        error: err => {
          this.error = 'Upload failed.';
          this.progress[file.name] = 0;
        }
      });
    });
  }

  constructor(private http: HttpClient) {}
}