import { baseRequestClient, requestClient } from '#/api/request';

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
 * 登录
 */
export async function loginApi(data: AuthApi.LoginParams) {
  return requestClient.post<AuthApi.LoginResult>('/auth/login', data);
}

/**
 * 获取所有菜单（Vben动态路由使用）
 */
export async function getAllMenusApi() {
  return await requestClient.get<any>('/auth/menus');
}

/**
 * 获取用户信息
 */
export async function getUserInfoApi() {
  return requestClient.get<any>('/auth/info');
}

/**
 * 更新个人资料
 */
export async function updateProfileApi(data: { avatar: string; email: string }) {
  return requestClient.post('/auth/profile', data);
}

/**
 * 修改密码
 */
export async function changePasswordApi(data: {
  old_password: string;
  new_password: string;
}) {
  return requestClient.post('/auth/password', data);
}

/**
 * 获取操作日志
 */
export async function getOperationLogApi(params?: {
  page?: number;
  page_size?: number;
}) {
  return requestClient.get('/auth/log', { params });
}

/**
 * 刷新accessToken
 */
export async function refreshTokenApi() {
  return baseRequestClient.post<AuthApi.RefreshTokenResult>('/auth/refresh', {
    withCredentials: true,
  });
}

/**
 * 退出登录
 */
export async function logoutApi() {
  return requestClient.post('/auth/logout', {
    withCredentials: true,
  });
}

/**
 * 获取用户权限码
 */
export async function getAccessCodesApi() {
  return requestClient.get<string[]>('/auth/permission');
}