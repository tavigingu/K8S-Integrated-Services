import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { FileService, FileUploadResponse } from '../services/file.service';

@Component({
  selector: 'app-file-upload',
  // Remove the 'standalone: true' property
  templateUrl: './file-upload.component.html',
  styleUrls: ['./file-upload.component.css'],
  standalone: false
})
export class FileUploadComponent {
  selectedFile: File | null = null;
  isDragging = false;
  uploading = false;
  uploadProgress = 0;
  uploadSuccess = false;
  errorMessage = '';
  uploadResponse: FileUploadResponse | null = null;
  
  constructor(private fileService: FileService, private router: Router) {}
  
  onDragOver(event: DragEvent): void {
    event.preventDefault();
    event.stopPropagation();
    this.isDragging = true;
  }
  
  onDragLeave(event: DragEvent): void {
    event.preventDefault();
    event.stopPropagation();
    this.isDragging = false;
  }
  
  onDrop(event: DragEvent): void {
    event.preventDefault();
    event.stopPropagation();
    this.isDragging = false;
    
    if (event.dataTransfer?.files.length) {
      const file = event.dataTransfer.files[0];
      this.validateAndSetFile(file);
    }
  }
  
  onFileSelected(event: Event): void {
    const target = event.target as HTMLInputElement;
    
    if (target.files && target.files.length > 0) {
      const file = target.files[0];
      this.validateAndSetFile(file);
    }
  }
  
  validateAndSetFile(file: File): void {
    // Reset previous errors
    this.errorMessage = '';
    
    // Check file type
    const validTypes = ['text/plain', 'application/json'];
    if (!validTypes.includes(file.type)) {
      this.errorMessage = 'Invalid file type. Please upload text or JSON files only.';
      return;
    }
    
    // Check file size (max 10MB)
    const maxSize = 10 * 1024 * 1024; // 10MB in bytes
    if (file.size > maxSize) {
      this.errorMessage = 'File size exceeds the maximum limit of 10MB.';
      return;
    }
    
    this.selectedFile = file;
  }
  
  resetFile(): void {
    this.selectedFile = null;
    this.errorMessage = '';
  }
  
  resetForm(): void {
    this.selectedFile = null;
    this.uploading = false;
    this.uploadProgress = 0;
    this.uploadSuccess = false;
    this.errorMessage = '';
    this.uploadResponse = null;
  }
  
  uploadFile(): void {
    if (!this.selectedFile) return;
    
    this.uploading = true;
    this.errorMessage = '';
    
    this.fileService.uploadFile(this.selectedFile, (progress) => {
      this.uploadProgress = progress;
    }).subscribe({
      next: (response) => {
        // Skip empty responses (progress updates)
        if (!response.fileId) return;
        
        this.uploadResponse = response;
        this.uploading = false;
        this.uploadSuccess = true;
      },
      error: (error) => {
        this.uploading = false;
        this.errorMessage = error.message || 'An error occurred during upload. Please try again.';
      }
    });
  }
  
  viewFile(): void {
    if (this.uploadResponse?.fileId) {
      this.router.navigate(['/file', this.uploadResponse.fileId]);
    }
  }
  
  formatFileSize(bytes: number): string {
    if (bytes === 0) return '0 Bytes';
    
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  }
}