# HairCut - 企业级连锁理发店数字化平台

## 项目概述

一套面向全国连锁理发店的**完整数字化管理平台**，类似「美团丽人 + 连锁门店管理 + 理发师工作台 + 企业ERP后台」的SaaS系统。

## 系统组成

| 子系统 | 技术栈 | 说明 |
|--------|--------|------|
| 用户端小程序 | Taro + React + TS | 面向顾客的微信小程序 |
| 理发师端小程序 | Taro + React + TS | 面向理发师的工具小程序 |
| 店长端小程序 | Taro + React + TS | 面向门店管理者的小程序 |
| 总部运营后台 | React18 + Vite + Shadcn UI | Web管理系统 |
| 官网宣传网站 | Next.js 14 (SSR) | 品牌展示与SEO |
| API服务端 | Golang + Gin + GORM | 统一业务后端 |

## 技术架构

### 后端技术栈
- **框架**: Gin (Golang)
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **缓存**: Redis 7
- **搜索**: Elasticsearch 8
- **文档存储**: MongoDB
- **消息队列**: RabbitMQ
- **实时通信**: WebSocket
- **文件存储**: MinIO
- **认证**: JWT + Casbin RBAC

### 前端技术栈
- **小程序**: Taro 3.x + React + TypeScript + TDesign
- **后台**: React 18 + Vite + Shadcn UI + TanStack Query + Zustand + TailwindCSS
- **官网**: Next.js 14 (App Router) + SSR + TailwindCSS

## 核心功能模块

### 用户端
- ✅ 附近门店查询（腾讯地图集成）
- ✅ 在线预约流程（选店→选理发师→选服务→选时间→支付）
- ✅ 订单管理与在线支付（微信支付）
- ✅ 5级会员体系（普通/银卡/金卡/黑金/钻石）
- ✅ 优惠券系统（新人券/满减券/折扣券等）
- ✅ 积分商城兑换
- ✅ AI发型推荐
- ✅ 用户评价（图片/视频/文字）

### 理发师端
- ✅ 工作台（预约/收入/客户概览）
- ✅ 预约管理（今日/未来/历史）
- ✅ 客户档案管理
- ✅ 排班与考勤
- ✅ 作品集管理
- ✅ 粉丝互动

### 店长端
- ✅ 门店数据总览（营业额/客流/订单/会员）
- ✅ 员工管理与排班
- ✅ 财务统计（收入/退款/提成）
- ✅ 库存管理
- ✅ 采购申请与审批

### 总部运营后台
- ✅ 组织架构管理（集团>区域>城市>门店）
- ✅ 全国门店管理
- ✅ 营销中心（优惠券/活动/拼团/秒杀）
- ✅ CRM系统（客户生命周期/流失预警）
- ✅ 数据分析大屏（GMV/复购率/留存率）

### 官网
- ✅ 品牌故事与文化展示
- ✅ 全国门店地图查询
- ✅ 明星理发师团队
- ✅ 新闻动态与招聘
- ✅ 加盟合作申请

## 快速开始

### 环境要求
- Node.js >= 18.x
- Go >= 1.21
- MySQL >= 8.0
- Redis >= 7.0
- Elasticsearch >= 8.0
- MongoDB >= 6.0
- RabbitMQ >= 3.12
- MinIO latest

### 一键启动（Docker Compose）
```bash
# 启动所有基础设施服务
docker-compose up -d mysql redis elasticsearch mongodb rabbitmq minio

# 初始化数据库
./scripts/init-db.sh

# 启动后端API
cd server && go run cmd/server/main.go

# 启动前端应用（选择一个）
cd apps/admin-web && npm run dev
# 或
cd apps/user-miniapp && npm run dev:weapp
```

### 单独启动各服务
详见各子目录下的 README.md

## 项目结构

```
haircut/
├── server/              # 后端API服务（Golang Gin）
├── apps/
│   ├── user-miniapp/    # 用户端小程序
│   ├── staff-miniapp/   # 理发师端小程序
│   ├── manager-miniapp/ # 店长端小程序
│   ├── admin-web/       # 运营后台Web
│   └── official-website/# 官网
├── docs/                # 设计文档
├── scripts/             # 运维脚本
└── docker/              # Docker配置
```

## 设计规范

**UI风格**: Apple HIG + Notion极简 + 美团丽人
- 主色调: 深玫瑰金 (#C8A882)
- 背景: 纯白 (#FFFFFF)
- 圆角卡片 + 微动效 + 暗黑模式支持

## 权限体系

RBAC多角色权限控制：
1. 超级管理员 - 全部权限
2. 总部管理员 - 运营管理
3. 区域经理 - 区域内数据
4. 店长 - 门店级别
5. 理发师 - 个人工作台
6. 客服 - 用户咨询
7. 会员用户 - 个人操作

## 部署方案

- 开发环境: Docker Compose
- 生产环境: Kubernetes + Helm Chart
- CI/CD: GitHub Actions / GitLab CI

## License

MIT License © 2024 HairCut Team
