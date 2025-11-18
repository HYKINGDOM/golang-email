# 前端服务说明（React / Vite）

## 框架与技术栈
- React 18 + TypeScript
- 路由：React Router
- 状态管理：React Query（服务端数据获取与缓存）
- HTTP 客户端：Axios（运行时读取 baseURL）
- UI 组件库：Ant Design
- 富文本：React-Quill
- 可视化：Chart.js + react-chartjs-2
- 构建工具：Vite
- 样式：Tailwind CSS

## 项目结构
```
frontend/
  public/              # 静态资源与运行时配置
    config.json        # 运行时配置（baseURL/ssePath/刷新间隔）
  src/
    pages/             # 页面组件（Dashboard/EmailConfigs/Tasks/Templates/Monitoring/DBTest）
    components/        # 通用组件（可扩展）
    services/          # API 客户端（Axios）
    hooks/             # 自定义 Hooks（可扩展）
    types/             # TS 类型（可扩展）
    config.ts          # 前端运行时配置加载
    App.tsx            # 路由与布局
```

## 路由配置
- `/` 仪表盘：统计图与概览
- `/configs` 邮箱配置：列表与创建弹窗
- `/tasks` 任务管理：创建任务与列表（可扩展）
- `/templates` 模板管理：富文本编辑、追踪开关
- `/monitor` 实时监控：SSE 订阅任务进度
- `/dbtest` 数据库测试：调用后端诊断接口显示结果

## 状态管理
- React Query 管理接口数据，统一 `queryKey` 与 `queryFn`，在仪表盘设置 `refetchInterval` 控制刷新间隔。

## UI 组件库
- Ant Design：表格、表单、弹窗、按钮等组件。遵循表单校验与交互规范。

## 构建工具与脚本
```bash
# 开发
npm run dev
# 构建
npm run build
# 预览
npm run preview
```

## 代码分割策略
- 按页面级进行懒加载（可选），结合路由动态导入（当前为直载，可按需优化）。
- 将第三方库按需引入，避免不必要的打包体积。

## 性能优化方案
- 使用 React Query 缓存与批量请求减少网络开销。
- Chart.js 在数据量大时启用抽样与降采样策略（可扩展）。
- Tailwind 原子化样式降低 CSS 负载，页面按需渲染。
- 对富文本编辑器与图表组件进行代码分割与懒加载（可选）。

## 测试方案
- 单元测试：建议使用 Jest + React Testing Library 测组件与逻辑。
- 端到端测试：建议使用 Cypress 测页面流程（SSE/表单/接口）。

## 多环境部署
- 运行时配置：`public/config.json` 提供 `baseURL/ssePath/chartRefreshMs`。
- 不同环境通过替换 `config.json` 即可，无需重新打包。

## 浏览器兼容性
- 现代浏览器（Chrome/Edge/Firefox/Safari 最新版）。如需 IE 兼容需额外 Polyfill。

## 版本日志
- v0.1.0：初始版本交付，完成页面与配置、监控与模板管理。
