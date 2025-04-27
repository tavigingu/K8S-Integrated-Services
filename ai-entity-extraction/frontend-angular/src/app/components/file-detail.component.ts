import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { FileService, FileRecord, ProcessResult, Entity } from '../services/file.service';

@Component({
  selector: 'app-file-detail',
  // Fix the path to reference the correct folder structure
  templateUrl: './file-detail.component.html',
  styleUrls: ['./file-detail.component.css'],
  standalone: false
})
export class FileDetailComponent implements OnInit {
  fileId: number | null = null;
  file: FileRecord | null = null;
  processResult: ProcessResult | null = null;
  loading = false;
  errorMessage = '';
  
  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private fileService: FileService
  ) {}
  
  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      const idParam = params.get('id');
      if (idParam) {
        this.fileId = +idParam;
        this.loadFileDetails();
      } else {
        this.errorMessage = 'Invalid file ID.';
      }
    });
  }
  
  loadFileDetails(): void {
    if (!this.fileId) return;
    
    this.loading = true;
    this.errorMessage = '';
    
    this.fileService.getFileById(this.fileId).subscribe({
      next: (file) => {
        this.file = file;
        
        if (file.processingResult) {
          this.processResult = this.fileService.parseProcessingResult(file.processingResult);
        }
        
        this.loading = false;
      },
      error: (error) => {
        this.loading = false;
        this.errorMessage = error.message || 'Failed to load file details.';
      }
    });
  }
  
  refreshFileDetails(): void {
    this.loadFileDetails();
  }
  
  navigateToHistory(): void {
    this.router.navigate(['/history']);
  }
  
  navigateToUpload(): void {
    this.router.navigate(['/upload']);
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
  
  getCategoryColor(category: string, opacity: number = 1): string {
    // Define the colors with an index signature
    const colors: {[key: string]: string} = {
      'person': '#ef4444',
      'location': '#3b82f6',
      'organization': '#8b5cf6',
      'datetime': '#f59e0b',
      'quantity': '#10b981',
      'email': '#6366f1',
      'phone': '#ec4899',
      'url': '#0ea5e9',
      'address': '#14b8a6'
    };
    
    // Convert category to lowercase to ensure matching
    const lowerCategory = category.toLowerCase();
    
    // Default color if category not found
    const baseColor = colors[lowerCategory] || '#6b7280';
    
    // Return with opacity if needed
    if (opacity < 1) {
      return this.hexToRgba(baseColor, opacity);
    }
    
    return baseColor;
  }
  
  hexToRgba(hex: string, opacity: number): string {
    const r = parseInt(hex.slice(1, 3), 16);
    const g = parseInt(hex.slice(3, 5), 16);
    const b = parseInt(hex.slice(5, 7), 16);
    
    return `rgba(${r}, ${g}, ${b}, ${opacity})`;
  }
  
  getConfidenceColor(confidence: number): string {
    if (confidence >= 0.8) {
      return '#10b981'; // Green for high confidence
    } else if (confidence >= 0.5) {
      return '#f59e0b'; // Yellow for medium confidence
    } else {
      return '#ef4444'; // Red for low confidence
    }
  }
}