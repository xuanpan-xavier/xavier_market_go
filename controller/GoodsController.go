package controller

// 会员管理系统内的商品管理
import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xmarket_gin/common"
	"xmarket_gin/model"
	"xmarket_gin/repository"
	"xmarket_gin/response"
	"xmarket_gin/vo"
)

type IGoodsController interface {
	RestController
	ShowAll(c *gin.Context)
}

type GoodsController struct {
	Repository repository.IGoodsRepository
}

func (g GoodsController) Create(c *gin.Context) {
	var requestGoods vo.CreateGoodsRequest
	if err := c.ShouldBind(&requestGoods); err != nil {
		response.Fail(c, gin.H{"data": requestGoods}, "数据验证错误，请填写完整信息")
		panic(err)
		return
	}
	price, _ := strconv.ParseFloat(requestGoods.Price, 64)
	goods, err := g.Repository.Create(requestGoods.Name, requestGoods.Category, float32(price))
	if err != nil {
		//panic(err)
		response.Fail(c, nil, "创建失败")
		return
	}

	response.Success(c, gin.H{"goods": goods}, "nil")
}

func (g GoodsController) Delete(c *gin.Context) {
	// 获取path中参数
	id, err := strconv.ParseUint(c.Params.ByName("ID"), 10, 64)
	if err != nil {
		response.Fail(c, nil, "转化uuid失败")
	}
	if err := g.Repository.DeleteById(uint(id)); err != nil {
		response.Fail(c, nil, "删除失败")
		return
	}
	response.Success(c, nil, "")
}

func (g GoodsController) Show(c *gin.Context) {
	// 获取path中参数
	id, err1 := strconv.ParseUint(c.Params.ByName("ID"), 10, 64)
	if err1 != nil {
		response.Fail(c, nil, "转化uuid失败")
	}
	goods, err2 := g.Repository.SelectById(uint(id))
	if err2 != nil {
		response.Fail(c, nil, "商品不存在")
		return
	}
	response.Success(c, gin.H{"goods": goods}, "")
}

func (g GoodsController) ShowAll(c *gin.Context) {
	query := c.DefaultQuery("query", "")
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	DB := common.GetDB()
	var goods []model.Goods
	var total int
	if query == "" {
		DB.Order("category, id").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&goods)
		DB.Model(model.Goods{}).Count(&total)
	} else {
		DB.Order("category, id").Offset((pageNum-1)*pageSize).Limit(pageSize).Where("category = ?", query).Find(&goods)
		DB.Model(model.Goods{}).Where("category = ?", query).Count(&total)
	}

	response.Success(c, gin.H{"data": goods, "total": total}, "查询成功")
}

func (g GoodsController) Update(c *gin.Context) {
	var requestGoods vo.UpdateGoodsRequest
	if err := c.ShouldBind(&requestGoods); err != nil {
		response.Fail(c, nil, "数据验证错误，请填写完整信息")
		panic(err)
		return
	}
	// 获取path中参数
	id, err1 := strconv.ParseUint(c.Params.ByName("ID"), 10, 64)
	if err1 != nil {
		response.Fail(c, nil, "转化uint失败")
	}
	updateGoods, err2 := g.Repository.SelectById(uint(id))
	if err2 != nil {
		response.Fail(c, nil, "商品不存在")
		return
	}
	// 更新
	price, _ := strconv.ParseFloat(requestGoods.Price, 64)
	goods, err3 := g.Repository.Update(*updateGoods, float32(price))
	if err3 != nil {
		response.Fail(c, nil, "更新失败")
	}
	response.Success(c, gin.H{"goods": goods}, "修改成功")

}

func NewGoodsController() IGoodsController {
	GoodsController := GoodsController{Repository: repository.NewGoodsRepository()}
	GoodsController.Repository.(repository.GoodsRepository).DB.AutoMigrate(model.Goods{})
	return GoodsController
}
