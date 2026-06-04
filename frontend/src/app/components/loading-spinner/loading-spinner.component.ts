import {Component, Input, OnInit} from '@angular/core';

@Component({
  standalone: true,
  selector: 'app-loading-spinner',
  templateUrl: './loading-spinner.component.html',
  styleUrls: ['./loading-spinner.component.scss']
})
export class LoadingSpinnerComponent implements OnInit {
  @Input() loadingTitle = "Please wait, loading..."
  @Input() loadingSubTitle = ""

  constructor() { }

  ngOnInit(): void {
  }

}
