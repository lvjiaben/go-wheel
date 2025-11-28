import { requestClient } from '#/api/request';

export namespace AdminRoleApi {
  /** 角色类型集合 */
  export const RoleTypes = [
    'normal', 
    'super',
  ] as const;

  /** 系统角色 */
  export interface AdminRole {
    [key: string]: any;
    /** 角色ID */
    id: number;
    /** 父级ID */
    pid: number;
    /** 角色名称 */
    name: string;
    /** 角色描述 */
    description: string;
    /** 是否超级管理员 */
    is_super: number;
    /** 状态 */
    status: number;
    /** 排序 */
    sort: number;
    /** 创建时间 */
    created_at: number;
    /** 更新时间 */
    updated_at: number;
    /** 菜单ID列表 */
    menu_ids?: number[];
  }

}

/**
 * 获取角色数据列表
 */
async function getRoleList() {
  return requestClient.get<Array<AdminRoleApi.AdminRole>>('/admin/role/list');
}

/**
 * 保存单个角色
 * @param data 角色数据
 */
async function saveRole(
  id: number,
  data: Omit<AdminRoleApi.AdminRole, 'id' | 'created_at' | 'updated_at'>,
) {
  if(id > 0){
    data.id = id;
  }
  return requestClient.post('/admin/role/save', data);
}

/**
 * 删除角色
 * @param id 角色 ID
 */
async function deleteRole(id: number) {
  return requestClient.delete(`/admin/role/delete/${id}`);
}

/**
 * 获取我的菜单权限
 */
async function getMyMenus() {
  return requestClient.get('/admin/role/my-menus');
}

export {
  saveRole,
  deleteRole,
  getRoleList,
  getMyMenus,
};
