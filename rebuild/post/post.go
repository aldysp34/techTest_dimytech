package post

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Customer struct {
	ID           uint   `gorm:"primaryKey"`
	CustomerName string `gorm:"not null"`
}

type CustomerAddress struct {
	ID         uint   `gorm:"primaryKey"`
	CustomerID uint   `gorm:"not null"`
	Address    string `gorm:"not null"`
}

type Transaction struct {
	gorm.Model
	CustomerID        uint
	CustomerAddressID uint
	TransactionDate   string
	Products          []Product       `gorm:"many2many:transaction_products;"`
	PaymentMethods    []PaymentMethod `gorm:"many2many:transaction_payment_methods;"`
}

type Product struct {
	gorm.Model
	Name  string
	Price float64
}

type PaymentMethod struct {
	gorm.Model
	Name     string
	IsActive bool
}

type Post struct {
	DB *gorm.DB
}

func NewPost(db *gorm.DB) *Post {
	return &Post{DB: db}
}

func (p *Post) Migrate() {
	err := p.DB.AutoMigrate(&Customer{}, &CustomerAddress{}, &Transaction{}, &Product{}, &PaymentMethod{})
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Post) CreateTransaction(c *gin.Context) {
	var transaction Transaction
	err := c.BindJSON(&transaction)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	// Create transaction
	result := p.DB.Create(&transaction)
	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(500, gin.H{"error": "Failed to create transaction"})
		return
	}

	c.JSON(201, gin.H{"transactionId": transaction.ID})
}

func (p *Post) Routes(r *gin.Engine) {
	r.POST("/transaction", p.CreateTransaction)
}
