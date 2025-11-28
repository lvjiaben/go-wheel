import { requestClient } from '#/api/request';

/**
 * 获取验证码
 */
export async function captchaApi() {
  return requestClient.post<any>('/common/captcha');
}
