import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { RouterModule, Routes } from '@angular/router';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { CommonModule } from '@angular/common';

import { AppComponent } from './app.component';
import { NavigationComponent } from './components/navigation.component';
import { HomeComponent } from './components/home.component';
import { FileUploadComponent } from './components/file-upload.component';
import { FileListComponent } from './components/file-list.component';
import { FileDetailComponent } from './components/file-detail.component';
import { LoadingSpinnerComponent } from './components/loading-spinner.component';

// Define routes
const routes: Routes = [
  { path: '', component: HomeComponent },
  { path: 'upload', component: FileUploadComponent },
  { path: 'history', component: FileListComponent },
  { path: 'file/:id', component: FileDetailComponent },
  { path: '**', redirectTo: '' }
];

@NgModule({
  declarations: [
    AppComponent,
    NavigationComponent,
    HomeComponent,
    FileUploadComponent,
    FileListComponent,
    FileDetailComponent,
    LoadingSpinnerComponent
  ],
  imports: [
    BrowserModule,
    CommonModule,
    HttpClientModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule.forRoot(routes)
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }