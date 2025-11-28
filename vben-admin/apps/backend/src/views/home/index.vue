<script lang="ts" setup>
import type { AnalysisOverviewItem } from '@vben/common-ui';
import type { EchartsUIType } from '@vben/plugins/echarts';
import type { TabOption } from '@vben/types';
import type { HomeApi } from '#/api/core/home';
import {
  Card,
  Button,
  ButtonGroup,
} from 'ant-design-vue';

import { markRaw, onMounted, ref } from 'vue';

import { AnalysisOverview } from '@vben/common-ui';
import {
  SvgBellIcon,
  SvgCakeIcon,
  SvgCardIcon,
  SvgDownloadIcon,
} from '@vben/icons';
import { EchartsUI, useEcharts } from '@vben/plugins/echarts';

import { message } from 'ant-design-vue';

import { getHomeData } from '#/api/core/home';

// 时间范围选择
const timeRange = ref<7 | 30>(7);
const loading = ref(false);

// 统计数据
const homeData = ref<HomeApi.HomeData | null>(null);

// 统计卡片数据
const overviewItems = ref<AnalysisOverviewItem[]>([
  {
    icon: "",
    title: '用户余额',
    totalTitle: '用户积分',
    totalValue: 0,
    value: 0,
  },
  {
    icon: "",
    title: '用户数量',
    totalTitle: '今日注册',
    totalValue: 0,
    value: 0,
  },
  {
    icon: "",
    title: '附件数量',
    totalTitle: '今日上传',
    totalValue: 0,
    value: 0,
  },
  {
    icon: "",
    title: '管理员数量',
    totalTitle: '今日封禁用户',
    totalValue: 0,
    value: 0,
  },
]);

// 图表
const trendsChartRef = ref<EchartsUIType>();
const { renderEcharts: renderTrends } = useEcharts(trendsChartRef);

// 加载数据
const loadData = async () => {
  loading.value = true;
  try {
    const data = await getHomeData(timeRange.value);
    homeData.value = data;

    // 更新统计卡片
    if (overviewItems.value[0]) {
      overviewItems.value[0].value = data.user_money;
      overviewItems.value[0].totalValue = data.user_score;
    }
    if (overviewItems.value[1]) {
      overviewItems.value[1].value = data.user_count;
      overviewItems.value[1].totalValue = data.user_today;
    }
    if (overviewItems.value[2]) {
      overviewItems.value[2].value = data.upload_count;
      overviewItems.value[2].totalValue = data.upload_today_count;
    }
    if (overviewItems.value[3]) {
      overviewItems.value[3].value = data.admin_count;
      overviewItems.value[3].totalValue = data.user_status;
    }

    // 渲染图表
    renderTrendsChart(data);
  } catch (error) {
    message.error('加载数据失败');
  } finally {
    loading.value = false;
  }
};

// 渲染图表
const renderTrendsChart = (data: HomeApi.HomeData) => {
  renderTrends({
    grid: {
      bottom: 0,
      containLabel: true,
      left: '1%',
      right: '1%',
      top: '2%',
    },
    legend: {
      data: ['用户注册', '余额变动', '积分变动'],
      top: 10,
    },
    series: [
      {
        areaStyle: {},
        data: data.line_chart.yAxis[0],
        itemStyle: {
          color: '#5ab1ef',
        },
        name: '用户注册',
        smooth: true,
        type: 'line',
      },
      {
        areaStyle: {},
        data: data.line_chart.yAxis[1],
        itemStyle: {
          color: '#019680',
        },
        name: '余额变动',
        smooth: true,
        type: 'line',
      },
      {
        areaStyle: {},
        data: data.line_chart.yAxis[2],
        itemStyle: {
          color: '#f56c6c',
        },
        name: '积分变动',
        smooth: true,
        type: 'line',
      },
    ],
    tooltip: {
      axisPointer: {
        lineStyle: {
          color: '#019680',
          width: 1,
        },
      },
      trigger: 'axis',
    },
    xAxis: {
      axisTick: {
        show: false,
      },
      boundaryGap: false,
      data: data.line_chart.xAxis,
      splitLine: {
        lineStyle: {
          type: 'solid',
          width: 1,
        },
        show: true,
      },
      type: 'category',
    },
    yAxis: [
      {
        axisTick: {
          show: false,
        },
        splitArea: {
          show: true,
        },
        splitNumber: 4,
        type: 'value',
      },
    ],
  });
};

// 切换时间范围
const handleTimeChange = (value: string) => {
    console.log(value);
  const time = Number.parseInt(value) as 7 | 30;
  timeRange.value = time;
  loadData();
};

onMounted(() => {
  loadData();
});
</script>

<template>
  <div class="p-5">
    <AnalysisOverview :items="overviewItems" />
    <Card class="mt-5">
        <template #title>
            图表统计
        </template>
        <template #extra>
            <Button-group>
                <Button @click="handleTimeChange('7')">7天</Button>
                <Button @click="handleTimeChange('30')">30天</Button>
            </Button-group>
        </template>
        <div>
            <EchartsUI ref="trendsChartRef" :loading="loading" />
        </div>
    </Card>

  </div>
</template>

