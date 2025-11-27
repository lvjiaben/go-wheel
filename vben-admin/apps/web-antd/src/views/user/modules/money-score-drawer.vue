<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';
import type { UserApi } from '#/api/user';

import { computed, nextTick, reactive, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';

import { breakpointsTailwind, useBreakpoints } from '@vueuse/core';

import { useVbenForm } from '#/adapter/form';
import { updateUserMoney, updateUserScore } from '#/api/user';
import { $t } from '#/locales';

const emit = defineEmits<{
  success: [];
}>();

const formData = ref<UserApi.User & { type: 'money' | 'score' }>();

const schema = computed((): VbenFormSchema[] => {
  const isMoney = formData.value?.type === 'money';
  
  return [
    {
      component: 'Input',
      componentProps: {
        disabled: true,
      },
      fieldName: 'username',
      label: $t('user.username'),
    },
    {
      component: 'InputNumber',
      componentProps: {
        disabled: true,
        precision: 2,
      },
      fieldName: isMoney ? 'money' : 'score',
      label: isMoney ? $t('user.currentMoney') : $t('user.currentScore'),
    },
    {
      component: 'RadioGroup',
      componentProps: {
        options: [
          { label: $t('user.add'), value: 'add' },
          { label: $t('user.sub'), value: 'sub' },
        ],
      },
      defaultValue: 'add',
      fieldName: 'type',
      label: $t('user.operationType'),
    },
    {
      component: 'InputNumber',
      componentProps: {
        min: 0.01,
        placeholder: isMoney ? $t('user.moneyAmountPlaceholder') : $t('user.scoreAmountPlaceholder'),
        precision: 2,
        step: 0.01,
      },
      fieldName: 'amount',
      label: isMoney ? $t('user.moneyAmount') : $t('user.scoreAmount'),
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: $t('user.notePlaceholder'),
      },
      fieldName: 'note',
      label: $t('user.note'),
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: $t('user.sourcePlaceholder'),
      },
      fieldName: 'source',
      label: $t('user.source'),
    },
  ];
});

const breakpoints = useBreakpoints(breakpointsTailwind);
const isHorizontal = computed(() => breakpoints.greaterOrEqual('md').value);

const [Form, formApi] = useVbenForm(
  reactive({
    commonConfig: {
      colon: true,
    },
    schema,
    showDefaultActions: false,
  }),
);

const [Drawer, drawerApi] = useVbenDrawer({
  onConfirm: onSubmit,
  async onOpenChange(isOpen) {
    if (isOpen) {
      const data = drawerApi.getData<UserApi.User & { type: 'money' | 'score' }>();
      formApi.resetForm();
      formData.value = data;

      await nextTick();
      if (data) {
        formApi.setValues({
          username: data.username,
          money: data.money,
          score: data.score,
          type: 'add',
          amount: undefined,
          note: '',
          source: '',
        });
      }
    }
  },
});

async function onSubmit() {
  // 直接获取表单值，不做前端验证（后端会验证）
  drawerApi.lock();
  const values = await formApi.getValues<{
    type: 'add' | 'sub';
    amount: number;
    note: string;
    source: string;
  }>();

  try {
    const params: UserApi.UpdateMoneyScoreParams = {
      id: formData.value!.id,
      type: values.type,
      note: values.note,
      source: values.source,
    };

    if (formData.value?.type === 'money') {
      params.money = values.amount;
      await updateUserMoney(params);
    } else {
      params.score = values.amount;
      await updateUserScore(params);
    }

    drawerApi.close();
    emit('success');
  } finally {
    drawerApi.unlock();
  }
}

const getDrawerTitle = computed(() =>
  formData.value?.type === 'money'
    ? $t('user.updateMoney')
    : $t('user.updateScore'),
);
</script>
<template>
  <Drawer class="w-[500px]" :title="getDrawerTitle">
    <Form class="mx-4" :layout="isHorizontal ? 'horizontal' : 'vertical'" />
  </Drawer>
</template>
