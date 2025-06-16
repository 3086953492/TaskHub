import { userRegister } from '../common/api/user';
import { UserInfo } from '../entity/userInfo';
import { StorageUtil } from '../common/utils/StorageUtil';
import promptAction from '@ohos.promptAction';
import router from '@ohos.router';
import common from '@ohos.app.ability.common';

@Entry
@Component
struct Register {
  @State username: string = '';
  @State email: string = '';
  @State password: string = '';
  @State confirmPassword: string = '';
  @State nickname: string = '';
  @State isLoading: boolean = false;
  @State showPassword: boolean = false;
  @State showConfirmPassword: boolean = false;
  @State cardScale: number = 0.9;
  private context: common.UIAbilityContext = getContext(this) as common.UIAbilityContext;

  // 验证邮箱格式
  private validateEmail(email: string): boolean {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  }

  // 注册处理函数
  async handleRegister() {
    // 验证输入
    if (!this.username.trim()) {
      promptAction.showToast({
        message: '请输入用户名',
        duration: 2000
      });
      return;
    }

    if (this.username.length < 3 || this.username.length > 15) {
      promptAction.showToast({
        message: '用户名长度必须在3-15个字符之间',
        duration: 2000
      });
      return;
    }

    if (!this.email.trim()) {
      promptAction.showToast({
        message: '请输入邮箱',
        duration: 2000
      });
      return;
    }

    if (!this.validateEmail(this.email)) {
      promptAction.showToast({
        message: '请输入有效的邮箱地址',
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

    if (this.password.length < 6) {
      promptAction.showToast({
        message: '密码长度至少6位',
        duration: 2000
      });
      return;
    }

    if (this.password !== this.confirmPassword) {
      promptAction.showToast({
        message: '两次输入的密码不一致',
        duration: 2000
      });
      return;
    }

    if (!this.nickname.trim()) {
      promptAction.showToast({
        message: '请输入昵称',
        duration: 2000
      });
      return;
    }

    if (this.nickname.length > 50) {
      promptAction.showToast({
        message: '昵称长度最多50个字符',
        duration: 2000
      });
      return;
    }

    this.isLoading = true;

    try {
      const userInfo = new UserInfo();
      userInfo.username = this.username;
      userInfo.email = this.email;
      userInfo.password = this.password;
      userInfo.nickname = this.nickname;

      const response = await userRegister(userInfo);
      
      console.info('注册响应数据：', JSON.stringify(response));
      console.info('响应code值：', response.code);
      
      if (response.code === 200) {
        promptAction.showToast({
          message: '注册成功！请登录',
          duration: 2000
        });
        console.info('注册成功');
        
        // 注册成功后跳转到登录页面
        router.replaceUrl({
          url: 'pages/Login'
        }).catch((error: Error) => {
          console.error('页面跳转失败：', error);
        });
      } else {
        promptAction.showToast({
          message: response.msg || '注册失败',
          duration: 2000
        });
      }
    } catch (error) {
      console.error('注册失败：' + error);
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
        // 注册卡片
        Column({ space: 24 }) {
          // 标题
          Text('创建新账户')
            .fontSize(28)
            .fontWeight(FontWeight.Bold)
            .fontColor('#1F2937')
            .margin({ bottom: 8 })

          Text('请填写以下信息完成注册')
            .fontSize(14)
            .fontColor('#6B7280')
            .margin({ bottom: 16 })

          // 滚动容器，包含所有输入框
          Scroll() {
            Column({ space: 20 }) {
              // 用户名输入框
              Column({ space: 12 }) {
                Row() {
                  Text('👤')
                    .fontSize(16)
                    .margin({ right: 8 })

                  Text('用户名')
                    .fontSize(14)
                    .fontColor('#374151')
                    .fontWeight(FontWeight.Medium)
                }
                .alignSelf(ItemAlign.Start)

                TextInput({ placeholder: '请输入用户名（3-15字符）' })
                  .height(52)
                  .borderRadius(12)
                  .backgroundColor('#F8FAFC')
                  .border({ width: 1, color: '#E2E8F0' })
                  .padding({ left: 20, right: 20 })
                  .placeholderColor('#94A3B8')
                  .onChange((value: string) => {
                    this.username = value;
                  })
              }

              // 邮箱输入框
              Column({ space: 12 }) {
                Row() {
                  Text('📧')
                    .fontSize(16)
                    .margin({ right: 8 })

                  Text('邮箱')
                    .fontSize(14)
                    .fontColor('#374151')
                    .fontWeight(FontWeight.Medium)
                }
                .alignSelf(ItemAlign.Start)

                TextInput({ placeholder: '请输入邮箱地址' })
                  .height(52)
                  .borderRadius(12)
                  .backgroundColor('#F8FAFC')
                  .border({ width: 1, color: '#E2E8F0' })
                  .padding({ left: 20, right: 20 })
                  .placeholderColor('#94A3B8')
                  .type(InputType.Email)
                  .onChange((value: string) => {
                    this.email = value;
                  })
              }

              // 昵称输入框
              Column({ space: 12 }) {
                Row() {
                  Text('✨')
                    .fontSize(16)
                    .margin({ right: 8 })

                  Text('昵称')
                    .fontSize(14)
                    .fontColor('#374151')
                    .fontWeight(FontWeight.Medium)
                }
                .alignSelf(ItemAlign.Start)

                TextInput({ placeholder: '请输入昵称（最多50字符）' })
                  .height(52)
                  .borderRadius(12)
                  .backgroundColor('#F8FAFC')
                  .border({ width: 1, color: '#E2E8F0' })
                  .padding({ left: 20, right: 20 })
                  .placeholderColor('#94A3B8')
                  .onChange((value: string) => {
                    this.nickname = value;
                  })
              }

              // 密码输入框
              Column({ space: 12 }) {
                Row() {
                  Text('🔒')
                    .fontSize(16)
                    .margin({ right: 8 })

                  Text('密码')
                    .fontSize(14)
                    .fontColor('#374151')
                    .fontWeight(FontWeight.Medium)
                }
                .alignSelf(ItemAlign.Start)

                TextInput({ placeholder: '请输入密码（至少6位）' })
                  .type(this.showPassword ? InputType.Normal : InputType.Password)
                  .height(52)
                  .borderRadius(12)
                  .backgroundColor('#F8FAFC')
                  .border({ width: 1, color: '#E2E8F0' })
                  .padding({ left: 20, right: 20 })
                  .placeholderColor('#94A3B8')
                  .onChange((value: string) => {
                    this.password = value;
                  })
              }

              // 确认密码输入框
              Column({ space: 12 }) {
                Row() {
                  Text('🔐')
                    .fontSize(16)
                    .margin({ right: 8 })

                  Text('确认密码')
                    .fontSize(14)
                    .fontColor('#374151')
                    .fontWeight(FontWeight.Medium)
                }
                .alignSelf(ItemAlign.Start)

                TextInput({ placeholder: '请再次输入密码' })
                  .type(this.showConfirmPassword ? InputType.Normal : InputType.Password)
                  .height(52)
                  .borderRadius(12)
                  .backgroundColor('#F8FAFC')
                  .border({ width: 1, color: '#E2E8F0' })
                  .padding({ left: 20, right: 20 })
                  .placeholderColor('#94A3B8')
                  .onChange((value: string) => {
                    this.confirmPassword = value;
                  })
              }
            }
            .width('100%')
          }
          .height(300)
          .scrollable(ScrollDirection.Vertical)
          .scrollBar(BarState.Auto)

          // 注册按钮
          Button(this.isLoading ? '注册中...' : '立即注册')
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
              this.handleRegister();
            })
            .animation({
              duration: 200,
              curve: Curve.EaseInOut
            })

          // 登录链接
          Row() {
            Text('已有账户？')
              .fontSize(14)
              .fontColor('#6B7280')

            Text(' 立即登录')
              .fontSize(14)
              .fontColor('#6366F1')
              .fontWeight(FontWeight.Medium)
              .onClick(() => {
                router.replaceUrl({
                  url: 'pages/Login'
                }).catch((error: Error) => {
                  console.error('页面跳转失败：', error);
                });
              })
          }
          .margin({ top: 16 })

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