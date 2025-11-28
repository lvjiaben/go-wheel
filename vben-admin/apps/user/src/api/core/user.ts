import type { Recordable } from '@vben/types';

import { requestClient } from '#/api/request';

export namespace UserApi {
  /** 登录请求参数 */
  export interface LoginParams {
    username: string;
    password: string;
    captcha: {
      id: string;
      code: string;
    };
  }

  /** 手机号登录请求参数 */
  export interface MobileLoginParams {
    mobile: string;
    code: string;
  }

  /** 注册请求参数 */
  export interface RegParams {
    mobile: string;
    password: string;
    code: string;
    invite_code?: string;
  }

  /** 重置密码请求参数 */
  export interface ResetPwdParams {
    mobile: string;
    code: string;
    new_password: string;
  }

  /** 发送短信验证码请求参数 */
  export interface SendSmsParams {
    mobile: string;
    event: string;
  }

  /** 登录返回结果 */
  export interface LoginResult {
    accessToken: string;
  }
}

/**
 * 账号密码登录
 */
export async function userLoginApi(data: Recordable<any>) {
  return requestClient.post<UserApi.LoginResult>('/user/login', data);
}

/**
 * 手机号登录
 */
export async function mobileLoginApi(data: UserApi.MobileLoginParams) {
  return requestClient.post<UserApi.LoginResult>('/user/mobile-login', data);
}

/**
 * 用户注册
 */
export async function userRegApi(data: UserApi.RegParams) {
  return requestClient.post<UserApi.LoginResult>('/user/reg', data);
}

/**
 * 重置密码
 */
export async function resetPwdApi(data: UserApi.ResetPwdParams) {
  return requestClient.post('/user/reset-pwd', data);
}

/**
 * 退出登录
 */
export async function userLogoutApi() {
  return requestClient.post('/user/logout', {
    withCredentials: true,
  });
}

/**
 * 发送短信验证码
 */
export async function sendSmsApi(data: UserApi.SendSmsParams) {
  return requestClient.post('/sms/send', data);
}

/**
 * 获取图形验证码
 */
export async function userCaptchaApi() {
  return requestClient.post<any>('/common/captcha');
}

/** 修改手机号请求参数 */
export interface ChangeMobileParams {
  mobile: string;
  new_mobile: string;
  code: string;
}

/** 修改密码请求参数 */
export interface ChangePwdParams {
  old_password: string;
  new_password: string;
}

/**
 * 修改手机号
 */
export async function changeMobileApi(data: ChangeMobileParams) {
  return requestClient.post('/user/change-mobile', data);
}

/**
 * 修改密码
 */
export async function changePwdApi(data: ChangePwdParams) {
  return requestClient.post('/user/change-pwd', data);
}

