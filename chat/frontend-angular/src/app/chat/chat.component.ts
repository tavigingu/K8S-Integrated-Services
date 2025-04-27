import { Component, OnInit, ViewChild, ElementRef, AfterViewChecked, OnDestroy, PLATFORM_ID, Inject } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { ChatService } from '../services/chat.service';
import { Message } from '../models/message.model';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-chat',
  standalone: true,
  imports: [FormsModule, CommonModule],
  templateUrl: './chat.component.html',
  styleUrls: ['./chat.component.css']
})
export class ChatComponent implements OnInit, AfterViewChecked, OnDestroy {
  messages: Message[] = [];
  newMessage: string = '';
  username: string = '';
  private subscription: Subscription;
  private isBrowser: boolean;

  @ViewChild('messageContainer') private messageContainer!: ElementRef;

  constructor(
    private chatService: ChatService,
    @Inject(PLATFORM_ID) platformId: Object
  ) {
    this.subscription = new Subscription();
    this.isBrowser = isPlatformBrowser(platformId);
  }

  ngOnInit(): void {
    if (this.isBrowser) {
      const storedUsername = localStorage.getItem('username');
      if (storedUsername) {
        this.username = storedUsername;
      } else {
        window.location.href = '/';
        return;
      }

      // Încărcăm mesajele inițiale folosind REST
      this.loadInitialMessages();

      // Apoi ascultăm pentru mesaje noi prin WebSocket
      this.subscription.add(
        this.chatService.messages$.subscribe(message => {
          // Verificăm dacă mesajul există deja pentru a evita duplicatele
          const exists = this.messages.some(m => 
            m.id === message.id || 
            (m.content === message.content && 
             m.username === message.username && 
             m.timestamp === message.timestamp)
          );
          
          if (!exists) {
            this.messages.push(message);
            setTimeout(() => this.scrollToBottom(), 100);
          }
        })
      );
    }
  }

  ngAfterViewChecked(): void {
    if (this.isBrowser) {
      this.scrollToBottom();
    }
  }

  ngOnDestroy(): void {
    if (this.subscription) {
      this.subscription.unsubscribe();
    }
  }

  loadInitialMessages(): void {
    this.chatService.getMessages().subscribe({
      next: (messages) => {
        this.messages = messages;
        if (this.isBrowser) {
          setTimeout(() => this.scrollToBottom(), 100);
        }
      },
      error: (err) => {
        console.error('Eroare la încărcarea mesajelor:', err);
      }
    });
  }

  sendMessage(): void {
    if (this.newMessage.trim()) {
      const message: Message = {
        content: this.newMessage,
        username: this.username
      };
      
      this.chatService.sendMessage(message);
      this.newMessage = '';
    }
  }

  private scrollToBottom(): void {
    if (!this.isBrowser) return;
    
    try {
      this.messageContainer!.nativeElement.scrollTop = this.messageContainer!.nativeElement.scrollHeight;
    } catch (err) {
      console.error('Eroare la derularea în jos:', err);
    }
  }
}