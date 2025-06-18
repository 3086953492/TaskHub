// 消息类型
export type MessageType = 'success' | 'error' | 'warning' | 'info'

// 消息配置接口
export interface MessageConfig {
  type: MessageType
  title: string
  message?: string
  duration?: number
  showClose?: boolean
}

// 消息实例接口
export interface MessageInstance {
  close: () => void
}

// 创建消息元素
const createMessageElement = (config: MessageConfig): HTMLElement => {
  const messageEl = document.createElement('div')
  messageEl.className = `message message-${config.type}`
  
  // 设置样式
  Object.assign(messageEl.style, {
    position: 'fixed',
    top: '20px',
    left: '50%',
    transform: 'translateX(-50%)',
    padding: '12px 20px',
    borderRadius: '8px',
    backgroundColor: getBackgroundColor(config.type),
    color: getTextColor(config.type),
    border: `1px solid ${getBorderColor(config.type)}`,
    boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)',
    zIndex: '9999',
    minWidth: '300px',
    maxWidth: '500px',
    fontSize: '14px',
    lineHeight: '1.5',
    fontFamily: 'system-ui, -apple-system, sans-serif',
    display: 'flex',
    alignItems: 'center',
    gap: '8px'
  })
  
  // 添加图标
  const iconEl = document.createElement('span')
  iconEl.innerHTML = getIcon(config.type)
  iconEl.style.fontSize = '16px'
  messageEl.appendChild(iconEl)
  
  // 添加文本内容
  const textEl = document.createElement('div')
  textEl.style.flex = '1'
  
  const titleEl = document.createElement('div')
  titleEl.style.fontWeight = '600'
  titleEl.textContent = config.title
  textEl.appendChild(titleEl)
  
  if (config.message) {
    const messageTextEl = document.createElement('div')
    messageTextEl.style.marginTop = '4px'
    messageTextEl.style.opacity = '0.9'
    messageTextEl.textContent = config.message
    textEl.appendChild(messageTextEl)
  }
  
  messageEl.appendChild(textEl)
  
  // 添加关闭按钮
  if (config.showClose !== false) {
    const closeEl = document.createElement('button')
    closeEl.innerHTML = '×'
    closeEl.style.cssText = `
      background: none;
      border: none;
      font-size: 18px;
      cursor: pointer;
      color: inherit;
      opacity: 0.6;
      padding: 0;
      margin-left: 8px;
      width: 20px;
      height: 20px;
      display: flex;
      align-items: center;
      justify-content: center;
    `
    closeEl.addEventListener('mouseenter', () => {
      closeEl.style.opacity = '1'
    })
    closeEl.addEventListener('mouseleave', () => {
      closeEl.style.opacity = '0.6'
    })
    messageEl.appendChild(closeEl)
    
    closeEl.addEventListener('click', () => {
      removeMessage(messageEl)
    })
  }
  
  return messageEl
}

// 获取背景颜色
const getBackgroundColor = (type: MessageType): string => {
  switch (type) {
    case 'success': return '#f0f9ff'
    case 'error': return '#fef2f2'
    case 'warning': return '#fffbeb'
    case 'info': return '#f8fafc'
    default: return '#f8fafc'
  }
}

// 获取文本颜色
const getTextColor = (type: MessageType): string => {
  switch (type) {
    case 'success': return '#065f46'
    case 'error': return '#991b1b'
    case 'warning': return '#92400e'
    case 'info': return '#374151'
    default: return '#374151'
  }
}

// 获取边框颜色
const getBorderColor = (type: MessageType): string => {
  switch (type) {
    case 'success': return '#34d399'
    case 'error': return '#f87171'
    case 'warning': return '#fbbf24'
    case 'info': return '#94a3b8'
    default: return '#94a3b8'
  }
}

// 获取图标
const getIcon = (type: MessageType): string => {
  switch (type) {
    case 'success': return '✓'
    case 'error': return '✕'
    case 'warning': return '⚠'
    case 'info': return 'ⓘ'
    default: return 'ⓘ'
  }
}

// 移除消息
const removeMessage = (messageEl: HTMLElement) => {
  messageEl.style.opacity = '0'
  messageEl.style.transform = 'translateX(-50%) translateY(-10px)'
  messageEl.style.transition = 'all 0.3s ease'
  
  setTimeout(() => {
    if (messageEl.parentElement) {
      messageEl.parentElement.removeChild(messageEl)
    }
  }, 300)
}

// 显示消息
export const showMessage = (config: MessageConfig): MessageInstance => {
  const messageEl = createMessageElement(config)
  document.body.appendChild(messageEl)
  
  // 添加显示动画
  messageEl.style.opacity = '0'
  messageEl.style.transform = 'translateX(-50%) translateY(-10px)'
  messageEl.style.transition = 'all 0.3s ease'
  
  setTimeout(() => {
    messageEl.style.opacity = '1'
    messageEl.style.transform = 'translateX(-50%) translateY(0)'
  }, 10)
  
  // 自动关闭
  const duration = config.duration || 3000
  if (duration > 0) {
    setTimeout(() => {
      removeMessage(messageEl)
    }, duration)
  }
  
  return {
    close: () => removeMessage(messageEl)
  }
}

// 便捷方法
export const message = {
  success: (title: string, message?: string, duration?: number) => 
    showMessage({ type: 'success', title, message, duration }),
  
  error: (title: string, message?: string, duration?: number) => 
    showMessage({ type: 'error', title, message, duration }),
  
  warning: (title: string, message?: string, duration?: number) => 
    showMessage({ type: 'warning', title, message, duration }),
  
  info: (title: string, message?: string, duration?: number) => 
    showMessage({ type: 'info', title, message, duration })
} 