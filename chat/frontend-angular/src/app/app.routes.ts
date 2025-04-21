import { Routes } from '@angular/router';
import { SetUsernameComponent } from './set-username/set-username.component';
import { ChatComponent } from './chat/chat.component';

export const routes: Routes = [
  { path: '', component: SetUsernameComponent },
  { path: 'chat', component: ChatComponent },
];