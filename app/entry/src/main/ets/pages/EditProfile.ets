import { updateUserInfo } from '../common/api/user';
import { StorageUtil, UserInfoData } from '../common/utils/StorageUtil';
import promptAction from '@ohos.promptAction';
import router from '@ohos.router';
import common from '@ohos.app.ability.common';

@Entry
@Component
struct EditProfile {
  @State userInfo: UserInfoData | null = null;
  @State isLoading: boolean = false;
  @State isSaving: boolean = false;
  
  // 编辑字段
  @State editUsername: string = '';
  @State editEmail: string = '';
  @State editNickname: string = '';
  @State editPassword: string = '';
  @State confirmPassword: string = '';
  
  // 控制是否修改密码
  @State changePassword: boolean = false;
  
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
      if (userInfo) {
        this.userInfo = userInfo;
        // 初始化编辑字段
        this.editUsername = userInfo.username || '';
        this.editEmail = userInfo.email || '';
        this.editNickname = userInfo.nickname || '';
        console.info('加载用户信息成功：', JSON.stringify(userInfo));
      } else {
        promptAction.showToast({
          message: '未找到用户信息',
          duration: 2000
        });
        this.goBack();
      }
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

  // 验证邮箱格式
  private validateEmail(email: string): boolean {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  }

  // 保存用户信息
  async saveUserInfo() {
    // 验证输入
    if (!this.editUsername.trim()) {
      promptAction.showToast({
        message: '请输入用户名',
        duration: 2000
      });
      return;
    }

    if (this.editUsername.length < 3 || this.editUsername.length > 15) {
      promptAction.showToast({
        message: '用户名长度必须在3-15个字符之间',
        duration: 2000
      });
      return;
    }

    if (!this.editEmail.trim()) {
      promptAction.showToast({
        message: '请输入邮箱',
        duration: 2000
      });
      return;
    }

    if (!this.validateEmail(this.editEmail)) {
      promptAction.showToast({
        message: '请输入有效的邮箱地址',
        duration: 2000
      });
      return;
    }

    if (!this.editNickname.trim()) {
      promptAction.showToast({
        message: '请输入昵称',
        duration: 2000
      });
      return;
    }

    if (this.editNickname.length > 50) {
      promptAction.showToast({
        message: '昵称长度最多50个字符',
        duration: 2000
      });
      return;
    }

    // 如果要修改密码，验证密码
    if (this.changePassword) {
      if (!this.editPassword.trim()) {
        promptAction.showToast({
          message: '请输入新密码',
          duration: 2000
        });
        return;
      }

      if (this.editPassword.length < 6) {
        promptAction.showToast({
          message: '密码长度至少6位',
          duration: 2000
        });
        return;
      }

      if (this.editPassword !== this.confirmPassword) {
        promptAction.showToast({
          message: '两次输入的密码不一致',
          duration: 2000
        });
        return;
      }
    }

    this.isSaving = true;

    try {
      // 构建更新数据
      const updateData: UserInfoData = {
        username: this.editUsername,
        email: this.editEmail,
        nickname: this.editNickname
      };

      // 如果要修改密码，添加密码字段
      if (this.changePassword && this.editPassword) {
        updateData.password = this.editPassword;
      }

      const response = await updateUserInfo(updateData);
      
      console.info('更新用户信息响应：', JSON.stringify(response));
      
      if (response.code === 200) {
        // 更新本地存储的用户信息
        const updatedUserInfo: UserInfoData = {
          id: this.userInfo?.id,
          username: this.editUsername,
          email: this.editEmail,
          nickname: this.editNickname,
          avatar: this.userInfo?.avatar,
          role: this.userInfo?.role
          // 注意：不保存password到本地存储
        };
        
        const storeSuccess = await StorageUtil.setUserInfo(updatedUserInfo);
        if (storeSuccess) {
          promptAction.showToast({
            message: '保存成功',
            duration: 2000
          });
          console.info('用户信息更新成功');
          this.goBack();
        } else {
          promptAction.showToast({
            message: '本地数据保存失败',
            duration: 2000
          });
        }
      } else {
        promptAction.showToast({
          message: response.msg || '保存失败',
          duration: 2000
        });
      }
    } catch (error) {
      console.error('更新用户信息失败：', error);
      promptAction.showToast({
        message: '网络错误，请重试',
        duration: 2000
      });
    } finally {
      this.isSaving = false;
    }
  }

  // 返回上一页
  goBack() {
    router.back();
  }

  build() {
    Column() {
      // 标题栏
      Row() {
        // 返回按钮
        Button('← 返回')
          .backgroundColor(Color.Transparent)
          .fontColor('#007DFF')
          .fontSize(16)
          .onClick(() => {
            this.goBack();
          })

        Text('编辑资料')
          .fontSize(20)
          .fontWeight(FontWeight.Bold)
          .fontColor('#182431')
          .layoutWeight(1)
          .textAlign(TextAlign.Center)

        // 保存按钮
        Button(this.isSaving ? '保存中...' : '保存')
          .backgroundColor(this.isSaving ? '#99007DFF' : '#007DFF')
          .fontColor(Color.White)
          .fontSize(14)
          .padding({ left: 16, right: 16, top: 8, bottom: 8 })
          .borderRadius(8)
          .enabled(!this.isSaving)
          .onClick(() => {
            this.saveUserInfo();
          })
      }
      .width('100%')
      .height(60)
      .padding({ left: 16, right: 16 })
      .justifyContent(FlexAlign.SpaceBetween)
      .alignItems(VerticalAlign.Center)
      .backgroundColor('#FFFFFF')

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
        .layoutWeight(1)
        .justifyContent(FlexAlign.Center)
        .alignItems(HorizontalAlign.Center)
      } else {
        // 编辑表单
        Scroll() {
          Column({ space: 24 }) {
            // 用户名编辑
            Column({ space: 12 }) {
              Row() {
                Text('👤')
                  .fontSize(16)
                  .margin({ right: 8 })

                Text('用户名')
                  .fontSize(16)
                  .fontColor('#182431')
                  .fontWeight(FontWeight.Medium)
              }
              .alignSelf(ItemAlign.Start)

              TextInput({ text: this.editUsername, placeholder: '请输入用户名（3-15字符）' })
                .height(48)
                .borderRadius(8)
                .backgroundColor('#F8FAFC')
                .border({ width: 1, color: '#E2E8F0' })
                .padding({ left: 16, right: 16 })
                .onChange((value: string) => {
                  this.editUsername = value;
                })
            }

            // 邮箱编辑
            Column({ space: 12 }) {
              Row() {
                Text('📧')
                  .fontSize(16)
                  .margin({ right: 8 })

                Text('邮箱')
                  .fontSize(16)
                  .fontColor('#182431')
                  .fontWeight(FontWeight.Medium)
              }
              .alignSelf(ItemAlign.Start)

              TextInput({ text: this.editEmail, placeholder: '请输入邮箱地址' })
                .height(48)
                .borderRadius(8)
                .backgroundColor('#F8FAFC')
                .border({ width: 1, color: '#E2E8F0' })
                .padding({ left: 16, right: 16 })
                .type(InputType.Email)
                .onChange((value: string) => {
                  this.editEmail = value;
                })
            }

            // 昵称编辑
            Column({ space: 12 }) {
              Row() {
                Text('✨')
                  .fontSize(16)
                  .margin({ right: 8 })

                Text('昵称')
                  .fontSize(16)
                  .fontColor('#182431')
                  .fontWeight(FontWeight.Medium)
              }
              .alignSelf(ItemAlign.Start)

              TextInput({ text: this.editNickname, placeholder: '请输入昵称（最多50字符）' })
                .height(48)
                .borderRadius(8)
                .backgroundColor('#F8FAFC')
                .border({ width: 1, color: '#E2E8F0' })
                .padding({ left: 16, right: 16 })
                .onChange((value: string) => {
                  this.editNickname = value;
                })
            }

            // 修改密码选项
            Column({ space: 16 }) {
              Row() {
                Checkbox({ name: 'changePassword', group: 'changePasswordGroup' })
                  .select(this.changePassword)
                  .selectedColor('#007DFF')
                  .onChange((value: boolean) => {
                    this.changePassword = value;
                    if (!value) {
                      this.editPassword = '';
                      this.confirmPassword = '';
                    }
                  })
                  .margin({ right: 12 })

                Text('修改密码')
                  .fontSize(16)
                  .fontColor('#182431')
                  .fontWeight(FontWeight.Medium)
              }
              .alignItems(VerticalAlign.Center)
              .alignSelf(ItemAlign.Start)

              if (this.changePassword) {
                Column({ space: 12 }) {
                  // 新密码
                  Column({ space: 8 }) {
                    Row() {
                      Text('🔒')
                        .fontSize(14)
                        .margin({ right: 8 })

                      Text('新密码')
                        .fontSize(14)
                        .fontColor('#182431')
                    }
                    .alignSelf(ItemAlign.Start)

                    TextInput({ text: this.editPassword, placeholder: '请输入新密码（至少6位）' })
                      .height(48)
                      .borderRadius(8)
                      .backgroundColor('#F8FAFC')
                      .border({ width: 1, color: '#E2E8F0' })
                      .padding({ left: 16, right: 16 })
                      .type(InputType.Password)
                      .onChange((value: string) => {
                        this.editPassword = value;
                      })
                  }

                  // 确认新密码
                  Column({ space: 8 }) {
                    Row() {
                      Text('🔐')
                        .fontSize(14)
                        .margin({ right: 8 })

                      Text('确认新密码')
                        .fontSize(14)
                        .fontColor('#182431')
                    }
                    .alignSelf(ItemAlign.Start)

                    TextInput({ text: this.confirmPassword, placeholder: '请再次输入新密码' })
                      .height(48)
                      .borderRadius(8)
                      .backgroundColor('#F8FAFC')
                      .border({ width: 1, color: '#E2E8F0' })
                      .padding({ left: 16, right: 16 })
                      .type(InputType.Password)
                      .onChange((value: string) => {
                        this.confirmPassword = value;
                      })
                  }
                }
                .width('100%')
                .padding({ left: 24 })
              }
            }

            // 说明文字
            Text('• 用户名修改后需要重新登录\n• 密码修改后需要重新登录\n• 请确保邮箱地址有效')
              .fontSize(12)
              .fontColor('#6B7280')
              .lineHeight(20)
              .margin({ top: 16 })
          }
          .width('100%')
          .padding(24)
        }
        .layoutWeight(1)
        .scrollable(ScrollDirection.Vertical)
        .scrollBar(BarState.Auto)
      }
    }
    .width('100%')
    .height('100%')
    .backgroundColor('#F1F3F5')
  }
} 