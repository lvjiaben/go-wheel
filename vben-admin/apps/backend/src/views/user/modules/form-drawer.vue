<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { computed, nextTick, reactive, ref } from 'vue';

import { useVbenDrawer, useVbenForm } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { message } from 'ant-design-vue';
import { useBreakpoints } from '@vueuse/core';

import { createUser, updateUser } from '#/api/user';

import type { UserApi } from '#/api/user';

const emit = defineEmits<{
  success: [];
}>();

const rowData = ref<UserApi.User | null>(null);
const loading = ref(false);

let handleSubmit: () => Promise<void>;

const [Drawer, drawerApi] = useVbenDrawer({
  onConfirm: () => {
    handleSubmit();
  },
  async onOpenChange(isOpen: boolean) {
    if (isOpen) {
      // 获取数据
      const data = drawerApi.getData<UserApi.User>();
      rowData.value = data;

      const isEdit = !!data?.id;
      drawerApi.setState({
        title: isEdit
          ? $t('ui.actionTitle.edit', [$t('user.name')])
          : $t('ui.actionTitle.create', [$t('user.name')]),
      });

      // 设置表单值
      await nextTick();
      if (isEdit && data) {
        formApi.setValues({
          id: data.id,
          username: data.username,
          email: data.email,
          mobile: data.mobile,
          avatar: data.avatar,
          code: data.code,
          pid: data.pid,
          tid: data.tid,
          status: data.status,
          status_text: data.status_text,
          password: '', // 编辑时密码留空
        });
      } else {
        formApi.resetForm();
      }
    }
  },
});

const breakpoints = useBreakpoints({
  md: 768,
});
const isHorizontal = breakpoints.greater('md');

const formSchema = computed((): VbenFormSchema[] => {
  const isEdit = !!rowData.value?.id;
  return [
    {
      component: 'AttachmentInput',
      componentProps: {
        placeholder: $t('user.avatarPlaceholder'),
        showInput: false,
        showPreview: true,
      },
      fieldName: 'avatar',
      label: $t('user.avatar'),
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: $t('user.codePlaceholder'),
      },
      fieldName: 'code',
      label: $t('user.code'),
    },
    {
      component: 'InputNumber',
      componentProps: {
        placeholder: $t('user.pidPlaceholder'),
      },
      fieldName: 'pid',
      label: $t('user.pid'),
    },
    {
      component: 'InputNumber',
      componentProps: {
        placeholder: $t('user.tidPlaceholder'),
      },
      fieldName: 'tid',
      label: $t('user.tid'),
    },
    {
      component: 'RadioGroup',
      componentProps: {
        options: [
          { label: $t('common.enable'), value: 1 },
          { label: $t('common.disable'), value: 0 },
        ],
      },
      defaultValue: 1,
      fieldName: 'status',
      label: $t('user.status'),
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: $t('user.statusTextPlaceholder'),
      },
      fieldName: 'status_text',
      label: $t('user.statusText'),
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: $t('user.usernamePlaceholder'),
      },
      fieldName: 'username',
      label: $t('user.username'),
    },
    {
      component: 'InputPassword',
      componentProps: {
        placeholder: isEdit
          ? $t('user.passwordEditPlaceholder')
          : $t('user.passwordPlaceholder'),
      },
      fieldName: 'password',
      label: $t('user.password'),
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: $t('user.mobilePlaceholder'),
      },
      fieldName: 'mobile',
      label: $t('user.mobile'),
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: $t('user.emailPlaceholder'),
      },
      fieldName: 'email',
      label: $t('user.email'),
    },
  ];
});

const [Form, formApi] = useVbenForm(
  reactive({
    commonConfig: {
      componentProps: {
        class: 'w-full',
      },
    },
    schema: formSchema,
    showDefaultActions: false,
    wrapperClass: 'grid-cols-1',
  }),
);

handleSubmit = async () => {
  loading.value = true;
  try {
    const values = await formApi.getValues() as any;
    const isEdit = !!rowData.value?.id;

    if (isEdit) {
      values.id = rowData.value!.id;
      if (!values.password) {
        delete values.password;
      }
      await updateUser(values);
      message.success($t('common.success'));
    } else {
      await createUser(values);
      message.success($t('common.success'));
    }
    drawerApi.close();
    emit('success');
  } catch (error) {
    console.error('Submit error:', error);
  } finally {
    loading.value = false;
  }
};

defineExpose({
  drawerApi,
});
</script>

<template>
  <Drawer
    :confirm-loading="loading"
    :loading="loading"
    class="w-[600px]"
  >
    <Form class="mx-4" :layout="isHorizontal ? 'horizontal' : 'vertical'" />
  </Drawer>
</template>



