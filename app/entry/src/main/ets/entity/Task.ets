/**
 * 任务数据模型
 */
export interface TaskItem {
  task_id: number;
  status: number;
  created_at: string;
  title: string;
  priority: number;
  due_date: string;
}

/**
 * 任务列表响应数据
 */
export interface TaskListResponse {
  code: number;
  msg: string;
  data: TaskItem[];
} 