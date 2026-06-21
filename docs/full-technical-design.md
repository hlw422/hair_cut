# HairCut 连锁理发店数字化平台 - 完整技术设计文档

> **版本**: v1.0.0  
> **最后更新**: 2024-01-20  
> **文档状态**: 正式版  
> **适用范围**: 开发团队 / 运维团队 / 产品团队

---

## 目录

1. [产品与系统架构](#一产品与系统架构)
2. [微服务拆分方案](#二微服务拆分方案)
3. [API接口设计规范](#三api接口设计规范)
4. [Redis缓存设计方案](#四redis缓存设计方案)
5. [Elasticsearch搜索方案](#五elasticsearch搜索方案)
6. [RabbitMQ消息队列方案](#六rabbitmq消息队列方案)
7. [安全设计方案](#七安全设计方案)
8. [高并发设计方案](#八高并发设计方案)
9. [全国连锁门店扩展方案](#九全国连锁门店扩展方案)
10. [多租户SaaS方案](#十多租户saas方案)
11. [Docker部署方案](#十一docker部署方案)
12. [Kubernetes部署方案](#十二kubernetes部署方案)
13. [CI/CD方案](#十三cicd方案)

---

## 一、产品与系统架构

### 1.1 系统组成（6大子系统）

| 子系统 | 技术栈 | 目标用户 | 核心功能 |
|--------|--------|----------|----------|
| **用户端小程序** | Taro + React + TS + TDesign | 消费顾客 | 预约/支付/会员/AI推荐 |
| **理发师端小程序** | Taro + React + TS | 理发师 | 工作台/排班/客户管理 |
| **店长端小程序** | Taro + React + TS | 门店管理者 | 运营数据/员工/库存 |
| **总部运营后台** | React18 + Vite + Shadcn UI | 集团运营人员 | 数据分析/营销/CRM |
| **官网宣传网站** | Next.js 14 SSR | 品牌展示 | SEO优化/门店查询/加盟 |
| **API服务端** | Golang + Gin + GORM | 所有前端 | 统一业务逻辑 |

### 1.2 整体技术栈

```
┌─────────────────────────────────────────────────────────────┐
│                        客户端层                              │
│  ┌───────────┐ ┌───────────┐ ┌───────────┐ ┌─────────────┐  │
│  │ 用户小程序 │ │ 理发师小程序│ │ 店长小程序 │ │ Web(后台+官网)│  │
│  └─────┬─────┘ └─────┬─────┘ └─────┬─────┘ └──────┬──────┘  │
│        │            │           │              │          │
│        └────────────┴───────────┘              │          │
│                         │                      │          │
├─────────────────────────▼──────────────────────▼──────────┤
│                    API网关层 (Nginx)                        │
│         HTTPS终止 / 负载均衡 / 静态资源 / CORS               │
├─────────────────────────────────────────────────────────────┤
│                     应用服务层 (Golang Gin)                   │
│  ┌────────┬────────┬────────┬────────┬────────┬────────┐   │
│  │ 用户   │ 门店   │ 预约   │ 支付   │ 营销   │ 分析   │   │
│  │ 模块   │ 模块   │ 订单   │ 会员   │ CRM    │ 报表   │   │
│  └────────┴────────┴────────┴────────┴────────┴────────┘   │
│                                                             │
│  中间件链: JWT认证 → Casbin权限 → 限流 → 日志 → 异常恢复      │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│                     基础设施层                               │
│  ┌──────┐ ┌──────┐ ┌────────┐ ┌──────┐ ┌──────┐ ┌──────┐  │
│  │MySQL │ │Redis │ │ Elastic │ │MongoDB│ │Rabbit│ │ MinIO│  │
│  │ 8.0  │ │ 7.x  │ │ search  │ │       │ │  MQ  │ │      │  │
│  └──────┘ └──────┘ └────────┘ └──────┘ └──────┘ └──────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 1.3 核心业务流程图

#### 预约下单流程（核心链路）
```
用户选门店 → 选理发师 → 选服务 → 选时间 → 提交订单
     ↓
创建预约记录(MySQL事务) 
     ↓
锁定时间段(Redis分布式锁: apt:{stylistId}:{date}:{time})
     ↓
计算金额(原价→会员折扣→优惠券→积分抵扣=实付金额)
     ↓
创建待支付订单(Order状态=0)
     ↓
调用微信支付(获取prepay_id返回前端)
     ↓
用户完成支付 ← 微信回调通知(RabbitMQ异步监听)
     ↓
更新订单状态(Paid) + 扣减优惠券 + 扣减积分
     ↓
发送预约确认通知(WebSocket实时推送 + 微信模板消息)
     ↓
ES同步更新搜索索引(异步)
```

---

## 二、微服务拆分方案

### 2.1 第一阶段：单体模块化（当前）

```
haircut-server (单进程)
├── internal/
│   ├── api/handler/     # HTTP处理器
│   ├── service/        # 业务逻辑层
│   ├── repository/     # 数据访问层
│   ├── model/          # 数据模型
│   ├── middleware/      # 中间件
│   └── pkg/            # 内部工具包
└── cmd/server/main.go  # 入口
```

**优势**: 
- 开发效率高（无需跨服务调用）
- 事务一致性强（本地ACID）
- 部署运维简单

**拆分准备**:
- 通过interface抽象Service和Repository层
- 每个领域包可独立测试

### 2.2 第二阶段：按领域拆分为独立微服务

```
                          API Gateway (Nginx/Kong)
                                  │
        ┌───────────────────────────┼───────────────────────────┐
        ▼                           ▼                           ▼
┌──────────────┐          ┌──────────────┐          ┌──────────────┐
│ user-service │          │ order-service │          │ store-service│
│ (用户/会员)   │◄────RPC─►│ (预约/订单)   │◄────RPC─►│ (门店/理发师)  │
└──────────────┘          └──────────────┘          └──────────────┘
        │                           │                           │
        ▼                           ▼                           ▼
┌──────────────┐          ┌──────────────┐          ┌──────────────┐
│ payment-svc  │          │ marketing-svc│          │ analytics-svc│
│ (支付/财务)   │          │ (营销/CRM)   │          │ (数据分析)     │
└──────────────┘          └──────────────┘          └──────────────┘
```

**触发条件**:
- 团队规模 > 15人
- 单个服务QPS > 1000
- 部署频率需求高（不同服务独立发布）

**通信方式**:
- 同步: gRPC（内部服务间）
- 异步: RabbitMQ事件驱动（解耦）

---

## 三、API接口设计规范

### 3.1 RESTful URL规范

```
基础路径: /api/v1

# 公开接口（无需认证）
POST   /api/v1/auth/login          # 登录
POST   /api/v1/auth/wechat         # 微信登录
GET    /api/v1/stores/nearby?lat=&lng=&radius=  # 附近门店
GET    /api/v1/stores/:id          # 门店详情
GET    /api/v1/stylists/:id        # 理发师详情
GET    /api/v1/services            # 服务项目列表
GET    /api/v1/health              # 健康检查

# 用户认证接口
GET    /api/v1/user/profile        # 个人信息
PUT    /api/v1/user/profile        # 更新信息
GET    /api/v1/user/member         # 会员信息
GET    /api/v1/user/appointments   # 我的预约列表
GET    /api/v1/user/orders         # 我的订单列表

# 预约接口
POST   /api/v1/appointments        # 创建预约
GET    /api/v1/appointments/:id/slots # 可用时间段
POST   /api/v1/appointments/:id/cancel  # 取消预约

# 订单接口
POST   /api/v1/orders              # 创建订单
POST   /api/v1/orders/:id/pay      # 发起支付
GET    /api/v1/orders/:id          # 订单详情
POST   /api/v1/orders/:id/refund   # 申请退款

# 优惠券接口
GET    /api/v1/coupons/available   # 我的可用优惠券
POST   /api/v1/coupons/:id/use     # 使用优惠券

# 管理后台接口（需管理员角色）
GET    /api/v1/admin/users         # 用户管理
GET    /api/v1/admin/stores        # 门店管理
GET    /api/v1/admin/analytics/dashboard  # 数据看板
...
```

### 3.2 统一响应格式

```json
// 成功响应
{
  "code": 0,
  "message": "success",
  "data": { ... },
  "meta": { "page": 1, "per_page": 20, "total": 100, "total_pages": 5 }
}

// 错误响应
{
  "code": 400,
  "message": "参数验证失败",
  "data": null,
  "error": {
    "field": "phone",
    "message": "手机号格式不正确"
  }
}
```

### 3.3 分页规范

```
GET /api/v1/stores?page=1&per_page=20&sort_by=rating&sort_dir=desc&keyword=

Query参数:
- page: 页码（默认1，最小1）
- per_page: 每页数量（默认20，最大100）
- sort_by: 排序字段（白名单限制防SQL注入）
- sort_dir: 排序方向 asc/desc
- keyword: 搜索关键词

响应Meta:
{
  "page": 1,
  "per_page": 20,
  "total": 128,        // 总记录数
  "total_pages": 7     // 总页数
}
```

---

## 四、Redis缓存设计方案

### 4.1 缓存层级策略

```
请求 → 本地缓存(LRU, 5s TTL) → Redis集群 → MySQL数据库
                ↓ miss             ↓ hit
            回源查询          直接返回
```

### 4.2 缓存Key命名规范

```go
// Key命名规则: {业务}:{标识}:{参数}
const (
    // 用户相关
    CacheUserBase      = "user:base:%d"           // 用户基础信息
    CacheUserMember     = "user:member:%d"          // 会员信息
    
    // 门店相关
    CacheStoreDetail    = "store:detail:%d"         // 门店详情(热点)
    CacheStoreNearby    = "store:nearby:%f:%f:%d"  // 附近门店(lat,lng,radius)
    CacheStoreListCity  = "store:list:%s:%d:%d"    // 城市门店列表(city,page,size)

    // 理发师相关
    CacheStylistDetail  = "stylist:detail:%d"
    CacheStylistStore   = "stylist:store:%d:%d"     // 门店的理发师列表

    // 预约锁（分布式锁）
    LockAppointmentSlot = "lock:apt:%d:%s:%s"       // {stylist_id}:{date}:{time}

    // 会话/Token
    SessionUserToken   = "session:token:%s"        // JWT Token黑名单(登出用)
    
    // 验证码
    CodeSMSVerify      = "code:sms:%s"              // 手机验证码
)
```

### 4.3 缓存过期策略

| 数据类型 | TTL策略 | 说明 |
|---------|---------|------|
| 用户基本信息 | 30分钟 | 变更频率低 |
| 门店详情 | 10分钟 | 较少变更 |
| 附近门店列表 | 2分钟 | 实时性要求高 |
| 热门服务项目 | 1小时 | 几乎不变 |
| 配置项/字典 | 24小时 | 极低变更 |
| 分布式锁 | 自动释放 | 30秒超时自动删除 |

### 4.4 缓存一致性保障

```
写操作流程:
1. 更新MySQL数据库
2. 删除对应Redis缓存（而非更新，防止脏读）
3. 发布缓存失效消息到MQ（可选，多实例时广播清除）

读操作流程:
1. 先查Redis
2. 命中直接返回
3. 未命中查MySQL
4. 写入Redis（设置TTL）
5. 返回数据
```

### 5.5 缓存穿透/击穿/雪崩防护

- **穿透**: 布隆过滤器（判断key是否存在）+ 缓空值（TTL短，如60s）
- **击穿**: 热点数据永不过期 + 互斥锁重建
- **雪崩**: TTL随机化（基础时间 ± 随机值） + 多级缓存 + 熔断降级

---

## 五、Elasticsearch搜索方案

### 5.1 ES索引设计

```json
// 门店搜索索引 (stores_index)
{
  "mappings": {
    "properties": {
      "store_id":     {"type": "long"},
      "name":         {"type": "text", "analyzer": "ik_max_word", "fields": {"keyword": {"type": "keyword"}}},
      "city":         {"type": "keyword"},
      "district":     {"type": "keyword"},
      "address":      {"type": "text", "analyzer": "ik_max_word"},
      "location":     {"type": "geo_point"},          // 地理坐标点
      "tags":         {"type": "keyword"},
      "avg_price":    {"type": "float"},
      "rating":       {"type": "float"},
      "status":       {"type": "integer"},
      "is_featured":  {"type": "boolean"},
      "updated_at":   {"type": "date"}
    }
  }
}

// 理发师搜索索引 (stylists_index)
{
  "mappings": {
    "properties": {
      "stylist_id":   {"type": "long"},
      "name":         {"type": "text", "analyzer": "ik_max_word"},
      "store_name":   {"type": "keyword"},
      "title":        {"type": "keyword"},
      "specialties":  {"type": "text", "analyzer": "ik_smart"},
      "experience":   {"type": "integer"},
      "rating":       {"type": "float"},
      "status":       {"type": "integer"}
    }
  }
}

// 服务项目搜索索引 (services_index) - 类似结构
```

### 5.2 地理位置搜索（附近门店）

```go
// 使用ES Geo Distance Query实现
func SearchNearbyStores(lat, lng float64, radiusKm int, page, size int) ([]StoreDoc, error) {
    query := map[string]interface{}{
        "bool": map[string]interface{}{
            "must": []map[string]interface{}{
                {
                    "geo_distance": map[string]interface{}{
                        "distance": fmt.Sprintf("%dkm", radiusKm),
                        "location": map[string]float64{"lat": lat, "lon": lng},
                    },
                },
                {
                    "term": map[string]interface{}{"status": 1}, // 只查营业中的门店
                },
            },
        },
    }

    // 按距离排序
    sort := []map[string]interface{}{
        {"_geo_distance": map[string]interface{}{
            "location":  map[string]float64{"lat": lat, "lon": lng},
            "order":    "asc",
            "unit":     "km",
        }},
        {"rating": map[string]string{"order": "desc"}},
    }

    return es.Search("stores_index", query, sort, page, size)
}
```

### 5.3 数据同步机制

```
MySQL写入 → Binlog Canal监听 → RabbitMQ消息 → ES同步写入
                                    或
手动触发 → Service层双写 → MQ异步同步 → ES Update
```

**同步策略**:
- 实时同步（门店状态变更、新开门店）
- 延迟同步（评价更新、评分变化，延迟5分钟批量处理）
- 全量重建（每天凌晨定时全量同步一次，保证一致性）

---

## 六、RabbitMQ消息队列方案

### 6.1 Exchange & Queue 设计

```
Exchange: haircut.direct (直连交换机, durable=true)

Queues:
├── queue.payment.callback     # 支付回调通知 (消费者: PaymentService)
├── queue.order.created        # 订单创建通知 (消费者: NotificationService, AnalyticsService)
├── queue.appointment.confirm   # 预约确认通知
├── queue.user.registered       # 新用户注册 (消费者: MarketingService, WelcomeEmail)
├── queue.es.sync              # ES索引同步 (消费者: SyncToESService)
├── cache.invalidation         # 缓存失效广播 (多消费者: 各API实例)
└── queue.sms.send             # 短信发送 (消费者: SMSService)
```

### 6.2 核心消息流转示例

#### 支付回调流程
```
微信支付服务器 → POST /api/v1/payment/wechat/notify
                       ↓
              PaymentHandler接收原始XML
                       ↓
              解析并验证签名
                       ↓
              发布消息到 queue.payment.callback:
              {
                "transaction_id": "420000xxxxx",
                "out_trade_no": "ORD20240120001",
                "trade_state": "SUCCESS",
                "amount": { "total": 19800, ... },
                "pay_time": "2024-01-20T14:35:22+08:00"
              }
                       ↓
              PaymentConsumer消费:
              ① 更新Order状态为Paid
              ② 创建PaymentRecord记录
              ③ 如果使用了优惠券 → 标记已使用
              ④ 如果使用积分 → 扣减积分余额
              ⑤ 发布 order.paid 事件（供其他服务订阅）
                       ↓
              NotificationConsumer消费:
              发送WebSocket消息给用户("支付成功")
              发送微信模板消息("预约提醒")
              
              AnalyticsConsumer消费:
              实时更新GMV统计
```

### 6.3 死信队列(DLQ)与重试机制

```yaml
# Queue配置
queue.payment.callback:
  durable: true
  arguments:
    x-dead-letter-exchange: ""        # 死信到默认exchange
    x-dead-letter-routing-key: "dlq.payment.callback"  # 死信路由key
    x-message-ttl: 86400000          # 消息最大存活时间24h
    x-max-retry-count: 3             # 最大重试次数

# 消费者重试策略:
# 第1次失败: 立即重试
# 第2次失败: 延迟5s重试  
# 第3次失败: 延迟30s重试
# 第4次失败: 进入DLQ死信队列，人工排查后重新投递
```

---

## 七、安全设计方案

### 7.1 认证与授权体系

```
用户登录 → 验证凭证 → 生成JWT Token(含userId/roles/tenantId/exp)
    ↓
前端存储(Token in localStorage / HttpOnly Cookie)
    ↓
每次请求携带(Authorization: Bearer <token>)
    ↓
JWT中间件解析验证 → 注入上下文(user_id, roles)
    ↓
Casbin中间件权限校验(roles × path × method)
    ↓
通过 → Handler执行业务逻辑
拒绝 → 返回403 Forbidden
```

### 7.2 JWT Token安全配置

```json
{
  "algorithm": "HS256",
  "secret_key": ">=32字节强随机密钥(生产环境从Vault/KMS获取)",
  "access_token_ttl": "24小时",
  "refresh_token_ttl": "7天",
  "issuer": "haircut-server",
  "audience": ["haircut-user-app", "haircut-admin-web"]
}
```

### 7.3 RBAC权限矩阵（部分）

| 角色 | 用户管理 | 门店管理 | 订单查看 | 财务报表 | 数据导出 | 系统设置 |
|------|---------|---------|---------|---------|---------|---------|
| super_admin | ✅ CRUD | ✅ CRUD | ✅ 查看 | ✅ 全部 | ✅ 全部 | ✅ 全部 |
| admin | ✅ 查看 | ✅ CRUD | ✅ 查看 | ✅ 区域内 | ✅ 区域内 | ❌ |
| regional_manager | ✅ 查看 | ✅ 区域内 | ✅ 区域内 | ✅ 区域内 | ✅ 区域内 | ❌ |
| store_manager | ❌ | ✅ 本店 | ✅ 本店 | ✅ 本店 | ✅ 本店 | ❌ |
| stylist | ❌ | ❌ | ✅ 本人 | ✅ 个人收入 | ❌ | ❌ |
| cs_staff | ✅ 查看 | ❌ | ✅ 查看 | ❌ | ❌ | ❌ |
| member_user | 仅个人 | ❌ | ✅ 个人 | ❌ | ❌ | ❌ |

### 7.4 数据脱敏与加密

| 字段类型 | 存储方式 | 展示方式 |
|---------|---------|---------|
| 密码 | BCrypt哈希(不可逆) | 不显示(仅修改入口) |
| 手机号 | 明文存储 | 列表显示: 138****8888 |
| 身份证号 | AES加密存储 | 显示: 310***********123X |
| 银行卡号 | AES加密存储 | 显示: **** **** **** 1234 |
| 微信SessionKey | 加密存储 | 不返回给前端 |
| 支付密钥 | Vault/KMS管理 | 应用启动时加载到内存 |

### 7.5 接口安全措施

- **HTTPS强制**: 生产环境所有API必须走TLS 1.2+
- **CORS严格配置**: 白名单Origin，禁止 `*`
- **限流**: IP维度 + 用户维度双重限流
- **参数校验**: GIN Validator + 自定义规则，防SQL注入/XSS
- **SQL注入防护**: 全部使用GORM Parameterized Query（杜绝拼接SQL）
- **CSRF保护**: 关键操作（支付/修改密码）需二次验证码
- **敏感操作审计**: 所有写操作记录operation_log（含IP/UA/参数）

### 7.6 文件上传安全

```
文件上传流程:
1. 前端上传至MinIO (预签名URL或直传)
2. 返回文件访问URL
3. 后台异步扫描病毒/检测文件类型(Magic Number校验)
4. 图片自动压缩/WebP转换(减少存储和带宽)
5. 敏感文件(如身份证)加密存储
6. URL有效期控制(Presigned URL, 默认1小时)
```

---

## 八、高并发设计方案

### 8.1 并发瓶颈分析

| 场景 | 预估QPS | 瓶颈点 | 优化方案 |
|-----|---------|--------|---------|
| 日常浏览 | 500-2000 | DB查询 | 多级缓存 + CDN |
| 周末高峰 | 5000-20000 | 预约时段竞争 | Redis分布式锁 + 异步削峰 |
| 新品发布/活动 | 10000-50000 | 下单/支付 | 消息队列削峰 + 库存预热 |
| 秒杀活动 | 50000+ | 库存扣减 | Redis原子操作 + 令牌桶限流 |

### 8.2 核心优化手段

#### 8.2.1 数据库层面
- **读写分离**: 主库写 + 从库读（GORM支持多个DB连接）
- **分库分表**: 订单表按月水平分表 `orders_202401`, `orders_202402`...
- **索引优化**: 复合索引覆盖高频查询，避免回表
- **连接池**: 合理配置MaxOpenConns(100) + MaxIdleConns(10)
- **慢SQL监控**: 开启慢查询日志(>100ms)，定期优化

#### 8.2.2 缓存层面
- **热点数据预热**: 启动时或定时任务将热点数据加载到Redis
- **本地缓存**: Go-SyncMap/LRU缓存极热数据(如系统配置，TTL=5s)，减少Redis压力
- **Cache-Aside模式**: 读时先查缓存，写时删缓存（非更新）
- **缓存击穿防护**: 热点Key永不过期 + 互斥锁重建

#### 8.2.3 应用层面
- **连接复用**: HTTP Keep-Alive + 数据库连接池
- **异步处理**: 非关键路径异步（如发送通知、更新ES、统计日志走MQ）
- **批量操作**: 批量查询替代N+1问题（GORM Preload）
- **限流降级**: 信号量/令牌桶算法限流，超载时返回友好提示

#### 8.2.4 架构层面
- **负载均衡**: Nginx upstream + 多个Server实例
- **水平扩展**: 无状态设计，可随时增加Pod副本
- **CDN加速**: 静态资源(JS/CSS/图片)走CDN，回源到MinIO
- **动静分离**: API走动态服务，静态资源走Nginx/CDN

### 8.3 预约时段冲突解决（核心难点）

```go
// 方案: Redis分布式锁 + Lua脚本保证原子性
func BookAppointment(stylistID uint64, date string, timeSlot string) error {
    lockKey := fmt.Sprintf("apt:lock:%d:%s:%s", stylistID, date, timeSlot)
    
    // 尝试获取锁（30秒自动释放，防止死锁）
    acquired, unlock := redis.SetNX(ctx, lockKey, 30*time.Second)
    if !acquired {
        return errors.New("该时间段已被预约，请选择其他时间")
    }
    defer unlock() // 无论成功失败都释放锁
    
    // 再次检查数据库（双重检查，防止并发竞态）
    exists, _ := db.Where("stylist_id=? AND date=? AND time_slot=? AND status IN (?)", 
        stylistID, date, timeSlot, []int{0, 1}).Exists()
    if exists {
        return errors.New("该时间段已被预约")
    }
    
    // 创建预约记录（MySQL事务）
    return db.Transaction(func(tx *gorm.DB) error {
        appointment := &Appointment{...}
        if err := tx.Create(appointment).Error; err != nil {
            return err
        }
        
        // 设置长期锁（直到预约结束或取消）
        lockKeyLong := fmt.Sprintf("apt:booked:%d:%s:%s", stylistID, date, timeSlot)
        redis.Set(ctx, lockKeyLong, "1", time.Until(date+" "+strings.Split(timeSlot, "-")[1]))
        
        return nil
    })
}
```

---

## 九、全国连锁门店扩展方案

### 9.1 多地域部署架构

```
                    用户请求 (DNS智能解析)
                         │
              ┌──────────┼──────────┐
              ▼          ▼          ▼
        华东区域      华南区域      华北区域
    (上海/杭州)   (广州/深圳)   (北京/天津)
         │            │            │
    ┌────┴────┐  ┌────┴────┐  ┌────┴────┐
    │API Server│  │API Server│  │API Server│
    └────┬────┘  └────┬────┘  └────┬────┘
         │            │            │
    ┌────┴────┐  ┌────┴────┐  ┌────┴────┐
    │MySQL主从│  │MySQL主从│  │MySQL主从│
    │(区域库) │  │(区域库) │  │(区域库) │
    └─────────┘  └─────────┘  └─────────┘
         │            │            │
         └───── ───────┘────────────┘
                    │
              总部数据中心
         (汇总数据/全局配置/统一认证)
              Redis Cluster (跨地域同步)
              Elasticsearch Cluster
              MongoDB Replica Set
              RabbitMQ Cluster
```

### 9.2 数据同步策略

| 数据类型 | 同步方向 | 同步方式 | 延迟要求 |
|---------|---------|---------|---------|
| 用户账号 | 总部→区域 | 实时同步(主从复制) | <1s |
| 门店基础数据 | 总部→区域 | 定时批量(每小时) | <1h |
| 订单/交易数据 | 区域→总部 | 实时MQ上报 | <5s |
| 统计报表数据 | 区域→总部 | T+1离线批处理 | <24h |
| 配置/字典 | 总部→区域 | 推送+版本控制 | <10min |

### 9.3 跨区域调度能力

- **就近接入**: DNS根据用户IP地理位置解析到最近节点
- **异地预约**: 支持用户在A城市预约B城市门店（数据可跨区读取）
- **故障切换**: 某区域故障时流量切到邻近健康节点（需考虑数据一致性窗口）

---

## 十、多租户SaaS方案

### 10.1 租户模型设计

```
三种租户级别:

Level 1: 平台租户 (TenantID = Platform)
  - HairCut官方自营门店
  - 拥有全部功能和最高配额

Level 2: 品牌加盟商 (TenantID = Brand_XXX)
  - 加盟品牌旗下所有门店共享一个租户ID
  - 可自定义部分UI、定价策略
  - 数据完全隔离

Level 3: 单店商户 (TenantID = Store_YYY)
  - 小型独立门店入驻
  - 使用平台标准功能
  - 数据隔离但可参与平台联合营销
```

### 10.2 数据隔离实现

```go
// 中间件自动注入租户上下文
func TenantMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 优先从JWT Token提取tenant_id
        tenantID := extractTenantFromJWT(c)
        
        // 其次从Header提取（开放API场景）
        if tenantID == 0 {
            tenantID = extractTenantFromHeader(c)
        }
        
        // 最后使用默认租户（兼容老版本）
        if tenantID == 0 {
            tenantID = DefaultTenantID // 平台默认
        }
        
        c.Set("tenant_id", tenantID)
        c.Set("db_scope", ByTenant(tenantID)) // GORM Scope
        
        c.Next()
    }
}

// GORM全局Scope（所有查询自动附加WHERE条件）
func ByTenant(tenantID uint64) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("tenant_id = ? OR tenant_id = 0", tenantID)
        // tenant_id=0 表示公共数据（如系统配置），所有租户可见
    }
}
```

### 10.3 私有化定制扩展点

- **主题配色**: 租户级CSS变量覆盖（品牌色、Logo等）
- **字段扩展**: 通用扩展字段JSON列（ext_fields JSON）
- **插件市场**: 可选功能模块（高级报表、AI推荐等额外付费）
- **API版本控制**: 保证旧租户不受升级影响

---

## 十一、Docker部署方案

### 11.1 已提供的Docker编排（见 docker-compose.yml）

包含基础设施容器：
- MySQL 8.0
- Redis 7
- Elasticsearch 8.11
- MongoDB 6.0
- RabbitMQ 3.12 (含Management UI)
- MinIO (对象存储)

### 11.2 应用容器化

```dockerfile
# 后端服务 (见 docker/Dockerfile.server)
# 多阶段构建: Build阶段(Go编译) → Run阶段(Alpine运行)

# 前端应用 (见 docker/Dockerfile.admin-web)
# Node构建 → Nginx托管静态资源
```

### 11.3 快速启动命令

```bash
# 1. 复制环境变量
cp docker/.env.example .env

# 2. 编辑.env填入真实配置（密钥、AppID等）
vim .env

# 3. 一键启动所有基础设施
docker-compose up -d mysql redis elasticsearch mongodb rabbitmq minio

# 4. 初始化数据库
chmod +x scripts/init-db.sh
./scripts/init-db.sh

# 5. 启动后端API
cd server && go run cmd/server/main.go

# 6. 启动前端（另开终端）
cd apps/admin-web && npm run dev
```

---

## 十二、Kubernetes部署方案（生产环境）

### 12.1 架构概览

```
Kubernetes Cluster v1.28+
├── Namespace: haircut-prod
│   ├── Deployment: haircut-api (3 replicas)
│   │   ├── Pod: api-xxx (Resource: CPU 500m, Mem 512Mi)
│   │   ├── Pod: api-yyy
│   │   └── Pod: api-zzz
│   ├── Deployment: haircut-admin-web (2 replicas)
│   ├── Deployment: haircut-website (2 replicas)
│   ├── Service: haircut-api-svc (ClusterIP, port 8080)
│   ├── Service: admin-web-svc (NodePort/LoadBalancer)
│   ├── ConfigMap: app-config (环境变量配置)
│   ├── Secret: db-credentials (数据库密码等敏感信息)
│   ├── Ingress: haircut-ingress (域名路由 + TLS证书)
│   └── HPA: api-hpa (自动扩缩容, 3~50 replicas based on CPU)
│
├── Namespace: haircut-infra
│   ├── StatefulSet: mysql (主从, PVC持久卷)
│   ├── StatefulSet: redis-cluster (6节点, 3主3从)
│   ├── StatefulSet: elasticsearch (3节点集群)
│   ├── StatefulSet: mongodb (副本集 3节点)
│   ├── Deployment: rabbitmq (集群模式)
│   ├── Deployment: minio (分布式, PVC)
│   └── PV/PVC: nfs-storage-class (网络存储)
```

### 12.2 关键资源配置示例

```yaml
# api-deployment.yaml (精简版)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: haircut-api
  namespace: haircut-prod
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 0
      maxSurge: 1
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
        image: registry.example.com/haircut/api:v1.0.0
        ports:
        - containerPort: 8080
        envFrom:
        - configMapRef:
            name: app-config
        - secretRef:
            name: db-credentials
        resources:
          requests:
            cpu: "250m"
            memory: "256Mi"
          limits:
            cpu: "1000m"
            memory: "1Gi"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 15
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
# hpa.yaml (自动扩缩容)
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: api-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: haircut-api
  minReplicas: 3
  maxReplicas: 50
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70  # CPU利用率超过70%时扩容
```

---

## 十三、CI/CD方案

### 13.1 CI流水线（GitHub Actions示例）

```yaml
# .github/workflows/ci.yml
name: CI Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  lint-and-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      
      - name: Lint
        run: |
          go vet ./...
          golangci-lint run ./...
      
      - name: Unit Test
        run: go test -v -coverprofile=coverage.out ./internal/...
      
      - name: Upload Coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
  
  build-docker:
    needs: lint-and-test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Build Docker Image
        run: docker build -t haircut-api:${{ github.sha }} -f docker/Dockerfile.server .
      
      - name: Push to Registry
        run: |
          echo ${{ secrets.DOCKER_PASSWORD }} | docker login registry.example.com -u ${{ secrets.DOCKER_USERNAME }} --password-stdin
          docker tag haircut-api:${{ github.sha }} registry.example.com/haircut/api:${{ github.sha }}
          docker push registry.example.com/haircut/api:${{ github.sha }}

  deploy-staging:
    needs: build-docker
    if: github.ref == 'refs/heads/develop'
    runs-on: ubuntu-latest
    environment: staging
    steps:
      - name: Deploy to Staging K8s
        run: |
          kubectl set image deployment/haircut-api api=registry.example.com/haircut/api:${{ github.sha }} -n haircut-staging
          kubectl rollout status deployment/haircut-api -n haircut-staging

  deploy-production:
    needs: build-docker
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    environment: production
    steps:
      - name: Deploy to Production K8s (Blue-Green)
        run: |
          # 蓝绿部署，零停机切换
          kubectl apply -f k8s/production-blue.yaml
          # 等待新Pod Ready
          kubectl rollout status deployment/api-blue -n haircut-prod --timeout=300s
          # 切换Service流量
          kubectl patch svc api-svc -p '{"spec":{"selector":{"version":"blue"}}}' -n haircut-prod
```

### 13.2 分支策略

```
main分支 (生产环境, 受保护)
  ↑ Merge Request (Code Review + CI通过必选)
  │
develop分支 (开发/测试环境)
  ↑ Feature分支 (功能开发)
  │
feature/user-auth (开发中)
feature/payment-wechat (开发中)
feature/ai-hairstyle (开发中)
```

### 13.3 发布流程

1. **开发完成**: Feature分支代码合并到develop
2. **CI验证**: 自动跑单元测试、集成测试、Lint、Build
3. **测试环境部署**: 自动部署到Staging环境，QA团队验收测试
4. **生产发布**: 创建Release PR (main←develop), 人工审核后合并
5. **蓝绿部署**: 自动构建镜像→部署到Blue环境→健康检查→切换流量
6. **监控告警**: 观察错误率、响应时间、业务指标，异常立即回滚

---

## 附录A：环境变量清单

| 变量名 | 说明 | 示例值 | 必填 |
|-------|------|--------|-----|
| MYSQL_HOST | MySQL地址 | localhost | 是 |
| MYSQL_PORT | MySQL端口 | 3306 | 否(默认3306) |
| MYSQL_USER | MySQL用户 | haircut_user | 是 |
| MYSQL_PASSWORD | MySQL密码 | xxx | 是 |
| REDIS_HOST | Redis地址 | localhost | 是 |
| REDIS_PASSWORD | Redis密码 | xxx | 否(无密码留空) |
| JWT_SECRET | JWT签名密钥 | >=32字符随机串 | **是(必须强)** |
| WECHAT_APP_ID | 微信小程序AppID | wx... | 是 |
| WECHAT_PAY_MCH_ID | 微信支付商户号 | 14xx... | 是 |
| TENCENT_MAP_SDK_KEY | 腾讯地图SDK Key | AKBZ-XXXXXX... | 是 |
| MINIO_ACCESS_KEY | MinIO AccessKey | minioadmin | 是 |
| SITE_DOMAIN | 站点域名 | https://www.haircut.com | 是 |

---

## 附录B：性能指标基线

| 指标 | 目标值 | 监控方式 |
|------|--------|---------|
| API平均响应时间 (P50) | < 100ms | Prometheus + Grafana |
| API响应时间 (P99) | < 500ms | 同上 |
| 错误率 | < 0.1% | 同上 |
| 系统可用性 (SLA) | > 99.9% | UptimeRobot / Pingdom |
| 并发支撑能力 | > 5000 QPS | JMeter压测 |
| 数据库查询耗时 (慢查询) | < 200ms | MySQL Slow Log |
| Redis命中率 | > 85% | Redis INFO stats |

---

**文档维护**: HairCut技术架构组  
**联系方式**: dev@haircut.com  
**最后更新**: 2024-01-20
