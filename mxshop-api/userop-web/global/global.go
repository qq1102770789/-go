package global

import (
	ut "github.com/go-playground/universal-translator"

	"mxshop-api/userop-web/config"
	"mxshop-api/userop-web/proto"
)

var (
	ServerConfig   *config.ServerConfig = &config.ServerConfig{}
	Trans          ut.Translator
	GoodsSrvClient proto.GoodsClient
	NacosConfig    = &config.NacosConfig{}
	MessageClient  proto.MessageClient
	AddressClient  proto.AddressClient
	UserFavClient  proto.UserFavClient
)
