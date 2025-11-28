import { requestClient } from '#/api/request';

export namespace AdminAdminApi {
  /** 管理员 */
  export interface Admin {
    [key: string]: any;
    /** 管理员ID */
    id: number;
    /** 上级管理员ID */
    pid: number;
    /** 用户名 */
    username: string;
    /** 真实姓名 */
    realname: string;
    /** 头像 */
    avatar: string;
    /** 邮箱 */
    email: string;
    /** 手机号 */
    mobile: string;
    /** 状态 */
    status: number;
    /** 最后登录时间 */
    last_login_time: number;
    /** 创建时间 */
    created_at: number;
    /** 更新时间 */
    updated_at: number;
    /** 角色ID列表 */
    role_ids?: number[];
    /** 密码（仅保存时使用） */
    password?: string;
  }
}

/**
 * 获取管理员数据列表
 */
async function getAdminList() {
  return requestClient.get<Array<AdminAdminApi.Admin>>('/admin/admin/list');
}

/**
 * 保存单个管理员
 * @param data 管理员数据
 */
async function saveAdmin(
  id: number,
  data: Omit<AdminAdminApi.Admin, 'id' | 'created_at' | 'updated_at' | 'last_login_time'>,
) {
  if (id > 0) {
    data.id = id;
  }
  return requestClient.post('/admin/admin/save', data);
}

/**
 * 删除管理员
 * @param id 管理员 ID
 */
async function deleteAdmin(id: number) {
  return requestClient.delete(`/admin/admin/delete/${id}`);
}

export { deleteAdmin, getAdminList, saveAdmin };

