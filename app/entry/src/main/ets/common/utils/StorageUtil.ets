import dataPreferences from '@ohos.data.preferences';

/**
 * 用户信息接口
 */
export interface UserInfoData {
  id?: number;
  username?: string;
  email?: string;
  nickname?: string;
  avatar?: string;
  role?: string;
}

/**
 * 本地存储工具类
 */
export class StorageUtil {
  private static readonly STORE_NAME = 'app_preferences';
  private static readonly TOKEN_KEY = 'user_token';
  private static readonly USER_INFO_KEY = 'user_info';
  
  private static preferences: dataPreferences.Preferences | null = null;
  
  /**
   * 初始化存储
   */
  static async init(context: Context): Promise<void> {
    try {
      StorageUtil.preferences = await dataPreferences.getPreferences(context, StorageUtil.STORE_NAME);
    } catch (error) {
      console.error('初始化存储失败：', error);
    }
  }
  
  /**
   * 存储token
   */
  static async setToken(token: string): Promise<boolean> {
    try {
      if (!StorageUtil.preferences) {
        console.error('存储未初始化');
        return false;
      }
      await StorageUtil.preferences.put(StorageUtil.TOKEN_KEY, token);
      await StorageUtil.preferences.flush();
      return true;
    } catch (error) {
      console.error('存储token失败：', error);
      return false;
    }
  }
  
  /**
   * 获取token
   */
  static async getToken(): Promise<string> {
    try {
      if (!StorageUtil.preferences) {
        console.error('存储未初始化');
        return '';
      }
      return await StorageUtil.preferences.get(StorageUtil.TOKEN_KEY, '') as string;
    } catch (error) {
      console.error('获取token失败：', error);
      return '';
    }
  }
  
  /**
   * 删除token
   */
  static async removeToken(): Promise<boolean> {
    try {
      if (!StorageUtil.preferences) {
        console.error('存储未初始化');
        return false;
      }
      await StorageUtil.preferences.delete(StorageUtil.TOKEN_KEY);
      await StorageUtil.preferences.flush();
      return true;
    } catch (error) {
      console.error('删除token失败：', error);
      return false;
    }
  }
  
  /**
   * 存储用户信息
   */
  static async setUserInfo(userInfo: UserInfoData): Promise<boolean> {
    try {
      if (!StorageUtil.preferences) {
        console.error('存储未初始化');
        return false;
      }
      await StorageUtil.preferences.put(StorageUtil.USER_INFO_KEY, JSON.stringify(userInfo));
      await StorageUtil.preferences.flush();
      return true;
    } catch (error) {
      console.error('存储用户信息失败：', error);
      return false;
    }
  }
  
  /**
   * 获取用户信息
   */
  static async getUserInfo(): Promise<UserInfoData | null> {
    try {
      if (!StorageUtil.preferences) {
        console.error('存储未初始化');
        return null;
      }
      const userInfoStr = await StorageUtil.preferences.get(StorageUtil.USER_INFO_KEY, '') as string;
      return userInfoStr ? JSON.parse(userInfoStr) as UserInfoData : null;
    } catch (error) {
      console.error('获取用户信息失败：', error);
      return null;
    }
  }
  
  /**
   * 删除用户信息
   */
  static async removeUserInfo(): Promise<boolean> {
    try {
      if (!StorageUtil.preferences) {
        console.error('存储未初始化');
        return false;
      }
      await StorageUtil.preferences.delete(StorageUtil.USER_INFO_KEY);
      await StorageUtil.preferences.flush();
      return true;
    } catch (error) {
      console.error('删除用户信息失败：', error);
      return false;
    }
  }
  
  /**
   * 检查用户是否已登录（有效token和用户信息）
   */
  static async isUserLoggedIn(): Promise<boolean> {
    try {
      const token = await StorageUtil.getToken();
      const userInfo = await StorageUtil.getUserInfo();
      return !!token && !!userInfo;
    } catch (error) {
      console.error('检查登录状态失败：', error);
      return false;
    }
  }
  
  /**
   * 用户退出登录（清除token和用户信息）
   */
  static async logout(): Promise<boolean> {
    try {
      const tokenRemoved = await StorageUtil.removeToken();
      const userInfoRemoved = await StorageUtil.removeUserInfo();
      return tokenRemoved && userInfoRemoved;
    } catch (error) {
      console.error('退出登录失败：', error);
      return false;
    }
  }
  
  /**
   * 清除所有存储数据
   */
  static async clear(): Promise<boolean> {
    try {
      if (!StorageUtil.preferences) {
        console.error('存储未初始化');
        return false;
      }
      await StorageUtil.preferences.clear();
      await StorageUtil.preferences.flush();
      return true;
    } catch (error) {
      console.error('清除存储失败：', error);
      return false;
    }
  }
} 