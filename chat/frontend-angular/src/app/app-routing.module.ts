import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { SetUsernameComponent } from './set-username/set-username.component';
import { ChatComponent } from './chat/chat.component';

const routes: Routes = [
  { path: '', component: SetUsernameComponent }, // Prima pagină: setare utilizator
  { path: 'chat', component: ChatComponent },    // Pagina de chat
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }