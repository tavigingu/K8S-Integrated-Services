// import { Injectable } from '@angular/core';
// import { HttpClient } from '@angular/common/http';
// import { Observable } from 'rxjs';

// export interface EnvironmentConfig {
//   production: boolean;
//   apiUrl: string;
//   wsUrl: string;
// }

// @Injectable({
//   providedIn: 'root'
// })
// export class ConfigService {
//   private config: EnvironmentConfig | null = null;

//   constructor(private http: HttpClient) {}

//   loadConfig(): Observable<EnvironmentConfig> {
//     return this.http.get<EnvironmentConfig>('/config.json');
//   }

//   setConfig(config: EnvironmentConfig): void {
//     this.config = config;
//   }

//   getConfig(): EnvironmentConfig {
//     if (!this.config) {
//       throw new Error('Configuration not loaded');
//     }
//     return this.config;
//   }
// }