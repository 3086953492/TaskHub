/**
 * HTTP API
 */
export class Domain {

  // 后端服务IP
  static readonly SERVER: string = 'http://192.168.101.1:8080'

}

/**
 * HTTP状态类
 */
export class HttpStatus {

  // 成功
  static readonly SUCCESS: number = 200

  // 失败
  static readonly ERROR: number = 500

}

export const enum ContentType {

  // 响应体类型
  JSON = 'application/json'

}
