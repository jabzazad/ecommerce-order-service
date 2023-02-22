package request

import "ecommerce-order/internal/models"

type BulkUpdateProducts struct {
	Products []*models.Product `json:"products"`
}
