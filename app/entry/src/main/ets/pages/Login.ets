// 这里需要引入第 6 步创建的用户API
import { userLogin } from '../common/api/user';
import { UserInfo } from '../entity/userInfo';
import { StorageUtil } from '../common/utils/StorageUtil';
import { LoginResponse } from '../common/response/LoginResponseData';
import promptAction from '@ohos.promptAction';
import router from '@ohos.router';
import common from '@ohos.app.ability.common';

@Entry
@Component
struct Login {
  @State username: string = '';
  @State password: string = '';
  @State isLoading: boolean = false;
  @State showPassword: boolean = false;
  @State cardScale: number = 1;
  private context: common.UIAbilityContext = getContext(this) as common.UIAbilityContext;

  // 登录处理函数
  async handleLogin() {
    // 验证输入
    if (!this.username.trim()) {
      promptAction.showToast({
        message: '请输入用户名',
        duration: 2000
      });
      return;
    }

    if (!this.password.trim()) {
      promptAction.showToast({
        message: '请输入密码',
        duration: 2000
      });
      return;
    }

    this.isLoading = true;

    try {
      const userInfo = new UserInfo();
      userInfo.username = this.username;
      userInfo.password = this.password;

      const response = await userLogin(userInfo);
      
      console.info('登录响应数据：', JSON.stringify(response));
      console.info('响应code值：', response.code);
      
      if (response.code === 200) {
        // 类型断言确保data是LoginResponseData类型
        const loginResponse = response as LoginResponse;
        // 存储token
        const token = loginResponse.data?.token;
        if (token) {
          const storeSuccess = await StorageUtil.setToken(token);
          if (storeSuccess) {
            console.info('Token存储成功');
            
            // 存储用户信息
            if (loginResponse.data?.user) {
              const userStoreSuccess = await StorageUtil.setUserInfo(loginResponse.data.user);
              if (userStoreSuccess) {
                console.info('用户信息存储成功：', JSON.stringify(loginResponse.data.user));
              } else {
                console.error('用户信息存储失败');
                // 即使用户信息存储失败，也继续登录流程，只是记录错误
              }
            } else {
              console.warn('服务器未返回用户信息');
            }
            
            promptAction.showToast({
              message: '登录成功',
              duration: 2000
            });
            console.info('登录成功，用户数据已存储：', JSON.stringify(loginResponse.data));
            
            // 跳转到主页
            router.replaceUrl({
              url: 'pages/Index'
            }).catch((error: Error) => {
              console.error('页面跳转失败：', error);
            });
          } else {
            promptAction.showToast({
              message: '登录信息保存失败',
              duration: 2000
            });
          }
        } else {
          promptAction.showToast({
            message: '服务器响应异常：缺少令牌',
            duration: 2000
          });
        }
      } else {
        promptAction.showToast({
          message: response.msg || '登录失败',
          duration: 2000
        });
      }
    } catch (error) {
      console.error('登录失败：' + error);
      promptAction.showToast({
        message: '网络错误，请重试',
        duration: 2000
      });
    } finally {
      this.isLoading = false;
    }
  }

  async aboutToAppear() {
    // 初始化存储
    await StorageUtil.init(this.context);
    
    // 检查是否已经登录
    const isLoggedIn = await StorageUtil.isUserLoggedIn();
    if (isLoggedIn) {
      // 如果用户已登录，直接跳转到主页
      console.info('用户已登录，直接跳转到主页');
      router.replaceUrl({
        url: 'pages/Index'
      }).catch((error: Error) => {
        console.error('页面跳转失败：', error);
      });
      return;
    }
    
    // 入场动画
    animateTo({
      duration: 600,
      curve: Curve.EaseOut
    }, () => {
      this.cardScale = 1;
    });
  }

  build() {
    Stack() {
      // 渐变背景
      Column()
        .width('100%')
        .height('100%')
        .linearGradient({
          direction: GradientDirection.Bottom,
          colors: [
            [0x667eea, 0.0],
            [0x764ba2, 0.5],
            [0x6366f1, 1.0]
          ]
        })

      // 装饰性圆形
      Circle({ width: 200, height: 200 })
        .fill(Color.White)
        .opacity(0.1)
        .position({ x: -50, y: 100 })

      Circle({ width: 150, height: 150 })
        .fill(Color.White)
        .opacity(0.08)
        .position({ x: '80%', y: '20%' })

      Circle({ width: 100, height: 100 })
        .fill(Color.White)
        .opacity(0.06)
        .position({ x: '70%', y: '70%' })

      // 主要内容
      Column() {

        // 登录卡片
        Column({ space: 24 }) {
          // 标题
          Text('欢迎回来')
            .fontSize(28)
            .fontWeight(FontWeight.Bold)
            .fontColor('#1F2937')
            .margin({ bottom: 8 })

          Text('请登录您的账户以继续使用')
            .fontSize(14)
            .fontColor('#6B7280')
            .margin({ bottom: 16 })

          // 输入框区域
          Column({ space: 20 }) {
            // 用户名输入框
            Column({ space: 12 }) {
              Row() {
                // 用户图标
                Text('👤')
                  .fontSize(16)
                  .margin({ right: 8 })

                Text('用户名')
                  .fontSize(14)
                  .fontColor('#374151')
                  .fontWeight(FontWeight.Medium)
              }
              .alignSelf(ItemAlign.Start)

              TextInput({ placeholder: '请输入您的用户名' })
                .height(52)
                .borderRadius(12)
                .backgroundColor('#F8FAFC')
                .border({ width: 1, color: '#E2E8F0' })
                .padding({ left: 20, right: 20 })
                .placeholderColor('#94A3B8')
                .onChange((value: string) => {
                  this.username = value;
                })
                .focusable(true)
                .onFocus(() => {
                  // 焦点动画效果可以在这里添加
                })
            }

            // 密码输入框
            Column({ space: 12 }) {
              Row() {
                // 密码图标
                Text('🔒')
                  .fontSize(16)
                  .margin({ right: 8 })

                Text('密码')
                  .fontSize(14)
                  .fontColor('#374151')
                  .fontWeight(FontWeight.Medium)
              }
              .alignSelf(ItemAlign.Start)

              Row() {
                TextInput({ placeholder: '请输入您的密码' })
                  .type(this.showPassword ? InputType.Normal : InputType.Password)
                  .height(52)
                  .layoutWeight(1)
                  .borderRadius(12)
                  .backgroundColor('#F8FAFC')
                  .border({ width: 1, color: '#E2E8F0' })
                  .padding({ left: 20, right: 20 })
                  .placeholderColor('#94A3B8')
                  .onChange((value: string) => {
                    this.password = value;
                  })

              }
            }
          }

          // 登录按钮
          Button(this.isLoading ? '登录中...' : '立即登录')
            .width('100%')
            .height(52)
            .backgroundColor(this.isLoading ? '#94A3B8' : '#6366F1')
            .borderRadius(12)
            .fontSize(16)
            .fontWeight(FontWeight.Medium)
            .fontColor(Color.White)
            .enabled(!this.isLoading)
            .shadow({
              radius: 12,
              color: this.isLoading ? '#40000000' : '#406366F1',
              offsetX: 0,
              offsetY: 4
            })
            .onClick(() => {
              this.handleLogin();
            })
            .animation({
              duration: 200,
              curve: Curve.EaseInOut
            })

          // 注册链接
          Row() {
            Text('还没有账户？')
              .fontSize(14)
              .fontColor('#6B7280')

            Text(' 立即注册')
              .fontSize(14)
              .fontColor('#6366F1')
              .fontWeight(FontWeight.Medium)
              .onClick(() => {
                promptAction.showToast({
                  message: '注册功能开发中',
                  duration: 2000
                });
              })
          }
          .margin({ top: 24 })

        }
        .width('100%')
        .padding(32)
        .backgroundColor(Color.White)
        .borderRadius(24)
        .shadow({
          radius: 30,
          color: '#20000000',
          offsetX: 0,
          offsetY: 15
        })
        .scale({ x: this.cardScale, y: this.cardScale })
        .animation({
          duration: 600,
          curve: Curve.EaseOut
        })

      }
      .width('100%')
      .height('100%')
      .padding({ left: 24, right: 24, top: 60, bottom: 40 })
      .justifyContent(FlexAlign.Center)
    }
    .width('100%')
    .height('100%')
  }
}
