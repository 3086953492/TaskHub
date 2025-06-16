import { StorageUtil, UserInfoData } from '../common/utils/StorageUtil';
import { BottomNavigation } from '../common/components/BottomNavigation';
import { getTaskList, TaskQueryParams } from '../common/api/task';
import { TaskItem, TaskListResponse } from '../entity/Task';
import promptAction from '@ohos.promptAction';
import router from '@ohos.router';
import common from '@ohos.app.ability.common';

@Entry
@Component
struct Index {
  private context: common.UIAbilityContext = getContext(this) as common.UIAbilityContext;
  
  // 任务列表数据
  @State taskList: TaskItem[] = [];
  @State isLoading: boolean = true;
  @State currentUser: UserInfoData | null = null;
  @State currentPage: number = 1;
  private pageSize: number = 10; // 每页显示数量

  // 获取状态文本
  getStatusText(status: number): string {
    switch (status) {
      case 1: return "待处理";
      case 2: return "进行中";
      case 3: return "已完成";
      case 4: return "已取消";
      default: return "未知";
    }
  }

  // 获取状态颜色
  getStatusColor(status: number): string {
    switch (status) {
      case 1: return '#FF6B35'; // 橙色 - 待处理
      case 2: return '#007DFF'; // 蓝色 - 进行中
      case 3: return '#00C853'; // 绿色 - 已完成
      case 4: return '#9E9E9E'; // 灰色 - 已取消
      default: return '#9E9E9E';
    }
  }

  // 获取优先级文本
  getPriorityText(priority: number): string {
    switch (priority) {
      case 1: return "高";
      case 2: return "中";
      case 3: return "低";
      default: return "未知";
    }
  }

  // 获取优先级颜色
  getPriorityColor(priority: number): string {
    switch (priority) {
      case 1: return '#F44336'; // 红色 - 高优先级
      case 2: return '#FF9800'; // 橙色 - 中优先级
      case 3: return '#4CAF50'; // 绿色 - 低优先级
      default: return '#9E9E9E';
    }
  }

  // 格式化日期
  formatDate(dateString: string): string {
    const date = new Date(dateString);
    return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')}`;
  }
  
  // 初始化
  async aboutToAppear() {
    await StorageUtil.init(this.context);
    await this.loadUserInfo();
    await this.loadTaskList();
  }

  // 加载用户信息
  async loadUserInfo() {
    try {
      this.currentUser = await StorageUtil.getUserInfo();
      console.info('当前用户信息：', JSON.stringify(this.currentUser));
    } catch (error) {
      console.error('获取用户信息失败：', error);
    }
  }

  // 加载任务列表
  async loadTaskList() {
    try {
      this.isLoading = true;

      // 根据用户角色构建查询参数
      const queryParams: TaskQueryParams = {
        page: this.currentPage,
        page_size: this.pageSize
      };

      // 根据用户角色决定查询条件
      if (this.currentUser?.role === 'admin') {
        // 管理员：查看所有任务，所以不传任何参数
      } else {
        // 普通用户：查看未分配的任务（assignee_id为0）
        queryParams.assignee_id = 0;
      }

      console.info('查询参数：', JSON.stringify(queryParams));

      const response = await getTaskList(queryParams) as TaskListResponse;
      
      console.info('任务列表响应：', JSON.stringify(response));

      if (response.code === 200) {
        this.taskList = response.data || [];
        console.info(`成功获取${this.taskList.length}个任务`);
      } else {
        console.error('获取任务列表失败：', response.msg);
        promptAction.showToast({
          message: response.msg || '获取任务列表失败',
          duration: 2000
        });
      }

    } catch (error) {
      console.error('加载任务列表异常：', error);
      promptAction.showToast({
        message: '网络错误，请重试',
        duration: 2000
      });
    } finally {
      this.isLoading = false;
    }
  }

  // 刷新任务列表
  async refreshTaskList() {
    this.currentPage = 1;
    await this.loadTaskList();
  }

  build() {
    Column() {
      // 内容区域
      Column() {
        // 标题栏
        Row() {
          Text('任务列表')
            .fontSize(24)
            .fontWeight(FontWeight.Bold)
            .fontColor('#182431')
            .layoutWeight(1)
          
          // 刷新按钮
          Button('刷新')
            .backgroundColor('#007DFF')
            .fontColor(Color.White)
            .fontSize(14)
            .padding({ left: 16, right: 16, top: 8, bottom: 8 })
            .borderRadius(8)
            .onClick(() => {
              this.refreshTaskList();
            })
        }
        .width('100%')
        .height(60)
        .padding({ left: 16, right: 16 })
        .justifyContent(FlexAlign.SpaceBetween)
        .alignItems(VerticalAlign.Center)
        .backgroundColor('#FFFFFF')
        .margin(10)

        if (this.isLoading) {
          // 加载状态
          Column() {
            Text('🔄')
              .fontSize(48)
              .margin({ bottom: 16 })
            Text('加载任务列表中...')
              .fontSize(16)
              .fontColor('#99182431')
          }
          .width('100%')
          .height(200)
          .justifyContent(FlexAlign.Center)
          .alignItems(HorizontalAlign.Center)
        } else if (this.taskList.length === 0) {
          // 空数据状态
          Column() {
            Text('📝')
              .fontSize(48)
              .margin({ bottom: 16 })
            Text('暂无任务数据')
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
                this.refreshTaskList();
              })
          }
          .width('100%')
          .height(200)
          .justifyContent(FlexAlign.Center)
          .alignItems(HorizontalAlign.Center)
        } else {
          // 任务列表
          List({ space: 12 }) {
            ForEach(this.taskList, (task: TaskItem) => {
              ListItem() {
                Column({ space: 8 }) {
                  // 任务标题和优先级
                  Row() {
                    Text(task.title)
                      .fontSize(16)
                      .fontWeight(FontWeight.Medium)
                      .fontColor('#182431')
                      .maxLines(2)
                      .textOverflow({ overflow: TextOverflow.Ellipsis })
                      .layoutWeight(1)

                    Text(this.getPriorityText(task.priority))
                      .fontSize(12)
                      .fontColor(Color.White)
                      .backgroundColor(this.getPriorityColor(task.priority))
                      .padding({ left: 8, right: 8, top: 4, bottom: 4 })
                      .borderRadius(12)
                  }
                  .width('100%')
                  .alignItems(VerticalAlign.Top)

                  // 状态和创建时间
                  Row() {
                    Text(this.getStatusText(task.status))
                      .fontSize(14)
                      .fontColor(this.getStatusColor(task.status))
                      .fontWeight(FontWeight.Medium)

                    Blank()

                    Text(`创建: ${this.formatDate(task.created_at)}`)
                      .fontSize(12)
                      .fontColor('#99182431')
                  }
                  .width('100%')

                  // 截止日期
                  Row() {
                    Text(`截止: ${this.formatDate(task.due_date)}`)
                      .fontSize(12)
                      .fontColor('#99182431')

                    Blank()

                    Text(`#${task.task_id}`)
                      .fontSize(12)
                      .fontColor('#66182431')
                  }
                  .width('100%')
                }
                .width('100%')
                .padding(16)
                .backgroundColor('#FFFFFF')
                .borderRadius(8)
                .shadow({
                  radius: 4,
                  color: '#1A000000',
                  offsetX: 0,
                  offsetY: 2
                })
              }
            })
          }
          .width('100%')
          .layoutWeight(1)
          .padding({ left: 16, right: 16, bottom: 16 })
          .scrollBar(BarState.Auto)
        }
      }
      .layoutWeight(1)
      .width('100%')
      .backgroundColor('#F1F3F5')

      // 底部导航栏
      BottomNavigation({ currentPageIndex: 0 })
    }
    .width('100%')
    .height('100%')
    .backgroundColor('#F1F3F5')
  }
}