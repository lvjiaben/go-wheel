import { requestClient } from '#/api/request';

export namespace UserApi {
	/** 用户表 */
	export interface User {
		[key: string]: any;
		/** Id */
		id?: number;
		/** 上级 */
		pid?: number;
		/** 顶级 */
		tid?: number;
		/** 状态 */
		status?: number;
		/** 状态信息 */
		status_text?: string;
		/** 邀请码 */
		code?: string;
		/** UNIONID */
		wechat_unionid?: string;
		/** OPENID */
		wechat_openid?: string;
		/** 版本 */
		version?: number;
		/** 头像 */
		avatar?: string;
		/** 账号 */
		username: string;
		/** 密码 */
		password?: string;
		/** 密码盐 */
		salt: string;
		/** 邮箱 */
		email?: string;
		/** 手机号码 */
		mobile?: string;
		/** 积分 */
		score?: number;
		/** 余额 */
		money?: number;
		/** TOKEN */
		token?: string;
		/** 创建时间 */
		created_at: number;
		/** 更新时间 */
		updated_at?: number;
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

	/** 操作字段参数 */
	export interface OperateParams {
		ids?: number[];
		field: string;
		value: number;
	}

	/** 余额/积分调整参数 */
	export interface UpdateMoneyScoreParams {
		id?: number;
		type: 'add' | 'sub';
		money?: number;
		score?: number;
		note?: string;
		source?: string;
	}

}

/**
 * 获取用户表列表
 */
async function getUserList(params: UserApi.ListParams) {
	return requestClient.get<UserApi.ListResponse>('/user/list', { params });
}

/**
 * 创建用户表
 */
async function createUser(
	data: Omit<UserApi.User, 'id' | 'created_at' | 'updated_at'>,
) {
	return requestClient.post('/user/create', data);
}

/**
 * 更新用户表
 */
async function updateUser(
	data: Partial<UserApi.User> & { id: number },
) {
	return requestClient.post('/user/update', data);
}

/**
 * 删除用户表
 */
async function deleteUser(data: any) {
	return requestClient.post('/user/delete', data);
}

/**
 * 操作用户表字段
 */
async function operateUser(data: UserApi.OperateParams) {
	return requestClient.post('/user/operate', data);
}

/**
 * 调整用户余额
 */
async function updateUserMoney(data: UserApi.UpdateMoneyScoreParams) {
	return requestClient.post('/user/money', data);
}

/**
 * 调整用户积分
 */
async function updateUserScore(data: UserApi.UpdateMoneyScoreParams) {
	return requestClient.post('/user/score', data);
}

export {
	getUserList,
	createUser,
	updateUser,
	deleteUser,
	operateUser,
	updateUserMoney,
	updateUserScore,
};
