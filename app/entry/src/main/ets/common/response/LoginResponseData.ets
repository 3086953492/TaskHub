import { UserInfoData } from '../utils/StorageUtil';

/**
 * 登录响应数据接口
 */
export interface LoginResponseData {
  token?: string;
  user?: UserInfoData;
}

/**
 * 带类型的登录响应结果
 */
export interface LoginResponse {
  code: number;
  msg: string;
  data: LoginResponseData;
} 