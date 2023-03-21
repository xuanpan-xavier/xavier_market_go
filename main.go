package main

import (
	"context"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"os"
	"time"

	//"github.com/go-echarts/go-echarts/v2/charts"
	//"github.com/go-echarts/go-echarts/v2/opts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	//"os"
)

//func makeFindOptions(length int64) *options.FindOptions {
//	//分页
//	var limit int64
//	limit = 1
//	skip := (length/2 - 1)
//	if skip < 0 {
//		skip = 0
//	}
//	sort := bson.D{
//		bson.E{"price", -1},
//	}
//	//查询条件
//	opts := &options.FindOptions{
//		Sort:  sort,
//		Limit: &limit,
//		Skip:  &skip,
//	}
//
//	return opts
//}
var (
	itemCntPie = 6
	mon        = []string{"before 2022", "2022-Jan.", "2022-Feb.", "2022-Mar.", "2022-Apr.", "2022-May."}
	data       = []int{8, 0, 69, 416, 535, 31529}
)

func generatePieItems() []opts.PieData {
	items := make([]opts.PieData, 0)
	for i := 0; i < itemCntPie; i++ {
		items = append(items, opts.PieData{Name: mon[i], Value: data[i]})
	}
	return items
}
func main() {
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI("mongodb://rwuser:mongo_SDU2022@123.249.29.154:8635")

	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("A:Connected fail to MongoDB!")
		log.Fatal(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println("B:Connected fail to MongoDB!")
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	// 指定获取要操作的数据集
	collection := client.Database("TH").Collection("rent")
	//// 查询多个
	//// 将选项传递给Find()
	//findOptions := options.Find()
	//findOptions.SetLimit(2)
	//type BaseInfo struct {
	//	Types string `bson:"types"`
	//}
	//type rent struct {
	//	City       string   `bson:"city"`
	//	Base_info  BaseInfo `bson:"base_info"`
	//	Lease_mode string   `bson:"lease_mode"`
	//	Price      int      `bson:"price"`
	//	//maintain   string `bson:"maintain"`
	//}

	//// 定义一个切片用来存储查询结果
	//var results []*rent
	//// 把bson.D{{}}作为一个filter来匹配所有文档
	//cur, err := collection.Find(context.TODO(), bson.D{}, findOptions)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("111")
	//// 查找多个文档返回一个光标
	//// 遍历游标允许我们一次解码一个文档
	//for cur.Next(context.TODO()) {
	//	fmt.Printf("222")
	//	// 创建一个值，将单个文档解码为该值
	//	var elem rent
	//	err := cur.Decode(&elem)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("city: %s, types: %s \n", elem.city, elem.types)
	//	results = append(results, &elem)
	//}
	//
	//if err := cur.Err(); err != nil {
	//	log.Fatal(err)
	//}
	//
	//// 完成后关闭游标
	//cur.Close(context.TODO())
	//fmt.Printf("Found multiple documents (array of pointers): %#v\n", results)
	// 查询name=Bob的文档
	//var c string
	//var t string
	//var l string
	//c = "北京"
	//l = "整租"
	//pipeLine := mongo.Pipeline{
	//	{{"$match", bson.D{{"city", "上海"}}}},
	//	{{"$group", bson.D{{"_id", bson.D{{"types", "$base_info.types"}, {"mode", "$lease_mode"}}}, {"avg", bson.D{{"$avg", "$price"}}}, {"max", bson.D{{"$max", "$price"}}}, {"min", bson.D{{"$min", "$price"}}}}}},
	//}
	//cursor, err := collection.Aggregate(context.Background(), pipeLine)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//type ID struct {
	//	Mode  string `bson:"mode"`
	//	Types string `bson:"types"`
	//}
	//type Info struct {
	//	ID  ID      `bson:"_id"`
	//	AVG float64 `bson:"avg"`
	//	MAX float64 `bson:"max"`
	//	MIN float64 `bson:"min"`
	//}
	//// 定义bson.M类型的文档数组，bson.M是一个map类型的键值数据结构
	//var results []Info
	//// 使用All函数获取所有查询结果，并将结果保存至results变量。
	//if err = cursor.All(context.TODO(), &results); err != nil {
	//	log.Fatal(err)
	//}
	//// 遍历结果数组
	//var modes string
	//var types string
	//var adds string
	//var ids []opts.BarData
	//type median struct {
	//	Price int `bson:"price"`
	//}
	//var meds []opts.BarData
	//for _, result := range results {
	//	//fmt.Println(result)
	//	modes = result.ID.Mode
	//	types = result.ID.Types
	//	adds = modes + types
	//	ids = append(ids, opts.BarData{Value: adds})
	//	count, _ := collection.CountDocuments(context.Background(), bson.D{{"city", "上海"}, {"base_info.types", types}, {"lease_mode", modes}})
	//	cur, err := collection.Find(context.Background(), bson.D{{"city", "上海"}, {"base_info.types", types}, {"lease_mode", modes}}, makeFindOptions(count))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	var med []median
	//	if err = cur.All(context.TODO(), &med); err != nil {
	//		log.Fatal(err)
	//	}
	//	for _, m := range med {
	//		fmt.Println(m)
	//		meds = append(meds, opts.BarData{Value: m.Price})
	//	}
	//}
	//
	////create a new bar instance
	//bar := charts.NewBar()
	//// set some global options like Title/Legend/ToolTip or anything else
	//bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
	//	Title:    "Nosql 实验",
	//	Subtitle: "city:上海 data:median",
	//}),
	//	charts.WithInitializationOpts(opts.Initialization{Width: "20000px", Height: "500px"}),
	//)
	//
	//// Put data into instance
	//bar.SetXAxis(ids).
	//	AddSeries("median", meds)
	////AddSeries("min", mins).
	////AddSeries("min", maxs)
	//// Where the magic happens
	//f, _ := os.Create("bar.html")
	//bar.Render(f)

	type rentTime struct {
		Maintain time.Time `bson:"maintain"`
	}
	type rent struct {
		Base_info rentTime `bson:"base_info"`
	}
	tt := [][]time.Time{
		{time.Date(2022, time.January, 1, 00, 0, 0, 0, time.UTC), time.Date(2022, time.February, 1, 00, 0, 0, 0, time.UTC)},
		{time.Date(2022, time.February, 1, 00, 0, 0, 0, time.UTC), time.Date(2022, time.March, 1, 00, 0, 0, 0, time.UTC)},
		{time.Date(2022, time.March, 1, 00, 0, 0, 0, time.UTC), time.Date(2022, time.April, 1, 00, 0, 0, 0, time.UTC)},
		{time.Date(2022, time.April, 1, 00, 0, 0, 0, time.UTC), time.Date(2022, time.May, 1, 00, 0, 0, 0, time.UTC)},
		{time.Date(2022, time.May, 1, 00, 0, 0, 0, time.UTC), time.Date(2022, time.June, 1, 00, 0, 0, 0, time.UTC)},
		//	{time.Date(2022, time.June, 1, 00, 0, 0, 0, time.UTC), time.Date(2022, time.July, 1, 00, 0, 0, 0, time.UTC)},
		//	{time.Date(2022, time.July, 1, 00, 0, 0, 0, time.UTC), time.Date(2022, time.August, 1, 00, 0, 0, 0, time.UTC)},
		//	{time.Date(2022, time.August, 1, 00, 0, 0, 0, time.UTC), time.Date(2022, time.September, 1, 00, 0, 0, 0, time.UTC)},
		//	{time.Date(2022, time.September, 1, 00, 0, 0, 0, time.UTC), time.Date(2022, time.October, 1, 00, 0, 0, 0, time.UTC)},
		//	{time.Date(2022, time.October, 1, 00, 0, 0, 0, time.UTC), time.Date(2022, time.November, 1, 00, 0, 0, 0, time.UTC)},
		//	{time.Date(2022, time.November, 1, 00, 0, 0, 0, time.UTC), time.Date(2022, time.December, 1, 00, 0, 0, 0, time.UTC)},
		//	{time.Date(2022, time.December, 1, 00, 0, 0, 0, time.UTC), time.Date(2023, time.January, 1, 00, 0, 0, 0, time.UTC)},
	}
	var c []int64
	cc, _ := collection.CountDocuments(context.Background(), bson.D{{"city", "上海"}, {"base_info.maintain", bson.D{{"$lte", time.Date(2022, time.January, 1, 00, 0, 0, 0, time.UTC)}}}})
	c = append(c, cc)
	for _, t := range tt {
		count, _ := collection.CountDocuments(context.Background(), bson.D{{"city", "上海"}, {"base_info.maintain", bson.D{{"$gte", t[0]}, {"$lte", t[1]}}}})
		c = append(c, count)
	}
	for _, cs := range c {
		fmt.Println(cs)
	}
	// 断开连接
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "月份&租房"}),
	)

	pie.AddSeries("pie", generatePieItems())
	// create a new bar instance
	// Where the magic happens
	f, _ := os.Create("pie.html")
	pie.Render(f)
}

//
//import (
//	"github.com/gin-gonic/gin"
//	"github.com/spf13/viper"
//	"os"
//	"xmarket_gin/common"
//	"xmarket_gin/routes"
//)
//
//func main() {
//	InitConfig()
//	db := common.InitDB()
//	defer db.Close()
//	rds := common.RedisPollInit().Get()
//	// redis-server.exe redis.windows.conf
//	defer rds.Close()
//
//	r := gin.Default()
//	r = routes.CollectRoute(r)
//	port := viper.GetString("server.port")
//	if port != "" {
//		panic(r.Run(":" + port))
//	}
//	panic(r.Run())
//}
//
//func InitConfig() {
//	workDir, _ := os.Getwd()
//	viper.SetConfigName("application")
//	viper.SetConfigType("yml")
//	viper.AddConfigPath(workDir + "/config")
//	err := viper.ReadInConfig()
//	if err != nil {
//		panic(err)
//	}
//}
