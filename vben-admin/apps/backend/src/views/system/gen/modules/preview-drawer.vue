<script lang="ts" setup>
import { ref, watch } from 'vue';

import { Drawer, Tabs } from 'ant-design-vue';

import type { GenApi } from '#/api/system/gen';

interface Props {
  visible: boolean;
  previewData: GenApi.GeneratedCode | null;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  'update:visible': [value: boolean];
}>();

const activeKey = ref('controller');

const handleClose = () => {
  emit('update:visible', false);
};

watch(
  () => props.visible,
  (val) => {
    if (val) {
      activeKey.value = 'controller';
    }
  },
);
</script>

<template>
  <Drawer
    :open="visible"
    :width="800"
    title="代码预览"
    @close="handleClose"
  >
    <Tabs v-if="previewData" v-model:activeKey="activeKey">
      <!-- 后端代码 -->
      <Tabs.TabPane key="controller" tab="Controller">
        <pre class="code-preview">{{ previewData.backend.controller }}</pre>
      </Tabs.TabPane>

      <Tabs.TabPane key="service" tab="Service">
        <pre class="code-preview">{{ previewData.backend.service }}</pre>
      </Tabs.TabPane>

      <Tabs.TabPane key="model" tab="Model">
        <pre class="code-preview">{{ previewData.backend.model }}</pre>
      </Tabs.TabPane>

      <Tabs.TabPane key="validate" tab="Validate">
        <pre class="code-preview">{{ previewData.backend.validate }}</pre>
      </Tabs.TabPane>

      <Tabs.TabPane key="route" tab="Route">
        <pre class="code-preview">{{ previewData.backend.route }}</pre>
      </Tabs.TabPane>

      <!-- 前端代码 -->
      <Tabs.TabPane key="api" tab="API">
        <pre class="code-preview">{{ previewData.frontend.api }}</pre>
      </Tabs.TabPane>

      <Tabs.TabPane key="list" tab="List.vue">
        <pre class="code-preview">{{ previewData.frontend.list_view }}</pre>
      </Tabs.TabPane>

      <Tabs.TabPane key="data" tab="Data.ts">
        <pre class="code-preview">{{ previewData.frontend.data_ts }}</pre>
      </Tabs.TabPane>

      <Tabs.TabPane key="form" tab="Form.vue">
        <pre class="code-preview">{{ previewData.frontend.form_vue }}</pre>
      </Tabs.TabPane>

      <!-- 菜单 SQL -->
      <Tabs.TabPane key="menu" tab="Menu SQL">
        <pre class="code-preview">{{ previewData.menu.sql }}</pre>
      </Tabs.TabPane>
    </Tabs>
  </Drawer>
</template>

<style scoped>
.code-preview {
  background-color: #f5f5f5;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  padding: 16px;
  overflow-x: auto;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.5;
  max-height: 600px;
}
</style>

