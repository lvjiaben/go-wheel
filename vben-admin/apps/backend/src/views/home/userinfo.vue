<script lang="ts" setup>
import { computed, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { $t } from '@vben/locales';
import { useUserStore } from '@vben/stores';

import { Menu } from 'ant-design-vue';

import { breakpointsTailwind, useBreakpoints } from '@vueuse/core';

import AccountInfo from './userinfo/account-info.vue';
import OperationLog from './userinfo/operation-log.vue';
import PasswordChange from './userinfo/password-change.vue';

const userStore = useUserStore();
const breakpoints = useBreakpoints(breakpointsTailwind);

// 当前选中的菜单
const selectedKey = ref<string>('account');

// 菜单项
const menuItems = computed(() => [
  {
    icon: 'user',
    key: 'account',
    label: $t('userinfo.menu.account'),
  },
  {
    icon: 'lock',
    key: 'password',
    label: $t('userinfo.menu.password'),
  },
  {
    icon: 'file-text',
    key: 'log',
    label: $t('userinfo.menu.log'),
  },
]);

// 是否是PC端
const isDesktop = computed(() => breakpoints.greaterOrEqual('md').value);

// 菜单模式
const menuMode = computed(() => (isDesktop.value ? 'inline' : 'horizontal'));

// 处理菜单选择
const handleMenuSelect = (info: any) => {
  selectedKey.value = info.key as string;
};
</script>

<template>
  <Page auto-content-height>
    <div class="flex flex-col gap-4 md:h-full md:flex-row">
      <!-- 左侧菜单 (PC) / 顶部菜单 (Mobile) -->
      <div
        class="w-full flex-shrink-0 overflow-hidden rounded-lg border bg-white md:w-56 md:max-h-full"
      >
        <!-- 用户头像和信息 -->
        <div class="border-b p-4 text-center md:p-6">
          <div class="mx-auto mb-3 h-20 w-20 overflow-hidden rounded-full">
            <img
              v-if="userStore.userInfo?.avatar"
              :alt="userStore.userInfo.realName"
              class="h-full w-full object-cover"
              :src="userStore.userInfo.avatar"
            />
            <div
              v-else
              class="flex h-full w-full items-center justify-center bg-gray-200 text-2xl text-gray-500"
            >
              {{
                userStore.userInfo?.realName?.charAt(0).toUpperCase() || 'U'
              }}
            </div>
          </div>
          <div class="font-semibold">
            {{ userStore.userInfo?.realName || $t('common.user') }}
          </div>
        </div>

        <!-- 菜单 -->
        <Menu
          :selected-keys="[selectedKey]"
          class="border-0"
          :mode="menuMode"
          @select="handleMenuSelect"
        >
          <Menu.Item v-for="item in menuItems" :key="item.key">
            {{ item.label }}
          </Menu.Item>
        </Menu>
      </div>

      <!-- 右侧内容区域 -->
      <div class="flex-1 overflow-y-auto rounded-lg border bg-white p-6">
        <AccountInfo v-show="selectedKey === 'account'" />
        <PasswordChange v-show="selectedKey === 'password'" />
        <OperationLog v-show="selectedKey === 'log'" />
      </div>
    </div>
  </Page>
</template>

