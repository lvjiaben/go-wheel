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
    layout: "header-sidebar-nav",
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
  breadcrumb: {
    styleType: "background"
  },
  navigation: {
    styleType: "plain"
  },
  sidebar: {
    collapsedButton: false,
    fixedButton: false
  },
  tabbar: {
    styleType: "plain"
  },
  theme: {
    mode: "light",
    radius: "0"
  },
  transition: {
    enable: false,
    loading: false,
    name: "fade",
    progress: false
  },
  widget: {
    globalSearch: false,
    lockScreen: false,
    themeToggle: false
  }
});
