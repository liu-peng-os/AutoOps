# 阶段 1 Bug 修复记录 (2026-04-28)

## 修复内容

### 1. 后端：Go nil slice 序列化为 JSON null 导致前端崩溃

**问题**：
- 数据库 `cmdb_group` 表为空时，后端返回 `{"code":200,"data":null}`
- 前端收到 `null` 后直接调用 `.map()` / `.find()` / `.length` 导致运行时错误

**根本原因**：
- Go 的 `var tree []T` 初始化为 `nil`，序列化成 JSON 是 `null` 而不是 `[]`

**修复位置**：
- `api/api/cmdb/model/cmdbGroup.go:66` — `var tree []CmdbGroup` → `tree := make([]CmdbGroup, 0)`
- `api/api/cmdb/model/cmdbGroupHost.go:48` — `var tree []CmdbGroupHostDto` → `tree := make([]CmdbGroupHostDto, 0)`

**影响**：
- 空分组时后端现在返回 `{"code":200,"data":[]}`，前端不再崩溃

---

### 2. 前端：缺少空值防御

**问题**：
- 多个组件直接使用 `res.data` 而不检查是否为 `null`

**修复位置**：
- `web/src/views/cmdb/cmdbHost.vue:534` — `this.groupList = res.data || []`
- `web/src/views/cmdb/Host/CreateExcel.vue:136` — `this.groupList = res.data || []`
- `web/src/views/cmdb/Host/HostSsh.vue:292` — `(response.data.data || []).map(...)`
- `web/src/views/cmdb/Host/SSH.vue:134` — `(response.data.data || []).map(...)`

**影响**：
- 即使后端返回 `null`，前端也能正常处理

---

### 3. 前端：Vue 2 遗留代码导致运行时错误

**问题**：
- `this.$set` 在 Vue 3 中已移除
- `slot="footer"` 是 Vue 2 语法，Vue 3 应使用 `<template #footer>`

**修复位置**：
- `web/src/views/system/Admin.vue:538` — 删除 `this.$set`，直接赋值
- `web/src/views/system/Admin.vue:262,335` — `slot="footer"` → `<template #footer>`

**影响**：
- 对话框按钮正常渲染
- 状态更新不再报错

---

### 4. 前端：密码明文显示在 toast 通知中

**问题**：
- 重置密码成功后，toast 显示 `"修改成功，新密码是：123456"`

**修复位置**：
- `web/src/views/system/Admin.vue:671` — 移除密码明文，只显示 `"修改成功"`

**影响**：
- 安全性提升

---

### 5. 构建配置修复

**问题 1**：`go.mod` 里 `go 1.25.0` 但 Docker 镜像是 `golang:1.24-alpine`
- **修复**：`api/go.mod:3` — `go 1.25.0` → `go 1.24.0`

**问题 2**：`goproxy.cn` 不稳定导致构建失败
- **修复**：`docker/api/Dockerfile:4` — 添加备用代理 `https://goproxy.io`

---

## 验证方式

1. 启动服务：
   ```bash
   cd docker
   docker compose build devops-api devops-web
   docker compose up -d --force-recreate devops-api devops-web
   ```

2. 访问主机管理页面：`http://localhost:8088`，登录后进入 CMDB → 主机管理

3. 浏览器控制台不再出现以下错误：
   - `Cannot read properties of null (reading 'map')`
   - `Cannot read properties of null (reading 'find')`
   - `Cannot read properties of null (reading 'length')`

---

## 未完成的前端优化（已识别但未修复）

以下问题已通过 `/frontend-design` 识别，但未在本次修复：

### 性能问题
- 全量导入 Element Plus 图标（应按需引入）
- 无路由懒加载（应使用 `() => import(...)`）
- 全量导入 ECharts（应使用 `echarts/core`）
- 重复的 xterm 包（`xterm` 和 `@xterm/xterm`）

### 可访问性问题（Critical）
- 所有 focus outline 被移除，键盘导航不可用
- 登录表单缺少 `autocomplete` 属性
- 验证码图片无 `alt` 文本

### 代码质量问题
- `HeadImage` 组件被挂载两次
- 使用已废弃的 `/deep/` CSS 选择器（应使用 `:deep()`）
- 死依赖：`vue-template-compiler`（Vue 2 包）

**建议**：这些问题应在阶段 2 早期统一处理，避免影响当前阶段的功能开发。

---

## 相关文件

- 后端修复：`api/api/cmdb/model/cmdbGroup.go`, `api/api/cmdb/model/cmdbGroupHost.go`
- 前端修复：`web/src/views/cmdb/cmdbHost.vue`, `web/src/views/cmdb/Host/CreateExcel.vue`, `web/src/views/cmdb/Host/HostSsh.vue`, `web/src/views/cmdb/Host/SSH.vue`, `web/src/views/system/Admin.vue`
- 构建配置：`api/go.mod`, `docker/api/Dockerfile`
