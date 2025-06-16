import { StorageUtil, UserInfoData } from '../common/utils/StorageUtil';
import { BottomNavigation } from '../common/components/BottomNavigation';
import promptAction from '@ohos.promptAction';
import router from '@ohos.router';
import common from '@ohos.app.ability.common';

@Entry
@Component
struct UserProfile {
  @State userInfo: UserInfoData | null = null;
  @State isLoading: boolean = true;
  private context: common.UIAbilityContext = getContext(this) as common.UIAbilityContext;

  async aboutToAppear() {
    await StorageUtil.init(this.context);
    await this.loadUserInfo();
  }

  // 加载用户信息
  async loadUserInfo() {
    try {
      this.isLoading = true;
      const userInfo = await StorageUtil.getUserInfo();
      this.userInfo = userInfo;
      console.info('获取到的用户信息：', JSON.stringify(userInfo));
    } catch (error) {
      console.error('加载用户信息失败：', error);
      promptAction.showToast({
        message: '加载用户信息失败',
        duration: 2000
      });
    } finally {
      this.isLoading = false;
    }
  }

  // 退出登录
  async logout() {
    const success = await StorageUtil.logout();
    if (success) {
      promptAction.showToast({
        message: '已退出登录',
        duration: 2000
      });
      
      // 跳转到登录页面
      router.replaceUrl({
        url: 'pages/Login'
      }).catch((error: Error) => {
        console.error('页面跳转失败：', error);
      });
    } else {
      promptAction.showToast({
        message: '退出登录失败',
        duration: 2000
      });
    }
  }

  // 编辑用户信息（暂时只显示提示）
  editUserInfo() {
    promptAction.showToast({
      message: '编辑功能开发中',
      duration: 2000
    });
  }

  // 获取角色显示文本
  getRoleText(role?: string): string {
    switch (role) {
      case 'admin': return '管理员';
      case 'user': return '普通用户';
      default: return '普通用户';
    }
  }

  // 构建头像文字
  @Builder
  buildAvatarText() {
    Text(this.userInfo?.nickname ? this.userInfo.nickname.charAt(0).toUpperCase() : 
         this.userInfo?.username ? this.userInfo.username.charAt(0).toUpperCase() : '?')
      .fontSize(32)
      .fontColor(Color.White)
      .fontWeight(FontWeight.Bold)
  }

  build() {
    Column() {
      // 内容区域
      Column() {
        // 标题栏
        Row() {
          Text('我的信息')
            .fontSize(24)
            .fontWeight(FontWeight.Bold)
            .fontColor('#182431')
            .layoutWeight(1)
          
          // 退出按钮
          Button('退出登录')
            .backgroundColor('#FF4444')
            .fontColor(Color.White)
            .fontSize(14)
            .padding({ left: 16, right: 16, top: 8, bottom: 8 })
            .borderRadius(8)
            .onClick(() => {
              this.logout();
            })
        }
        .width('100%')
        .height(60)
        .padding({ left: 16, right: 16 })
        .justifyContent(FlexAlign.SpaceBetween)
        .alignItems(VerticalAlign.Center)
        .backgroundColor('#FFFFFF')
        .margin({ bottom: 16 })

        if (this.isLoading) {
          // 加载状态
          Column() {
            Text('🔄')
              .fontSize(48)
              .margin({ bottom: 16 })
            Text('加载中...')
              .fontSize(16)
              .fontColor('#99182431')
          }
          .width('100%')
          .height(200)
          .justifyContent(FlexAlign.Center)
          .alignItems(HorizontalAlign.Center)
        } else if (this.userInfo) {
          // 用户信息卡片
          Column({ space: 16 }) {
            // 头像和基本信息
            Column({ space: 16 }) {
              // 头像
              Circle({ width: 80, height: 80 })
                .fill('#007DFF')
                .overlay(this.buildAvatarText())

              // 用户名和昵称
              Column({ space: 8 }) {
                Text(this.userInfo.nickname || this.userInfo.username || '未知用户')
                  .fontSize(20)
                  .fontWeight(FontWeight.Bold)
                  .fontColor('#182431')

                Text(`@${this.userInfo.username || '未知'}`)
                  .fontSize(14)
                  .fontColor('#99182431')
              }
            }
            .width('100%')
            .padding(24)
            .alignItems(HorizontalAlign.Center)

            // 详细信息列表
            Column({ space: 12 }) {
              // 用户ID
              Row() {
                Row({ space: 12 }) {
                  Text('🆔')
                    .fontSize(20)
                  Text('用户ID')
                    .fontSize(16)
                    .fontColor('#182431')
                    .fontWeight(FontWeight.Medium)
                }
                .layoutWeight(1)

                Text(this.userInfo.id?.toString() || '未知')
                  .fontSize(16)
                  .fontColor('#99182431')
              }
              .width('100%')
              .padding({ left: 16, right: 16, top: 12, bottom: 12 })

              Divider().color('#F1F3F5')

              // 邮箱
              Row() {
                Row({ space: 12 }) {
                  Text('📧')
                    .fontSize(20)
                  Text('邮箱')
                    .fontSize(16)
                    .fontColor('#182431')
                    .fontWeight(FontWeight.Medium)
                }
                .layoutWeight(1)

                Text(this.userInfo.email || '未设置')
                  .fontSize(16)
                  .fontColor('#99182431')
              }
              .width('100%')
              .padding({ left: 16, right: 16, top: 12, bottom: 12 })

              Divider().color('#F1F3F5')

              // 角色
              Row() {
                Row({ space: 12 }) {
                  Text('👑')
                    .fontSize(20)
                  Text('角色')
                    .fontSize(16)
                    .fontColor('#182431')
                    .fontWeight(FontWeight.Medium)
                }
                .layoutWeight(1)

                Text(this.getRoleText(this.userInfo.role))
                  .fontSize(16)
                  .fontColor('#99182431')
              }
              .width('100%')
              .padding({ left: 16, right: 16, top: 12, bottom: 12 })

              Divider().color('#F1F3F5')

            }
            .width('100%')
            .backgroundColor('#FFFFFF')
            .borderRadius(12)
            .padding({ top: 8, bottom: 8 })

            // 编辑按钮
            Button('编辑资料')
              .width('100%')
              .height(48)
              .backgroundColor('#007DFF')
              .fontColor(Color.White)
              .fontSize(16)
              .fontWeight(FontWeight.Medium)
              .borderRadius(12)
              .margin({ top: 16 })
              .onClick(() => {
                this.editUserInfo();
              })
          }
          .width('100%')
          .padding({ left: 16, right: 16 })
        } else {
          // 无用户信息状态
          Column() {
            Text('😕')
              .fontSize(48)
              .margin({ bottom: 16 })
            Text('未找到用户信息')
              .fontSize(16)
              .fontColor('#99182431')
              .margin({ bottom: 16 })
            
            Button('重新加载')
              .backgroundColor('#007DFF')
              .fontColor(Color.White)
              .fontSize(14)
              .padding({ left: 20, right: 20, top: 10, bottom: 10 })
              .borderRadius(8)
              .onClick(() => {
                this.loadUserInfo();
              })
          }
          .width('100%')
          .height(200)
          .justifyContent(FlexAlign.Center)
          .alignItems(HorizontalAlign.Center)
        }
      }
      .layoutWeight(1)
      .width('100%')
      .backgroundColor('#F1F3F5')

      // 底部导航栏
      BottomNavigation({ currentPageIndex: 1 })
    }
    .width('100%')
    .height('100%')
    .backgroundColor('#F1F3F5')
  }
} 