import type { Recordable, UserInfo } from '@vben/types';

import { ref } from 'vue';
import { useRouter } from 'vue-router';

import { LOGIN_PATH } from '@vben/constants';
import { preferences } from '@vben/preferences';
import { resetAllStores, useAccessStore, useUserStore } from '@vben/stores';

import { message, notification } from 'ant-design-vue';
import { defineStore } from 'pinia';

import {
  getAccessCodesApi,
  getUserInfoApi,
  userLogoutApi as logoutApi,
  mobileLoginApi,
  resetPwdApi,
  sendSmsApi,
  userLoginApi,
  userRegApi,
} from '#/api';
import { $t } from '#/locales';

export const useAuthStore = defineStore('auth', () => {
  const accessStore = useAccessStore();
  const userStore = useUserStore();
  const router = useRouter();

  const loginLoading = ref(false);

  /**
   * 异步处理登录操作
   * Asynchronously handle the login process
   * @param params 登录表单数据
   */
  async function authLogin(
    params: Recordable<any>,
    onSuccess?: () => Promise<void> | void,
  ) {
    // 异步处理用户登录操作并获取 accessToken
    let userInfo: null | UserInfo = null;
    try {
      loginLoading.value = true;
      const { accessToken } = await userLoginApi(params);

      // 如果成功获取到 accessToken
      if (accessToken) {
        accessStore.setAccessToken(accessToken);

        // 获取用户信息并存储到 accessStore 中
        const [fetchUserInfoResult, accessCodes] = await Promise.all([
          fetchUserInfo(),
          getAccessCodesApi(),
        ]);

        userInfo = fetchUserInfoResult;

        userStore.setUserInfo(userInfo);
        accessStore.setAccessCodes(accessCodes);

        if (accessStore.loginExpired) {
          accessStore.setLoginExpired(false);
        } else {
          onSuccess
            ? await onSuccess?.()
            : await router.push(
                userInfo?.homePath || preferences.app.defaultHomePath,
              );
        }

        if (userInfo?.realName) {
          notification.success({
            description: `${$t('authentication.loginSuccessDesc')}:${userInfo?.realName}`,
            duration: 3,
            message: $t('authentication.loginSuccess'),
          });
        }
      }
    } catch (error) {
      loginLoading.value = false;
    } finally {
      loginLoading.value = false;
    }

    return {
      userInfo,
    };
  }

  async function logout(redirect: boolean = true) {
    if(redirect){
      try {
        await logoutApi();
      } catch {
        // 不做任何处理
      }
    }
    resetAllStores();
    accessStore.setLoginExpired(false);

    // 回登录页带上当前路由地址
    await router.replace({
      path: LOGIN_PATH,
      query: redirect
        ? {
            redirect: encodeURIComponent(router.currentRoute.value.fullPath),
          }
        : {},
    });
  }

  async function fetchUserInfo() {
    let userInfo: null | UserInfo = null;
    userInfo = await getUserInfoApi();
    userStore.setUserInfo(userInfo);
    return userInfo;
  }

  /**
   * 手机号验证码登录
   * @param params 手机号和验证码
   */
  async function authMobileLogin(
    params: { mobile: string; code: string },
    onSuccess?: () => Promise<void> | void,
  ) {
    let userInfo: null | UserInfo = null;
    try {
      loginLoading.value = true;
      const { accessToken } = await mobileLoginApi(params);

      if (accessToken) {
        accessStore.setAccessToken(accessToken);

        const [fetchUserInfoResult, accessCodes] = await Promise.all([
          fetchUserInfo(),
          getAccessCodesApi(),
        ]);

        userInfo = fetchUserInfoResult;

        userStore.setUserInfo(userInfo);
        accessStore.setAccessCodes(accessCodes);

        if (accessStore.loginExpired) {
          accessStore.setLoginExpired(false);
        } else {
          onSuccess
            ? await onSuccess?.()
            : await router.push(
                userInfo?.homePath || preferences.app.defaultHomePath,
              );
        }

        if (userInfo?.realName) {
          notification.success({
            description: `${$t('authentication.loginSuccessDesc')}:${userInfo?.realName}`,
            duration: 3,
            message: $t('authentication.loginSuccess'),
          });
        }
      }
    } finally {
      loginLoading.value = false;
    }

    return { userInfo };
  }

  /**
   * 用户注册
   * @param params 注册参数
   */
  async function authRegister(params: {
    mobile: string;
    password: string;
    code: string;
    invite_code?: string;
  }) {
    try {
      loginLoading.value = true;
      const { accessToken } = await userRegApi(params);

      if (accessToken) {
        accessStore.setAccessToken(accessToken);

        const [fetchUserInfoResult, accessCodes] = await Promise.all([
          fetchUserInfo(),
          getAccessCodesApi(),
        ]);

        const userInfo = fetchUserInfoResult;

        userStore.setUserInfo(userInfo);
        accessStore.setAccessCodes(accessCodes);

        notification.success({
          description: $t('authentication.loginSuccessDesc'),
          duration: 3,
          message: $t('page.auth.registerSuccess'),
        });

        await router.push(userInfo?.homePath || preferences.app.defaultHomePath);
      }
    } finally {
      loginLoading.value = false;
    }
  }

  /**
   * 重置密码
   * @param params 重置密码参数
   */
  async function authResetPwd(params: {
    mobile: string;
    code: string;
    new_password: string;
  }) {
    try {
      loginLoading.value = true;
      await resetPwdApi(params);
      message.success($t('page.auth.resetPwdSuccess'));
      await router.push(LOGIN_PATH);
    } finally {
      loginLoading.value = false;
    }
  }

  /**
   * 发送短信验证码
   * @param mobile 手机号
   * @param event 事件类型
   */
  async function sendSmsCode(mobile: string, event: string) {
    await sendSmsApi({ mobile, event });
  }

  function $reset() {
    loginLoading.value = false;
  }

  return {
    $reset,
    authLogin,
    authMobileLogin,
    authRegister,
    authResetPwd,
    fetchUserInfo,
    loginLoading,
    logout,
    sendSmsCode,
  };
});
