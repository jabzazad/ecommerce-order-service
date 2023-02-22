package request

import "ecommerce-order/internal/models"

// CreateOrderRequest create order request
type CreateOrderRequest struct {
	AddressID     string  `json:"address_id"`
	TotalPrice    float64 `json:"total_price"`
	NetPrice      float64 `json:"net_price"`
	TotalQuantity float64 `json:"total_quantity"`
	Discount      float64 `json:"discount"`
	ShippingPrice float64 `json:"shipping_price"`
	TrackingCode  string  `json:"tracking_code,omitempty"`
	VoucherID     uint    `json:"voucher_id,omitempty"`
	OrderDetails  []struct {
		ProductID uint    `json:"product_id"`
		Quantity  float64 `json:"quantity"`
	} `json:"order_details"`
}

// FindOrderRequest find order request
type FindOrderRequest struct {
	UserID uint `json:"user_id"`
	models.PageForm
}
