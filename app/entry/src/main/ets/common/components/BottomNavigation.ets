import router from '@ohos.router';

/**
 * 底部导航栏项目接口
 */
export interface NavItem {
  id: number;
  text: string;
  icon: string;
  activeIcon: string;
  route: string;
}

/**
 * 底部导航栏组件
 */
@Component
export struct BottomNavigation {
  @Prop currentPageIndex: number = 0; // 通过属性接收当前页面索引
  @State currentIndex: number = 0;
  private navItems: NavItem[] = [
    {
      id: 0,
      text: '首页',
      icon: '🏠',
      activeIcon: '🏠',
      route: 'pages/Index'
    },
    {
      id: 1,
      text: '我的',
      icon: '👤',
      activeIcon: '👤',
      route: 'pages/UserProfile'
    }
  ];

  // 初始化当前选中状态
  aboutToAppear() {
    this.currentIndex = this.currentPageIndex;
    console.info('BottomNavigation初始化，当前页面索引：', this.currentPageIndex);
  }

  // 监听属性变化
  onPageShow() {
    this.currentIndex = this.currentPageIndex;
    console.info('BottomNavigation页面显示，当前页面索引：', this.currentPageIndex);
  }

  // 导航项点击处理
  private onNavItemClick(index: number, route: string) {
    console.info('导航点击：', index, route, '当前索引：', this.currentIndex);
    
    if (this.currentIndex === index) {
      console.info('已在当前页面，跳过跳转');
      return; // 如果是当前页面，不执行跳转
    }
    
    this.currentIndex = index;
    
    // 页面跳转
    router.replaceUrl({
      url: route
    }).then(() => {
      console.info('页面跳转成功：', route);
    }).catch((error: Error) => {
      console.error('页面跳转失败：', error);
      // 跳转失败时恢复原来的索引
      this.currentIndex = this.currentPageIndex;
    });
  }

  build() {
    Row() {
      ForEach(this.navItems, (item: NavItem, index: number) => {
        Column({ space: 4 }) {
          // 图标
          Text(item.icon)
            .fontSize(24)
            .fontColor(this.currentIndex === index ? '#007DFF' : '#99182431')
            .animation({
              duration: 200,
              curve: Curve.EaseInOut
            })
          
          // 文本
          Text(item.text)
            .fontSize(12)
            .fontColor(this.currentIndex === index ? '#007DFF' : '#99182431')
            .fontWeight(this.currentIndex === index ? FontWeight.Medium : FontWeight.Normal)
            .animation({
              duration: 200,
              curve: Curve.EaseInOut
            })
        }
        .width('50%')
        .height('100%')
        .justifyContent(FlexAlign.Center)
        .alignItems(HorizontalAlign.Center)
        .onClick(() => {
          this.onNavItemClick(index, item.route);
        })
        .backgroundColor(this.currentIndex === index ? '#F0F7FF' : Color.Transparent)
        .borderRadius(12)
        .animation({
          duration: 200,
          curve: Curve.EaseInOut
        })
      })
    }
    .width('100%')
    .height(70)
    .backgroundColor('#FFFFFF')
    .padding({ left: 16, right: 16, top: 8, bottom: 8 })
    .justifyContent(FlexAlign.SpaceEvenly)
    .alignItems(VerticalAlign.Center)
    .shadow({
      radius: 8,
      color: '#1A000000',
      offsetX: 0,
      offsetY: -2
    })
  }
} 