import { requestClient } from '#/api/request';

export namespace UserApi {
  /** 用户 */
  export interface User {
    [key: string]: any;
    /** 用户ID */
    id: number;
    /** 上级用户ID */
    pid: number;
    /** 顶级用户ID */
    tid: number;
    /** 状态 */
    status: number;
    /** 状态文本 */
    status_text: string;
    /** 邀请码 */
    code: string;
    /** 版本号（乐观锁） */
    version: number;
    /** 头像 */
    avatar: string;
    /** 用户名 */
    username: string;
    /** 邮箱 */
    email: string;
    /** 手机号 */
    mobile: string;
    /** 积分 */
    score: number;
    /** 余额 */
    money: number;
    /** 创建时间 */
    created_at: number;
    /** 更新时间 */
    updated_at: number;
    /** 密码（仅保存时使用） */
    password?: string;
  }

  /** 列表请求参数 */
  export interface ListParams {
    page: number;
    page_size: number;
    search?: string;
    filter?: string;
    sort_by?: string;
    sort_order?: 'asc' | 'desc';
  }

  /** 列表响应 */
  export interface ListResponse {
    list: User[];
    total: number;
    page: number;
    limit: number;
  }

  /** 更新余额/积分参数 */
  export interface UpdateMoneyScoreParams {
    id: number;
    type: 'add' | 'sub';
    money?: number;
    score?: number;
    note: string;
    source: string;
  }

  /** 操作字段参数 */
  export interface OperateParams {
    ids?: number[];
    field: string;
    value: number;
  }

}

/**
 * 获取用户列表
 */
async function getUserList(params: UserApi.ListParams) {
  return requestClient.get<UserApi.ListResponse>('/user/list', { params });
}

/**
 * 创建用户
 */
async function createUser(
  data: Omit<UserApi.User, 'id' | 'created_at' | 'updated_at' | 'deleted_at' | 'version'>,
) {
  return requestClient.post('/user/create', data);
}

/**
 * 更新用户
 */
async function updateUser(
  data: Partial<UserApi.User> & { id: number },
) {
  return requestClient.post('/user/update', data);
}

/**
 * 删除用户
 */
async function deleteUser(data: any) {
  return requestClient.post('/user/delete', data);
}

/**
 * 更新用户余额
 */
async function updateUserMoney(data: UserApi.UpdateMoneyScoreParams) {
  return requestClient.post('/user/update-money', data);
}

/**
 * 更新用户积分
 */
async function updateUserScore(data: UserApi.UpdateMoneyScoreParams) {
  return requestClient.post('/user/update-score', data);
}

/**
 * 操作用户字段（status等开关字段）
 */
async function operateUser(data: UserApi.OperateParams) {
  return requestClient.post('/user/operate', data);
}

export {
  createUser,
  deleteUser,
  getUserList,
  operateUser,
  updateUser,
  updateUserMoney,
  updateUserScore,
};

