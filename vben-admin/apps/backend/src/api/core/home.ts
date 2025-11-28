import { requestClient } from '#/api/request';

export namespace HomeApi {
  /** 首页统计数据 */
  export interface HomeData {
    /** 用户余额 */
    user_money: number;
    /** 用户积分 */
    user_score: number;
    /** 用户总数 */
    user_count: number;
    /** 今日新增用户 */
    user_today: number;
    /** 禁用用户数 */
    user_status: number;
    /** 管理员数量 */
    admin_count: number;
    /** 附件数量 */
    upload_count: number;
    /** 今日上传附件数量 */
    upload_today_count: number;
    /** 折线图数据 */
    line_chart: {
      xAxis: string[];
      yAxis: [number[], number[], number[]];
    };
  }
}

/**
 * 获取首页数据
 * @param time 时间范围：7天或30天
 */
async function getHomeData(time: 7 | 30 = 7) {
  return requestClient.get<HomeApi.HomeData>('/home/index', {
    params: { time },
  });
}

export { getHomeData };

