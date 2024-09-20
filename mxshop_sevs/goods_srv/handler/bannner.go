package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop_sevs/goods_srv/global"
	"mxshop_sevs/goods_srv/model"
	"mxshop_sevs/goods_srv/model/proto"
)

// BannerList(context.Context, *Empty) (*BannerListResponse, error)
// CreateBanner(context.Context, *BannerRequest) (*BannerResponse, error)
// DeleteBanner(context.Context, *BannerRequest) (*Empty, error)
// UpdateBanner(context.Context, *BannerRequest) (*Empty, error)
func (s *GoodsServer) BannerList(ctx context.Context, req *proto.Empty) (*proto.BannerListResponse, error) {
	bannerListResponse := proto.BannerListResponse{}

	var banners []model.Banner
	result := global.DB.Find(&banners)
	bannerListResponse.Total = int32(result.RowsAffected)
	var bannerResponses []*proto.BannerResponse

	for _, banner := range banners {
		bannerResponses = append(bannerResponses, &proto.BannerResponse{
			Id:    banner.ID,
			Image: banner.Image,
			Index: banner.Index,
			Url:   banner.Url,
		})
	}
	bannerListResponse.Data = bannerResponses
	return &bannerListResponse, nil
}
func (s *GoodsServer) CreateBanner(ctx context.Context, req *proto.BannerRequest) (*proto.BannerResponse, error) {
	banner := model.Banner{}
	banner.Image = req.Image
	banner.Index = req.Index
	banner.Url = req.Url
	global.DB.Save(&banner)
	return &proto.BannerResponse{
		Id: banner.ID,
	}, nil
}
func (s *GoodsServer) DeleteBanner(ctx context.Context, req *proto.BannerRequest) (*proto.Empty, error) {
	if result := global.DB.Delete(&model.Banner{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")

	}
	return &proto.Empty{}, nil
}
func (s *GoodsServer) UpdateBanner(ctx context.Context, req *proto.BannerRequest) (*proto.Empty, error) {
	var banner model.Banner
	if result := global.DB.First(&banner, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}
	if req.Url != "" {
		banner.Url = req.Url
	}
	if req.Image != "" {
		banner.Image = req.Image
	}
	if req.Index != 0 {
		banner.Index = req.Index
	}
	global.DB.Save(&banner)
	return &proto.Empty{}, nil
}
