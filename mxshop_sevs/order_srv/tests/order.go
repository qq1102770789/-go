package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mxshop_sevs/order_srv/model/proto"
)

var orderClient proto.OrderClient
var conn *grpc.ClientConn

func TestCreateCartItem(userId, nums, goodsId int32) {
	rsp, err := orderClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  userId,
		Nums:    nums,
		GoodsId: goodsId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id)
}

func TestCartItemList(userId int32) {
	rsp, err := orderClient.CartItemList(context.Background(), &proto.UserInfo{
		Id: userId,
	})
	if err != nil {
		panic(err)
	}
	for _, item := range rsp.Data {
		fmt.Println(item.Id, item.GoodsId, item.Nums)
	}
}

func TestUpdateCartItem(id int32) {
	_, err := orderClient.UpdateCartItem(context.Background(), &proto.CartItemRequest{
		Id:      id,
		Checked: true,
	})
	if err != nil {
		zap.S().Error(err.Error())
	}
}

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	orderClient = proto.NewOrderClient(conn)
}

func TestCreateOrder() {
	_, err := orderClient.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  1,
		Address: "北京市",
		Name:    "bobby",
		Mobile:  "18996071730",
		Post:    "请尽快发货",
	})
	if err != nil {
		panic(err)

	}
}

func TestGetOrderDetail(orderId int32) {
	rsp, err := orderClient.OrderDetail(context.Background(), &proto.OrderRequest{
		Id: orderId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.OrderInfo.OrderSn)
	for _, good := range rsp.Goods {
		fmt.Println(good.GoodsName)
	}

}

func TestOrderList() {
	rsp, err := orderClient.OrderList(context.Background(), &proto.OrderFilterRequest{
		UserId: 1,
	})
	if err != nil {
		panic(err)
	}

	for _, order := range rsp.Data {
		fmt.Println(order.OrderSn)
	}
}

func main() {
	Init()
	//启动两个依赖的微服务 商品 和库存
	TestCreateCartItem(1, 1, 422)
	//TestCartItemList(1)
	//TestUpdateCartItem(1)
	//TestCreateOrder()
	//TestGetOrderDetail(20)
	//TestOrderList()
	conn.Close()
}
