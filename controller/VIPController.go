package controller

// 会员管理系统内的用户管理
import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"xmarket_gin/common"
	"xmarket_gin/model"
	"xmarket_gin/repository"
	"xmarket_gin/response"
	"xmarket_gin/vo"
)

type IVIPController interface {
	RestController
	ShowALL(c *gin.Context)
	UpdateVIP(c *gin.Context)
	User2VIP(c *gin.Context)
	UpdatePoint(c *gin.Context)
	ReadVouchers(c *gin.Context)
	UpdateVouchers(c *gin.Context)
	ReadVouchersByTelephone(c *gin.Context)
}

type VIPController struct {
	Repository repository.IVIPRepository
}

func NewVIPController() IVIPController {
	VIPController := VIPController{Repository: repository.NewVIPRepository()}
	VIPController.Repository.(repository.VIPRepository).DB.AutoMigrate(model.User{}, model.VipCard{})
	return VIPController
}

// Create 管理系统新建用户
func (V VIPController) Create(c *gin.Context) {
	var requestVIP vo.CreateVIPRequest
	if err := c.ShouldBind(&requestVIP); err != nil {
		response.Fail(c, nil, "数据验证错误，请填写完整信息")
		return
	}

	VIP, err := V.Repository.Create(requestVIP.Name, requestVIP.Telephone, requestVIP.Password)
	if err != nil {
		//panic(err)
		response.Fail(c, nil, "创建失败")
		return
	}

	response.Success(c, gin.H{"VIP": VIP}, "nil")
}

func (V VIPController) Delete(c *gin.Context) {
	// 获取path中参数
	if err := V.Repository.DeleteByTelephone(c.Params.ByName("Telephone")); err != nil {
		response.Fail(c, nil, "删除失败")
		return
	}
	response.Success(c, nil, "")
}

func (V VIPController) Update(c *gin.Context) {
	var requestVIP vo.CVIPRequest
	if err := c.ShouldBind(&requestVIP); err != nil {
		response.Fail(c, nil, "数据验证错误，请填写完整信息")
		return
	}
	//hasedPassword, err := bcrypt.GenerateFromPassword([]byte(requestVIP.Password), bcrypt.DefaultCost)
	// 获取path中参数
	updateVIP, err := V.Repository.SelectByTelephone(c.Params.ByName("Telephone"))
	if err != nil {
		response.Fail(c, gin.H{"err": err}, "用户不存在")
		return
	}
	// 更新
	VIP, err := V.Repository.Update(*updateVIP, requestVIP.Name)
	if err != nil {
		panic(err)
	}
	response.Success(c, gin.H{"VIP": VIP}, "修改成功")

}

func (V VIPController) UpdateVIP(c *gin.Context) {
	var requestVIP vo.UpdateVIPRequest
	if err := c.ShouldBind(&requestVIP); err != nil {
		response.Fail(c, nil, "数据验证错误，请填写完整信息")
		panic(err)
		return
	}
	// 获取path中参数
	updateVIP, err := V.Repository.SelectByTelephone(c.Params.ByName("Telephone"))
	if err != nil {
		response.Fail(c, gin.H{"err": err}, "用户不存在")
		return
	}
	// 更新
	VIP, err := V.Repository.UpdateVIPTime(*updateVIP, requestVIP.StartTime, requestVIP.EndTime)
	if err != nil {
		panic(err)
	}
	response.Success(c, gin.H{"VIP": VIP}, "修改成功")

}

func (V VIPController) User2VIP(c *gin.Context) {
	var user vo.User2VIPRequest
	if err := c.ShouldBind(&user); err != nil {
		response.Fail(c, nil, "数据验证错误，请填写完整信息")
		panic(err)
		return
	}
	// 获取path中参数
	updateVIP, err := V.Repository.SelectByTelephone(c.Params.ByName("Telephone"))
	if err != nil {
		response.Fail(c, gin.H{"err": err}, "用户不存在")
		return
	}
	// 更新
	var updateTime model.Time
	y, _ := strconv.ParseInt(user.Year, 10, 64)
	if user.IsVIP == "1" {
		updateTime = updateVIP.VipCard.EndTime.AddDate(int(y), 0, 0)
		VIP, err := V.Repository.UpdateVIPTime(*updateVIP, updateVIP.VipCard.StartTime, updateTime)
		if err != nil {
			panic(err)
		}
		response.Success(c, gin.H{"VIP": VIP}, "修改成功")
	} else {
		updateTime = model.Time(time.Now().AddDate(int(y), 0, 0))
		VIP, err := V.Repository.UpdateVIPTime(*updateVIP, model.Time(time.Now()), updateTime)
		if err != nil {
			panic(err)
		}
		response.Success(c, gin.H{"VIP": VIP}, "修改成功")
	}
}

func (V VIPController) UpdatePoint(c *gin.Context) {
	var minus vo.UpdatePointRequest
	if err := c.ShouldBind(&minus); err != nil {
		response.Fail(c, nil, "数据验证错误，请填写完整信息")
		panic(err)
		return
	}
	// 获取path中参数
	updateVIP, err := V.Repository.SelectByTelephone(minus.Telephone)
	if err != nil {
		response.Fail(c, gin.H{"err": err}, "用户不存在")
		return
	}
	// 更新
	VIP, err := V.Repository.UpdatePoint(*updateVIP, minus.Minus)
	if err != nil {
		panic(err)
	}
	voucher := model.Voucher{
		Telephone: VIP.Telephone,
		Point:     minus.Minus,
		IsUsed:    0,
	}
	DB := common.GetDB()
	if err := DB.Create(&voucher).Error; err != nil {
		response.Fail(c, gin.H{"err": err}, "代金券兑换失败")
		return
	}

	response.Success(c, gin.H{"VIP": VIP}, "修改成功")
}

func (V VIPController) Show(c *gin.Context) {
	// 获取path中参数
	VIP, err := V.Repository.SelectByTelephone(c.Params.ByName("Telephone"))
	if err != nil {
		response.Fail(c, gin.H{"err": err}, "用户不存在")
		return
	}
	var isVIP string
	timeNow := time.Now().Format("2006-01-02")
	timeVIP := VIP.VipCard.EndTime.String()
	year1, _ := strconv.ParseUint(timeNow[0:4], 10, 64)
	year2, _ := strconv.ParseUint(timeVIP[0:4], 10, 64)
	month1, _ := strconv.ParseUint(timeNow[5:7], 10, 64)
	month2, _ := strconv.ParseUint(timeVIP[5:7], 10, 64)
	day1, _ := strconv.ParseUint(timeNow[8:10], 10, 64)
	day2, _ := strconv.ParseUint(timeVIP[8:10], 10, 64)
	if year1 < year2 || month1 < month2 || day1 < day2 {
		isVIP = "1"
	} else {
		isVIP = "0"
	}
	response.Success(c, gin.H{"VIP": VIP, "IsVIP": isVIP}, "")

}

func (V VIPController) ShowALL(c *gin.Context) {
	query := c.DefaultQuery("query", "")
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	DB := common.GetDB()
	var VIP []model.User
	var total int
	if query == "" {
		DB.Order("Telephone").Offset((pageNum - 1) * pageSize).Limit(pageSize).Preload("VipCard").Find(&VIP)
		DB.Model(model.User{}).Count(&total)
	} else {
		DB.Order("Telephone").Where("Telephone LIKE ?", query).Preload("VipCard").Find(&VIP)
		DB.Model(model.User{}).Where("Telephone LIKE ?", query).Count(&total)
	}
	response.Success(c, gin.H{"data": VIP, "total": total}, "查询成功")
}

func (V VIPController) ReadVouchers(c *gin.Context) {
	query := c.DefaultQuery("query", "")
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	DB := common.GetDB()
	var vouchers []model.Voucher
	var total int
	if query == "" {
		DB.Order("ID").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&vouchers)
		DB.Model(model.Voucher{}).Count(&total)
	} else {
		DB.Order("id").Where("Telephone = ?", query).Find(&vouchers)
	}
	response.Success(c, gin.H{"data": vouchers, "total": total}, "查询成功")
}

func (V VIPController) ReadVouchersByTelephone(c *gin.Context) {
	DB := common.GetDB()
	// 获取path中参数
	var v []model.Voucher
	if err := DB.Order("id").Where("Telephone = ?", c.Params.ByName("Telephone")).Find(&v).Error; err != nil {
		response.Fail(c, gin.H{"err": err}, "获取代金券信息失败")
		return
	}
	response.Success(c, gin.H{"data": v}, "查询成功")
}
func (V VIPController) UpdateVouchers(c *gin.Context) {
	//var voucher vo.UpdateVouchersRequest
	//if err := c.ShouldBind(&voucher); err != nil {
	//	response.Fail(c, nil, "数据验证错误，请填写完整信息")
	//	panic(err)
	//	return
	//}
	//
	//uid, err := uuid.FromString(voucher.ID)
	//if err != nil {

	//	response.Fail(c, gin.H{"err": err, "uid": voucher.ID}, "转换uuid失败")
	//	return
	//}
	DB := common.GetDB()
	// 获取path中参数
	var v model.Voucher
	vid := c.Params.ByName("ID")
	if err := DB.Where("ID = ?", vid).First(&v).Error; err != nil {
		response.Fail(c, gin.H{"err": err, "vid": vid}, "代金券不存在")
		return
	}
	// 更新
	if v.IsUsed == 1 {
		response.Success(c, gin.H{"data": v, "vid": vid}, "代金券已使用")
		DB.Model(&v).Update("IsUsed", int(1))
		return
	}
	if err := DB.Model(&v).Update("IsUsed", int(1)).Error; err != nil {
		response.Fail(c, gin.H{"err": err}, "代金券使用失败")
		return
	}

	response.Success(c, gin.H{"voucher": v}, "修改成功")
}
