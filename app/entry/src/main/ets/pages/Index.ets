import { StorageUtil } from '../common/utils/StorageUtil';
import { BottomNavigation } from '../common/components/BottomNavigation';
import promptAction from '@ohos.promptAction';
import router from '@ohos.router';
import common from '@ohos.app.ability.common';

// 任务数据模型
interface TaskItem {
  task_id: number;
  status: number;
  created_at: string;
  title: string;
  priority: number;
  due_date: string;
}

@Entry
@Component
struct Index {
  private context: common.UIAbilityContext = getContext(this) as common.UIAbilityContext;
  
  // 静态任务数据
  @State taskList: TaskItem[] = [
    {
      task_id: 1,
      status: 1,
      created_at: "2024-01-15T09:00:00Z",
      title: "完成产品需求文档撰写",
      priority: 1,
      due_date: "2024-02-01T18:00:00Z"
    },
    {
      task_id: 2,
      status: 2,
      created_at: "2024-01-12T14:30:00Z",
      title: "优化系统性能监控",
      priority: 2,
      due_date: "2024-01-30T12:00:00Z"
    },
    {
      task_id: 3,
      status: 3,
      created_at: "2024-01-10T11:15:00Z",
      title: "用户界面设计评审",
      priority: 1,
      due_date: "2024-01-25T16:00:00Z"
    },
    {
      task_id: 4,
      status: 1,
      created_at: "2024-01-14T16:45:00Z",
      title: "数据库备份策略制定",
      priority: 3,
      due_date: "2024-02-15T10:00:00Z"
    },
    {
      task_id: 5,
      status: 2,
      created_at: "2024-01-13T13:20:00Z",
      title: "移动端适配测试",
      priority: 2,
      due_date: "2024-01-28T14:00:00Z"
    }
  ];

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
  
  // 初始化存储
  async aboutToAppear() {
    await StorageUtil.init(this.context);
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
        }
        .width('100%')
        .height(60)
        .padding({ left: 16, right: 16 })
        .justifyContent(FlexAlign.SpaceBetween)
        .alignItems(VerticalAlign.Center)
        .backgroundColor('#FFFFFF')
        .margin(10)

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