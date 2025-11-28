import { requestClient } from '#/api/request';

export namespace ConfigApi {
  /** 配置项 */
  export interface Config {
    id: number;
    dir: string;
    key: string;
    name: string;
    tip: string;
    type: string;
    value: string;
    variable: string;
    created_at: number;
  }

  /** 配置分组 */
  export interface ConfigGroup {
    id: number;
    dir: string;
    children: Config[];
  }

  /** 列表响应 */
  export interface ListResponse {
    list: ConfigGroup[];
  }
}

/**
 * 获取配置列表
 */
export async function getConfigList() {
  return requestClient.get<ConfigApi.ListResponse>('/system/config/list');
}

/**
 * 创建配置
 */
export async function createConfig(data: Partial<ConfigApi.Config>) {
  return requestClient.post('/system/config/create', data);
}

/**
 * 批量更新配置
 */
export async function updateConfigs(data: Record<string, any>) {
  return requestClient.post('/system/config/update', data);
}

/**
 * 删除配置
 */
export async function deleteConfig(id: number) {
  return requestClient.delete(`/system/config/delete/${id}`);
}

