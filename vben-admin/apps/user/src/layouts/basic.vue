<script lang="ts" setup>

import { computed, watch } from 'vue';
import { useRouter } from 'vue-router';

import { AuthenticationLoginExpiredModal } from '@vben/common-ui';
import { useWatermark } from '@vben/hooks';
import { Settings } from '@vben/icons';
import {
  BasicLayout,
  LockScreen,
  UserDropdown,
} from '@vben/layouts';
import { preferences } from '@vben/preferences';
import { useAccessStore, useUserStore } from '@vben/stores';

import { $t } from '#/locales';
import { useAuthStore } from '#/store';
import LoginForm from '#/views/_core/authentication/login.vue';


const userStore = useUserStore();
const authStore = useAuthStore();
const accessStore = useAccessStore();
const router = useRouter();
const { destroyWatermark, updateWatermark } = useWatermark();

const menus = computed(() => [
  {
    handler: () => {
      router.push('/userinfo');
    },
    icon: Settings,
    text: $t('userinfo.title'),
  },
]);

const avatar = computed(() => {
  return userStore.userInfo?.avatar ?? preferences.app.defaultAvatar;
});

async function handleLogout() {
  await authStore.logout(false);
}
watch(
  () => preferences.app.watermark,
  async (enable) => {
    if (enable) {
      await updateWatermark({
        content: `${userStore.userInfo?.username} - ${userStore.userInfo?.realName}`,
      });
    } else {
      destroyWatermark();
    }
  },
  {
    immediate: true,
  },
);
</script>

<template>
  <BasicLayout @clear-preferences-and-logout="handleLogout">
    <template #user-dropdown>
      <UserDropdown
        :avatar
        :menus
        :text="userStore.userInfo?.realName"
        @logout="handleLogout"
      />
    </template>
    <template #extra>
      <AuthenticationLoginExpiredModal
        v-model:open="accessStore.loginExpired"
        :avatar
      >
        <LoginForm />
      </AuthenticationLoginExpiredModal>
    </template>
    <template #lock-screen>
      <LockScreen :avatar @to-login="handleLogout" />
    </template>
  </BasicLayout>
</template>
