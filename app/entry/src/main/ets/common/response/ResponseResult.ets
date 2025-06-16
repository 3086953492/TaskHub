export class ResponseResult {

  code: number;

  msg: string | Resource;

  data: ArrayBuffer | Object | string

  constructor() {
    this.code = 0;
    this.msg = '';
    this.data = '';
  }
}
