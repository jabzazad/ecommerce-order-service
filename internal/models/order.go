package models

import (
	"time"
)

const (
	// OrderPrefix order prefix
	OrderPrefix = "OR"
	// TimeFormat time format for running number
	TimeFormat = "20060102150405"
)

// OrderStatus order status
type OrderStatus uint

const (
	// OrderStatusUnknown order status unknown
	OrderStatusUnknown OrderStatus = iota
	// OrderStatusDraft order status draft
	OrderStatusDraft
	// OrderStatusConfirm order status confirm
	OrderStatusConfirm
	// OrderStatusPick order status pick
	OrderStatusPick
	// OrderStatusPack order status pack
	OrderStatusPack
	// OrderStatusShip order status ship
	OrderStatusShip
	// OrderStatusCancel order status cancel
	OrderStatusCancel
	// OrderStatusComplete order status complete
	OrderStatusComplete
	// OrderStatusReturn order status return
	OrderStatusReturn
)

// Type type
func (t OrderStatus) Type() *Type {
	name := "-"
	switch t {
	case OrderStatusDraft:
		name = "สร้างออเดอร์"
	case OrderStatusConfirm:
		name = "ยืนยันออเดอร์"
	case OrderStatusPick:
		name = "หยิบสินค้าจากคลัง"
	case OrderStatusPack:
		name = "แพ็คสินค้า"
	case OrderStatusShip:
		name = "จัดส่งสินค้า"
	case OrderStatusCancel:
		name = "ยกเลิก"
	case OrderStatusComplete:
		name = "สำเร็จ"
	}

	return &Type{
		Name:  name,
		Value: int(t),
	}
}

// ShippingMethod shipping method
type ShippingMethod uint

const (
	// ShippingMethodUnknown shipping method unknown
	ShippingMethodUnknown ShippingMethod = iota
	// Kerry shipping method kerry
	Kerry
	// Flash shipping method flash
	Flash
	// Thaipost shipping method thaipost
	Thaipost
	// ScgExpress scg express
	ScgExpress
	// NinjaVan ninja van
	NinjaVan
	// NimExpress nim express
	NimExpress
	// JAndT j&t
	JAndT
	// InterExpress inter express
	InterExpress
	// JWDExpress jwd express
	JWDExpress
	// BestExpress best express
	BestExpress
	// DHLExpress dhl express
	DHLExpress
	// TrueELogistic true express
	TrueELogistic
	// Lalamove la la move
	Lalamove
	// CJLogistics cj logistics
	CJLogistics
	// TP2 thaipost register
	TP2
	// JWDC JWD Chilled Express
	JWDC
	// JWDF JWD Frozen Express
	JWDF
	// SCGF SCG Yamato Express Frozen
	SCGF
	// SCGC SCG Yamato Express Chilled
	SCGC
	// SelfService self service
	SelfService
	// JAndTPickUp j&t pickup
	JAndTPickUp
	// JAndTDropoff j&t dropoff
	JAndTDropoff
	// NinjaPickup Ninja pick up
	NinjaPickup
	// NinjaDropoff Ninja drop off
	NinjaDropoff
)

// Type type
func (t ShippingMethod) Type() *Type {
	name := "-"
	switch t {
	case Kerry:
		name = "Kerry"

	case Flash:
		name = "Flash"

	case Thaipost:
		name = "Thaipost EMS"

	case ScgExpress:
		name = "Scg Express"

	case NinjaVan:
		name = "NinjaVan"

	case NimExpress:
		name = "Nim Express"

	case JAndT:
		name = "J&T"

	case InterExpress:
		name = "Inter Express"

	case JWDExpress:
		name = "JWD Express"

	case BestExpress:
		name = "Best Express"

	case DHLExpress:
		name = "DHL"

	case TrueELogistic:
		name = "True e-Logistics"

	case TP2:
		name = "Thaipost register"

	case JWDC:
		name = "JWD Chilled Express"

	case JWDF:
		name = "JWD Frozen Express"

	case SCGF:
		name = "SCG Yamato Express Frozen"

	case SCGC:
		name = "SCG Yamato Express Chilled"
	}

	return &Type{
		Name:  name,
		Value: int(t),
	}
}

// Order order struct
type Order struct {
	Model
	Message
	Status          OrderStatus    `json:"status"`
	PaymentAt       *time.Time     `json:"payment_at,omitempty"`
	CreatedByUserID uint           `json:"created_by_user_id"`
	UpdatedByUserID *uint          `json:"updated_by_user_id,omitempty"`
	DeletedByUserID *uint          `json:"deleted_by_user_id,omitempty"`
	OrderNumber     string         `json:"order_number"`
	AddressID       string         `json:"address_id"`
	TotalPrice      float64        `json:"total_price"`
	NetPrice        float64        `json:"net_price"`
	TotalQuantity   float64        `json:"total_quantity"`
	Discount        float64        `json:"discount"`
	ShippingPrice   float64        `json:"shipping_price"`
	TrackingCode    string         `json:"tracking_code,omitempty"`
	VoucherID       uint           `json:"voucher_id,omitempty"`
	UserID          uint           `json:"user_id"`
	OrderDetails    []*OrderDetail `json:"order_details,omitempty" gorm:"foreignKey:OrderID;references:ID" copier:"-"`
}

// OrderDetail order detail struct
type OrderDetail struct {
	Model
	OrderID         uint     `json:"order_id"`
	ProductID       uint     `json:"product_id"`
	Quantity        float64  `json:"quantity"`
	Price           float64  `json:"price"`
	CreatedByUserID uint     `json:"created_by_user_id"`
	UpdatedByUserID *uint    `json:"updated_by_user_id,omitempty"`
	DeletedByUserID *uint    `json:"deleted_by_user_id,omitempty"`
	Product         *Product `json:"product,omitempty" gorm:"ForeignKey:ProductID;references:ID"`
}
