export interface Message {
    id?: string; // id-ul va fi generat de backend
    content: string;
    username: string;
    timestamp?: string; // timestamp-ul va fi setat de backend
  }