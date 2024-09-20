package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"mxshop_sevs/user_srv/handler"
	"mxshop_sevs/user_srv/model"
	"mxshop_sevs/user_srv/model/proto"
	"os"
	"time"
)

func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))

}

func main() {
	dsn := "root:root@tcp(192.168.43.170:3306)/wu_yi_gao_su?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	//options := &password.Options{16, 100, 32, sha512.New}
	//salt, encodedPwd := password.Encode("admin123", options)
	//newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	db.AutoMigrate(&model.User{})

	//var users []model.User
	//db.Find(&users)
	//
	//// 处理查询结果
	//for _, user := range users {
	//	fmt.Println(user.Mobile)
	//}
	for i := 0; i < 10; i++ {
		var user handler.UserServer
		rsp, err := user.CreateUser(context.Background(), &proto.CreateUserInfo{
			NickName: fmt.Sprintf("bobbby%d", i),
			Mobile:   fmt.Sprintf("1899607173%d", i),
			PassWord: fmt.Sprintf("admin123"),
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("%S", rsp.Mobile)
		db.Save(&rsp)
	}

}

// 自动迁移表结构

//options := &password.Options{16, 100, 32, sha512.New}
//salt, encodedPwd := password.Encode("generic password", options)
//newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
//fmt.Println(newPassword)
//
//passwordInfo := strings.Split(newPassword, "$")
//fmt.Println(passwordInfo)
//check := password.Verify("generic password", passwordInfo[2], passwordInfo[3], options)
//fmt.Println(check)
