import { ComponentFixture, TestBed } from '@angular/core/testing';

import { BucketDetail } from './bucket-detail';

describe('BucketDetail', () => {
  let component: BucketDetail;
  let fixture: ComponentFixture<BucketDetail>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [BucketDetail]
    })
    .compileComponents();

    fixture = TestBed.createComponent(BucketDetail);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
