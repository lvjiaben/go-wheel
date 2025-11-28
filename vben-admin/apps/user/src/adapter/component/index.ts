/**
 * 通用组件共同的使用的基础组件，原先放在 adapter/form 内部，限制了使用范围，这里提取出来，方便其他地方使用
 * 可用于 vben-form、vben-modal、vben-drawer 等组件使用,
 */
import ArrayFormTable from '#/components/array-form-table/index.vue';
import AttachmentInput from '#/components/attachment-input/index.vue';
import ImageCaptcha from '#/components/form/image-captcha.vue';
import RichEditor from '#/components/rich-editor/index.vue';
import TableSelect from '#/components/table-select/index.vue';

import type { Component } from 'vue';

import type { BaseFormComponentType } from '@vben/common-ui';
import type { Recordable } from '@vben/types';

import { defineAsyncComponent, defineComponent, h, ref, watch } from 'vue';

import { ApiComponent, globalShareState, IconPicker } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { notification } from 'ant-design-vue';
import dayjs from 'dayjs';

const AutoComplete = defineAsyncComponent(
  () => import('ant-design-vue/es/auto-complete'),
);
const Button = defineAsyncComponent(() => import('ant-design-vue/es/button'));
const Checkbox = defineAsyncComponent(
  () => import('ant-design-vue/es/checkbox'),
);
const CheckboxGroup = defineAsyncComponent(() =>
  import('ant-design-vue/es/checkbox').then((res) => res.CheckboxGroup),
);
const DatePicker = defineAsyncComponent(
  () => import('ant-design-vue/es/date-picker'),
);
const Divider = defineAsyncComponent(() => import('ant-design-vue/es/divider'));
const Input = defineAsyncComponent(() => import('ant-design-vue/es/input'));
const InputNumber = defineAsyncComponent(
  () => import('ant-design-vue/es/input-number'),
);
const InputPassword = defineAsyncComponent(() =>
  import('ant-design-vue/es/input').then((res) => res.InputPassword),
);
const Mentions = defineAsyncComponent(
  () => import('ant-design-vue/es/mentions'),
);
const Radio = defineAsyncComponent(() => import('ant-design-vue/es/radio'));
const RadioGroup = defineAsyncComponent(() =>
  import('ant-design-vue/es/radio').then((res) => res.RadioGroup),
);
const RangePicker = defineAsyncComponent(() =>
  import('ant-design-vue/es/date-picker').then((res) => res.RangePicker),
);
const Rate = defineAsyncComponent(() => import('ant-design-vue/es/rate'));
const Select = defineAsyncComponent(() => import('ant-design-vue/es/select'));
const Space = defineAsyncComponent(() => import('ant-design-vue/es/space'));
const Switch = defineAsyncComponent(() => import('ant-design-vue/es/switch'));
const Textarea = defineAsyncComponent(() =>
  import('ant-design-vue/es/input').then((res) => res.Textarea),
);
const TimePicker = defineAsyncComponent(
  () => import('ant-design-vue/es/time-picker'),
);
const TreeSelect = defineAsyncComponent(
  () => import('ant-design-vue/es/tree-select'),
);
const Upload = defineAsyncComponent(() => import('ant-design-vue/es/upload'));

const withDefaultPlaceholder = <T extends Component>(
  component: T,
  type: 'input' | 'select',
  componentProps: Recordable<any> = {},
) => {
  return defineComponent({
    name: component.name,
    inheritAttrs: false,
    setup: (props: any, { attrs, expose, slots }) => {
      const placeholder =
        props?.placeholder ||
        attrs?.placeholder ||
        $t(`ui.placeholder.${type}`);
      // 透传组件暴露的方法
      const innerRef = ref();
      expose(
        new Proxy(
          {},
          {
            get: (_target, key) => innerRef.value?.[key],
            has: (_target, key) => key in (innerRef.value || {}),
          },
        ),
      );
      return () =>
        h(
          component,
          { ...componentProps, placeholder, ...props, ...attrs, ref: innerRef },
          slots,
        );
    },
  });
};

// 默认的日期范围快捷选项
const getDefaultRangePresets = () => [
  {
    label: $t('common.components.rangePicker.today'),
    value: [dayjs().startOf('day'), dayjs().endOf('day')],
  },
  {
    label: $t('common.components.rangePicker.yesterday'),
    value: [
      dayjs().subtract(1, 'day').startOf('day'),
      dayjs().subtract(1, 'day').endOf('day'),
    ],
  },
  {
    label: $t('common.components.rangePicker.last7Days'),
    value: [dayjs().subtract(7, 'days').startOf('day'), dayjs().endOf('day')],
  },
  {
    label: $t('common.components.rangePicker.last30Days'),
    value: [dayjs().subtract(30, 'days').startOf('day'), dayjs().endOf('day')],
  },
  {
    label: $t('common.components.rangePicker.thisWeek'),
    value: [dayjs().startOf('week'), dayjs().endOf('week')],
  },
  {
    label: $t('common.components.rangePicker.thisMonth'),
    value: [dayjs().startOf('month'), dayjs().endOf('month')],
  },
  {
    label: $t('common.components.rangePicker.lastMonth'),
    value: [
      dayjs().subtract(1, 'month').startOf('month'),
      dayjs().subtract(1, 'month').endOf('month'),
    ],
  },
];

// 这里需要自行根据业务组件库进行适配，需要用到的组件都需要在这里类型说明
export type ComponentType =
  | 'ApiSelect'
  | 'ApiTreeSelect'
  | 'ArrayFormTable'
  | 'AttachmentInput'
  | 'AutoComplete'
  | 'Checkbox'
  | 'CheckboxGroup'
  | 'DatePicker'
  | 'DefaultButton'
  | 'Divider'
  | 'IconPicker'
  | 'ImageCaptcha'
  | 'Input'
  | 'InputNumber'
  | 'InputPassword'
  | 'Mentions'
  | 'PrimaryButton'
  | 'Radio'
  | 'RadioGroup'
  | 'RangePicker'
  | 'Rate'
  | 'RichEditor'
  | 'Select'
  | 'Space'
  | 'Switch'
  | 'TableSelect'
  | 'Textarea'
  | 'TimePicker'
  | 'TreeSelect'
  | 'Upload'
  | BaseFormComponentType;

async function initComponentAdapter() {
  const components: Partial<Record<ComponentType, Component>> = {
    // 如果你的组件体积比较大，可以使用异步加载
    // Button: () =>
    // import('xxx').then((res) => res.Button),
    ApiSelect: withDefaultPlaceholder(
      {
        ...ApiComponent,
        name: 'ApiSelect',
      },
      'select',
      {
        component: Select,
        loadingSlot: 'suffixIcon',
        visibleEvent: 'onDropdownVisibleChange',
        modelPropName: 'value',
      },
    ),
    ApiTreeSelect: withDefaultPlaceholder(
      {
        ...ApiComponent,
        name: 'ApiTreeSelect',
      },
      'select',
      {
        component: TreeSelect,
        fieldNames: { label: 'label', value: 'value', children: 'children' },
        loadingSlot: 'suffixIcon',
        modelPropName: 'value',
        optionsPropName: 'treeData',
        visibleEvent: 'onVisibleChange',
      },
    ),
    ArrayFormTable,
    AttachmentInput,
    AutoComplete,
    Checkbox,
    CheckboxGroup,
    DatePicker,
    // 自定义默认按钮
    DefaultButton: (props, { attrs, slots }) => {
      return h(Button, { ...props, attrs, type: 'default' }, slots);
    },
    Divider,
    IconPicker: withDefaultPlaceholder(IconPicker, 'select', {
      iconSlot: 'addonAfter',
      inputComponent: Input,
      modelValueProp: 'value',
    }),
    Input: withDefaultPlaceholder(Input, 'input'),
    InputNumber: withDefaultPlaceholder(InputNumber, 'input'),
    InputPassword: withDefaultPlaceholder(InputPassword, 'input'),
    Mentions: withDefaultPlaceholder(Mentions, 'input'),
    // 自定义主要按钮
    PrimaryButton: (props, { attrs, slots }) => {
      return h(Button, { ...props, attrs, type: 'primary' }, slots);
    },
    Radio,
    RadioGroup,
    // RangePicker 添加默认的快捷选项和时间戳格式支持
    RangePicker: defineComponent({
      name: 'CustomRangePicker',
      inheritAttrs: false,
      setup(_, { attrs, slots }) {
        // 如果用户没有传入 presets，则使用默认的快捷选项
        const presets = attrs?.presets || getDefaultRangePresets();

        // 支持自定义的 timestampFormat 属性
        // timestampFormat: true - 转换为 10 位时间戳数组
        // timestampFormat: 'ms' - 转换为 13 位毫秒时间戳数组
        const timestampFormat = attrs?.timestampFormat;

        // 将时间戳转换为 dayjs 对象（用于显示）
        const convertTimestampToDayjs = (value: any) => {
          if (!timestampFormat || !value) return value;
          if (Array.isArray(value)) {
            return value.map((v: any) => {
              if (!v) return undefined;
              // 如果已经是 dayjs 对象，直接返回
              if (dayjs.isDayjs(v)) return v;
              // 如果是时间戳，转换为 dayjs 对象
              if (typeof v === 'number') {
                // timestampFormat === 'ms' 表示毫秒时间戳
                return timestampFormat === 'ms' ? dayjs(v) : dayjs.unix(v);
              }
              return v;
            });
          }
          return value;
        };

        // 创建内部值的 ref（用于显示的 dayjs 对象）
        const internalValue = ref(convertTimestampToDayjs(attrs.value));

        // 监听外部值变化（外部传入的是时间戳）
        watch(() => attrs.value, (newVal: any) => {
          // 将时间戳转换为 dayjs 对象用于显示
          internalValue.value = convertTimestampToDayjs(newVal);
        });

        // 处理值变化（用户选择日期时）
        const handleChange = (dates: any) => {
          // 更新内部显示值
          internalValue.value = dates;

          if (timestampFormat && dates && Array.isArray(dates)) {
            // 转换为时间戳数组传给外部
            const timestamps = dates.map((d: any) => {
              if (!d) return undefined;
              // timestampFormat === 'ms' 返回 13 位毫秒时间戳
              // 否则返回 10 位秒时间戳
              return timestampFormat === 'ms' ? d.valueOf() : d.unix();
            });
            // 触发 update:value 事件，传递时间戳
            const updateHandler = attrs['onUpdate:value'] as any;
            updateHandler?.(timestamps);
          } else {
            // 不转换，直接传递 dayjs 对象
            const updateHandler = attrs['onUpdate:value'] as any;
            updateHandler?.(dates);
          }
        };

        return () => h(
          RangePicker as any,
          {
            ...attrs,
            presets,
            value: internalValue.value,
            'onUpdate:value': handleChange,
          },
          slots,
        );
      },
    }),
    Rate,
    RichEditor,
    Select: withDefaultPlaceholder(Select, 'select'),
    Space,
    Switch,
    TableSelect,
    Textarea: withDefaultPlaceholder(Textarea, 'input'),
    TimePicker,
    TreeSelect: withDefaultPlaceholder(TreeSelect, 'select'),
    Upload,
    ImageCaptcha,
  };

  // 将组件注册到全局共享状态中
  globalShareState.setComponents(components);

  // 定义全局共享状态中的消息提示
  globalShareState.defineMessage({
    // 复制成功消息提示
    copyPreferencesSuccess: (title, content) => {
      notification.success({
        description: content,
        message: title,
        placement: 'bottomRight',
      });
    },
  });
}

export { initComponentAdapter };
