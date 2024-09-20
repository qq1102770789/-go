package model

import (
	"context"
	"gorm.io/gorm"
	"mxshop_sevs/goods_srv/global"
	"strconv"
)

// 分类模型
type Category struct {
	BaseModel                    // 基础模型
	Name             string      `gorm:"type:varchar(20);not null;comment:'分类名称'" json:"name"` // 分类名称
	ParentCategoryID int32       `gorm:"type:int" json:"parent"`                               // 父分类ID
	ParentCategory   *Category   `json:"-"`                                                    // 父分类，自引用
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
	Level            int32       `gorm:"type:int;not null;default:1" json:"level"` // 分类级别，默认为1
	IsTab            bool        `gorm:"not null;default:false" json:"is_tab"`     // 是否标签
	Url              string      `gorm:"type:varchar(255)" json:"url"`             // 分类链接
}

// 品牌模型
type Brands struct {
	BaseModel        // 基础模型
	Name      string `gorm:"type:varchar(30);comment:'品牌名称'"`                        // 品牌名称
	Logo      string `gorm:"type:varchar(255);not null;default:'';comment:'品牌logo'"` // 品牌logo
}

// 商品分类品牌关联模型
type GoodsCategoryBrand struct {
	BaseModel           // 基础模型
	Category   Category // 分类
	CategoryID int32    `gorm:"type:int;idx_category_brand,unique"` // 分类ID，与分类关联
	BrandsID   int32    `gorm:"type:int;idx_category_brand,unique"` // 品牌ID，与品牌关联
	Brands     Brands   // 品牌
}

// 设置商品分类品牌关联模型的表名
func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

// 轮播图模型
type Banner struct {
	BaseModel        // 基础模型
	Image     string `gorm:"type:varchar(255);not null;comment:'轮播图'"`  // 轮播图路径
	Url       string `gorm:"type:varchar(255);not null;comment:'跳转链接'"` // 跳转链接
	Index     int32  `gorm:"type:int;not null"`                         // 轮播图索引，默认为1
}
type Goods struct {
	BaseModel
	Category        Category // 分类
	CategoryID      int32    `gorm:"type:int"`                          // 分类ID，与分类关联
	BrandsID        int32    `gorm:"type:int;not null;column:brand_id"` // 品牌ID，与品牌关联
	Brands          Brands   // 品牌
	OnSale          bool     `gorm:"default:false;not null"`                    // 是否上架
	ShipFee         bool     `gorm:"default:false;not null"`                    // 是否包邮
	IsNew           bool     `gorm:"default:false;not null"`                    // 是否新品
	IsHot           bool     `gorm:"default:false;not null"`                    // 是否热销
	Name            string   `gorm:"type:varchar(255);not null;comment:'商品名称'"` // 商品名称
	GoodsSn         string   `gorm:"type:varchar(50);not null;comment:'商品货号'"`  // 商品货号
	ClickNum        int32    `gorm:"type:int;not null;default:0"`               // 点击数
	SoldNum         int32    `gorm:"type:int;not null;default:0"`               // 销售量
	FavNum          int32    `gorm:"type:int;not null;default:0"`               // 收藏数
	MarketPrice     float32  `gorm:"not null;comment:'市场价格'"`                   // 市场价格
	ShopPrice       float32  `gorm:"not null;comment:'本店价格'"`                   // 本店价格
	ShipFree        bool     `gorm:"default:false;not null"`
	GoodsBrief      string   `gorm:"type:varchar(255);not null;comment:'商品简介'"` // 商品简介
	Images          GormList `gorm:"type:varchar(2000);not null"`
	DescImages      GormList `gorm:"type:varchar(4000);not null"`
	GoodsFrontImage string   `gorm:"type:varchar(255);not null;comment:'商品封面图'"`
	Stock           int32    `gorm:"type:int;not null;default:0;column:stocks"` // 库存
}

func (g *Goods) AfterCreate(tx *gorm.DB) (err error) {
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandsID:    g.BrandsID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}

	_, err = global.EsClient.Index().Index(esModel.GetIndexName()).BodyJson(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (g *Goods) AfterUpdate(tx *gorm.DB) (err error) {
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandsID:    g.BrandsID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}

	_, err = global.EsClient.Update().Index(esModel.GetIndexName()).
		Doc(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (g *Goods) AfterDelete(tx *gorm.DB) (err error) {
	_, err = global.EsClient.Delete().Index(EsGoods{}.GetIndexName()).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
