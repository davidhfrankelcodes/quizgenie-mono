import { ComponentFixture, TestBed } from '@angular/core/testing';

import { BucketList } from './bucket-list';

describe('BucketList', () => {
  let component: BucketList;
  let fixture: ComponentFixture<BucketList>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [BucketList]
    })
    .compileComponents();

    fixture = TestBed.createComponent(BucketList);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
