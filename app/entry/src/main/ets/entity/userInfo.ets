export class UserInfo {

  username: string;

  password: string;

  // 添加注册所需的字段
  email?: string;

  nickname?: string;

  constructor() {
    this.username = ''
    this.password = ''
    this.email = ''
    this.nickname = ''
  }
}
