package lesson01

import (
	"base_go/03_gorm/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
)

func Run(db *gorm.DB) {
	// 创建表
	// db.AutoMigrate(&models.User{})
	// 新增
	//birthday := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	//user := models.User{
	//	Username: "john_doe",
	//	Email:    "john@example.com",
	//	Password: "securepassword",
	//	Age:      30,
	//	Birthday: &birthday,
	//	Status:   true,
	//}
	//// 新增
	//add(db, &user)
	// 删
	//softDeleteUser(db, 1)
	//permanentDeleteUser(db, 1)
	//
	//// 查单个
	//u := getUserByID(db, 1)
	//fmt.Println("user= ", u)
	// 查一个列表
	list := getUserList(db, 1)
	for i, user := range list {
		fmt.Println("userSortID= ", i, "user=", user)
	}
	// 更新
	//newEmail := "new_john@example.com"
	//updateUser(db, 3, newEmail)

}

func add(db *gorm.DB, user *models.User) {
	result := db.Create(user)
	if result.Error != nil {
		log.Println("Failed to create user:", result.Error)
		return
	}
	fmt.Printf("User created successfully with ID: %d\n", user.ID)
}

func getUserByID(db *gorm.DB, id uint) *models.User {
	var user models.User
	result := db.Where(fmt.Print("id=", id)).First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Println("User not found")
		} else {
			log.Println("Failed to get user:", result.Error)
		}
	}
	fmt.Printf("User found: %+v\n", user)
	return &user
}

func getUserList(db *gorm.DB, id uint) []models.User {
	var users []models.User
	// db.Select("id,username,email")
	// 分页
	page, pageSize := 1, 10
	db = db.Where("ID >= ? ", 3)
	db = db.Where("username= ?", "john_doe2")
	result := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("id asc").Find(&users)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Println("User not found")
		} else {
			log.Println("Failed to get user:", result.Error)
		}
	}
	return users
}

func softDeleteUser(db *gorm.DB, id uint) {
	result := db.Delete(&models.User{}, id)
	if result.Error != nil {
		log.Println("Failed to delete user:", result.Error)
		return
	}
	fmt.Printf("User with ID %d soft deleted successfully\n", id)
}

func permanentDeleteUser(db *gorm.DB, id uint) {
	result := db.Unscoped().Delete(&models.User{}, id)
	if result.Error != nil {
		log.Println("Failed to permanently delete user:", result.Error)
		return
	}
	fmt.Printf("User with ID %d permanently deleted\n", id)
}

func updateUser(db *gorm.DB, id uint, newEmail string) {
	// 先查询用户
	var user models.User
	result := db.First(&user, id)
	if result.Error != nil {
		log.Println("Failed to find user:", result.Error)
		return
	}
	// 更新用户信息
	result = db.Model(&user).Updates(models.User{
		Email:  newEmail,
		Status: false,
	})
	if result.Error != nil {
		log.Println("Failed to update user:", result.Error)
		return
	}

	fmt.Println("User updated successfully")
}
