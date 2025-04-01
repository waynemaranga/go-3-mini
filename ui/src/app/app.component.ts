import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ChatService, ChatMessage } from './services/chat.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
  standalone: true,
  imports: [CommonModule, FormsModule]
})


export class AppComponent implements OnInit {
  title(title: any) {
    throw new Error('Method not implemented.');
  }
  messages: ChatMessage[] = [];
  userMessage = '';
  loading = false;
  
  @ViewChild('chatHistory') private chatHistoryRef!: ElementRef;

  constructor(private chatService: ChatService) {}

  ngOnInit(): void {
    this.loadChatHistory();
  }

  loadChatHistory(): void {
    this.loading = true;
    this.chatService.getHistory().subscribe({
      next: (history) => {
        this.messages = history;
        this.loading = false;
        this.scrollToBottom();
      },
      error: (error) => {
        console.error('Failed to load chat history:', error);
        this.loading = false;
      }
    });
  }

  sendMessage(): void {
    if (!this.userMessage.trim() || this.loading) {
      return;
    }

    const userMessageText = this.userMessage.trim();
    this.userMessage = '';
    
    // Add user message to the UI immediately
    this.messages.push({
      role: 'user',
      content: userMessageText
    });
    
    this.scrollToBottom();
    this.loading = true;

    // Send to the server
    this.chatService.sendMessage(userMessageText).subscribe({
      next: (response) => {
        this.messages.push({
          role: 'assistant',
          content: response.response
        });
        this.loading = false;
        this.scrollToBottom();
      },
      error: (error) => {
        console.error('Error sending message:', error);
        this.messages.push({
          role: 'assistant',
          content: 'Sorry, there was an error processing your request.'
        });
        this.loading = false;
        this.scrollToBottom();
      }
    });
  }

  private scrollToBottom(): void {
    setTimeout(() => {
      if (this.chatHistoryRef) {
        const element = this.chatHistoryRef.nativeElement;
        element.scrollTop = element.scrollHeight;
      }
    }, 100);
  }
}