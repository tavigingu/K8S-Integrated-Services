import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
//import { ConfigService } from './services/config.service';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet],
  standalone: true,
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent {
  title = 'frontend-angular';

  // constructor(private configService: ConfigService) {
  //   this.config = this.configService.getConfig();
  // }
}