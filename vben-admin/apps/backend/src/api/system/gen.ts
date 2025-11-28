import { requestClient } from '#/api/request';

export namespace GenApi {
  /** 表信息 */
  export interface TableInfo {
    table_name: string;
    table_comment: string;
  }

  /** 字段配置 */
  export interface FieldConfig {
    column_name: string;
    column_type: string;
    column_comment: string;
    is_nullable: boolean;
    is_primary_key: boolean;
    is_auto_increment: boolean;
    field_name: string;
    field_type: string;
    json_tag: string;
    gorm_tag: string;
    in_create: boolean;
    in_update: boolean;
    is_required: boolean;
    validate_rules: string;
    show_in_table: boolean;
    table_sortable: boolean;
    table_searchable: boolean;
    table_search_type: string;
    table_display_type: string;
    show_in_form: boolean;
    form_component: string;
    form_component_props: Record<string, any>;
    search_form_type: string;
    is_operate_field: boolean;
    is_sort_field: boolean;
    is_time_field: boolean;
    is_relation_field: boolean;
    is_enum_field: boolean;
    is_set_field: boolean;
    is_text_field: boolean;
    is_bool_field: boolean;
    is_image_field: boolean;
    is_images_field: boolean;
  }

  /** 生成配置 */
  export interface GenConfig {
    table_name: string;
    table_comment: string;
    module_name: string;
    struct_name: string;
    package_name: string;
    frontend_src_path: string;
    methods: {
      list: boolean;
      create: boolean;
      update: boolean;
      delete: boolean;
      operate: boolean;
    };
    fields: FieldConfig[];
    search_fields: string[];
    operate_fields: string[];
    default_sort_field: string;
    default_sort_order: string;
    menu_config: {
      parent_menu_name: string;
      menu_name: string;
      menu_icon: string;
      menu_sort: number;
    };
  }

  /** 生成的代码 */
  export interface GeneratedCode {
    backend: {
      controller: string;
      service: string;
      model: string;
      validate: string;
      route: string;
    };
    frontend: {
      api: string;
      list_view: string;
      data_ts: string;
      form_vue: string;
    };
    menu: {
      sql: string;
    };
  }

  /** 生成历史 */
  export interface GenHistory {
    id: number;
    table_name: string;
    table_comment: string;
    module_name: string;
    struct_name: string;
    package_name: string;
    frontend_src_path: string;
    config: string;
    created_at: number;
    updated_at: number;
  }
}

/**
 * 获取数据库表列表
 */
async function getTableList(params?: { search?: string }) {
  return requestClient.get<GenApi.TableInfo[]>('/system/gen/table-list', {
    params,
  });
}

/**
 * 获取表详细信息
 */
async function getTableInfo(params: { table_name: string }) {
  return requestClient.get<any>('/system/gen/table-info', { params });
}

/**
 * 获取表的默认配置
 */
async function getTableConfig(params: { table_name: string }) {
  return requestClient.get<GenApi.GenConfig>('/system/gen/table-config', {
    params,
  });
}

/**
 * 预览生成的代码
 */
async function previewCode(data: { config: GenApi.GenConfig }) {
  return requestClient.post<GenApi.GeneratedCode>('/system/gen/preview', data);
}

/**
 * 生成代码
 */
async function generateCode(data: { config: GenApi.GenConfig }) {
  return requestClient.post('/system/gen/generate', data);
}

/**
 * 获取生成历史
 */
async function getHistory() {
  return requestClient.get<GenApi.GenHistory[]>('/system/gen/history');
}

/**
 * 删除生成的代码
 */
async function deleteGenerated(data: { id: number }) {
  return requestClient.post('/system/gen/delete', data);
}

/**
 * 下载生成的代码
 */
async function downloadCode(data: { config: GenApi.GenConfig }): Promise<Blob> {
  // 使用 POST 请求但需要特殊处理 blob 响应
  const response = await requestClient.post<any>('/system/gen/download', data, {
    responseType: 'blob',
    responseReturn: 'body', // 直接返回 body (Blob)
  });
  return response;
}

export {
  deleteGenerated,
  downloadCode,
  generateCode,
  getHistory,
  getTableConfig,
  getTableInfo,
  getTableList,
  previewCode,
};

