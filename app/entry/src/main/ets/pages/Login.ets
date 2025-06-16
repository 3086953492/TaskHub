@Entry
@Component
struct Login {
  @State username: string = '';
  @State password: string = '';
  @State isLoading: boolean = false;
  @State errorMessage: string = '';

  // 验证表单
  validateForm(): boolean {
    if (!this.username.trim()) {
      this.errorMessage = '请输入用户名';
      return false;
    }
    if (this.username.trim().length < 3 || this.username.trim().length > 15) {
      this.errorMessage = '用户名长度应为3-15字符';
      return false;
    }
    if (!this.password) {
      this.errorMessage = '请输入密码';
      return false;
    }
    if (this.password.length < 6) {
      this.errorMessage = '密码长度至少6位';
      return false;
    }
    this.errorMessage = '';
    return true;
  }

  // 登录处理
  async handleLogin() {
    if (!this.validateForm()) {
      return;
    }

    this.isLoading = true;
    this.errorMessage = '';

    try {
      // 模拟API调用
      setTimeout(() => {
        if (this.username === 'admin' && this.password === '123456') {
          // 登录成功
          console.log('登录成功');
          this.errorMessage = '';
          // TODO: 保存token，跳转到主页
        } else {
          this.errorMessage = '用户名或密码错误';
        }
        this.isLoading = false;
      }, 1000);
    } catch (error) {
      this.errorMessage = '登录失败，请重试';
      this.isLoading = false;
    }
  }

  build() {
    Column() {
      // 顶部装饰区域
      Column() {
        Image($r('app.media.app_icon'))
          .width(80)
          .height(80)
          .margin({ bottom: 16 })

        Text('TaskHub')
          .fontSize(32)
          .fontWeight(FontWeight.Bold)
          .fontColor('#182431')
          .margin({ bottom: 8 })

        Text('任务管理系统')
          .fontSize(16)
          .fontColor('#99182431')
          .margin({ bottom: 40 })
      }
      .width('100%')
      .padding({ top: 60 })

      // 登录表单
      Column({ space: 20 }) {
        // 用户名输入框
        Column({ space: 8 }) {
          Row() {
            Text('用户名')
              .fontSize(14)
              .fontColor('#182431')
              .fontWeight(FontWeight.Medium)
          }
          .width('100%')

          TextInput({ placeholder: '请输入用户名', text: this.username })
            .width('100%')
            .height(48)
            .backgroundColor('#F5F5F5')
            .borderRadius(8)
            .padding({ left: 16, right: 16 })
            .fontSize(16)
            .placeholderColor('#99182431')
            .onChange((value: string) => {
              this.username = value;
              this.errorMessage = '';
            })
        }

        // 密码输入框
        Column({ space: 8 }) {
          Row() {
            Text('密码')
              .fontSize(14)
              .fontColor('#182431')
              .fontWeight(FontWeight.Medium)
          }
          .width('100%')

          TextInput({ placeholder: '请输入密码', text: this.password })
            .width('100%')
            .height(48)
            .backgroundColor('#F5F5F5')
            .borderRadius(8)
            .padding({ left: 16, right: 16 })
            .fontSize(16)
            .type(InputType.Password)
            .placeholderColor('#99182431')
            .onChange((value: string) => {
              this.password = value;
              this.errorMessage = '';
            })
        }

        // 错误信息
        if (this.errorMessage) {
          Text(this.errorMessage)
            .fontSize(14)
            .fontColor('#F44336')
            .textAlign(TextAlign.Center)
            .width('100%')
        }

        // 登录按钮
        Button(this.isLoading ? '登录中...' : '登录')
          .width('100%')
          .height(48)
          .backgroundColor(this.isLoading ? '#CCCCCC' : '#007DFF')
          .borderRadius(8)
          .fontSize(16)
          .fontColor(Color.White)
          .fontWeight(FontWeight.Medium)
          .enabled(!this.isLoading)
          .onClick(() => {
            this.handleLogin();
          })

        // 注册链接
        Row() {
          Text('还没有账号？')
            .fontSize(14)
            .fontColor('#99182431')

          Text('立即注册')
            .fontSize(14)
            .fontColor('#007DFF')
            .fontWeight(FontWeight.Medium)
            .onClick(() => {
              // TODO: 跳转到注册页面
              console.log('跳转到注册页面');
            })
        }
        .justifyContent(FlexAlign.Center)
        .margin({ top: 20 })
      }
      .width('100%')
      .padding({ left: 32, right: 32 })

      Blank()

      // 底部信息
      Column() {
        Text('© 2025 TaskHub. All rights reserved.')
          .fontSize(12)
          .fontColor('#66182431')
          .textAlign(TextAlign.Center)
      }
      .width('100%')
      .padding({ bottom: 30 })
    }
    .width('100%')
    .height('100%')
    .backgroundColor('#FFFFFF')
    .justifyContent(FlexAlign.SpaceBetween)
  }
} 