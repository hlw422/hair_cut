package mysql

import (
	"time"
)

// Inventory 库存表 - 门店商品/物料库存
type Inventory struct {
	BaseModel
	StoreID     uint64 `json:"store_id" gorm:"not null;index;comment:所属门店ID"`
	SupplierID  *uint64 `json:"supplier_id,omitempty" gorm:"index;comment:供应商ID"`

	// 商品信息
	Name        string  `json:"name" gorm:"type:varchar(100);not null;comment:商品名称"`
	Category    string  `json:"category" gorm:"type:varchar(50);index;comment:分类(洗护/染烫/工具/耗材)"`
	SKU         string  `json:"sku" gorm:"type:varchar(50);uniqueIndex;comment:SKU编码"`
	Unit        string  `json:"unit" gorm:"type:varchar(20);default:'个';comment:计量单位(瓶/盒/包/个)"`

	// 库存数量
	Quantity    int     `json:"quantity" gorm:"type:int;not null;default:0;comment:当前库存量"`
	MinQuantity int     `json:"min_quantity" gorm:"type:int;default:10;comment:最低库存预警值"`
	MaxQuantity int     `json:"max_quantity" gorm:"type:int;default:500;comment:最高库存上限"`

	// 价格信息
	PurchasePrice float64 `json:"purchase_price" gorm:"type:decimal(10,2);default:0.00;comment:进货单价(元)"`
	SalePrice     float64 `json:"sale_price" gorm:"type:decimal(10,2);default:0.00;comment:销售单价(元)(可选)"`

	// 其他
	Specification string     `json:"specification" gorm:"type:varchar(200);comment:规格说明(容量/尺寸等)"`
	ImageURL      string     `json:"image_url" gorm:"type:text;comment:商品图片"`
	Status        int8       `json:"status" gorm:"type:tinyint;default:1;comment:状态: 0停用 1正常"`
	LastRestockAt *time.Time `json:"last_restock_at,omitempty" gorm:"comment:最近入库时间"`

	// 关联关系
	Store          *Store           `json:"store,omitempty" gorm:"foreignKey:StoreID"`
	Supplier       *Supplier        `json:"supplier,omitempty" gorm:"foreignKey:SupplierID"`
	PurchaseItems  []PurchaseItem   `json:"-" gorm:"foreignKey:InventoryID"` // 采购项关联
}

func (Inventory) TableName() string {
	return "inventories"
}

// IsLowStock 检查是否低库存（低于预警值）
func (inv *Inventory) IsLowStock() bool {
	return inv.Quantity <= inv.MinQuantity
}

// Supplier 供应商表 - 采购供应商管理
type Supplier struct {
	BaseModel
	Name        string `json:"name" gorm:"type:varchar(100);not null;comment:供应商名称"`
	ContactName string `json:"contact_name" gorm:"type:varchar(50);comment:联系人姓名"`
	Phone       string `json:"phone" gorm:"type:varchar(20);comment:联系电话"`
	Address     string `json:"address" gorm:"type:varchar(200);comment:地址"`
	Email       string `json:"email" gorm:"type:varchar(100);comment:邮箱"`
	BankName    string `json:"bank_name" gorm:"type:varchar(50);comment:开户银行"`
	BankAccount string `json:"bank_account" gorm:"type:varchar(30);comment:银行账号"`
	TaxNumber   string `json:"tax_number" gorm:"type:varchar(30);comment:税号"`
	Description string `json:"description" gorm:"type:text;comment:备注/合作说明"`
	Rating      float32 `json:"rating" gorm:"type:decimal(2,1);default:5.0;comment:评分(1-5)"`
	Status      int8   `json:"status" gorm:"type:tinyint;default:1;comment:状态: 0禁用 1合作中 2停止合作"`

	// 关联
	Inventories []Inventory `json:"inventories,omitempty" gorm:"foreignKey:SupplierID"`
	PurchaseOrders []PurchaseOrder `json:"purchase_orders,omitempty" gorm:"foreignKey:SupplierID"` // 采购单
}

func (Supplier) TableName() string {
	return "suppliers"
}

// PurchaseOrder 采购单表 - 门店采购申请与审批
type PurchaseOrder struct {
	BaseModel
	OrderNo    string `json:"order_no" gorm:"type:varchar(64);uniqueIndex;not null;comment:采购单号"`
	StoreID    uint64 `json:"store_id" gorm:"not null;index;comment:申请门店ID"`
	SupplierID uint64 `json:"supplier_id" gorm:"not null;index;comment:供应商ID"`

	// 金额汇总
	TotalAmount float64 `json:"total_amount" gorm:"type:decimal(12,2);default:0.00;comment:采购总金额(元)"`

	// 状态机
	Status int8 `json:"status" gorm:"type:tinyint;default:0;index;comment:状态: 0待提交 1待审批 2已审批 3采购中 4已完成 5已取消 6已驳回"`

	// 时间节点
	SubmittedAt   *time.Time `json:"submitted_at,omitempty" gorm:"comment:提交时间"`
	ApprovedAt    *time.Time `json:"approved_at,omitempty" gorm:"comment:审批时间"`
	ApprovedBy    *uint64    `json:"approved_by,omitempty" gorm:"comment:审批人ID"`
	RejectedReason string     `json:"rejected_reason,omitempty" gorm:"type:varchar(200);comment:驳回原因"`
	ExpectedDate  *time.Time `json:"expected_date,omitempty" gorm:"type:date;comment:预计到货日期"`
	ReceivedAt    *time.Time `json:"received_at,omitempty" gorm:"comment:实际到货时间"`

	// 备注
	ApplicantRemark string `json:"applicant_remark" gorm:"type:text;comment:申请人备注"`
	ApproverRemark string `json:"approver_remark" gorm:"type:text;comment:审批人备注"`

	// 关联关系
	Store    *Store          `json:"store,omitempty" gorm:"foreignKey:StoreID"`
	Supplier *Supplier       `json:"supplier,omitempty" gorm:"foreignKey:SupplierID"`
	Items    []PurchaseItem  `json:"items,omitempty" gorm:"foreignKey:PurchaseOrderID"` // 采购项列表
}

func (PurchaseOrder) TableName() string {
	return "purchase_orders"

const (
	PurchaseStatusDraft       = 0 // 待提交
	PurchaseStatusPending     = 1 // 待审批
	PurchaseStatusApproved    = 2 // 已审批
	PurchaseStatusPurchasing  = 3 // 采购中
	PurchaseStatusCompleted   = 4 // 已完成
	PurchaseStatusCancelled   = 5 // 已取消
	PurchaseStatusRejected    = 6 // 已驳回
)
}

// PurchaseItem 采购单项表 - 采购单中的具体商品明细
type PurchaseItem struct {
	BaseModel
	PurchaseOrderID uint64 `json:"purchase_order_id" gorm:"not null;index;comment:采购单ID"`
	InventoryID     uint64 `json:"inventory_id" gorm:"not null;comment:库存商品ID"`
	ProductName     string `json:"product_name" gorm:"type:varchar(100);not null;comment:商品名称(冗余)"`
	Specification   string `json:"specification" gorm:"type:varchar(200);comment:规格"`
	Unit            string `json:"unit" gorm:"type:varchar(20);comment:单位"`
	Quantity        int     `json:"quantity" gorm:"type:int;not null;comment:采购数量"`
	UnitPrice       float64 `json:"unit_price" gorm:"type:decimal(10,2);not null;comment:采购单价(元)"`
	TotalPrice      float64 `json:"total_price" gorm:"type:decimal(10,2);not null;comment:小计金额(元)"`
	ReceivedQty     int     `json:"received_qty" gorm:"type:int;default:0;comment:已收货数量"`

	// 关联
	PurchaseOrder *PurchaseOrder `json:"purchase_order,omitempty" gorm:"foreignKey:PurchaseOrderID"`
	Inventory     *Inventory     `json:"inventory,omitempty" gorm:"foreignKey:InventoryID"`
}

func (PurchaseItem) TableName() string {
	return "purchase_items"
}
