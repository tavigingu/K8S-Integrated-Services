import { Injectable } from '@angular/core';
import { HttpClient, HttpEventType } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { environment } from '../environment/environment';

export interface FileRecord {
  id: number;
  fileName: string;
  blobUrl: string;
  uploadedAt: string;
  processingResult: string;
  status: string;
}

export interface FileUploadResponse {
  fileId: number;
  fileName: string;
  blobUrl: string;
  message: string;
}

export interface FilesListResponse {
  files: FileRecord[];
  count: number;
}

export interface ProcessResult {
  entities: Entity[];
  success: boolean;
  message?: string;
}

export interface Entity {
  name: string;
  category: string;
  confidence: number;
  offset: number;
  length: number;
  subType?: string;
  matches?: string[];
}

@Injectable({
  providedIn: 'root'
})
export class FileService {
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient) { }

  uploadFile(file: File, onProgress?: (percentage: number) => void): Observable<FileUploadResponse> {
    const formData = new FormData();
    formData.append('file', file);

    return this.http.post<FileUploadResponse>(`${this.apiUrl}/files`, formData, {
      reportProgress: true,
      observe: 'events'
    }).pipe(
      map(event => {
        switch (event.type) {
          case HttpEventType.UploadProgress:
            if (event.total && onProgress) {
              const progress = Math.round(100 * event.loaded / event.total);
              onProgress(progress);
            }
            return {} as FileUploadResponse;
          case HttpEventType.Response:
            return event.body as FileUploadResponse;
          default:
            return {} as FileUploadResponse;
        }
      }),
      catchError(error => {
        console.error('Upload error:', error);
        return throwError(() => new Error('Error uploading file. Please try again.'));
      })
    );
  }

  getFiles(limit: number = 10, offset: number = 0): Observable<FilesListResponse> {
    return this.http.get<FilesListResponse>(`${this.apiUrl}/files?limit=${limit}&offset=${offset}`).pipe(
      catchError(error => {
        console.error('Error fetching files:', error);
        return throwError(() => new Error('Error loading files. Please try again.'));
      })
    );
  }

  getFileById(id: number): Observable<FileRecord> {
    return this.http.get<FileRecord>(`${this.apiUrl}/files/${id}`).pipe(
      catchError(error => {
        console.error('Error fetching file:', error);
        return throwError(() => new Error('Error loading file details. Please try again.'));
      })
    );
  }

  parseProcessingResult(result: string): ProcessResult | null {
    if (!result) return null;
    
    try {
      return JSON.parse(result) as ProcessResult;
    } catch (error) {
      console.error('Error parsing processing result:', error);
      return null;
    }
  }
}