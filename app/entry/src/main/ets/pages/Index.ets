// 这里需要引入第 6 步创建的用户API
import { userLogin } from '../common/api/user';

@Entry
@Component
struct Index {

  // 在进入页面时调用 userLogin 接口进行测试
  aboutToAppear(): void {
    userLogin({username: 'user', password: 'useruser'}).then(res => {
      console.info('页面成功获取到用户信息：' + res.data)
    })
  }

  build() {
    RelativeContainer() {

    }
  }
}
