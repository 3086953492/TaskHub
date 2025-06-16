// 这里引入第 4 步创建的请求工具类
import HttpRequest from '../request/HttpRequest'

// 这里引入第 3 步创建的响应体
import { ResponseResult } from '../response/ResponseResult'

// 这里引入第 5 步创建的自定义用户对象，主要用于定义各参数的类型
import { UserInfo } from '../../entity/userInfo'

// 这里引入用户信息数据类型
import { UserInfoData } from '../utils/StorageUtil'

// 用户登录 API
export function userLogin(params: UserInfo): Promise<ResponseResult> {
  return HttpRequest.request({
    url: '/user/login',
    data: params,
    method:'POST'
  })
}

// 用户注册 API
export function userRegister(params: UserInfo): Promise<ResponseResult> {
  return HttpRequest.request({
    url: '/user/register',
    data: params,
    method:'POST'
  })
}

// 更新用户信息 API
export function updateUserInfo(params: UserInfoData): Promise<ResponseResult> {
  return HttpRequest.request({
    url: '/user/info',
    data: params,
    method:'PUT'
  })
}
