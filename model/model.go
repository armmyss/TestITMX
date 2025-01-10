package model

import (
	"errors"
	"gorm.io/gorm"
)

type Customers struct {
	ID   uint
	Name string `json:"name"`
	Age  uint   `json:"age"`
}

// เป็น func ที่ init data ของ customer เข้าไปโดยมีการเช็คก่อนว่าใน db ไม่มีข้อมูล
func InitData(db *gorm.DB) error {
	var Count int64
	db.Model(&Customers{}).Count(&Count)
	if Count == 0 {
		customers := []Customers{
			{Name: "Arm", Age: 21},
			{Name: "Bob", Age: 22},
			{Name: "Alice", Age: 23},
		}
		result := db.Create(&customers)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func GetCustomers(db *gorm.DB) ([]Customers, error) {
	var customers []Customers
	result := db.Find(&customers)
	if result.Error != nil {
		return []Customers{}, result.Error
	}
	return customers, nil
}

func GetCustomer(db *gorm.DB, customerID uint) (*Customers, error) {
	var customer Customers
	result := db.First(&customer, customerID)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return nil, errors.New("customer not found")
		}
	}
	return &customer, nil
}

func AddCustomers(db *gorm.DB, customers *Customers) error {
	result := db.Create(customers)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateCustomer(db *gorm.DB, customer *Customers) error {
	var existingCustomer Customers
	if err := db.First(&existingCustomer, customer.ID).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return errors.New("customer not found")
		}
		return err
	}

	result := db.Model(&Customers{}).Where("id = ?", customer.ID).Updates(customer)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteCustomer(db *gorm.DB, customerID uint) error {
	var customer Customers
	result := db.Delete(&customer, customerID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("customer not found")
	}
	return nil
}
