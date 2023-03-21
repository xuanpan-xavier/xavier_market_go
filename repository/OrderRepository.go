package repository

import (
	"github.com/jinzhu/gorm"
	"xmarket_gin/common"
	"xmarket_gin/model"
)

type IOrderRepository interface {
	// Create 管理系统创建商品
	Create(telephone string, goods []model.Goods, total float32, num string) (*model.Order, error)
	Update(order model.Order, goods []model.Goods, total float32) (*model.Order, error)
	SelectById(id uint) (*model.Order, error)
	SelectALL() ([]model.Order, int, error)
	DeleteById(id uint) error
}
type OrderRepository struct {
	DB *gorm.DB
}

func (o OrderRepository) Create(telephone string, goods []model.Goods, total float32, num string) (*model.Order, error) {
	var Order = model.Order{
		Telephone: telephone,
		Goods:     goods,
		Total:     total,
		Number:    num,
	}
	if err := o.DB.Create(&Order).Error; err != nil {
		return nil, err
	}
	return &Order, nil
}

func (o OrderRepository) Update(order model.Order, goods []model.Goods, total float32) (*model.Order, error) {
	if err := o.DB.Model(&order).Association("Goods").Replace(&goods).Error; err != nil {
		return nil, err
	}
	if err := o.DB.Model(&order).Update("Total", total).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (o OrderRepository) SelectById(id uint) (*model.Order, error) {
	var order model.Order
	if err := o.DB.Model(&order).Where("ID = ?", id).Preload("Goods").Find(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (o OrderRepository) SelectALL() ([]model.Order, int, error) {
	var order []model.Order
	if err := o.DB.Preload("Goods").Find(&order).Error; err != nil {
		return nil, 0, err
	}
	var total int
	if err := o.DB.Model(model.Order{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return order, total, nil
}

func (o OrderRepository) DeleteById(id uint) error {
	var order model.Order
	order.ID = id
	if err := o.DB.Delete(&order).Error; err != nil {
		return err
	}
	return nil
}

func NewOrderRepository() IOrderRepository {
	db := common.GetDB()
	return OrderRepository{DB: db}
}
