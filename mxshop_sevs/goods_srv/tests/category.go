package main

import (
	"context"
	"fmt"
	"mxshop_sevs/goods_srv/model/proto"
)

func TestGetCategoryList() {
	rsp, err := brandClient.GetAllCategorysList(context.Background(), &proto.Empty{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	fmt.Println(rsp.JsonData)
}
func TestGetSubCategoryList() {
	rsp, err := brandClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id: 130358,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(rsp.SubCategorys)

}
