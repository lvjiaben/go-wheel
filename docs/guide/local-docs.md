# 本地运行文档

如果你想在本地查看和编辑文档，可以按照以下步骤操作。

## 前置要求

- Node.js 18+
- pnpm

## 运行步骤

```bash
# 进入文档目录
cd docs

# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev
```

启动后访问终端显示的地址（通常是 `http://localhost:5173`）即可预览文档。

## 其他命令

```bash
# 构建生产版本
pnpm build

# 预览构建后的版本
pnpm preview
```

