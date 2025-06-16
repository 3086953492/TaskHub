// 这里需要引入第 2 步和第 3 步创建的请求和响应实体
import { ResponseResult } from '../response/ResponseResult'
import { RequestBody } from '../response/RequestBody'

// 这里需要引入第 1 步公共常量类
import { Domain, HttpStatus, ContentType } from '../constant/CommonConstant'

// 引入存储工具类获取token
import { StorageUtil } from '../utils/StorageUtil'

// 这里引用 HTTP 库
import http from '@ohos.net.http'

// HTTP响应处理时候需要的库
import { BusinessError } from '@kit.BasicServicesKit';

class HttpRequest {
  async request(requestBody: RequestBody) {
    const promise: Promise<ResponseResult> = new Promise(async (resolve: Function, reject: Function) => {

      console.log('创建HTTP请求，请求地址：' + Domain.SERVER + requestBody.url + '，请求方式：' + requestBody.method + '，请求体：' + JSON.stringify(requestBody.data))

      // 创建HTTP请求
      const httpRequest = http.createHttp()

      // 获取token
      const token = await StorageUtil.getToken();
      console.info('携带token进行请求：', token ? token.substring(0, 20) + '...' : '无token');

      // 构建请求头
      const headers: Record<string, string> = {
        'Content-Type': ContentType.JSON
      };

      // 如果有token，添加Authorization头 - 直接使用token值，不加Bearer前缀
      if (token) {
        headers['Authorization'] = token;
        console.info('设置Authorization头为直接token值');
      }

      // 在这里发送 HTTP 请求，并且需要设置Url、Header、Body等信息
      httpRequest.request(
        Domain.SERVER + requestBody.url,
        {
          method: getHttpMethod(requestBody.method),
          readTimeout: 10000,
          header: headers,
          connectTimeout: 10000,
          extraData: requestBody.data
        },
        (err: BusinessError, data: http.HttpResponse) => {
          // 回调之后判断是否提示 Error
          if (!err) {

            let result = `${data.result}`
            console.info('HTTP响应，状态码：' + data.responseCode + '，响应数据：' + result)

            try {
              let resultJSON: ResponseResult = JSON.parse(result)

              // 判断响应码是否成功
              if (data.responseCode === HttpStatus.SUCCESS) {
                // 统一返回响应结果，让调用方根据业务code进行判断
                console.info('HTTP请求成功，响应数据：' + result)
                resolve(resultJSON)
              } else {
                // 对于非200状态码，也尝试解析JSON响应
                console.error('HTTP请求失败，状态码：' + data.responseCode + '，响应：' + result)
                resolve(resultJSON) // 仍然返回解析的JSON，让调用方处理错误
              }
            } catch (parseError) {
              console.error('JSON解析失败：' + parseError + '，原始响应：' + result)
              // JSON解析失败时，构造一个标准的错误响应
              const errorResponse: ResponseResult = new ResponseResult()
              errorResponse.code = data.responseCode
              errorResponse.msg = '服务器响应格式错误'
              errorResponse.data = result
              resolve(errorResponse)
            }

            // 当该请求使用完毕时，调用destroy方法主动销毁
            httpRequest.destroy();
          } else {
            console.error('ERROR:' + JSON.stringify(err));

            if (err.code === 2300007) {
              errorDialog('无法连接至服务器，请稍后重试')
            }

            if (err.code === 2300028) {
              errorDialog('当前网络环境不稳定，请稍后重试')
            }

            // 当该请求使用完毕时，调用destroy方法主动销毁
            httpRequest.destroy();

            reject(new Error('当前网络环境不稳定，请稍后重试'))
          }
        }
      )
    })

    return promise;
  }
}

const httpRequest = new HttpRequest();
export default httpRequest as HttpRequest;

// 获取HTTP方法
function getHttpMethod(method: string): http.RequestMethod {
  switch (method.toUpperCase()) {
    case 'GET':
      return http.RequestMethod.GET;
    case 'POST':
      return http.RequestMethod.POST;
    case 'PUT':
      return http.RequestMethod.PUT;
    case 'DELETE':
      return http.RequestMethod.DELETE;
    default:
      return http.RequestMethod.GET;
  }
}

// 判断是否是GET方法
function isGet(method: string): boolean | undefined {

  if (method === http.RequestMethod.GET) {
    return true;
  }

  return false;
}

// 判断是否是POST方法
function isPost(method: string): boolean | undefined {

  if (method === http.RequestMethod.POST) {
    return true;
  }

  return false;
}

// 提示错误
function errorDialog(msg: string) {
  console.error(msg)
}
