<template>
  <component
    :is="inputComponent"
    :model-value="getDisplayValue()"
    size="large"
    :placeholder="$t('page.auth.imageCaptchaTip')"
    :maxlength="6"
    @update:model-value="handleInput"
    @input="handleInput"
    @change="handleInput"
    class="captcha-input"
  >
    <template #addonAfter>
      <div class="captcha-addon">
        <img 
          v-if="captchaImage" 
          :src="captchaImage" 
          class="captcha-image"
          :alt="$t('page.auth.imageCaptcha')"
          @click="refreshCaptcha"
        />
        <span v-else class="text-gray-400 pl-2">loading...</span>
      </div>
    </template>
  </component>
</template>

<style scoped>
/* 确保输入框组件的addon区域与输入框高度完全匹配 */
:deep(.ant-input-group-addon) {
  padding: 0 !important;
  background: transparent !important;
}

.captcha-addon {
  max-width: none !important;
  width: 130px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  margin: 0;
  border-left: 1px solid #d9d9d9;
  background: #fafafa;
}

.captcha-image {
  width: 120px;
  height: 32px;
  cursor: pointer;
  padding: 0 4px;
  object-fit: cover;
  display: block;
}
</style>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { $t } from '@vben/locales';

// 使用adapter中注册的VbenInput以保持样式一致
import { globalShareState } from '@vben/common-ui';

const components = globalShareState.getComponents();
const inputComponent = components.Input;

interface CaptchaValue {
  id: string;
  code: string;
}

interface Props {
  modelValue?: CaptchaValue | string;
  captchaApi?: () => Promise<any>;
}

const props = withDefaults(defineProps<Props>(), {
  captchaApi: undefined,
});

// 确保组件不会继承不需要的属性
defineOptions({
  inheritAttrs: false,
});

const emit = defineEmits<{
  'update:modelValue': [value: CaptchaValue];
}>();

const loading = ref(false);
const captchaImage = ref('');
const captchaId = ref('');

// 生成验证码
const generateCaptcha = async () => {
  if (!props.captchaApi) {
    console.error('captchaApi未提供');
    return;
  }
  
  loading.value = true;
   try {
     const response = await props.captchaApi();
     captchaImage.value = response.captcha_data;
     captchaId.value = response.captcha_id;
     // 验证码生成后，只有在已有输入值时才发出
     const currentCode = getDisplayValue();
     if (currentCode) {
       const initialValue: CaptchaValue = {
         id: captchaId.value,
         code: currentCode,
       };
       emit('update:modelValue', initialValue);
     }
     
  } catch (error) {
    loading.value = false;
    console.error('生成验证码失败:', error);
  } finally {
    loading.value = false;
  }
};

// 获取显示值 - 确保传递给输入组件的是字符串
const getDisplayValue = () => {
  if (typeof props.modelValue === 'string') {
    return props.modelValue;
  }
  if (props.modelValue && typeof props.modelValue === 'object') {
    return props.modelValue.code || '';
  }
  return '';
};

// 处理输入
const handleInput = (value: string | Event) => {
  // 处理不同类型的事件参数
  let inputValue = '';
  if (typeof value === 'string') {
    inputValue = value;
  } else if (value && typeof value === 'object' && 'target' in value) {
    inputValue = (value.target as HTMLInputElement).value;
  }
  const captchaValue: CaptchaValue = {
    id: captchaId.value || '',
    code: inputValue || '',
  };
  emit('update:modelValue', captchaValue);
};

// 刷新验证码
const refreshCaptcha = () => {
  generateCaptcha();
};

// 组件挂载时生成验证码
onMounted(() => {
  generateCaptcha();
});

// 暴露验证码ID给外部使用
defineExpose({
  getCaptchaId: () => captchaId.value,
  refresh: refreshCaptcha,
});
</script>