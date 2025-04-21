import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SetUsernameComponent } from './set-username.component';

describe('SetUsernameComponent', () => {
  let component: SetUsernameComponent;
  let fixture: ComponentFixture<SetUsernameComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SetUsernameComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SetUsernameComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
