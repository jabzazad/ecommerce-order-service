package repositories

import (
	"ecommerce-order/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// OrderDetailRepository repo interface
type OrderDetailRepository interface {
	Create(db *gorm.DB, i interface{}) error
	Update(db *gorm.DB, i interface{}) error
	FindOneObjectByIDUInt(database *gorm.DB, id uint, i interface{}) error
	FindAllByUintIDs(db *gorm.DB, ids []uint, i interface{}) error
	BulkUpsert(db *gorm.DB, uniqueKey string, columns []string, i interface{}, batchSize int) error
	FindAllByOrderID(database *gorm.DB, id uint) ([]*models.OrderDetail, error)
	FindIDsByOrderID(database *gorm.DB, id uint) ([]uint, error)
	FindAllByOrderIDs(database *gorm.DB, ids []uint) ([]*models.OrderDetail, error)
}

type orderDetailRepository struct {
	Repository
}

// OrderDetailNewRepository new sql repository
func OrderDetailNewRepository() OrderDetailRepository {
	return &orderDetailRepository{
		NewRepository(),
	}
}

// FindAllByOrderID find all by order id
func (repo *orderDetailRepository) FindAllByOrderID(database *gorm.DB, id uint) ([]*models.OrderDetail, error) {
	entities := []*models.OrderDetail{}
	err := database.Where("order_id = ?", id).Preload(clause.Associations).Find(&entities).Error
	if err != nil {
		return nil, err
	}

	return entities, nil
}

// FindAllByIDs find all by id
func (repo *orderDetailRepository) FindAllByOrderIDs(database *gorm.DB, ids []uint) ([]*models.OrderDetail, error) {
	entities := []*models.OrderDetail{}
	err := database.Where("order_id in ?", ids).Find(&entities).Error
	if err != nil {
		return nil, err
	}

	return entities, nil
}

// FindIDsByOrderID find all by order id
func (repo *orderDetailRepository) FindIDsByOrderID(database *gorm.DB, id uint) ([]uint, error) {
	entities := []uint{}
	err := database.Select("id").Where("order_id = ?", id).Preload(clause.Associations).Find(&entities).Error
	if err != nil {
		return nil, err
	}

	return entities, nil
}
