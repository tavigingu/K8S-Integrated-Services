import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { FileService, FileRecord } from '../services/file.service';

@Component({
  selector: 'app-file-list',
  templateUrl: './file-list.component.html',
  styleUrls: ['./file-list.component.css'],
  standalone: false
})
export class FileListComponent implements OnInit {
  files: FileRecord[] = [];
  loading = false;
  errorMessage = '';
  currentOffset = 0;
  pageSize = 10;
  totalFiles = 0;
  
  constructor(private fileService: FileService, private router: Router) {}
  
  ngOnInit(): void {
    this.loadFiles();
  }
  
  loadFiles(): void {
    this.loading = true;
    this.errorMessage = '';
    
    this.fileService.getFiles(this.pageSize, this.currentOffset).subscribe({
      next: (response) => {
        this.files = response.files;
        this.totalFiles = response.count;
        this.loading = false;
      },
      error: (error) => {
        this.loading = false;
        this.errorMessage = error.message || 'Failed to load files. Please try again.';
      }
    });
  }
  
  navigateToUpload(): void {
    this.router.navigate(['/upload']);
  }
  
  viewFile(id: number): void {
    this.router.navigate(['/file', id]);
  }
  
  previousPage(): void {
    if (this.currentOffset > 0) {
      this.currentOffset = Math.max(0, this.currentOffset - this.pageSize);
      this.loadFiles();
    }
  }
  
  nextPage(): void {
    if (this.currentOffset + this.files.length < this.totalFiles) {
      this.currentOffset += this.pageSize;
      this.loadFiles();
    }
  }
  
  getStatusClass(status: string): string {
    switch (status.toLowerCase()) {
      case 'uploaded':
        return 'status-uploaded';
      case 'processed':
        return 'status-processed';
      case 'error':
        return 'status-error';
      default:
        return 'status-uploaded';
    }
  }
}