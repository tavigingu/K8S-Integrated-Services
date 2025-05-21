import { bootstrapApplication } from '@angular/platform-browser';
import { AppComponent } from './app/app.component';
import { provideRouter } from '@angular/router';
import { routes } from './app/app.routes';
import { provideHttpClient, withFetch } from '@angular/common/http';
import { importProvidersFrom, APP_INITIALIZER } from '@angular/core';
import { FormsModule } from '@angular/forms';
//import { ConfigService } from './app/services/config.service';
import { firstValueFrom } from 'rxjs';

export function initializeApp() {

}

// bootstrapApplication(AppComponent, {
//   providers: [
//     provideRouter(routes),
//     provideHttpClient(withFetch()),
//     importProvidersFrom(FormsModule),
//     {
//       provide: APP_INITIALIZER,
//       useFactory: initializeApp,
//       deps: [],
//       multi: true
//     }
//   ]
// }).catch(err => console.error(err));
bootstrapApplication(AppComponent, {
  providers: [
    provideRouter(routes),
    provideHttpClient(withFetch()),
    importProvidersFrom(FormsModule)
  ]
}).catch(err => console.error(err));