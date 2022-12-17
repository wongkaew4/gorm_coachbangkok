package main

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SqlLogger struct {
	logger.Interface
}

// function สำหรับทำ dry run คือทดสอบการ create table โดยจะไม่สร้าง database ใหม่
func (l SqlLogger) Trace(ct context.Context, begin time.Time, fc func() (sal string, rowsAffected int64), err error) {
	sql, _ := fc()
	fmt.Printf("%v\n===================================================\n", sql)
}

var db *gorm.DB

func main() {
	dsn := "root@P@ssw0rd@tcp(13.16.163.73:3306)/fe?parseTime=true" //fe คือชื่อ database ที่สร้างใน server ต้นทางปิด server แล้ว จึงใช้ไม่ได้
	dial := mysql.Open(dsn)                                         //mysql.open return Dialector สร้างตัวแร dial มารับ
	var err error
	db, err = gorm.Open(dial, &gorm.Config{
		Logger: &SqlLogger{},
		DryRun: false, //true คือ create database table แต่ไม่สร้างจริง ส่วน false คือสร้าง table ใน database จริง
	})
	if err != nil {
		panic(err)
	}

	// db.AutoMigrate(Gender{}, Test{})
	// CreateGender("Male")
	// GetGenders()
	// GetGender(1)
	// UpdateGender(3, "xxxx")
	DeleteGender(4)
}

func DeleteGender(id uint) {
	tx := db.Delete(&Gender{}, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println("Deleted")
	GetGender(id)

}

func UpdateGender(id uint, name string) {
	gender := Gender{}
	tx := db.First(&gender, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	gender.Name = name
	tx = db.Save(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	GetGender(id)

}

func GetGenders() { //ขอดู gender ทั้งหมด
	genders := []Gender{}
	tx := db.Find(&genders)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(genders)
}
func GetGender(id uint) { //ขอดู gender ตาม id
	gender := Gender{}
	tx := db.First(&gender, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)

}

func CreateGender(name string) {
	gender := Gender{Name: name}
	tx := db.Create(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

type Gender struct {
	ID   uint
	Name string `gorm:"unigue;size(10)"`
}

type Test struct {
	gorm.Model        //gorm ให้มาอยู่แล้วสามารถเลือกใช้ได้
	Code       uint   `gorm:"Comment:This is Code"`
	Name       string `gorm:"column:myname;size:20;unigue;default:Hello;not null"`
}
