package handler

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop_sevs/goods_srv/global"
	"mxshop_sevs/goods_srv/model"
	"mxshop_sevs/goods_srv/model/proto"
)

// GetAllCategorysList(context.Context, *Empty) (*CategoryListResponse, error)
func (s *GoodsServer) GetAllCategorysList(ctx context.Context, req *proto.Empty) (*proto.CategoryListResponse, error) {
	var categorys []model.Category
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys) //用SubCategory只能查询到第一次层级的分类的子分类，所以要用SubCategory.SubCategory来查询所有子分类
	b, _ := json.Marshal(categorys)

	return &proto.CategoryListResponse{JsonData: string(b)}, nil
}

func (s *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	categoryListResponse := proto.SubCategoryListResponse{}
	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	categoryListResponse.Info = &proto.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		Level:          category.Level,
		IsTab:          category.IsTab,
		ParentCategory: category.ParentCategoryID,
	}
	var subCategorys []model.Category
	var subCategoryResponse []*proto.CategoryInfoResponse
	//preloads := "SubCategory"
	//if category.Level == 1 {
	//	preloads = "SubCategory.SubCategory"
	//}
	global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Find(&subCategorys)
	for _, subCategory := range subCategorys {
		subCategoryResponse = append(subCategoryResponse, &proto.CategoryInfoResponse{
			Id:             subCategory.ID,
			Name:           subCategory.Name,
			Level:          subCategory.Level,
			IsTab:          subCategory.IsTab,
			ParentCategory: subCategory.ParentCategoryID,
		})
	}
	categoryListResponse.SubCategorys = subCategoryResponse
	return &categoryListResponse, nil

}
func (s *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	category := model.Category{}
	category.Name = req.Name
	category.Level = req.Level
	if category.Level != 1 {
		category.ParentCategoryID = req.ParentCategory
	}

	category.IsTab = req.IsTab

	global.DB.Save(&category)
	return &proto.CategoryInfoResponse{Id: int32(category.ID)}, nil
}

func (s *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*proto.Empty, error) {
	if result := global.DB.Delete(&model.Category{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	return &proto.Empty{}, nil
}

func (s *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.Empty, error) {
	var category model.Category

	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentCategory != 0 {
		category.ParentCategoryID = req.ParentCategory
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	if req.IsTab {
		category.IsTab = req.IsTab
	}

	global.DB.Save(&category)

	return &proto.Empty{}, nil
}

//CreateCategory(context.Context, *CategoryInfoRequest) (*CategoryInfoResponse, error)
//DeleteCategory(context.Context, *DeleteCategoryRequest) (*Empty, error)
//UpdateCategory(context.Context, *CategoryInfoRequest) (*Empty, error)
