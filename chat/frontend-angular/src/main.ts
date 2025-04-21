import { bootstrapApplication } from '@angular/platform-browser';
import { AppComponent } from './app/app.component';
import { provideRouter } from '@angular/router';
import { routes } from './app/app.routes';
import { provideHttpClient } from '@angular/common/http'; // Pentru HttpClient
import { importProvidersFrom } from '@angular/core';
import { FormsModule } from '@angular/forms'; // Pentru FormsModule

bootstrapApplication(AppComponent, {
  providers: [
    provideRouter(routes),
    provideHttpClient(), // Adaugă suport pentru HttpClient
    importProvidersFrom(FormsModule), // Adaugă suport pentru FormsModule
  ]
}).catch(err => console.error(err));