import { requestClient } from '#/api/request';

export namespace AttachmentApi {
  /** 附件 */
  export interface Attachment {
    [key: string]: any;
    /** 附件ID */
    id: number;
    /** 文件名 */
    filename: string;
    /** 文件URL */
    url: string;
    /** 文件路径 */
    path: string;
    /** 文件大小（字节） */
    size: number;
    /** MIME类型 */
    mediatype: string;
    /** 所属目录 */
    parent: string;
    /** 文件扩展名 */
    extension: string;
    /** 创建时间 */
    created_at: number;
    /** 更新时间 */
    updated_at: number;
    /** 类型 */
    type?: string;
    /** 管理员ID */
    admin_id?: number;
    /** 用户ID */
    user_id?: number;
  }

  /** 目录信息 */
  export interface Directory {
    name: string;
    path: string;
    count: number;
  }

  /** 分页参数 */
  export interface ListParams {
    page?: number;
    page_size?: number;
    parent?: string;
    search?: string;
  }

  /** 分页响应 */
  export interface ListResponse {
    list: Attachment[];
    total: number;
    page: number;
    pageSize: number;
  }
}

/**
 * 获取目录列表
 */
async function getDirectories() {
  return requestClient.get<AttachmentApi.Directory[]>(
    '/system/attachment/directories',
  );
}

/**
 * 获取附件列表
 */
async function getAttachmentList(params?: AttachmentApi.ListParams) {
  return requestClient.get<AttachmentApi.ListResponse>(
    '/system/attachment/list',
    { params },
  );
}

/**
 * 上传附件
 */
async function uploadAttachment(file: File, parent?: string) {
  const formData = new FormData();
  formData.append('file', file);
  if (parent) {
    formData.append('parent', parent);
  }
  return requestClient.post<AttachmentApi.Attachment>(
    '/system/attachment/upload',
    formData,
    {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    },
  );
}

/**
 * 删除附件
 */
async function deleteAttachment(ids: number[]) {
  return requestClient.post('/system/attachment/delete', { ids });
}

export { deleteAttachment, getAttachmentList, getDirectories, uploadAttachment };

