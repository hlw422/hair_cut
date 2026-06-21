package mysql

// Role 角色表 - RBAC角色定义
type Role struct {
	BaseModel
	Name        string `json:"name" gorm:"type:varchar(50);not null;uniqueIndex;comment:角色名称"`
	Code        string `json:"code" gorm:"type:varchar(50);uniqueIndex;not null;comment:角色编码(如super_admin/store_manager)"`
	Description string `json:"description" gorm:"type:varchar(200);comment:角色描述"`
	IsSystem    bool   `json:"is_system" gorm:"type:tinyint;default:false;comment:是否系统内置角色(不可删除)"`
	Status      int8   `json:"status" gorm:"type:tinyint;default:1;comment:状态: 0禁用 1启用"`

	// 关联关系
	Permissions  []Permission     `json:"permissions,omitempty" gorm:"many2many:role_permissions;"` // 角色权限多对多
	UserRoles    []UserRole       `json:"-" gorm:"foreignKey:RoleID"`                            // 用户关联(不序列化)
}

func (Role) TableName() string {
	return "roles"
}

// Permission 权限表 - 具体的功能权限点
type Permission struct {
	BaseModel
	Name        string `json:"name" gorm:"type:varchar(100);not null;comment:权限名称"`
	Code        string `json:"code" gorm:"type:varchar(100);uniqueIndex;not null;comment:权限编码(如store:create/order:delete)"`
	Type        int8   `json:"type" gorm:"type:tinyint;default:1;comment:类型: 1菜单 2按钮 3API接口"`
	ParentID    *uint64 `json:"parent_id,omitempty" gorm:"index;comment:父级权限ID(NULL=顶级)"	Path       string `json:"path" gorm:"type:varchar(200);comment:路径(用于前端路由/菜单)"`
	Icon        string `json:"icon" gorm:"type:varchar(50);comment:图标(菜单用)"`
	SortOrder   int    `json:"sort_order" gorm:"type:int;default:0;comment:排序权重"`
	Status      int8   `json:"status" gorm:"type:tinyint;default:1;comment:状态: 0禁用 1启用"`

	// 关联关系
	Parent     *Permission `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children   []Permission `json:"children,omitempty" gorm:"foreignKey:ParentID"` // 子权限
	Roles      []Role       `json:"-" gorm:"many2many:role_permissions;"`         // 角色(不序列化)
}

func (Permission) TableName() string {
	return "permissions"
}

// RolePermission 角色权限关联表（多对多中间表）
type RolePermission struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleID      uint64 `json:"role_id" gorm:"not null;uniqueIndex:idx_role_perm;index;comment:角色ID"`
	PermissionID uint64 `json:"permission_id" gorm:"not null;uniqueIndex:idx_role_perm;index;comment:权限ID"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

// UserRole 用户角色关联表（用户可以拥有多个角色）
type UserRole struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `json:"user_id" gorm:"not null;uniqueIndex:idx_user_role;index;comment:用户ID"`
	RoleID    uint64    `json:"role_id" gorm:"not null;uniqueIndex:idx_user_role;index;comment:角色ID"`
	StoreID   *uint64   `json:"store_id,omitempty" gorm:"comment:门店ID(数据权限范围, NULL=全部)"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;comment:授予时间"`

	// 关联关系
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Role *Role `json:"role,omitempty" gorm:"foreignKey:RoleID"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
