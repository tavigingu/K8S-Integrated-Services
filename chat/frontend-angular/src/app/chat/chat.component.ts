import { Component, OnInit, ViewChild, ElementRef, AfterViewChecked } from '@angular/core';
import { ChatService } from '../services/chat.service';
import { Message } from '../models/message.model';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-chat',
  standalone: true,
  imports: [FormsModule, CommonModule],
  templateUrl: './chat.component.html',
  styleUrls: ['./chat.component.css']
})
export class ChatComponent implements OnInit, AfterViewChecked {
  messages: Message[] = [];
  newMessage: string = '';
  username: string = '';

  @ViewChild('messageContainer') private messageContainer!: ElementRef;

  constructor(private chatService: ChatService) {}

  ngOnInit(): void {
    const storedUsername = localStorage.getItem('username');
    if (storedUsername) {
      this.username = storedUsername;
    } else {
      window.location.href = '/';
      return;
    }
    this.loadMessages();
    setInterval(() => this.loadMessages(), 5000);
  }

  ngAfterViewChecked(): void {
    this.scrollToBottom();
  }

  loadMessages(): void {
    this.chatService.getMessages().subscribe({
      next: (messages) => {
        this.messages = messages;
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
      this.chatService.sendMessage(message).subscribe({
        next: () => {
          this.newMessage = '';
          this.loadMessages();
        },
        error: (err) => {
          console.error('Eroare la trimiterea mesajului:', err);
        }
      });
    }
  }

  private scrollToBottom(): void {
    try {
      this.messageContainer.nativeElement.scrollTop = this.messageContainer.nativeElement.scrollHeight;
    } catch (err) {
      console.error('Eroare la derularea în jos:', err);
    }
  }
}