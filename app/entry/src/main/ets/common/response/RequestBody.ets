export class RequestBody {

  url: string;

  method: string;

  data: Object;

  constructor() {
    this.url = ''
    this.method = 'GET'
    this.data = new Object
  }
}
