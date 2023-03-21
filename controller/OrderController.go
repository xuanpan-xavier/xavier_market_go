package controller

// 会员管理系统内的订单列表
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
	"xmarket_gin/common"
	"xmarket_gin/model"
	"xmarket_gin/repository"
	"xmarket_gin/response"
	"xmarket_gin/vo"
)

type IOrderController interface {
	RestController
	ShowAll(c *gin.Context)
	ReadPie(c *gin.Context)
	ReadBar(c *gin.Context)
}
type OrderController struct {
	Repository repository.IOrderRepository
}

func (o OrderController) Create(c *gin.Context) {
	var requestOrder vo.CreateOrderRequest
	if err := c.ShouldBind(&requestOrder); err != nil {
		response.Fail(c, gin.H{"data": requestOrder}, "数据验证错误，请填写完整信息")
		panic(err)
		return
	}
	DB := common.GetDB()
	var goods []model.Goods
	var num []int
	var number string
	var total float32
	var exist int
	for _, id := range requestOrder.GoodsID {
		g, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.Fail(c, nil, "转化uint失败")
			return
		}
		var good model.Goods
		if err := DB.Model(&good).Where("ID = ?", uint(g)).Find(&good).Error; err != nil {
			response.Fail(c, gin.H{"data": gin.H{"goods[]": goods, "good": good, "g": g, "err": err}}, "获取商品信息失败")
			return
		}
		for i, g := range goods {
			if good.ID == g.ID {
				num[i] += 1
				exist = 1
				break
			}
		}
		if exist != 1 {
			goods = append(goods, good)
			num = append(num, 1)
			exist = 0
		}
		total += good.Price
	}
	fmt.Println(num)
	for _, n := range num {
		number = number + strconv.FormatInt(int64(n), 10)
	}
	var order *model.Order
	var errO error
	DB.Transaction(func(tx *gorm.DB) error {
		order, errO = o.Repository.Create(requestOrder.Telephone, goods, total, number)
		//for i, g := range goods {
		//	DB.Table("order_goods").Debug().Where("order_id = ?", order.ID).Where("goods_id = ?", g.ID).Update("goods_number", num[i])
		//}
		if errO != nil {
			panic(errO)
			response.Fail(c, gin.H{"data": errO}, "创建失败")
			return nil
		}
		var card model.VipCard
		DB.Model(&card).Where("Telephone = ?", requestOrder.Telephone).Find(&card)
		DB.Model(&card).Update("Point", card.Point+uint(total))
		return nil
	})
	response.Success(c, gin.H{"order": order}, "nil")
}

func (o OrderController) Delete(c *gin.Context) {
	// 获取path中参数
	id, _ := strconv.ParseUint(c.Params.ByName("OrderID"), 10, 64)

	DB := common.GetDB()
	var card model.VipCard
	var order model.Order
	DB.Transaction(func(tx *gorm.DB) error {
		DB.Model(&order).Where("ID = ?", uint(id)).Find(&order)
		DB.Model(&card).Where("Telephone = ?", order.Telephone).Find(&card)
		DB.Model(&card).Update("Point", card.Point-uint(order.Total))

		if err := o.Repository.DeleteById(uint(id)); err != nil {
			response.Fail(c, nil, "删除失败")
			return nil
		}
		return nil
	})
	response.Success(c, nil, "")
}

func (o OrderController) Show(c *gin.Context) {
	// 获取path中参数
	id, _ := strconv.ParseUint(c.Params.ByName("OrderID"), 10, 64)
	order, err2 := o.Repository.SelectById(uint(id))
	if err2 != nil {
		response.Fail(c, nil, "订单不存在")
		return
	}
	response.Success(c, gin.H{"order": order}, "")
}
func (o OrderController) Update(c *gin.Context) {
	var requestOrder vo.UpdateOrderRequest
	if err := c.ShouldBind(&requestOrder); err != nil {
		response.Fail(c, nil, "数据验证错误，请填写完整信息")
		panic(err)
		return
	}
	// 获取path中参数
	id, _ := strconv.ParseUint(c.Params.ByName("OrderID"), 10, 64)
	updateOrder, err2 := o.Repository.SelectById(uint(id))
	if err2 != nil {
		response.Fail(c, nil, "订单不存在")
		return
	}
	// 更新
	DB := common.GetDB()
	var card model.VipCard
	DB.Model(&card).Where("Telephone = ?", updateOrder.Telephone).Find(&card)
	DB.Model(&card).Update("Point", card.Point-uint(updateOrder.Total))
	var goods []model.Goods
	var total float32
	for _, id := range requestOrder.GoodsID {
		g, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.Fail(c, nil, "转化uint失败")
			return
		}
		var good model.Goods
		if err := DB.First(&good, g).Error; err != nil {
			response.Fail(c, gin.H{"data": err}, "获取商品信息失败")
			return
		}
		goods = append(goods, good)
		total += good.Price
	}
	order, err3 := o.Repository.Update(*updateOrder, goods, total)
	if err3 != nil {
		response.Fail(c, nil, "更新失败")
		return
	}

	DB.Model(&card).Update("Point", card.Point+uint(total))
	response.Success(c, gin.H{"order": order}, "修改成功")

}

func (o OrderController) ShowAll(c *gin.Context) {
	query := c.DefaultQuery("query", "")
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	DB := common.GetDB()
	var order []model.Order
	var total int
	if query == "" {
		DB.Order("id").Offset((pageNum - 1) * pageSize).Limit(pageSize).Preload("Goods").Find(&order)
		DB.Model(model.Order{}).Count(&total)
	} else {
		DB.Order("id").Where("Telephone = ?", query).Offset((pageNum - 1) * pageSize).Limit(pageSize).Preload("Goods").Find(&order)
		DB.Model(model.Order{}).Where("Telephone = ?", query).Count(&total)
	}

	response.Success(c, gin.H{"data": order, "total": total}, "查询成功")
}

func (o OrderController) ReadPie(c *gin.Context) {
	var requestOrder vo.ReadOrderRequest
	if err := c.ShouldBind(&requestOrder); err != nil {
		response.Fail(c, gin.H{"data": requestOrder}, "数据验证错误，请填写完整信息")
		panic(err)
		return
	}
	DB := common.GetDB()
	var orders []model.Order
	var goodsMap map[string]float32
	goodsMap = make(map[string]float32)
	DB.Model(&orders).Where("Telephone = ?", requestOrder.Telephone).Order("id").Preload("Goods").Find(&orders)
	for _, order := range orders {
		for i, goods := range order.Goods {
			n, _ := strconv.ParseFloat(order.Number[i:i+1], 64)
			goodsMap[goods.Category] += goods.Price * float32(n)
		}
	}
	var gs []string
	var ps []float32
	for g, p := range goodsMap {
		gs = append(gs, g)
		ps = append(ps, p)
	}
	response.Success(c, gin.H{"Category": gs, "Number": ps}, "查询成功")
}

func (o OrderController) ReadBar(c *gin.Context) {
	var requestOrder vo.ReadOrderRequest
	if err := c.ShouldBind(&requestOrder); err != nil {
		response.Fail(c, gin.H{"data": requestOrder}, "数据验证错误，请填写完整信息")
		panic(err)
		return
	}
	DB := common.GetDB()
	var orders []model.Order
	var Date []string
	var Total []float32
	//var results [][]sql.Result
	DB.Model(&orders).Where("Telephone = ?", requestOrder.Telephone).Order("id desc").Limit(10).Preload("Goods").Find(&orders)
	var count int
	DB.Model(&orders).Where("Telephone = ?", requestOrder.Telephone).Count(&count)
	Category := [11][10]float32{}
	for index, order := range orders {
		Date = append(Date, order.CreatedAt.Format("2006-01-02"))
		Total = append(Total, order.Total)
		var goodsMap map[string]float32
		goodsMap = make(map[string]float32)
		//var result []sql.Result
		//DB.Model(&order.Goods).Select("Category, sum(Price) as total").Group("Category").Scan(&result)
		//results = append(results, result)
		for i, goods := range order.Goods {
			n, _ := strconv.ParseFloat(order.Number[i:i+1], 64)
			goodsMap[goods.Category] += goods.Price * float32(n)
		}
		for category, price := range goodsMap {
			switch category {
			case "果蔬生鲜":
				Category[0][index] = price
			case "熟食冻品":
				Category[1][index] = price
			case "粮油调味":
				Category[2][index] = price
			case "休闲零食":
				Category[3][index] = price
			case "水饮冲调":
				Category[4][index] = price
			case "美妆个护":
				Category[5][index] = price
			case "家居百货":
				Category[6][index] = price
			case "日用品":
				Category[7][index] = price
			case "文体玩具":
				Category[8][index] = price
			case "服装首饰":
				Category[9][index] = price
			case "电子产品":
				Category[10][index] = price

			}
		}
	}
	response.Success(c, gin.H{"Date": Date, "Total": Total, "guo": Category[0][:count], "shu": Category[1][:count], "liang": Category[2][:count], "xiu": Category[3][:count], "shui": Category[4][:count], "mei": Category[5][:count], "jia": Category[6][:count], "ri": Category[7][:count], "wen": Category[8][:count], "fu": Category[9][:count], "dian": Category[10][:count]}, "查询成功")
}

func NewOrderController() IOrderController {
	OrderController := OrderController{Repository: repository.NewOrderRepository()}
	OrderController.Repository.(repository.OrderRepository).DB.AutoMigrate(model.Order{}, model.Voucher{})
	return OrderController
}
