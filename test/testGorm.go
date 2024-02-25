package main

import (
	"fmt"
	"ginchat/models"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {


	// Migrate the schema
	db.AutoMigrate(&models.UserBasic{})

	// Create
	user := &models.UserBasic{}
	user.Name = "Jack"
	user.LoginTime = time.Date(2021, 7, 1, 0, 0, 0, 0, time.Local)
	user.HeartbeatTime = time.Date(2021, 7, 1, 0, 0, 0, 0, time.Local)
	user.LoginOutTime = time.Date(2021, 7, 1, 0, 0, 0, 0, time.Local)
	db.Create(user)

	// Read
	fmt.Println(db.First(user, 1)) // find product with integer primary key
	// db.First(&user, "code = ?", "D42") // find product with code D42

	// Update - update product's price to 200
	db.Model(user).Update("Password", "1234")
	// Update - update multiple fields
	// db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	// db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	// db.Delete(&product, 1)
}
