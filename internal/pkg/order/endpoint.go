// Package order is a generated GoMock package.
package order

import (
	"ecommerce-order/internal/core/config"
	"ecommerce-order/internal/handlers"
	"ecommerce-order/internal/models"
	"ecommerce-order/internal/request"

	"github.com/gofiber/fiber/v2"
)

// Endpoint endpoint interface
type Endpoint interface {
	Create(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	FindAll(c *fiber.Ctx) error
	FindOne(c *fiber.Ctx) error
}

type endpoint struct {
	config  *config.Configs
	result  *config.ReturnResult
	service Service
}

// NewEndpoint new endpoint
func NewEndpoint() Endpoint {
	return &endpoint{
		config:  config.CF,
		result:  config.RR,
		service: NewService(),
	}
}

// Create create
// @Tags Order
// @Summary Create
// @Description Create
// @Accept json
// @Produce json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Param request body request.CreateOrderRequest true "request body"
// @Success 200 {object} models.Order
// @Failure 400 {object} models.Message
// @Failure 401 {object} models.Message
// @Failure 404 {object} models.Message
// @Failure 410 {object} models.Message
// @Security ApiKeyAuth
// @Router /orders [post]
func (ep *endpoint) Create(c *fiber.Ctx) error {
	return handlers.ResponseObject(c, ep.service.Create, &request.CreateOrderRequest{})
}

// Delete delete
// @Tags Order
// @Summary Delete
// @Description Delete
// @Accept json
// @Produce json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Param id path uint true "ID"
// @Param request body request.GetOneAndPermission true "query permission"
// @Success 200 {object} models.Message
// @Failure 400 {object} models.Message
// @Failure 401 {object} models.Message
// @Failure 404 {object} models.Message
// @Failure 410 {object} models.Message
// @Security ApiKeyAuth
// @Router /orders/{id} [delete]
func (ep *endpoint) Delete(c *fiber.Ctx) error {
	return handlers.ResponseSuccess(c, ep.service.Delete, &request.GetOneAndPermission{})
}

// FindAll find all
// @Tags Order
// @Summary FindAll
// @Description FindAll
// @Accept json
// @Produce json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Param request query models.PageForm true "query for get all"
// @Success 200 {object} models.Page
// @Failure 400 {object} models.Message
// @Failure 401 {object} models.Message
// @Failure 404 {object} models.Message
// @Failure 410 {object} models.Message
// @Security ApiKeyAuth
// @Router /orders [get]
func (ep *endpoint) FindAll(c *fiber.Ctx) error {
	return handlers.ResponseObject(c, ep.service.FindAll, &models.PageForm{})
}

// FindOne find one
// @Tags Order
// @Summary FindOne
// @Description FindOne
// @Accept json
// @Produce json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Param id path uint true "ID"
// @Param request query request.GetOneAndPermission true "query permission"
// @Success 200 {object} models.Order
// @Failure 400 {object} models.Message
// @Failure 401 {object} models.Message
// @Failure 404 {object} models.Message
// @Failure 410 {object} models.Message
// @Security ApiKeyAuth
// @Router /orders/{id} [get]
func (ep *endpoint) FindOne(c *fiber.Ctx) error {
	return handlers.ResponseObject(c, ep.service.FindOne, &request.GetOneAndPermission{})
}
