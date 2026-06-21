# HairCut 连锁理发店平台 - 数据库设计文档

> **版本**: v1.0.0  
> **更新时间**: 2024-01-01  
> **技术栈**: MySQL 8.0 + GORM (Go ORM)  
> **字符集**: utf8mb4_unicode_ci (支持完整Unicode，包括emoji)  
> **表数量**: 38 张核心业务表

---

## 一、ER图（实体关系图）

```
┌─────────────┐       ┌──────────────┐
│    Users     │───1:1─│   Members    │
│  (用户基础)   │       │ (会员信息)    │
└──────┬──────┘       └──────┬───────┘
       │                     │
       │1:N                  │1:N
       ▼                     ▼
┌──────────────┐     ┌────────────────┐
│ UserRoles    │     │ PointsRecords  │
│(用户-角色关联) │     │  (积分变动记录)  │
└──────┬───────┘     └────────────────┘
       │
       │N:1
       ▼
┌──────────────┐     ┌──────────────────┐
│    Roles      │◄─M:N─│   Permissions    │
│   (角色定义)   │     │   (权限点)        │
└──────────────┘     └──────────────────┘

┌─────────────┐       ┌──────────────┐
│    Stores    │───1:N─│ StorePhotos   │
│   (门店)      │       │ (门店照片)     │
├─────────────┤       ├──────────────┤
│ org_id(FK)  │──N:1─│Organizations │
└──────┬──────┘       │(组织架构树)   │
       │             └──────────────┘
       │1:N
       ├──────────────┬──────────────┬──────────────┐
       ▼              ▼              ▼              ▼
┌────────────┐ ┌───────────┐ ┌────────────┐ ┌────────────┐
│  Stylists  │ │ServiceItems│ │Employees   │ │Inventories │
│  (理发师)   │ │ (服务项目)  │ │  (员工)     │ │  (库存)     │
├────────────┤ └─────┬─────┘ ├────────────┤ └──────┬─────┘
│store_id(FK)│       │       │store_id(FK)│        │
└──────┬─────┘       │       └──────┬─────┘        │
       │1:N          │              │1:N           │N:1
       ▼             │              ▼              ▼
┌──────────────┐     │       ┌────────────┐  ┌──────────┐
│StylistPortfolio│    │       │Attendances  │  │Suppliers │
│  (作品集)      │     │       │  (考勤记录)  │  │ (供应商)  │
└──────────────┘     │       └────────────┘  └────┬─────┘
                     │                           │
       ┌─────────────┴───────────────┐            │N:1
       │                             │            ▼
       ▼                             ▼    ┌──────────────┐
┌──────────────┐               ┌────────────┐  │PurchaseOrders│
│Appointments  │───1:1?───────│   Orders    │  │ (采购单)      │
│  (预约记录)   │               │  (订单)     │  ├──────────────┤
├──────────────┤               ├────────────┤  │supplier_id(FK)│
│status(状态机) │               │status(状态机)│  └──────┬───────┘
└──────┬───────┘               └──────┬─────┘         │
       │                              │1:N            │1:N
       │N:1                           ▼               ▼
       │                      ┌────────────┐   ┌──────────────┐
       │                      │ OrderItems  │  │PurchaseItems  │
       │                      │ (订单明细)   │  │ (采购明细)    │
       │                      └────────────┘  └──────────────┘
       │1:N                          │
       ▼                              │
┌──────────────┐                     │N:1
│ PaymentRecords│                    ▼
│  (支付流水)   │             ┌──────────────┐
└──────────────┘             │    Reviews    │
                            │   (评价)      │
                            ├──────────────┤
┌──────────────┐             │order_id(UK)  │
│CouponTemplates│            └──────┬───────┘
│ (优惠券模板)  │                   │1:N
└──────┬───────┘                   ▼
       │1:N                 ┌──────────────┐
       ▼                    │ ReviewMedia   │
┌──────────────┐            │ (评价附件)     │
│   Coupons    │            └──────────────┘
│ (用户优惠券)  │
└──────────────┘

其他辅助表:
- Messages / Notifications (消息通知)
- FanRelations (粉丝关注)
- CustomerProfiles (客户档案)
- Campaigns (营销活动)
- FinanceRecords / StylistCommissions (财务/提成)
- OperationLogs / LoginLogs (日志审计)
- SystemConfigs (系统配置)
```

---

## 二、数据表详细说明

### 表分类统计

| 分类 | 表数量 | 核心表 |
|------|--------|--------|
| **用户与会员** | 3 | users, members, points_records |
| **门店与位置** | 2 | stores, store_photos |
| **理发师** | 2 | stylists, stylist_portfolios |
| **服务项目** | 2 | service_categories, service_items |
| **预约排班** | 3 | appointments, schedules, attendances |
| **订单支付** | 3 | orders, order_items, payment_records |
| **优惠券** | 2 | coupon_templates, coupons |
| **评价** | 2 | reviews, review_media |
| **消息通知** | 2 | messages, notifications |
| **社交关系** | 2 | fan_relations, customer_profiles |
| **组织员工** | 2 | organizations, employees |
| **库存采购** | 4 | inventories, suppliers, purchase_orders, purchase_items |
| **营销活动** | 1 | campaigns |
| **财务** | 2 | finance_records, stylist_commissions |
| **权限系统(RBAC)** | 4 | roles, permissions, role_permissions, user_roles |
| **日志配置** | 3 | operation_logs, login_logs, system_configs |
| **总计** | **38** | - |

---

### 1. 用户相关表

#### 1.1 users (用户表)

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGINT UNSIGNED | PK, AUTO_INCREMENT | 主键ID |
| tenant_id | BIGINT UNSIGNED | INDEX | 租户ID（多租户隔离） |
| openid | VARCHAR(100) | UNIQUE INDEX | 微信OpenID |
| union_id | VARCHAR(100) | INDEX | 微信UnionID（跨应用唯一） |
| session_key | VARCHAR(100) | - | 微信SessionKey（加密存储） |
| phone | VARCHAR(20) | UNIQUE INDEX | 手机号 |
| password | VARCHAR(200) | - | 密码（BCrypt加密） |
| nickname | VARCHAR(50) | NOT NULL | 昵称 |
| avatar_url | TEXT | - | 头像URL |
| gender | TINYINT | DEFAULT 0 | 性别：0未知 1男 2女 |
| birthday | DATE | - | 生日 |
| city_code | VARCHAR(10) | INDEX | 城市代码 |
| status | TINYINT | DEFAULT 1 | 状态：0禁用 1正常 2待审核 |
| created_at | TIMESTAMP | NOT NULL | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL | 更新时间 |
| deleted_at | TIMESTAMP | NULLABLE, INDEX | 软删除时间 |

**索引策略**:
- `UNIQUE idx_openid` - 微信登录唯一性
- `UNIQUE idx_phone` - 手机号唯一注册
- `INDEX idx_tenant` - 多租户查询优化

---

### 2. 门店相关表

#### 2.1 stores (门店表)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | BIGINT UNSIGNED | 主键 |
| tenant_id | BIGINT UNSIGNED | 租户ID |
| org_id | BIGINT INDEX | 组织架构节点ID |
| name | VARCHAR(100) | 门店名称 |
| logo_url | TEXT | Logo图片URL |
| cover_images | JSON | 封面图片列表(JSON数组) |
| province/city/district | VARCHAR(50) | 省/市/区 |
| address | VARCHAR(200) | 详细地址 |
| latitude | DECIMAL(10,7) INDEX | 纬度（腾讯地图坐标）|
| longitude | DECIMAL(10,7) INDEX | 经度 |
| phone | VARCHAR(20) | 联系电话 |
| open_time/close_time | VARCHAR(20) | 营业时间 |
| avg_price | DECIMAL(10,2) | 人均消费价格 |
| rating | DECIMAL(2,1) | 综合评分(1-5星) |
| star_level | TINYINT | 星级(1-5) |
| status | TINYINT INDEX | 状态：0停业 1营业 2筹备中 |

**关键索引**:
- `idx_city` - 按城市筛选门店
- `idx_location` (lat, lng) - 地理位置联合索引（附近门店查询）
- `idx_status` - 状态过滤

---

### 3. 预约状态机

```
待确认(0) → 已确认(1) → 进行中(2) → 已完成(3)
                ↓         ↓
            已取消(4)   爽约(5)
```

**状态转换规则**:
- 待确认 → 可取消（超时自动取消）
- 已确认 → 可取消（提前24小时）
- 进行中 → 不可取消
- 完成 → 终态，不可变更
- 取消/爽约 → 终态

---

### 4. 订单状态机

```
待支付(0) → 已支付(1) → 服务中(2) → 已完成(3)
    ↓           ↓           ↓          ↓
  取消(6)   部分退款(4)   部分退款(4)  全额退款(5)
                         全额退款(5)
```

---

### 5. 权限系统（RBAC）

**角色预设**:

| 角色编码 | 名称 | 权限范围 |
|----------|------|----------|
| super_admin | 超级管理员 | 全部功能+系统设置 |
| admin | 总部管理员 | 运营管理+数据分析 |
| regional_manager | 区域经理 | 区域内所有门店 |
| store_manager | 店长 | 单门店全部管理 |
| stylist | 理发师 | 个人工作台 |
| cs_staff | 客服 | 客户咨询处理 |
| member_user | 会员用户 | 个人操作 |

**权限类型**:
- 菜单权限：控制侧边栏显示
- 按钮权限：控制操作按钮显隐（新增/编辑/删除/导出等）
- 数据权限：按租户/区域/门店过滤数据

---

## 三、多租户设计方案

### 设计原则
1. **共享数据库，共享表结构**（Shared DB + Shared Schema）
2. 每张业务表增加 `tenant_id` 字段
3. 中间件自动注入当前租户上下文
4. 所有查询自动附加 `WHERE tenant_id = ?`

### 隔离策略
```go
// 中间件示例
func TenantMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tenantID := extractTenantID(c) // 从JWT/Header获取
        c.Set("tenant_id", tenantID)
        
        // GORM Scope 自动注入
        db.Session(&gorm.Session{}).Scopes(ByTenant(tenantID))
    }
}

func ByTenant(tenantID uint64) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("tenant_id = ?", tenantID)
    }
}
```

---

## 四、性能优化建议

### 1. 索引策略
- **单列索引**: 外键字段、状态字段、时间字段
- **复合索引**: 高频查询条件组合（如 user_id + status + created_at）
- **地理空间索引**: latitude + longitude 用于附近门店查询

### 2. 分库分表（预留扩展）
- **水平拆分**: 订单表/支付记录表按月分表（orders_202401, orders_202402...）
- **垂直拆分**: 日志类表独立数据库（operation_logs, login_logs）

### 3. 读写分离
- 主库：写操作 + 实时读
- 从库：报表查询 + 数据分析
- 通过GORM插件或中间件实现路由

### 4. 缓存层
- **Redis缓存热点**: 门店列表、理发师信息、服务项目
- **本地缓存(LRU)**: 系统配置、权限数据
- **CDN静态资源**: 图片/视频文件

---

## 五、数据字典（部分关键字段枚举值）

### 性别 (gender)
```
0 = 未知
1 = 男
2 = 女
```

### 会员等级 (member.level)
```
1 = 普通会员（消费满0元）
2 = 银卡会员（消费满2000元）
3 = 金卡会员（消费满5000元）
4 = 黑金会员（消费满15000元）
5 = 钻石会员（消费满30000元）
```

### 预约状态 (appointment.status)
```
0 = 待确认
1 = 已确认
2 = 进行中
3 = 已完成
4 = 已取消
5 = 爽约
```

### 订单状态 (order.status)
```
0 = 待支付
1 = 已支付（待服务）
2 = 服务中
3 = 已完成
4 = 已部分退款
5 = 已全额退款
6 = 已取消
```

### 支付方式 (payment_method)
```
1 = 微信支付
2 = 余额支付
3 = 积分全额兑换
4 = 混合支付（微信+余额+积分）
```

### 优惠券类型 (coupon_template.type)
```
1 = 满减券
2 = 折扣券
3 = 新人券
4 = 节日券
5 = 门店券
6 = 全国通用券
```

---

**文档维护**: 后端开发团队  
**最后更新**: 2024-01-01
