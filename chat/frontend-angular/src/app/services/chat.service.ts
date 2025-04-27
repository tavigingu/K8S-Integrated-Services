import { Injectable, PLATFORM_ID, Inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, Subject } from 'rxjs';
import { Message } from '../models/message.model';
import { isPlatformBrowser } from '@angular/common';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class ChatService {
  private socket: WebSocket | null = null;
  // private apiUrl = 'http://localhost:8080';
  // private wsUrl = 'ws://localhost:8080/ws';
  private apiUrl = environment.apiUrl; // Folosește variabila de mediu
  private wsUrl = environment.wsUrl;
  private messagesSubject = new Subject<Message>();
  public messages$ = this.messagesSubject.asObservable();
  private isBrowser: boolean;

  constructor(
    private http: HttpClient,
    @Inject(PLATFORM_ID) platformId: Object
  ) {
    this.isBrowser = isPlatformBrowser(platformId);
    
    if (this.isBrowser) {
      this.connectWebSocket();
    }
  }

  private connectWebSocket() {
    if (!this.isBrowser) return;
    
    this.socket = new WebSocket(this.wsUrl);

    this.socket.onopen = () => {
      console.log('WebSocket connected');
    };

    this.socket.onmessage = (event) => {
      const message: Message = JSON.parse(event.data);
      this.messagesSubject.next(message);
    };

    this.socket.onclose = () => {
      console.log('WebSocket closed, attempting to reconnect...');
      setTimeout(() => this.connectWebSocket(), 5000); // Reconnect after 5 seconds
    };

    this.socket.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
  }

  sendMessage(message: Message): void {
    if (!this.isBrowser) return;
    
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify(message));
    } else {
      console.error('WebSocket is not connected. Message not sent.');
    }
  }

  //Păstrăm metoda REST pentru încărcarea inițială sau backup
  getMessages(): Observable<Message[]> {
    return this.http.get<Message[]>(`${this.apiUrl}/messages`);
  }
}