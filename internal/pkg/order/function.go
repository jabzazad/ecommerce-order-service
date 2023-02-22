package order

import (
	"ecommerce-order/internal/models"
	"sync"
)

// MapProductToOrderDetail map product to order detail
func MapProductToOrderDetails(orderDetails <-chan []*models.OrderDetail, products []*models.Product, wg *sync.WaitGroup) {
	for orderDetails := range orderDetails {
		for _, orderDetail := range orderDetails {
			for _, product := range products {
				if orderDetail.ProductID == product.ID {
					orderDetail.Product = product
				}
			}
			wg.Done()
		}
	}
}

// MapProductToOrderDetail map product to order detail
func MapProductToOrderDetail(orderDetail <-chan *models.OrderDetail, products []*models.Product, wg *sync.WaitGroup) {
	for orderDetail := range orderDetail {
		for _, product := range products {
			if orderDetail.ProductID == product.ID {
				orderDetail.Product = product
			}
		}
		wg.Done()
	}
}
