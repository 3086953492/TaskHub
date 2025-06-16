import HttpRequest from '../request/HttpRequest'
import { ResponseResult } from '../response/ResponseResult'
import { RequestBody } from '../response/RequestBody'

/**
 * 任务查询参数接口
 */
export interface TaskQueryParams {
  page: number;
  page_size: number;
  assignee_id?: number;
  creator_id?: number;
}

/**
 * 构建查询字符串
 * @param params 查询参数
 * @returns 查询字符串
 */
function buildQueryString(params: TaskQueryParams): string {
  console.info('开始构建查询字符串，输入参数：', JSON.stringify(params));
  
  // 使用数组收集所有有效的查询参数
  const validParams: string[] = [];
  
  // 添加必填参数
  validParams.push(`page=${params.page}`);
  validParams.push(`page_size=${params.page_size}`);
  
  // 检查可选参数
  if (typeof params.assignee_id === 'number') {
    validParams.push(`assignee_id=${params.assignee_id}`);
    console.info('添加 assignee_id：', params.assignee_id);
  }
  
  if (typeof params.creator_id === 'number') {
    validParams.push(`creator_id=${params.creator_id}`);
    console.info('添加 creator_id：', params.creator_id);
  }

  const result = validParams.join('&');
  console.info('最终查询字符串：', result);
  console.info('有效参数数组：', JSON.stringify(validParams));
  
  return result;
}

/**
 * 获取任务列表
 * @param params 查询参数
 * @returns Promise<ResponseResult>
 */
export function getTaskList(params: TaskQueryParams): Promise<ResponseResult> {
  // 构建查询字符串
  const queryString = buildQueryString(params);
  const fullUrl = `/task?${queryString}`;

  console.info('完整请求URL：', fullUrl);

  // 创建RequestBody对象
  const requestBody = new RequestBody();
  requestBody.url = fullUrl;
  requestBody.method = 'GET';
  requestBody.data = new Object(); // GET请求data为空对象

  return HttpRequest.request(requestBody);
} 