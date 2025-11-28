import { defineOverridesPreferences } from '@vben/preferences';

/**
 * @description 项目配置文件
 * 只需要覆盖项目中的一部分配置，不需要的配置不用覆盖，会自动使用默认配置
 * !!! 更改配置后请清空缓存，否则可能不生效
 */
export const overridesPreferences = defineOverridesPreferences({
  // overrides
  app: {
    name: import.meta.env.VITE_APP_TITLE,
    enableCheckUpdates: false,
    layout: "sidebar-mixed-nav",
    authPageLayout: "panel-center",
    enablePreferences: false,
    accessMode: "backend",
  },
  logo: {
    source: "/logo.png",
    fit: "contain",
    enable: true
  },
  shortcutKeys: {
    enable: false
  },
  sidebar: {
    collapsed: true,
    width: 200
  },
  tabbar: {
    styleType: "brisk"
  },
  theme: {
    builtinType: "violet",
    colorPrimary: "hsl(245 82% 67%)",
    mode: "light"
  },
  widget: {
    globalSearch: false,
    lockScreen: false,
    themeToggle: false,
  }
});
