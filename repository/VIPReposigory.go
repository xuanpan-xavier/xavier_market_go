package repository

import (
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"xmarket_gin/common"
	"xmarket_gin/model"
)

type IVIPRepository interface {
	// Create 管理系统创建用户
	Create(name string, telephone string, password string) (*model.User, error)
	Update(VIP model.User, name string) (*model.User, error)
	UpdateVIPTime(VIP model.User, StartTime model.Time, EndTime model.Time) (*model.User, error)
	SelectByTelephone(telephone string) (*model.User, error)
	SelectALL() ([]model.User, int, error)
	DeleteByTelephone(telephone string) error
	UpdatePoint(VIP model.User, minus uint) (*model.User, error)
}
type VIPRepository struct {
	DB  *gorm.DB
	rds redis.Conn
}

func NewVIPRepository() IVIPRepository {
	return VIPRepository{DB: common.GetDB(), rds: common.GetRedis()}
}

func (V VIPRepository) Create(name string, telephone string, password string) (*model.User, error) {
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	VIP := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
		VipCard:   model.VipCard{Telephone: telephone},
	}
	if err := V.DB.Create(&VIP).Error; err != nil {
		return nil, err
	}
	//if err := V.DB.Create(&VIP.VipCard).Error; err != nil {
	//	return nil, err
	//}
	return &VIP, nil
}

func (V VIPRepository) Update(VIP model.User, name string) (*model.User, error) {
	if err := V.DB.Model(&VIP).Update("name", name).Error; err != nil {
		return nil, err
	}
	//if err := V.DB.Model(&VIP).Update("password", password).Error; err != nil {
	//	return nil, err
	//}
	rdsData, err := V.rds.Do("GET", VIP.Telephone)
	if err != nil {
		return nil, err
	}
	if rdsData != nil {
		_, err := V.rds.Do("DEL", VIP.Telephone)
		if err != nil {
			return nil, err
		}
	}
	return &VIP, nil
}

func (V VIPRepository) UpdateVIPTime(VIP model.User, StartTime model.Time, EndTime model.Time) (*model.User, error) {
	var VIPCard = VIP.VipCard
	if err := V.DB.Model(&VIPCard).Update("StartTime", StartTime).Error; err != nil {
		return nil, err
	}
	if err := V.DB.Model(&VIPCard).Update("EndTime", EndTime).Error; err != nil {
		return nil, err
	}
	//rdsData, err := V.rds.Do("GET", VIP.Telephone)
	//if err != nil {
	//	return nil, err
	//}
	//if rdsData != nil {
	//	_, err := V.rds.Do("DEL", VIP.Telephone)
	//	if err != nil {
	//		return nil, err
	//	}
	//}
	return &VIP, nil
}
func (V VIPRepository) UpdatePoint(VIP model.User, minus uint) (*model.User, error) {

	if err := V.DB.Model(&VIP.VipCard).Update("Point", VIP.VipCard.Point-minus).Error; err != nil {
		return nil, err
	}
	//rdsData, err := V.rds.Do("GET", VIP.Telephone)
	//if err != nil {
	//	return nil, err
	//}
	//if rdsData != nil {
	//	_, err := V.rds.Do("DEL", VIP.Telephone)
	//	if err != nil {
	//		return nil, err
	//	}
	//}
	return &VIP, nil
}

func (V VIPRepository) SelectByTelephone(telephone string) (*model.User, error) {
	var VIP model.User
	//rdsDataI, err := V.rds.Do("Get", telephone)
	//if err != nil {
	//	panic(0)
	//	return nil, err
	//}
	//if rdsDataI == nil {
	if err := V.DB.Preload("VipCard").First(&VIP, telephone).Error; err != nil {
		return nil, err
	}
	//	b, mErr := json.Marshal(&VIP)
	//	if mErr != nil {
	//		panic(2)
	//		return nil, mErr
	//	}
	//	_, mErr = V.rds.Do("SET", telephone, string(b))
	//	if mErr != nil {
	//		panic(3)
	//		return nil, mErr
	//	}
	//} else {
	//	var rdsData string
	//	rdsData, err = redis.String(rdsDataI, err)
	//	if err != nil {
	//		panic(1)
	//		return nil, err
	//	}
	//	b := []byte(rdsData)
	//	err = json.Unmarshal(b, &VIP)
	//	if err != nil {
	//		panic(4)
	//	}
	//}
	return &VIP, nil
}

func (V VIPRepository) SelectALL() ([]model.User, int, error) {
	var VIP []model.User
	if err := V.DB.Preload("VipCard").Find(&VIP).Error; err != nil {
		return nil, 0, err
	}
	var total int
	if err := V.DB.Model(model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return VIP, total, nil
}

func (V VIPRepository) DeleteByTelephone(telephone string) error {
	if err := V.DB.Delete(model.User{}, telephone).Error; err != nil {
		return err
	}
	if err := V.DB.Delete(model.VipCard{}, telephone).Error; err != nil {
		return err
	}
	//rdsData, err := V.rds.Do("GET", telephone)
	//if err != nil {
	//	return err
	//}
	//if rdsData != nil {
	//	_, err := V.rds.Do("DEL", telephone)
	//	if err != nil {
	//		return err
	//	}
	//}

	return nil
}
