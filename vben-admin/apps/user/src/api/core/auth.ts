import { requestClient } from '#/api/request';

export namespace AuthApi {
  /** 登录接口参数 */
  export interface LoginParams {
    password?: string;
    username?: string;
    captcha?:any;
  }

  /** 登录接口返回值 */
  export interface LoginResult {
    accessToken: string;
  }

  export interface RefreshTokenResult {
    data: string;
    status: number;
  }
}

/**
 * 获取所有菜单（Vben动态路由使用）
 */
export async function getAllMenusApi() {
  return await requestClient.get<any>('/menu/all');
}

/**
 * 获取用户信息
 */
export async function getUserInfoApi() {
  return requestClient.get<any>('/user/info');
}
/**
 * 刷新accessToken
 */
export async function refreshTokenApi() {
  return <AuthApi.RefreshTokenResult>{
    status: 500,
  };
}
/**
 * 获取用户权限码
 */
export async function getAccessCodesApi() {
  return [];
}