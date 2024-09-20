package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop_sevs/goods_srv/global"
	"mxshop_sevs/goods_srv/model"
	"mxshop_sevs/goods_srv/model/proto"
)

// 品牌和轮播图
func (s *GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	brandListResponse := proto.BrandListResponse{}
	var brands []model.Brands
	result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}
	var total int64
	global.DB.Model(&model.Brands{}).Count(&total)
	brandListResponse.Total = int32(total)
	var brandResponses []*proto.BrandInfoResponse
	for _, brand := range brands {
		brandResponses = append(brandResponses, &proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})

	}
	brandListResponse.Data = brandResponses
	return &brandListResponse, nil

}
func (s *GoodsServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	//新建品牌

	if result := global.DB.First(&model.Brands{}); result.RowsAffected == 1 {
		return nil, status.Error(codes.AlreadyExists, "品牌已存在")
	}
	brand := model.Brands{
		Name: req.Name,
		Logo: req.Logo,
	}
	global.DB.Save(brand)
	return &proto.BrandInfoResponse{
		Id: brand.ID,
	}, nil
}
func (s *GoodsServer) DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*proto.Empty, error) {
	if result := global.DB.Delete(&model.Brands{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "品牌不存在")
	}
	return &proto.Empty{}, nil
}

func (s *GoodsServer) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.Empty, error) {
	brands := model.Brands{}
	if result := global.DB.First(&brands); result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "品牌不存在")
	}
	if req.Name != "" {
		brands.Name = req.Name
	}
	if req.Logo != "" {
		brands.Logo = req.Logo
	}
	global.DB.Save(&brands)
	return &proto.Empty{}, nil
}
