import { requestClient } from '#/api/request';

export namespace AdminMenuApi {
  /** 菜单类型集合 */
  export const MenuTypes = [
    'menu', 
    'button',
    'iframe',
    'link',
  ] as const;

  /** 系统菜单 */
  export interface AdminMenu {
    [key: string]: any;
    /** 子级 */
    children?: AdminMenu[];
    /** 组件 */
    component?: string;
    /** 菜单ID */
    id: number;
    /** 菜单名称 */
    name: string;
    /** 菜单名称 */
    enname: string;
    /** 路由路径 */
    path: string;
    /** 父级ID */
    pid: number;
    /** 菜单类型 */
    type: (typeof MenuTypes)[number];
    /** 权限标识 */
    permission?: string;
    /** 外部链接 */
    external: string;
    /** 固定在标签栏 */
    fixed_tag: number;
    /** 菜单图标 */
    icon: string;
    /** 内嵌Iframe的URL */
    iframe: string;
    /** 路由路径 */
    route: string;
    /** 显示标签 */
    show_tag: number;
    /** 是否可见 */
    visible: number;
    /** 排序 */
    sort: number;
  }
}

/**
 * 获取菜单数据列表
 */
async function getMenuList() {
  return requestClient.get<Array<AdminMenuApi.AdminMenu>>('/admin/menu/list');
}

/**
 * 保存单个菜单
 * @param data 菜单数据
 */
async function saveMenu(
  id: any,
  data: Omit<AdminMenuApi.AdminMenu, 'children' | 'id'>,
) {
  if(typeof id === 'number' && id>0){
    data.id = id;
  }
  return requestClient.post('/admin/menu/save', data);
}


/**
 * 删除菜单
 * @param id 菜单 ID
 */
async function deleteMenu(id: number) {
  return requestClient.post('/admin/menu/delete', { id });
}

export {
  saveMenu,
  deleteMenu,
  getMenuList,
};