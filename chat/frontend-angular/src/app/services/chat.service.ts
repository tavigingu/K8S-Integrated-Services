import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, Subject } from 'rxjs';
import { Message } from '../models/message.model';

@Injectable({
  providedIn: 'root'
})
export class ChatService {
  private socket: WebSocket = new WebSocket('');
  private apiUrl = 'http://localhost:8080';
  private wsUrl = 'ws://localhost:8080/ws';
  private messagesSubject = new Subject<Message>();
  public messages$ = this.messagesSubject.asObservable();

  constructor(private http: HttpClient) {
    this.connectWebSocket();
  }

  private connectWebSocket() {
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
    if (this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify(message));
    } else {
      console.error('WebSocket is not connected. Message not sent.');
    }
  }

  // Păstrăm metoda REST pentru încărcarea inițială sau backup
  getMessages(): Observable<Message[]> {
    return this.http.get<Message[]>(`${this.apiUrl}/messages`);
  }
}