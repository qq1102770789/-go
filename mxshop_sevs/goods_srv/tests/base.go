package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mxshop_sevs/goods_srv/model/proto"
)

func Init() {
	var err error
	conn, err = grpc.Dial("192.168.43.127:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)

	}

	brandClient = proto.NewGoodsClient(conn)
}
func TestGetGoodsList() {
	rsp, err := brandClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
		TopCategory: 130361,
		PriceMin:    90,
		//KeyWords: "深海速冻",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, good := range rsp.Data {
		fmt.Println(good.Name, good.ShopPrice)
	}
}

func main() {
	Init()
	TestGetGoodsList()
	conn.Close()

}
