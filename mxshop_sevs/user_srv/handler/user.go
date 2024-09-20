package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"mxshop_sevs/user_srv/global"
	"mxshop_sevs/user_srv/model"
	"mxshop_sevs/user_srv/model/proto"
	"strings"
	"time"
)

// UserServer 结构体
type UserServer struct{}

// ModelToResponse 函数将用户模型转换为用户信息响应
// 在grpc的message中字段有默认值，你不能随便赋值nil，容易出错
// 因此需要清楚，那些字段是有默认值的
func ModelToResponse(user model.User) proto.UserInfoResponse {
	// 创建一个用户信息响应并将用户模型中的数据填充进去
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,          // 用户ID
		PassWord: user.Password,    // 用户密码
		NickName: user.NickName,    // 用户昵称
		Gender:   user.Gender,      // 用户性别
		Role:     int32(user.Role), // 用户角色
		Mobile:   user.Mobile,
	}
	// 如果用户的生日不为空，则将其转换为Unix时间戳并填充到响应中
	if user.Birthday != nil {
		userInfoRsp.BirthDay = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}

// Paginate 函数处理分页逻辑
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 如果页面为0，则默认为第一页
		if page == 0 {
			page = 1
		}

		// 对页面大小进行判断和处理，保证在合理范围内
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		// 计算偏移量并进行分页查询
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// GetUserList 函数获取带有分页的用户列表
func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	// 从数据库中查找所有用户
	result := global.DB.Find(&users)
	// 如果查找过程中出现错误，则返回错误
	if result.Error != nil {
		return nil, result.Error
	}
	fmt.Println("用户列表")
	// 创建一个用户列表响应并填充总数
	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected) //返回找到的记录数
	// 使用分页作用域查找用户并填充到响应中
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)
	for _, user := range users {
		userInfoRsp := ModelToResponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}

	return rsp, nil
}
func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	// 从数据库通过手机号查找用户
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error

	}
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}
func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}
func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	// 这里可以实现用户创建逻辑
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}
	//这里的一个小细节，使用first方法，如果找到会填充user，没找到就不会填充user，user仍然是个安全的变量可以继续使用
	user.Mobile = req.Mobile
	user.NickName = req.NickName
	//密码
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(req.PassWord, options)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil

}
func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*proto.Emptymessage, error) {
	//个人中心更新用户信息
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	birthDay := time.Unix(int64(req.BirthDay), 0)
	user.NickName = req.NickName
	user.Birthday = &birthDay
	user.Gender = req.Gender
	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &proto.Emptymessage{}, nil

}
func (s *UserServer) CheckPassWord(ctx context.Context, req *proto.PassWordCheckInfo) (*proto.CheckResponse, error) {
	// 密码校验
	options := &password.Options{16, 100, 32, sha512.New}
	passwordInfo := strings.Split(req.EncryptedPassword, "$")
	fmt.Println(passwordInfo)
	check := password.Verify(req.Password, passwordInfo[2], passwordInfo[3], options)
	return &proto.CheckResponse{Success: check}, nil
}
