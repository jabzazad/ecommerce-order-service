package order

import (
	"ecommerce-order/internal/core/config"
	"ecommerce-order/internal/core/context"
	"ecommerce-order/internal/models"
	"ecommerce-order/internal/pkg/client"
	"ecommerce-order/internal/repositories"
	"ecommerce-order/internal/request"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/imroc/req"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Service service interface
type Service interface {
	Create(c *context.Context, request *request.CreateOrderRequest) (*models.Order, error)
	Delete(c *context.Context, request *request.GetOne) error
	FindAll(c *context.Context, request *models.PageForm) (*models.Page, error)
	FindOne(c *context.Context, request *request.GetOne) (*models.Order, error)
}

type service struct {
	config                *config.Configs
	result                *config.ReturnResult
	repository            repositories.OrderRepository
	orderDetailRepository repositories.OrderDetailRepository
	clientService         client.Service
	mutex                 sync.Mutex
}

// NewService new service
func NewService() Service {
	return &service{
		config:                config.CF,
		result:                config.RR,
		repository:            repositories.OrderNewRepository(),
		orderDetailRepository: repositories.OrderDetailNewRepository(),
		clientService:         client.NewService(),
	}
}

// Create create
func (s *service) Create(c *context.Context, form *request.CreateOrderRequest) (*models.Order, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	db := c.GetDatabase()
	productIDs := []uint{}
	for _, product := range form.OrderDetails {
		productIDs = append(productIDs, product.ProductID)
	}

	header := req.Header{
		"Accept-Language": c.AcceptsLanguages(),
	}

	param := req.Param{
		"product_ids": productIDs,
	}

	url := fmt.Sprintf("%s%s/ids", s.config.Product.URL, s.config.Product.Path.Product)
	response := []*models.Product{}
	err := s.clientService.GetRequest(url, header, param, &response)
	if err != nil {
		return nil, err
	}

	for _, od := range form.OrderDetails {
		for _, product := range response {
			if product.ID == od.ProductID && product.AvailableStock < od.Quantity {
				return nil, s.result.OutOfStock
			}
		}
	}

	order := &models.Order{}
	order.OrderNumber = fmt.Sprintf("%s%s", models.OrderPrefix, time.Now().Format(models.TimeFormat))
	_ = copier.Copy(order, form)
	order.Status = models.OrderStatusDraft
	order.CreatedByUserID = c.GetUserID()
	err = s.repository.Create(db, order)
	if err != nil {
		logrus.Errorf("create order error: %s", err)
		return nil, err
	}

	orderDetails := []*models.OrderDetail{}
	for _, orderDetail := range form.OrderDetails {
		for _, product := range response {
			if product.ID == orderDetail.ProductID {
				if product.ProductID > 0 {
					for _, productMaster := range response {
						if product.ProductID == productMaster.ID {
							productMaster.AvailableStock -= orderDetail.Quantity
						}
					}
				}

				product.TotalSell += orderDetail.Quantity
				product.AvailableStock -= orderDetail.Quantity
			}

			order.TotalPrice += (product.Price * float64(orderDetail.Quantity))
			order.TotalQuantity += orderDetail.Quantity
			orderDetails = append(orderDetails, &models.OrderDetail{
				OrderID:         order.ID,
				ProductID:       orderDetail.ProductID,
				Quantity:        orderDetail.Quantity,
				Price:           product.Price,
				CreatedByUserID: c.GetUserID(),
			})
		}

	}

	err = s.repository.Update(db, order)
	if err != nil {
		logrus.Errorf("update order error: %s", err)
		return nil, err
	}

	err = s.orderDetailRepository.BulkUpsert(db, "id", []string{}, orderDetails, 100)
	if err != nil {
		logrus.Errorf("bulk order details error: %s", err)
		return nil, err
	}

	url = fmt.Sprintf("%s%s/bulk", s.config.Product.URL, s.config.Product.Path.Product)
	body := &request.BulkUpdateProducts{
		Products: response,
	}

	messageResponse := &models.Message{}
	err = s.clientService.PostRequest(url, header, nil, body, messageResponse)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// Delete delete
func (s *service) Delete(c *context.Context, form *request.GetOne) error {
	db := c.GetDatabase()
	tn := time.Now()
	userID := c.GetUserID()
	order := &models.Order{}
	err := s.repository.FindOneObjectByIDUInt(db, form.ID, order)
	if err != nil {
		logrus.Errorf("find order by id=%d error: %s", form.ID, err)
		return s.result.Internal.DatabaseNotFound
	}

	if order.CreatedByUserID != userID {
		return s.result.InvalidPermissionRole
	}

	orderDetails, err := s.orderDetailRepository.FindAllByOrderID(db, order.ID)
	if err != nil {
		return s.result.Internal.DatabaseNotFound
	}

	productIDs := []uint{}
	for _, orderDetail := range orderDetails {
		productIDs = append(productIDs, orderDetail.ProductID)
	}

	header := req.Header{
		"Accept-Language": c.AcceptsLanguages(),
	}

	param := req.Param{
		"product_ids": productIDs,
	}

	url := fmt.Sprintf("%s%s/ids", s.config.Product.URL, s.config.Product.Path.Product)
	response := []*models.Product{}
	err = s.clientService.GetRequest(url, header, param, &response)
	if err != nil {
		return err
	}

	for _, orderDetail := range orderDetails {
		for _, product := range response {
			if product.ID == orderDetail.ProductID {
				if product.ProductID > 0 {
					for _, productMaster := range response {
						if productMaster.ID == product.ProductID {
							productMaster.TotalSell -= orderDetail.Quantity
							productMaster.AvailableStock += orderDetail.Quantity
						}
					}
				}

				product.TotalSell -= orderDetail.Quantity
				product.AvailableStock += orderDetail.Quantity
			}

		}

		orderDetail.DeletedByUserID = &userID
		orderDetail.DeletedAt = gorm.DeletedAt{
			Time:  tn,
			Valid: true,
		}
	}

	order.DeletedAt = gorm.DeletedAt{
		Time:  tn,
		Valid: true,
	}

	order.Status = models.OrderStatusCancel
	order.DeletedByUserID = &userID
	err = s.repository.Update(db, order)
	if err != nil {
		return err
	}

	err = s.orderDetailRepository.BulkUpsert(db, "id", []string{}, orderDetails, 100)
	if err != nil {
		return err
	}

	url = fmt.Sprintf("%s%s/bulk", s.config.Product.URL, s.config.Product.Path.Product)
	body := &request.BulkUpdateProducts{
		Products: response,
	}

	messageResponse := &models.Message{}
	err = s.clientService.PostRequest(url, header, nil, body, messageResponse)
	if err != nil {
		return err
	}

	return nil
}

// FindAll find all
func (s *service) FindAll(c *context.Context, form *models.PageForm) (*models.Page, error) {
	db := c.GetDatabase()
	request := &request.FindOrderRequest{
		UserID: c.GetUserID(),
	}

	page, err := s.repository.FindAll(db, request)
	if err != nil {
		message := fmt.Sprintf("find order error: %s", err)
		logrus.Errorf(message)
		return nil, s.result.Internal.DatabaseNotFound
	}

	productIDs := []uint{}
	totalLength := len(page.Entities.([]*models.Order))
	for _, order := range page.Entities.([]*models.Order) {
		if order.OrderDetails != nil {
			for _, orderDetail := range order.OrderDetails {
				productIDs = append(productIDs, orderDetail.ProductID)
				totalLength++
			}
		}
	}

	header := req.Header{
		"Accept-Language": c.AcceptsLanguages(),
	}

	param := req.Param{
		"product_ids": productIDs,
	}

	url := fmt.Sprintf("%s%s/ids", s.config.Product.URL, s.config.Product.Path.Product)
	response := []*models.Product{}
	err = s.clientService.GetRequest(url, header, param, &response)
	if err != nil {
		return nil, err
	}

	jobs := make(chan []*models.OrderDetail, totalLength)
	wg := new(sync.WaitGroup)
	for w := 0; w < runtime.NumCPU()*2; w++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			MapProductToOrderDetails(jobs, response, wg)
		}()
	}

	for _, order := range page.Entities.([]*models.Order) {
		jobs <- order.OrderDetails
	}
	close(jobs)
	wg.Wait()

	return page, nil
}

// FindOne find one by id
func (s *service) FindOne(c *context.Context, request *request.GetOne) (*models.Order, error) {
	db := c.GetDatabase()
	order := &models.Order{}
	err := s.repository.FindOneObjectByIDUInt(db, request.ID, order)
	if err != nil {
		logrus.Errorf("find order by id=%d error: %s", request.ID, err)
		return nil, s.result.Internal.DatabaseNotFound
	}

	if order.CreatedByUserID != c.GetUserID() {
		return nil, s.result.InvalidPermissionRole
	}

	productIDs := []uint{}
	if order.OrderDetails != nil {
		for _, orderDetail := range order.OrderDetails {
			productIDs = append(productIDs, orderDetail.ProductID)
		}
	}

	header := req.Header{
		"Accept-Language": c.AcceptsLanguages(),
	}

	param := req.Param{
		"product_ids": productIDs,
	}

	url := fmt.Sprintf("%s%s/ids", s.config.Product.URL, s.config.Product.Path.Product)
	response := []*models.Product{}
	err = s.clientService.GetRequest(url, header, param, &response)
	if err != nil {
		return nil, err
	}

	jobs := make(chan *models.OrderDetail, len(order.OrderDetails))
	wg := new(sync.WaitGroup)
	for w := 0; w < runtime.NumCPU()*2; w++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			MapProductToOrderDetail(jobs, response, wg)
		}()
	}

	for _, orderDetail := range order.OrderDetails {
		jobs <- orderDetail
	}
	close(jobs)
	wg.Wait()

	return order, nil
}
