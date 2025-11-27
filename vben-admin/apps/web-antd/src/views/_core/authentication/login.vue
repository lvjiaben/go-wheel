<script lang="ts" setup>
import type { VbenFormSchema } from '@vben/common-ui';

import { computed } from 'vue';

import { AuthenticationLogin, z } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { useAuthStore } from '#/store';
import { captchaApi } from '#/api/core/common';

defineOptions({ name: 'Login' });

const authStore = useAuthStore();

const formSchema = computed((): VbenFormSchema[] => {
  return [
    {
      component: 'VbenInput',
      componentProps: {
        placeholder: $t('authentication.usernameTip'),
      },
      fieldName: 'username',
      label: $t('authentication.username'),
      rules: z.string().min(1, { message: $t('authentication.usernameTip') }),
    },
    {
      component: 'VbenInputPassword',
      componentProps: {
        placeholder: $t('authentication.password'),
      },
      fieldName: 'password',
      label: $t('authentication.password'),
      rules: z.string().min(1, { message: $t('authentication.passwordTip') }),
    },
    {
      component: 'ImageCaptcha',
      fieldName: 'captcha',
      label: $t('page.auth.imageCaptcha'),
      componentProps: {
        captchaApi: captchaApi,
      },
      rules: z.object({
        id: z.string(),
        code: z.string(),
      }).refine(
        (data) => data.id.length > 0 && data.code.length > 0,
        {
          message: $t('page.auth.imageCaptchaTip'),
        }
      ),
    },
  ];
});
</script>

<template>
  <AuthenticationLogin
    :form-schema="formSchema"
    :loading="authStore.loginLoading"
    @submit="authStore.authLogin"
    :showCodeLogin="false"
    :showForgetPassword="false"
    :showQrcodeLogin="false"
    :showRegister="false"
    :showThirdPartyLogin="false"
  />
</template>
