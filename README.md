# 💇‍♂️ HairCut - 企业级连锁理发店数字化平台

<p align="center">
  <strong>一套面向全国连锁理发店的完整 SaaS 数字化管理平台</strong>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Golang-1.21+-00ADD8?style=flat&logo=go" alt="Go" />
  <img src="https://img.shields.io/badge/React-18-61DAFB?style=flat&logo=react" alt="React" />
  <img src="https://img.shields.io/badge/Next.js-14-000000?style=flat&logo=next.js" alt="Next.js" />
  <img src="https://img.shields.io/badge/MySQL-8.0-4479A1?style=flat&logo=mysql" alt="MySQL" />
  <img src="https://img.shields.io/badge/Redis-7-DC382D?style=flat&logo=redis" alt="Redis" />
  <img src="https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker" alt="Docker" />
  <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License" />
</p>

---

## 📖 目录

- [项目简介](#项目简介)
- [系统架构](#系统架构)
- [功能特性](#功能特性)
- [技术栈](#技术栈)
- [快速开始](#快速开始)
- [项目结构](#项目结构)
- [数据库设计](#数据库设计)
- [API 接口文档](#api-接口文档)
- [部署方案](#部署方案)
- [开发指南](#开发指南)
- [性能优化](#性能优化)
- [安全设计](#安全设计)
- [常见问题](#常见问题)

---

## 🎯 项目简介

HairCut 是一个**企业级连锁理发店数字化管理平台**，整合了「美团丽人 + 连锁门店管理 + 理发师工作台 + 企业ERP后台」的全链路能力。

### 核心价值

| 维度 | 解决方案 |
|------|---------|
| **顾客体验** | 在线预约、AI发型推荐、会员权益、积分商城 |
| **门店效率** | 智能排班、库存管理、财务自动化、数据分析 |
| **总部管控** | 多门店统一运营、营销中台、CRM客户管理、决策大屏 |
| **理发师赋能** | 个人工作台、作品展示、粉丝经济、收入透明 |

### 适用场景

- ✅ 单店独立经营 → 数字化升级
- ✅ 区域连锁 (5-50家店) → 统一管理
- ✅ 全国连锁 (100+店) → SaaS多租户部署
- ✅ 美业集团 → 定制化私有部署

---

## 🏗️ 系统架构

### 整体架构图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              用户接入层                                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌───────────┐  ┌───────────┐   │
│  │ 用户小程序 │  │ 理发师端  │  │ 店长端    │  │ 运营后台   │  │ 官网SSR   │   │
│  │ Taro+React│  │ Taro+React│  │ Taro+React│  │ React+Vite│  │ Next.js  │   │
│  └─────┬─────┘  └─────┬────┘  └─────┬─────┘  └─────┬─────┘  └─────┬─────┘   │
└────────┼──────────────┼──────────────┼──────────────┼───────────────┼─────────┘
         └──────────────┴──────────────┴──────────────┴───────────────┘
                                    │
                          ┌─────────▼─────────┐
                          │   API Gateway     │
                          │  (Nginx/LB)       │
                          └─────────┬─────────┘
                                    │
                    ────────────────┼─────────────── 分界线
                                    │
                          ┌─────────▼─────────┐
                          │   后端服务层       │
                          │  Golang Gin       │
                          │                   │
                          │  ┌─────────────┐  │
                          │  │  Router      │  │
                          │  │  Middleware  │  │
                          │  │  Handler     │  │
                          │  │  Service     │  │
                          │  │  Repository  │  │
                          │  └─────────────┘  │
                          └─────────┬─────────┘
                                    │
        ┌───────────┬───────────────┼───────────────┬───────────┬──────────┐
        ▼           ▼               ▼               ▼           ▼          ▼
   ┌─────────┐ ┌─────────┐  ┌─────────────┐  ┌─────────┐ ┌────────┐ ┌────────┐
   │ MySQL   │ │ Redis   │  │Elasticsearch│  │ MongoDB  │ │RabbitMQ│ │ MinIO  │
   │ 主数据库│ │ 缓存/锁  │  │ 全文搜索    │  │ 日志存储  │ │消息队列│ │对象存储│
   └─────────┘ └─────────┘  └─────────────┘  └─────────┘ └────────┘ └────────┘
```

### 六大子系统

| 子系统 | 技术栈 | 目标用户 | 核心职责 |
|--------|--------|----------|----------|
| **用户端小程序** | Taro + React + TS | 顾客 | 预约下单、会员服务、社交互动 |
| **理发师端小程序** | Taro + React + TS | 理发师 | 工作管理、客户维护、收入查看 |
| **店长端小程序** | Taro + React + TS | 店长 | 门店运营、员工调度、财务管理 |
| **总部运营后台** | React18 + Vite + Shadcn UI | 运营人员 | 全局管控、数据分析、营销配置 |
| **官网宣传网站** | Next.js 14 (SSR) | 公众用户 | 品牌展示、门店查询、加盟招商 |
| **API 服务端** | Golang + Gin + GORM | 系统内部 | 业务逻辑、数据存储、第三方对接 |

---

## ✨ 功能特性

### 🔥 用户端核心流程

```
首页浏览 → 选择门店(地图) → 选择理发师(作品) → 选择服务项目 
    → 选择时间段(实时库存) → 支付下单(微信支付/余额/积分抵扣)
    → 到店核销 → 服务完成 → 评价晒图 → 积分奖励
```

#### 功能清单

| 模块 | 功能点 | 说明 |
|------|--------|------|
| **🏠 首页** | 城市定位、Banner轮播、快捷入口、推荐门店、明星理发师 | 腾讯地图LBS定位 |
| **🏪 门店** | 列表/详情/附近搜索/地图导航/营业状态/评分评价 | ES全文搜索 |
| **💇 理发师** | 详情页/作品集/擅长标签/从业年限/粉丝数/预约按钮 | 作品MinIO存储 |
| **📅 预约** | 5步预约流程/可选时间段/并发锁/自动取消/提醒通知 | Redis分布式锁 |
| **💰 订单** | 待付款/待服务/进行中/已完成/已退款/售后 | RabbitMQ异步处理 |
| **👤 会员** | 5级体系(普通/银卡/金卡/黑金/钻石)/成长值/等级特权 | 自动化规则引擎 |
| **🎫 优惠券** | 新人券/满减券/折扣券/兑换券/叠加使用 | 6种优惠类型 |
| **🏪 积分商城** | 商品兑换/积分获取规则/积分流水 | 积分抵扣支付 |
| **🤖 AI发型** | 人脸识别/脸型分析/风格推荐/案例匹配 | AI算法集成 |
| **💬 消息** | 预约提醒/订单状态/营销推送/系统公告 | WebSocket实时 |

### 👨‍🎨 理发师端功能

| 模块 | 说明 |
|------|------|
| **工作台** | 今日预约数、本月收入、待处理事项、客户概览 |
| **预约管理** | 预约确认/拒绝/开始服务/完成服务 |
| **客户档案** | 消费记录、偏好标签、历史发型、备注信息 |
| **收入统计** | 日/周/月收入趋势、服务量统计、提成明细 |
| **排班管理** | 设置可预约时间、休息日、特殊安排 |
| **作品管理** | 上传作品照片/视频、添加描述和标签 |
| **粉丝管理** | 粉丝列表、互动消息、粉丝画像 |

### 🏢 店长端功能

| 模块 | 说明 |
|------|------|
| **数据总览** | GMV、客流量、订单量、会员增长、转化率 |
| **员工管理** | 理发师CRUD、角色权限、绩效评估 |
| **预约总览** | 全店预约日历、资源冲突检测 |
| **财务报表** | 营收明细、退款统计、成本核算、利润分析 |
| **库存管理** | 产品入库/出库、库存预警、损耗记录 |
| **采购管理** | 供应商管理、采购申请、审批流程 |

### 📊 总部运营后台

| 模块 | 说明 |
|------|------|
| **Dashboard** | 全国GMV、订单趋势、活跃用户、区域对比 |
| **组织架构** | 集团→区域→城市→门店四级树形结构 |
| **用户管理** | C端用户列表、行为轨迹、标签分组 |
| **门店管理** | 开店/闭店审核、资质管理、评分监控 |
| **理发师管理** | 跨门店调动、职级评定、黑名单 |
| **订单中心** | 全平台订单查询、退款审批、异常处理 |
| **营销中心** | 优惠券模板、满减活动、拼团、秒杀、CRM触达 |
| **CRM系统** | 客户生命周期、流失预警、RFM模型、自动化营销 |
| **财务中心** | 对账单、提现审核、发票管理、财务报表 |
| **数据分析** | 可视化大屏、自定义报表、数据导出 |

### 🌐 官网功能

| 模块 | 说明 |
|------|------|
| **品牌首页** | Hero视频、品牌故事、核心数据展示 |
| **关于我们** | 发展历程、企业文化、荣誉资质 |
| **门店查询** | 地图选店、城市筛选、距离计算、导航跳转 |
| **理发师团队** | 明星理发师展示、擅长领域、在线预约入口 |
| **新闻中心** | 行业资讯、品牌动态、活动预告 |
| **招聘中心** | 职位发布、在线投递、面试邀请 |
| **加盟系统** | 加盟政策、费用说明、在线申请、进度跟踪 |

---

## 🛠️ 技术栈

### 后端技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| Golang | 1.21+ | 主开发语言 |
| Gin | v1.9+ | Web框架 |
| GORM | v1.25+ | ORM框架 |
| MySQL | 8.0+ | 关系型数据库 |
| Redis | 7.0+ | 缓存/分布式锁/会话 |
| Elasticsearch | 8.x | 全文搜索引擎 |
| MongoDB | 6.0+ | 文档存储(日志) |
| RabbitMQ | 3.12+ | 异步消息队列 |
| MinIO | Latest | 对象存储(文件) |
| JWT | - | 身份认证令牌 |
| Casbin | v2.9+ | RBAC权限控制 |
| WebSocket | - | 实时双向通信 |

### 前端技术栈

| 应用 | 技术 | 用途 |
|------|------|------|
| **小程序三端** | Taro 3.x + React + TypeScript + NutUI/TDesign | 跨端统一开发 |
| **运营后台** | React 18 + Vite 5 + Shadcn UI + TanStack Query + Zustand + TailwindCSS | 高效管理后台 |
| **官网** | Next.js 14 (App Router) + SSR + TailwindCSS | SEO优化品牌站 |

### 第三方服务

| 服务 | 用途 |
|------|------|
| 微信开放平台 | 小程序登录/支付/模板消息 |
| 腾讯地图 | LBS定位/附近门店/距离计算/路线导航 |
| 微信支付 | JSAPI/H5/APP多端支付 |
| 短信服务 (阿里云/腾讯云) | 验证码/营销短信 |
| AI人脸识别 (百度/阿里) | 发型推荐脸型分析 |

---

## 🚀 快速开始

### 环境要求

| 工具 | 最低版本 | 推荐版本 |
|------|---------|---------|
| Node.js | >= 18.x | 20 LTS |
| Go | >= 1.21 | 1.22+ |
| MySQL | >= 8.0 | 8.0.36+ |
| Redis | >= 7.0 | 7.2+ |
| Elasticsearch | >= 8.x | 8.11+ |
| MongoDB | >= 6.0 | 7.0+ |
| RabbitMQ | >= 3.12 | 3.13+ |
| Docker | >= 24.0 | 24.0+ |
| Docker Compose | >= 2.20 | 2.24+ |

> 💡 **提示**: 如果不想手动安装所有中间件，可以使用 Docker Compose 一键启动全部基础设施。

### 方式一：Docker Compose 一键部署（推荐）

```bash
# 1. 克隆仓库
git clone https://github.com/hlw422/hair_cut.git
cd hair_cut

# 2. 复制环境变量配置
cp docker/.env.example .env

# 3. 编辑 .env 文件，填入真实密钥
# 必填项：
#   - WECHAT_APP_ID / WECHAT_APP_SECRET (微信小程序)
#   - WECHAT_MCH_ID / WECHAT_API_KEY (微信支付)
#   - TENCENT_MAP_KEY (腾讯地图)
vim .env

# 4. 启动基础设施（7个容器）
docker-compose up -d mysql redis elasticsearch mongodb rabbitmq minio nginx

# 等待服务就绪（约30秒）
docker-compose ps

# 5. 初始化数据库（建表 + 种子数据）
chmod +x scripts/init-db.sh
./scripts/init-db.sh

# 6. 启动后端API服务
cd server
go mod download
go run cmd/server/main.go

# 新开终端 - 启动前端应用（选择需要的一个）
cd apps/admin-web && npm install && npm run dev
# 或
cd apps/user-miniapp && npm install && npm run dev:weapp
# 或
cd apps/official-website && npm install && npm run dev
```

### 方式二：本地开发环境

#### 1. 安装中间件

```bash
# 使用 Homebrew (macOS)
brew install mysql@8.0 redis elasticsearch mongodb-community rabbitmq

# 或使用 apt (Ubuntu)
sudo apt-get install mysql-server redis-server elasticsearch mongodb rabbitmq-server
```

#### 2. 配置后端服务

```bash
cd server

# 安装依赖
go mod download

# 复制并编辑配置文件
cp config/config.yaml config/config.yaml.local
# 修改 database / redis / 其他连接信息

# 运行服务
go run cmd/server/main.go

# 或编译运行
go build -o bin/haircut-api cmd/server/main.go
./bin/haircut-api
```

#### 3. 启动前端应用

```bash
# 运营后台
cd apps/admin-web
npm install
npm run dev
# 访问: http://localhost:5173

# 用户端小程序
cd apps/user-miniapp
npm install
npm run dev:weapp
# 使用微信开发者工具导入 dist 目录

# 官网
cd apps/official-website
npm install
npm run dev
# 访问: http://localhost:3000
```

### 验证安装

```bash
# 检查后端健康状态
curl http://localhost:8080/health

# 预期响应:
# {
#   "code": 200,
#   "message": "success",
#   "data": {
#     "status": "UP",
#     "version": "1.0.0",
#     "uptime": 3600
#   }
# }
```

---

## 📁 项目结构

```
haircut/
│
├── 📂 server/                          # 【后端API服务】Golang Gin
│   ├── cmd/
│   │   └── server/main.go             # 服务入口 & 路由注册
│   ├── config/
│   │   └── config.yaml                # 配置文件 (YAML)
│   ├── internal/
│   │   ├── api/handler/               # HTTP处理器 (Controller层)
│   │   │   ├── user_handler.go        #   用户相关接口
│   │   │   ├── store_handler.go       #   门店相关接口
│   │   │   ├── appointment_handler.go #   预约相关接口
│   │   │   └── order_handler.go       #   订单支付接口
│   │   ├── middleware/                # 中间件
│   │   │   ├── auth.go               #   JWT认证中间件
│   │   │   ├── rbac.go               #   Casbin RBAC权限校验
│   │   │   ├── cors.go               #   跨域处理
│   │   │   ├── ratelimit.go          #   IP限流 (滑动窗口)
│   │   │   ├── request_id.go         #   请求追踪UUID
│   │   │   ├── recovery.go           #   Panic恢复
│   │   │   └── logger.go             #   请求日志
│   │   ├── model/mysql/              # GORM数据模型 (38张表)
│   │   │   ├── base_model.go         #   公共基类
│   │   │   ├── user.go               #   用户表
│   │   │   ├── member.go             #   会员+积分表
│   │   │   ├── store.go              #   门店+门店照片表
│   │   │   ├── stylist.go            #   理发师+作品集表
│   │   │   ├── service_item.go       #   服务项目+分类表
│   │   │   ├── appointment.go        #   预约+排班+考勤表
│   │   │   ├── order.go              #   订单+订单项+支付记录
│   │   │   ├── coupon.go             #   优惠券模板+用户券表
│   │   │   ├── review.go             #   评价+媒体附件表
│   │   │   ├── message.go            #   消息+通知表
│   │   │   ├── employee.go           #   员工+组织架构表
│   │   │   ├── inventory.go          #   库存+采购+供应商表
│   │   │   ├── campaign.go           #   营销活动表
│   │   │   ├── finance.go            #   财务+提成表
│   │   │   ├── role.go               #   角色+权限+用户角色表
│   │   │   ├── system_log.go         #   操作日志+系统配置表
│   │   │   └── register.go           #   AutoMigrate注册
│   │   ├── config/config.go          # 配置加载器
│   │   └── pkg/                     # 内部工具包
│   │       ├── jwt/jwt.go           #   JWT Token封装
│   │       ├── casbin/casbin_enforcer.go # Casbin权限引擎
│   │       ├── redis/redis_client.go    #   Redis操作封装
│   │       └── minio/minio_client.go    #   MinIO文件上传
│   └── pkg/                        # 公共工具包
│       ├── response/response.go    # 统一响应格式
│       ├── pagination/pagination.go # 分页工具
│       └── logger/logger.go        # Zap日志
│
├── 📂 apps/                           # 【前端应用】
│   ├── shared/                       # 共享代码库 (类型定义/工具函数/组件)
│   │   └── package.json
│   │
│   ├── user-miniapp/                 # 【用户端微信小程序】Taro+React
│   │   ├── config/                  #   Taro构建配置 (dev/prod)
│   │   ├── src/
│   │   │   ├── app.config.ts       #   小程序全局配置 (路由/TabBar)
│   │   │   ├── pages/              #   页面目录
│   │   │   │   ├── index/          #   首页 (Banner/推荐/明星)
│   │   │   │   ├── store/          #   门店列表/详情
│   │   │   │   ├── stylist/        #   理发师详情
│   │   │   │   ├── appointment/    #   预约流程
│   │   │   │   ├── order/          #   订单管理
│   │   │   │   ├── member/         #   会员中心
│   │   │   │   └── ...            #   更多页面
│   │   │   ├── components/         #   公共组件
│   │   │   ├── services/           #   API请求封装
│   │   │   ├── store/              #   状态管理
│   │   │   └── utils/              #   工具函数
│   │   └── package.json
│   │
│   ├── staff-miniapp/               # 【理发师端微信小程序】Taro+React
│   │   ├── src/pages/
│   │   │   ├── workspace/          #   工作台
│   │   │   ├── appointment/        #   预约管理
│   │   │   ├── customer/           #   客户档案
│   │   │   ├── income/             #   收入统计
│   │   │   ├── schedule/           #   排班管理
│   │   │   ├── portfolio/          #   作品管理
│   │   │   └── fan/               #   粉丝管理
│   │   └── package.json
│   │
│   ├── manager-miniapp/             # 【店长端微信小程序】Taro+React
│   │   ├── src/pages/
│   │   │   ├── dashboard/          #   门店总览
│   │   │   ├── employee/           #   员工管理
│   │   │   ├── appointment/        #   全店预约
│   │   │   ├── finance/            #   财务报表
│   │   │   ├── inventory/          #   库存管理
│   │   │   └── purchase/           #   采购管理
│   │   └── package.json
│   │
│   ├── admin-web/                  # 【总部运营后台】React18+Vite+Shadcn UI
│   │   ├── src/
│   │   │   ├── pages/
│   │   │   │   ├── login/LoginPage.tsx       # 登录页
│   │   │   │   ├── dashboard/DashboardPage.tsx # 数据仪表盘
│   │   │   │   ├── users/                   # 用户管理
│   │   │   │   ├── stores/                  # 门店管理
│   │   │   │   ├── stylists/               # 理发师管理
│   │   │   │   ├── orders/                 # 订单中心
│   │   │   │   ├── marketing/              # 营销中心
│   │   │   │   ├── crm/                    # CRM系统
│   │   │   │   ├── finance/                # 财务中心
│   │   │   │   ├── analytics/              # 数据分析
│   │   │   │   └── system/                 # 系统设置
│   │   │   ├── components/        #   公共组件 (Layout/Table/Form...)
│   │   │   ├── hooks/             #   自定义Hooks
│   │   │   ├── store/             #   Zustand状态
│   │   │   ├── services/          #   API请求
│   │   │   └── types/             #   TypeScript类型
│   │   ├── vite.config.ts         # Vite构建配置
│   │   ├── tailwind.config.ts     # Tailwind主题 (品牌色系)
│   │   └── package.json
│   │
│   └── official-website/          # 【品牌官网】Next.js 14 SSR
│       ├── src/
│       │   └── app/
│       │       ├── page.tsx       #   首页 (SEO优化)
│       │       ├── about/         #   关于我们
│       │       ├── stores/        #   门店查询
│       │       ├── stylists/      #   理发师团队
│       │       ├── news/          #   新闻中心
│       │       ├── careers/       #   招聘中心
│       │       └── join/          #   加盟申请
│       ├── next.config.js         # Next.js配置 (图片优化/i18n)
│       └── package.json
│
├── 📂 docs/                          # 【设计文档】
│   ├── database-design.md           # 数据库ER图 + 38张表DDL
│   └── full-technical-design.md    # ⭐ 完整技术设计 (13章节)
│
├── 📂 scripts/                       # 【运维脚本】
│   ├── init-db.sh                   # 数据库初始化脚本
│   ├── deploy.sh                    # 一键部署脚本
│   └── backup.sh                    # 数据备份脚本 (MySQL+Redis+MinIO)
│
├── 📂 docker/                        # 【Docker配置】
│   ├── Dockerfile.server            # 后端服务镜像 (多阶段构建)
│   ├── Dockerfile.admin-web         # 前端Nginx托管镜像
│   ├── nginx.conf                   # 反向代理配置 (Gzip/CORS)
│   └── .env.example                 # Docker环境变量模板
│
├── .gitignore                       # Git忽略规则
├── docker-compose.yml               # 7个基础设施容器编排
├── README.md                        # 📖 项目说明文档 (本文件)
└── LICENSE                          # 开源协议 (MIT)
```

---

## 🗄️ 数据库设计

### ER 图概览

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           核心业务域                                     │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   ┌──────────┐     ┌──────────┐     ┌───────────┐     ┌───────────┐    │
│   │   User   │────<│  Member  │────>│MemberLevel │     │ PointLog  │    │
│   │  用户表  │     │  会员表  │     │  等级表    │     │ 积分记录  │    │
│   └────┬─────┘     └──────────┘     └───────────┘     └───────────┘    │
│        │                                                                │
│        │ 1:N                                                            │
│        ├──────────────────────────────────────────────────────────┐     │
│        ▼                                                          ▼     │
│   ┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────────┐  │
│   │ Coupon   │     │  Order   │     │Appointment│    │    Review     │  │
│   │ 优惠券   │     │  订单    │     │  预约     │     │    评价      │  │
│   └──────────┘     └────┬─────┘     └─────┬────┘     └──────────────┘  │
│                         │                 │                             │
│                         │ N:1             │ N:1                         │
│                         ▼                 ▼                             │
│                    ┌──────────┐     ┌──────────┐     ┌──────────────┐   │
│                    │OrderItem │     │ Stylist  │     │StorePhoto    │   │
│                    │ 订单项   │     │ 理发师   │     │ 门店照片     │   │
│                    └──────────┘     └────┬─────┘     └──────┬───────┘   │
│                                       │ N:1               │ N:1        │
│                                       ▼                   ▼            │
│                                  ┌──────────┐       ┌──────────┐       │
│                                  │Portfolio │       │  Store   │       │
│                                  │  作品集  │       │  门店    │       │
│                                  └──────────┘       └──────────┘       │
│                                                                         │
├─────────────────────────────────────────────────────────────────────────┤
│                          管理 & 扩展域                                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   ┌──────────┐     ┌──────────┐     ┌──────────┐     ┌───────────┐    │
│   │ Employee │     │Inventory │     │ Campaign │     │  Finance   │    │
│   │  员工    │     │  库存    │     │ 营销活动 │     │  财务     │    │
│   └────┬─────┘     └────┬─────┘     └──────────┘     └───────────┘    │
│        │               │                                                   │
│        │ N:1           │ N:1                                              │
│        ▼               ▼                                                  │
│   ┌──────────┐   ┌──────────┐   ┌──────────┐                              │
│   │Organization│  │Supplier  │   │Purchase  │                              │
│   │ 组织架构  │  │ 供应商   │   │  采购单   │                              │
│   └──────────┘   └──────────┘   └──────────┘                              │
│                                                                         │
├─────────────────────────────────────────────────────────────────────────┤
│                            权限 & 系统域                                 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   ┌──────────┐     ┌──────────┐     ┌──────────┐     ┌───────────┐    │
│   │   Role   │     │Permission│     │UserRole │     │ SystemLog │    │
│   │  角色    │     │  权限    │     │ 用户角色 │     │ 操作日志  │    │
│   └──────────┘     └──────────┘     └──────────┘     └───────────┘    │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 数据表清单（38张）

| 分类 | 表名 | 说明 | 关键字段 |
|------|------|------|----------|
| **用户域** | `users` | 用户基础信息 | openid/nickname/avatar/phone |
| | `members` | 会员信息 | level/growth_points/balance |
| | `member_level_rules` | 等级规则 | level_name/min_points/privileges |
| | `point_logs` | 积分流水 | type/amount/balance/source |
| **门店域** | `stores` | 门店信息 | name/address/lat/lng/status |
| | `store_photos` | 门店照片 | store_id/url/type/sort |
| **理发师域** | `stylists` | 理发师信息 | store_id/name/title/rating |
| | `stylist_portfolios` | 作品集 | stylist_id/media_url/description |
| **服务域** | `service_categories` | 服务分类 | name/icon/sort_order |
| | `service_items` | 服务项目 | category_id/name/price/duration |
| **预约域** | `appointments` | 预约记录 | user_id/stylist_id/time/status |
| | `stylist_schedules` | 排班表 | stylist_id/date/time_slots |
| | `stylist_attendances` | 考勤记录 | stylist_id/date/check_in/out |
| **订单域** | `orders` | 订单主表 | user_id/store_id/total_amount/status |
| | `order_items` | 订单子项 | order_id/service_id/price/quantity |
| | `payment_records` | 支付记录 | order_id/method/amount/transaction_id |
| **优惠券域** | `coupon_templates` | 优惠券模板 | type/discount/min_spend/stock |
| | `user_coupons` | 用户优惠券 | template_id/user_id/status/expired_at |
| **评价域** | `reviews` | 评价记录 | order_id/rating/content |
| | `review_medias` | 评价附件 | review_id/url/type |
| **消息域** | `messages` | 站内消息 | sender/receiver/type/content |
| | `notifications` | 系统通知 | user_id/title/content/is_read |
| **关系域** | `fan_relations` | 粉丝关系 | follower/following |
| | `customer_profiles` | 客户档案 | stylist_id/user_id/tags/notes |
| **组织域** | `organizations` | 组织架构 | parent_id/name/type/level |
| | `employees` | 员工信息 | org_id/user_id/position/status |
| **库存域** | `inventory_items` | 库存商品 | name/stock/alert_quantity/unit_price |
| | `suppliers` | 供应商 | name/contact/phone/rating |
| | `purchase_orders` | 采购单 | supplier_id/total/status |
| | `purchase_items` | 采购子项 | order_id/item_id/quantity/price |
| **营销域** | `campaigns` | 营销活动 | type/start/end/rules |
| **财务域** | `financial_records` | 财务记录 | type/amount/reference_id |
| | `stylist_commissions` | 提成记录 | stylist_id/order_id/rate/amount |
| **权限域** | `roles` | 角色定义 | name/description/is_system |
| | `permissions` | 权限定义 | resource/action/module |
| | `user_roles` | 用户角色关联 | user_id/role_id |
| **系统域** | `operation_logs` | 操作日志 | user_id/action/ip/detail |
| | `system_configs` | 系统配置 | key/value/group |

### 核心状态机

**预约状态流转**:
```
待确认(PENDING) → 已确认(CONFIRMED) → 进行中(IN_PROGRESS) → 已完成(COMPLETED)
       │                │
       ▼                ▼
   已取消(CANCELLED)   已拒绝(REJECTED)
```

**订单状态流转**:
```
待支付(PENDING) → 已支付(PAID) → 待服务(PENDING_SERVICE) → 进行中(SERVICE_IN_PROGRESS) → 已完成(COMPLETED)
     │               │              │                      │
     ▼               ▼              ▼                      ▼
  已关闭(CLOSED)  已退款(REFUNDED)  已取消(CANCELLED)      售后中(AFTER_SALE)
```

---

## 📡 API 接口文档

### 接口规范

| 规范项 | 说明 |
|--------|------|
| **Base URL** | `http://localhost:8080/api/v1` |
| **认证方式** | `Authorization: Bearer <jwt_token>` |
| **Content-Type** | `application/json` |
| **分页参数** | `page=1&page_size=20` |
| **时间格式** | RFC3339: `2024-01-15T10:30:00+08:00` |

### 统一响应格式

```json
// 成功响应
{
  "code": 200,
  "message": "success",
  "data": { ... }
}

// 错误响应
{
  "code": 40001,
  "message": "参数验证失败",
  "data": null,
  "errors": [
    { "field": "phone", "message": "手机号格式不正确" }
  ]
}

// 分页响应
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [ ... ],
    "pagination": {
      "total": 100,
      "page": 1,
      "page_size": 20,
      "total_pages": 5
    }
  }
}
```

### 核心接口列表

#### 用户模块 (`/api/v1/users`)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/users/me` | 获取当前用户信息 | 是 |
| PUT | `/users/me` | 更新个人信息 | 是 |
| GET | `/users/me/member` | 获取会员信息 | 是 |
| POST | `/users/login/wechat` | 微信登录 | 否 |
| POST | `/users/logout` | 退出登录 | 是 |

#### 门店模块 (`/api/v1/stores`)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/stores/nearby` | 附近门店列表 | 否 |
| GET | `/stores/:id` | 门店详情 | 否 |
| GET | `/stores/search` | 搜索门店 | 否 |
| GET | `/stores/:id/stylists` | 门店理发师列表 | 否 |
| GET | `/stores/:id/reviews` | 门店评价列表 | 否 |

#### 预约模块 (`/api/v1/appointments`)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/appointments` | 创建预约 | 是 |
| GET | `/appointments/my` | 我的预约列表 | 是 |
| GET | `/appointments/:id` | 预约详情 | 是 |
| PUT | `/appointments/:id/cancel` | 取消预约 | 是 |
| GET | `/stylists/:id/schedule` | 理发师可用时间段 | 是 |
| POST | `/stylists/:id/schedule` | 设置排班 | 是(理发师) |

#### 订单模块 (`/api/v1/orders`)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/orders` | 创建订单 | 是 |
| GET | `/orders/my` | 我的订单列表 | 是 |
| GET | `/orders/:id` | 订单详情 | 是 |
| POST | `/orders/:id/pay` | 发起支付 | 是 |
| POST | `/orders/:id/refund` | 申请退款 | 是 |

#### 会员模块 (`/api/v1/members`)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/members/levels` | 会员等级列表 | 否 |
| GET | `/members/balance` | 余额查询 | 是 |
| GET | `/members/points/log` | 积分流水 | 是 |
| POST | `/coupons/receive` | 领取优惠券 | 是 |
| GET | `/coupons/my` | 我的优惠券 | 是 |

#### 消息模块 (`/api/v1/messages`)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/messages` | 消息列表 | 是 |
| PUT | `/messages/:id/read` | 标记已读 | 是 |
| GET | `/notifications` | 通知列表 | 是 |
| WebSocket | `/ws` | 实时消息连接 | 是 |

---

## 🐳 部署方案

### Docker Compose 开发部署

项目已提供完整的 `docker-compose.yml`，一键启动全部基础设施：

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 停止所有服务
docker-compose down

# 清理数据卷（慎用）
docker-compose down -v
```

**包含的服务容器**:

| 服务 | 版本 | 端口 | 说明 |
|------|------|------|------|
| MySQL | 8.0 | 3306 | 主数据库 |
| Redis | 7-alpine | 6379 | 缓存/会话/锁 |
| Elasticsearch | 8.11 | 9200/9300 | 全文搜索 |
| MongoDB | 7 | 27017 | 日志存储 |
| RabbitMQ | 3.13-management | 5672/15672 | 消息队列+管理界面 |
| MinIO | latest | 9000/9001 | 对象存储+控制台 |
| Nginx | alpine | 80/443 | 反向代理 |

### Kubernetes 生产部署

生产环境推荐使用 Kubernetes 进行容器编排：

```yaml
# 示例: Deployment 配置
apiVersion: apps/v1
kind: Deployment
metadata:
  name: haircut-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: haircut-api
  template:
    metadata:
      labels:
        app: haircut-api
    spec:
      containers:
      - name: api
        image: haircut-api:v1.0.0
        ports:
        - containerPort: 8080
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: haircut-api-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: haircut-api
  minReplicas: 3
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70
```

### CI/CD 流水线

```yaml
# .github/workflows/deploy.yml (GitHub Actions)
name: Deploy Pipeline

on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Build Backend
        run: cd server && go build -o bin/haircut-api cmd/server/main.go
      
      - name: Build Frontend
        run: cd apps/admin-web && npm ci && npm run build
      
      - name: Docker Build & Push
        run: |
          docker build -t haircut-api:${{ github.sha }} -f docker/Dockerfile.server .
          docker push registry.example.com/haircut-api:${{ github.sha }}

  deploy-staging:
    needs: build
    runs-on: ubuntu-latest
    environment: staging
    steps:
      - name: Deploy to Staging
        run: kubectl set image deployment/haircut-api haircut-api=registry.example.com/haircut-api:${{ github.sha }}
```

---

## 🛡️ 安全设计

### 认证与授权

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   微信登录   │────▶│  JWT Token  │────▶│ Casbin RBAC │
│  OpenID/OAuth│     │ Access+Refresh│    │  权限校验    │
└─────────────┘     └─────────────┘     └─────────────┘
```

- **JWT**: 双Token机制 (Access Token 2h + Refresh Token 7d)
- **RBAC**: Casbin 权限模型，支持角色继承和数据隔离
- **多租户**: TenantID 数据行级隔离

### 数据安全

| 措施 | 实现 |
|------|------|
| **密码加密** | bcrypt (cost=12) |
| **敏感数据脱敏** | 手机号/身份证号中间位掩码 |
| **HTTPS强制** | Nginx SSL Termination + HSTS |
| **SQL注入防护** | GORM 参数化查询 + 输入验证 |
| **XSS防护** | 前端输入过滤 + CSP策略 |
| **CSRF防护** | SameSite Cookie + Token验证 |

### 接口安全

| 措施 | 说明 |
|------|------|
| **IP限流** | 滑动窗口算法，默认100请求/分钟 |
| **签名验证** | 敏感接口HMAC-SHA256签名 |
| **防重放攻击** | 请求时间戳 + Nonce机制 |
| **白名单** | CORS允许域名严格配置 |

---

## ⚡ 性能优化

### 缓存策略

| 层级 | 技术 | TTL | 场景 |
|------|------|-----|------|
| L1 本地缓存 | sync.Map | 5s | 热点配置数据 |
| L2 分布式缓存 | Redis | 5min-24h | 用户信息/门店列表/理发师详情 |
| L3 CDN缓存 | Nginx/CloudFlare | 1h-7d | 静态资源/官网页面 |

### Redis Key设计规范

```
haircut:user:{userId}:info           → 用户信息 Hash
haircut:store:{storeId}:detail      → 门店详情 Hash  
haircut:store:nearby:{lat}:{lng}    → 附近门店 Geo
haircut:appointment:{slot}:lock      → 预约时段锁 String (SETNX)
haircut:coupon:template:{id}:stock  → 优惠券库存 String (DECR)
haircut:user:{userId}:token         → JWT黑名单 Set
```

### Elasticsearch索引设计

```
# 门店索引
{
  "mappings": {
    "properties": {
      "name": { "type": "text", "analyzer": "ik_max_word" },
      "address": { "type": "text" },
      "location": { "type": "geo_point" },
      "tags": { "type": "keyword" },
      "status": { "type": "keyword" }
    }
  }
}

# 地理位置搜索 (附近门店)
GET /stores/_search
{
  "query": {
    "bool": {
      "must": [{ "match": { "name": "理发" } }],
      "filter": [
        { "geo_distance": { "distance": "10km", "location": { "lat": 39.9, "lng": 116.4 } } },
        { "term": { "status": "OPEN" } }
      ]
    }
  },
  "sort": [
    { "_geo_distance": { "location": { "lat": 39.9, "lng": 116.4 }, "order": "asc" } }
  ]
}
```

---

## 📊 监控告警

### 健康检查端点

```bash
# 服务健康状态
GET /health

# 详细信息 (需要管理员权限)
GET /health/detail
```

### 关键指标监控

| 指标类型 | 监控项 | 告警阈值 |
|----------|--------|----------|
| **系统** | CPU使用率 | > 80% |
| **内存** | 内存占用 | > 85% |
| **数据库** | MySQL连接数 | > 最大连接数80% |
| **缓存** | Redis内存使用 | > 90% |
| **API** | 平均响应时间(P99) | > 500ms |
| **API** | 错误率 | > 1% |
| **业务** | 支付成功率 | < 95% |

---

## 📋 开发指南

### 代码规范

```bash
# 后端代码检查
cd server && golangci-lint run ./...

# 前端代码检查
cd apps/admin-web && npx eslint . --ext .ts,.tsx

# 代码格式化
cd server && go fmt ./...
cd apps/admin-web && npx prettier --write .
```

### Git提交规范

```
feat: 新功能
fix: Bug修复
docs: 文档更新
style: 代码格式调整
refactor: 重构
perf: 性能优化
test: 测试相关
chore: 构建/工具链
```

### 分支策略

```
main (保护分支) ← 发布版本
  ↑
develop ← 开发集成分支
  ↑
feature/* ← 功能开发分支
bugfix/* ← 紧急修复
release/* ← 版本发布准备
```

---

## ❓ 常见问题

<details>
<summary><b>Q1: 如何配置微信小程序？</b></summary>

1. 在 [微信公众平台](https://mp.weixin.qq.com) 注册小程序账号
2. 获取 AppID 和 AppSecret
3. 配置服务器域名（需HTTPS）
4. 编辑 `.env` 文件填入密钥：
```env
WECHAT_APP_ID=your_app_id
WECHAT_APP_SECRET=your_app_secret
```

</details>

<details>
<summary><b>Q2: 微信支付如何对接？</b></summary>

1. 在 [微信支付商户平台](https://pay.weixin.qq.com) 申请商户号
2. 获取 MchID 和 API Key
3. 配置支付回调URL
4. 编辑 `.env` 文件：
```env
WECHAT_MCH_ID=your_merchant_id
WECHAT_API_KEY=your_api_key
WECHAT_NOTIFY_URL=https://your-domain.com/api/v1/payments/wechat/notify
```

</details>

<details>
<summary><b>Q3: 腾讯地图如何配置？</b></summary>

1. 在 [腾讯地图开放平台](https://lbs.qq.com) 创建应用
2. 申请 WebService API Key
3. 编辑 `.env` 文件：
```env
TENCENT_MAP_KEY=your_map_key
```
</details>

<details>
<summary><b>Q4: 如何扩展新的门店？</b></summary>

1. 登录运营后台
2. 进入「门店管理」→「新建门店」
3. 填写门店信息、上传资质
4. 提交审核后自动上线
</details>

<details>
<summary><b>Q5: 数据备份策略？</b></summary>

项目提供自动备份脚本：

```bash
# 手动执行备份
./scripts/backup.sh

# 备份内容：
# - MySQL 全量dump + Binlog
# - Redis RDB快照
# - MinIO 对象文件
# 自动保留最近7天备份
```

建议配合 Cron 定时执行：
```crontab
0 2 * * * /path/to/scripts/backup.sh >> /var/log/backup.log 2>&1
```

</details>

---

## 🗺️ 项目路线图

### v1.0 (当前版本)
- ✅ 基础架构搭建
- ✅ 38张数据库表设计
- ✅ 核心 API 实现
- ✅ 前端框架初始化

### v1.1 (计划中)
- 🔲 小程序完整页面开发
- 🔲 运营后台完整功能
- 🔲 微信支付全流程对接
- 🔲 腾讯地图深度集成

### v1.2 (规划中)
- 🔲 AI发型推荐功能
- 🔲 数据分析大屏
- 🔲 CRM自动化营销
- 🔲 移动端适配优化

### v2.0 (远期规划)
- 🔲 微服务拆分
- 🔲 多租户SaaS化
- 🔲 全国连锁多地域部署
- 🔲 国际化支持 (i18n)

---

## 🤝 贡献指南

我们欢迎任何形式的贡献！请遵循以下步骤：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'feat: add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

### 贡献者行为准则

- 尊重他人，保持专业和建设性讨论
- 代码必须通过 lint 检查和单元测试
- 提交信息符合 Conventional Commits 规范
- 大功能变更前先提 Issue 讨论

---

## 📄 License

本项目采用 MIT License 开源协议。详见 [LICENSE](./LICENSE) 文件。

```
MIT License

Copyright (c) 2024 HairCut Team

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software...
```

---

## 💬 联系我们

- **项目地址**: [https://github.com/hlw422/hair_cut](https://github.com/hlw422/hair_cut)
- **问题反馈**: 请提 [Issue](https://github.com/hlw422/hair_cut/issues)
- **技术交流**: 欢迎 Pull Request 和 Discussion

---

<p align="center">
  <strong>Made with ❤️ by HairCut Team | Powered by Golang & React</strong>
</p>

<p align="center">
  <sub>⭐ 如果这个项目对您有帮助，请给一个 Star 支持一下！</sub>
</p>