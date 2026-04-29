# AutoOps 前端设计系统

## 颜色体系

所有颜色使用 OKLCH 色彩空间，中性色向 indigo 品牌色微微倾斜（chroma 0.005–0.01）。

### CSS 变量（定义于 `web/src/assets/css/global.css`）

| 变量 | 值 | 用途 |
|---|---|---|
| `--color-accent` | `oklch(58% 0.18 265)` | 主品牌色（indigo） |
| `--color-accent-hover` | `oklch(52% 0.18 265)` | 按钮 hover |
| `--color-accent-subtle` | `oklch(96% 0.02 265)` | 浅色背景强调 |
| `--color-accent-muted` | `oklch(90% 0.05 265)` | 边框强调 |
| `--color-sidebar-bg` | `oklch(22% 0.06 265)` | 侧边栏背景 |
| `--color-sidebar-deep` | `oklch(18% 0.06 265)` | 侧边栏渐变终点 |
| `--color-bg` | `oklch(96% 0.008 265)` | 页面背景 |
| `--color-surface` | `oklch(99% 0.004 265)` | 卡片/面板背景 |
| `--color-text-primary` | `oklch(22% 0.02 265)` | 主文字 |
| `--color-text-secondary` | `oklch(48% 0.02 265)` | 次要文字 |
| `--color-text-muted` | `oklch(62% 0.015 265)` | 辅助文字 |
| `--color-success` | `oklch(58% 0.16 155)` | 成功状态 |
| `--color-warning` | `oklch(72% 0.16 75)` | 警告状态 |
| `--color-danger` | `oklch(58% 0.20 25)` | 危险/错误状态 |

### 颜色策略

**Restrained**（克制）：浅色内容区 + 一个 accent 色。适合运维工具的高信息密度界面。

禁止：
- `#000` / `#fff` 纯黑纯白
- 蓝紫渐变背景（`linear-gradient(135deg, #667eea 0%, #764ba2 100%)`）
- 渐变文字（`background-clip: text`）

## 主题

**浅色主题**（内容区）+ **深色侧边栏**。

场景：运维工程师在公司内网工位，白天使用为主，需要快速扫描大量信息。浅色内容区减少视觉疲劳，深色侧边栏提供清晰的导航层次。

## 排版

- 字体栈：`-apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', system-ui, sans-serif`
- 正文行长：不超过 75ch
- 层级：通过字号 + 字重对比（≥1.25 倍）建立，不依赖颜色

## 间距

| 变量 | 值 | 用途 |
|---|---|---|
| `--radius-sm` | `4px` | 标签、小按钮 |
| `--radius-md` | `6px` | 输入框、普通按钮 |
| `--radius-lg` | `8px` | 卡片 |
| `--radius-xl` | `12px` | 弹窗、大卡片 |

页面根容器 padding：`12px 16px`（`el-main` 默认 padding 已归零）。

## 阴影

| 变量 | 用途 |
|---|---|
| `--shadow-xs` | 头部栏、细微分隔 |
| `--shadow-sm` | 卡片默认状态 |
| `--shadow-md` | 卡片 hover |
| `--shadow-lg` | 弹出层、下拉菜单 |

## 动画

- 时长：`--duration-fast: 150ms`，`--duration-normal: 200ms`
- 缓动：`--ease-out: cubic-bezier(0.16, 1, 0.3, 1)`（ease-out-quint）
- 禁止：动画 CSS 布局属性（`transform: translateX` 用于 hover 位移是布局动画，已移除）
- 禁止：bounce / elastic 缓动

## 组件规范

### 侧边栏菜单

- 背景：`linear-gradient(135deg, --color-sidebar-bg, --color-sidebar-deep)`
- 菜单项 hover：`--color-sidebar-hover`（白色 10% 透明度）
- 激活项：`linear-gradient(135deg, --color-sidebar-active-from, --color-sidebar-active-to)`
- focus 样式：`focus-visible` 显示 outline，`focus`（鼠标点击）不显示

### Tags 导航栏

- 原生 div 实现，不依赖 Element Plus el-tag
- 当前页：accent 色小圆点 + accent 色文字 + accent-muted 边框
- 关闭按钮：hover 时显示

### 登录卡片

- 实色深色卡片：`oklch(14% 0.04 265)`
- 无 glassmorphism（禁止 `backdrop-filter` 装饰性用法）
- 输入框深色风格，与卡片背景协调

### 统计卡片（Dashboard）

- 图标统一使用 `--color-accent`（accent-subtle 背景）
- 告警类图标使用 `--color-warning`（warning 浅色背景）
- 禁止：每张卡片用不同颜色（蓝/绿/橙/粉各自为政）

## 禁止模式（全局）

1. **渐变文字**：`background-clip: text` + gradient，已全局清除
2. **蓝紫页面背景**：`linear-gradient(135deg, #667eea 0%, #764ba2 100%)`，已全局清除
3. **Glassmorphism**：`backdrop-filter: blur` 装饰性用法
4. **侧边框强调**：`border-left > 1px` 作为彩色装饰
5. **相同卡片网格**：icon + 标题 + 文字无限重复
6. **布局属性动画**：hover 时 `translateX/Y` 位移
