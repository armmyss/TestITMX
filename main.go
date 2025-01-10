package main

import (
	"TestITMX/handler"
	"TestITMX/model"
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	db, err := gorm.Open(sqlite.Open("customers.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.Customers{}) //สร้าง db โดยใช้ feature AutoMigrate
	model.InitData(db)                 //หลังจากที่ได้ db แล้วจะเรียกใช้ func initData ที่สร้างไว้หากไม่มี data ก็จะทำการ create แต่ถ้ามี data อนู่แล้วก็จะไม่มีอะไรเกิดขึ้น

	app := fiber.New()
	handler.SetupRoutes(app, db)
	app.Listen(":3000")
}
