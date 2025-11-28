import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Go-Wheel',
  description: '企业级后台管理系统文档',
  lang: 'zh-CN',
  // 如果部署到 https://<username>.github.io/<repo>/ 需要设置 base
  base: '/go-wheel/',
  
  themeConfig: {
    logo: '/logo.svg',
    
    nav: [
      { text: '开始', link: '/guide/introduction' },
      { text: '前端', link: '/frontend/introduction' },
      { text: '后端', link: '/backend/structure' },
    ],
    
    sidebar: {
      '/guide/': [
        {
          text: '开始',
          items: [
            { text: '介绍', link: '/guide/introduction' },
            { text: '安装部署', link: '/guide/installation' },
            { text: '常见问题', link: '/guide/faq' },
            { text: '联系我们', link: '/guide/contact' },
          ]
        }
      ],
      '/frontend/': [
        {
          text: '前端',
          items: [
            { text: '介绍', link: '/frontend/introduction' },
            { text: '功能概览', link: '/frontend/features' },
          ]
        },
        {
          text: '自定义组件',
          items: [
            { text: 'ArrayFormTable 数组表单', link: '/frontend/components/array-form-table' },
            { text: 'AttachmentInput 附件输入', link: '/frontend/components/attachment-input' },
            { text: 'RichEditor 富文本编辑器', link: '/frontend/components/rich-editor' },
            { text: 'TableSelect 表格选择器', link: '/frontend/components/table-select' },
            { text: 'ImageCaptcha 图形验证码', link: '/frontend/components/image-captcha' },
          ]
        },
        {
          text: '开发模板',
          items: [
            { text: 'CRUD 模板', link: '/frontend/templates/crud' },
          ]
        }
      ],
      '/backend/': [
        {
          text: '基础',
          items: [
            { text: '目录结构', link: '/backend/structure' },
            { text: '路由说明', link: '/backend/routes' },
            { text: '日志系统', link: '/backend/logger' },
            { text: '多语言', link: '/backend/i18n' },
            { text: 'RBAC 权限', link: '/backend/rbac' },
          ]
        },
        {
          text: '应用模块',
          items: [
            { text: 'API 接口', link: '/backend/modules/api' },
            { text: 'Backend 后台', link: '/backend/modules/backend' },
            { text: 'Cron 定时任务', link: '/backend/modules/cron' },
            { text: 'Queue 消息队列', link: '/backend/modules/queue' },
            { text: 'Views 视图', link: '/backend/modules/views' },
            { text: 'WebSocket', link: '/backend/modules/websocket' },
          ]
        },
        {
          text: '通用模块',
          items: [
            { text: 'CRUDBuilder 构建器', link: '/backend/common/builder' },
            { text: 'Validator 验证器', link: '/backend/common/validator' },
            { text: 'Service 通用服务', link: '/backend/common/service' },
          ]
        },
        {
          text: '核心组件',
          items: [
            { text: 'Container 容器', link: '/backend/pkg/container' },
            { text: 'Config 配置', link: '/backend/pkg/config' },
            { text: 'Captcha 验证码', link: '/backend/pkg/captcha' },
            { text: 'JWT 认证', link: '/backend/pkg/jwt' },
            { text: 'HTTPClient 客户端', link: '/backend/pkg/httpclient' },
            { text: 'Middleware 中间件', link: '/backend/pkg/middleware' },
            { text: 'Utils 工具函数', link: '/backend/pkg/utils' },
            { text: 'Redis 缓存', link: '/backend/pkg/redis' },
            { text: 'HTTP 响应', link: '/backend/pkg/response' },
            { text: 'Constants 常量', link: '/backend/pkg/constants' },
          ]
        },
        {
          text: '开发模板',
          items: [
            { text: 'CRUD 模板', link: '/backend/templates/crud' },
          ]
        },
        {
          text: '部署运维',
          items: [
            { text: '环境安装', link: '/backend/deploy/install' },
            { text: 'Makefile 命令', link: '/backend/deploy/makefile' },
            { text: '代码生成', link: '/backend/deploy/codegen' },
            { text: '热更新', link: '/backend/deploy/hotreload' },
          ]
        }
      ]
    },
    
    socialLinks: [
      { icon: 'github', link: 'https://github.com' }
    ],
    
    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2024'
    },
    
    search: {
      provider: 'local'
    },
    
    outline: {
      level: [2, 3],
      label: '页面导航'
    },
    
    lastUpdated: {
      text: '最后更新于'
    },
    
    docFooter: {
      prev: '上一页',
      next: '下一页'
    }
  }
})

