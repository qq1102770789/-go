package global

import (
	ut "github.com/go-playground/universal-translator"

	"mxshop-api/order-web/config"
	"mxshop-api/order-web/proto"
)

var (
	ServerConfig       *config.ServerConfig = &config.ServerConfig{}
	Trans              ut.Translator
	GoodsSrvClient     proto.GoodsClient
	NacosConfig        = &config.NacosConfig{}
	OrderSrvClient     proto.OrderClient
	InventorySrvClient proto.InventoryClient
)
