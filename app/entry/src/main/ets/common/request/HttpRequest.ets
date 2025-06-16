// 这里需要引入第 2 步和第 3 步创建的请求和响应实体
import { ResponseResult } from '../response/ResponseResult'
import { RequestBody } from '../response/RequestBody'

// 这里需要引入第 1 步公共常量类
import { Domain, HttpStatus, ContentType } from '../constant/CommonConstant'

// 这里引用 HTTP 库
import http from '@ohos.net.http'

// HTTP响应处理时候需要的库
import { BusinessError } from '@kit.BasicServicesKit';

class HttpRequest {
  request(requestBody: RequestBody) {
    const promise: Promise<ResponseResult> = new Promise((resolve: Function, reject: Function) => {

      console.log('创建HTTP请求，请求地址：' + Domain.SERVER + requestBody.url + '，请求方式：' + requestBody.method + '，请求体：' + JSON.stringify(requestBody.data))

      // 创建HTTP请求
      const httpRequest = http.createHttp()

      // 在这里发送 HTTP 请求，并且需要设置Url、Header、Body等信息
      httpRequest.request(
        Domain.SERVER + requestBody.url,
        {
          method: isPost(requestBody.method) ? http.RequestMethod.POST : http.RequestMethod.GET,
          readTimeout: 10000,
          header: {
            'Content-Type': ContentType.JSON
          },
          connectTimeout: 10000,
          extraData: requestBody.data
        },
        (err: BusinessError, data: http.HttpResponse) => {
          // 回调之后判断是否提示 Error
          if (!err) {

            let result = `${data.result}`
            let resultJSON: ResponseResult = JSON.parse(result)

            // 判断响应码是否成功
            if (data.responseCode === HttpStatus.SUCCESS) {
              // 统一返回响应结果，让调用方根据业务code进行判断
              console.info('HTTP请求成功，响应数据：' + result)
              resolve(resultJSON)
            } else {
              reject(new Error('HTTP请求失败，状态码：' + data.responseCode))
              errorDialog('HTTP请求失败，状态码：' + data.responseCode)
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
