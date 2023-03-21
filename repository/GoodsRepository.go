package repository

import (
	"github.com/jinzhu/gorm"
	"xmarket_gin/common"
	"xmarket_gin/model"
)

type IGoodsRepository interface {
	// Create 管理系统创建商品
	Create(name string, category string, price float32) (*model.Goods, error)
	Update(goods model.Goods, price float32) (*model.Goods, error)
	SelectById(id uint) (*model.Goods, error)
	SelectALL() ([]model.Goods, int, error)
	DeleteById(id uint) error
}
type GoodsRepository struct {
	DB *gorm.DB
}

func (g GoodsRepository) Create(name string, category string, price float32) (*model.Goods, error) {
	Goods := model.Goods{
		Name:     name,
		Category: category,
		Price:    price,
	}
	if err := g.DB.Create(&Goods).Error; err != nil {
		return nil, err
	}
	return &Goods, nil
}

func (g GoodsRepository) Update(goods model.Goods, price float32) (*model.Goods, error) {
	if err := g.DB.Model(&goods).Update("price", price).Error; err != nil {
		return nil, err
	}
	return &goods, nil
}

func (g GoodsRepository) SelectById(id uint) (*model.Goods, error) {
	var goods model.Goods
	if err := g.DB.Where("ID = ?", id).Find(&goods).Error; err != nil {
		return nil, err
	}
	return &goods, nil
}

func (g GoodsRepository) SelectALL() ([]model.Goods, int, error) {
	var goods []model.Goods
	if err := g.DB.Find(&goods).Error; err != nil {
		return nil, 0, err
	}
	var total int
	if err := g.DB.Model(model.Goods{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return goods, total, nil
}

func (g GoodsRepository) DeleteById(id uint) error {
	var goods model.Goods
	goods.ID = id
	if err := g.DB.Delete(&goods).Error; err != nil {
		return err
	}
	return nil
}

func NewGoodsRepository() IGoodsRepository {
	db := common.GetDB()
	return GoodsRepository{DB: db}
}
