package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop_sevs/goods_srv/global"
	"mxshop_sevs/goods_srv/model"
	"mxshop_sevs/goods_srv/model/proto"
)

// UserServer 结构体
type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

func ModelToResponse(goods model.Goods) proto.GoodsInfoResponse {
	return proto.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brands.ID,
			Name: goods.Brands.Name,
			Logo: goods.Brands.Logo,
		},
	}
}

func (s *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	// 使用 Elasticsearch (ES) 的目的是搜索出商品的 ID，通过 ID 拿到具体的字段信息是通过 MySQL 来完成
	// 我们使用 ES 是用来做搜索的，是否应该将所有的 MySQL 字段全部在 ES 中保存一份？
	// ES 用来做搜索，这个时候我们一般只把搜索和过滤的字段信息保存到 ES 中
	// ES 可以用来当做 MySQL 使用，但是实际上 MySQL 和 ES 之间是互补的关系，一般 MySQL 用来做存储使用，ES 用来做搜索使用
	// ES 想要提高性能，就要将 ES 的内存设置得够大，比如 1k 2k

	// 关键词搜索、查询新品、查询热门商品、通过价格区间筛选，通过商品分类筛选
	goodsListResponse := &proto.GoodsListResponse{}

	// 使用 match bool 复合查询
	q := elastic.NewBoolQuery()               // 创建一个布尔查询对象，用于组合多个查询条件
	localDB := global.DB.Model(model.Goods{}) // 获取 MySQL 数据库的 Goods 表模型

	// 关键词搜索
	if req.KeyWords != "" {
		q = q.Must(elastic.NewMultiMatchQuery(req.KeyWords, "name", "goods_brief")) // 使用 MultiMatchQuery 进行关键词搜索
	}
	//must会算分数，filter不会算分数，filter只用来过滤，must用来匹配，所以filter可以用来过滤掉一些不符合条件的结果，must用来匹配符合条件的结果

	// 查询热门商品
	if req.IsHot {
		localDB = localDB.Where(model.Goods{IsHot: true})       // MySQL 查询条件：是否热门
		q = q.Filter(elastic.NewTermQuery("is_hot", req.IsHot)) // ES 过滤条件：是否热门
	}

	// 查询新品
	if req.IsNew {
		q = q.Filter(elastic.NewTermQuery("is_new", req.IsNew)) // ES 过滤条件：是否新品
	}

	// 通过价格区间筛选
	if req.PriceMin > 0 {
		q = q.Filter(elastic.NewRangeQuery("shop_price").Gte(req.PriceMin)) // ES 过滤条件：最低价格
	}
	if req.PriceMax > 0 {
		q = q.Filter(elastic.NewRangeQuery("shop_price").Lte(req.PriceMax)) // ES 过滤条件：最高价格
	}

	// 通过品牌筛选
	if req.Brand > 0 {
		q = q.Filter(elastic.NewTermQuery("brands_id", req.Brand)) // ES 过滤条件：品牌
	}

	// 通过分类查询商品
	var subQuery string
	categoryIds := make([]interface{}, 0) // 存储所有符合条件的分类 ID
	if req.TopCategory > 0 {
		var category model.Category
		if result := global.DB.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在") // 如果分类不存在，返回错误
		}

		// 根据分类级别生成相应的子查询 SQL 语句
		if category.Level == 1 {
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category WHERE parent_category_id=%d)", req.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("select id from category WHERE parent_category_id=%d", req.TopCategory)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("select id from category WHERE id=%d", req.TopCategory)
		}

		// 执行子查询获取符合条件的分类 ID
		type Result struct {
			ID int32
		}
		var results []Result
		global.DB.Model(model.Category{}).Raw(subQuery).Scan(&results)
		for _, re := range results {
			categoryIds = append(categoryIds, re.ID)
		}

		// 生成 terms 查询，过滤符合条件的分类 ID
		q = q.Filter(elastic.NewTermsQuery("category_id", categoryIds...))
	}

	// 分页处理
	if req.Pages == 0 {
		req.Pages = 1
	}

	// 设置每页显示数量，最大为 100，默认 10
	switch {
	case req.PagePerNums > 100:
		req.PagePerNums = 100
	case req.PagePerNums <= 0:
		req.PagePerNums = 10
	}

	// 在 ES 中执行搜索查询
	result, err := global.EsClient.Search().Index(model.EsGoods{}.GetIndexName()).Query(q).From(int(req.Pages)).Size(int(req.PagePerNums)).Do(context.Background())
	if err != nil {
		return nil, err
	}

	// 获取商品 ID 列表
	goodsIds := make([]int32, 0)
	goodsListResponse.Total = int32(result.Hits.TotalHits.Value) // 总命中数量
	for _, value := range result.Hits.Hits {
		goods := model.EsGoods{}
		_ = json.Unmarshal(value.Source, &goods) // 反序列化 ES 搜索结果
		goodsIds = append(goodsIds, goods.ID)
	}

	// 在 MySQL 中查询具体商品信息
	var goods []model.Goods
	re := localDB.Preload("Category").Preload("Brands").Find(&goods, goodsIds)
	if re.Error != nil {
		return nil, re.Error
	}

	// 将查询结果转换为响应格式
	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good) // 将模型转换为响应结构
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}

	return goodsListResponse, nil
}

// 现在用户提交订单有多个商品，你得批量查询商品的信息吧
func (s *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	goodsListResponse := &proto.GoodsListResponse{}
	var goods []model.Goods

	//调用where并不会真正执行sql 只是用来生成sql的 当调用find， first才会去执行sql，
	result := global.DB.Where(req.Id).Find(&goods)
	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}
	goodsListResponse.Total = int32(result.RowsAffected)
	return goodsListResponse, nil
}
func (s *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var goods model.Goods

	if result := global.DB.Preload("Category").Preload("Brands").First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	goodsInfoResponse := ModelToResponse(goods)
	return &goodsInfoResponse, nil
}
func (s *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}
	//先检查redis中是否有这个token
	//防止同一个token的数据重复插入到数据库中，如果redis中没有这个token则放入redis
	//这里没有看到图片文件是如何上传， 在微服务中 普通的文件上传已经不再使用
	goods := model.Goods{
		Brands:          brand,
		BrandsID:        req.BrandId,
		Category:        category,
		CategoryID:      category.ID,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.GoodsBrief,
		ShipFree:        req.ShipFree,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		OnSale:          req.OnSale,
		Stock:           req.Stocks,
	}

	//srv之间互相调用了
	tx := global.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &proto.GoodsInfoResponse{
		Id: goods.ID,
	}, nil
}

func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*proto.Empty, error) {
	if result := global.DB.Delete(&model.Goods{BaseModel: model.BaseModel{ID: req.Id}}, req.Id); result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	return &proto.Empty{}, nil
}

func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.Empty, error) {
	var goods model.Goods

	if result := global.DB.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	goods.Brands = brand
	goods.BrandsID = brand.ID
	goods.Category = category
	goods.CategoryID = category.ID
	goods.Name = req.Name
	goods.GoodsSn = req.GoodsSn
	goods.MarketPrice = req.MarketPrice
	goods.ShopPrice = req.ShopPrice
	goods.GoodsBrief = req.GoodsBrief
	goods.ShipFree = req.ShipFree
	goods.Images = req.Images
	goods.DescImages = req.DescImages
	goods.GoodsFrontImage = req.GoodsFrontImage
	goods.IsNew = req.IsNew
	goods.IsHot = req.IsHot
	goods.OnSale = req.OnSale

	tx := global.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &proto.Empty{}, nil
}

// ModelToResponse 函数将用户模型转换为用户信息响应
// 在grpc的message中字段有默认值，你不能随便赋值nil，容易出错
// 因此需要清楚，那些字段是有默认值的

// 商品接口
//func (s *GoodsServer)  GoodsList(context context.Context, *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error){
//
//}
//BatchGetGoods(context.Context, *BatchGoodsIdInfo) (*GoodsListResponse, error)
//CreateGoods(context.Context, *CreateGoodsInfo) (*GoodsInfoResponse, error)
//DeleteGoods(context.Context, *DeleteGoodsInfo) (*Empty, error)
//UpdateGoods(context.Context, *CreateGoodsInfo) (*Empty, error)
//GetGoodsDetail(context.Context, *GoodInfoRequest) (*GoodsInfoResponse, error)
