package repositories

import (
	"ecommerce-order/internal/models"
	"ecommerce-order/internal/request"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// OrderRepository repo interface
type OrderRepository interface {
	Create(db *gorm.DB, i interface{}) error
	Update(db *gorm.DB, i interface{}) error
	FindOneObjectByIDUInt(db *gorm.DB, id uint, i interface{}) error
	FindOneByIDFullAssociations(db *gorm.DB, id uint64, i interface{}) error
	BulkUpsert(db *gorm.DB, uniqueKey string, columns []string, i interface{}, batchSize int) error
	FindAll(database *gorm.DB, request *request.FindOrderRequest) (*models.Page, error)
}

type orderRepository struct {
	Repository
}

// orderNewRepository new sql repository
func OrderNewRepository() OrderRepository {
	return &orderRepository{
		NewRepository(),
	}
}

// FindAll find all
func (repo *orderRepository) FindAll(database *gorm.DB, request *request.FindOrderRequest) (*models.Page, error) {
	entities := []*models.Order{}
	page, err := repo.FindAllAndPageInformation(
		repo.query(database, request).Preload(clause.Associations), &request.PageForm, &entities,
	)
	if err != nil {
		return nil, err
	}

	return models.NewPage(page, entities), nil
}

func (repo *orderRepository) query(database *gorm.DB, request *request.FindOrderRequest) *gorm.DB {
	query := database.Table("orders").Where("user_id = ?", request.UserID).Order("created_at desc")
	return query
}
