package post

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestCreateTransaction(t *testing.T) {
	db := initDB()
	pst := NewPost(db)

	r := gin.Default()
	pst.Routes(r)

	transaction := Transaction{
		CustomerID:        1,
		CustomerAddressID: 1,
		TransactionDate:   "14-06-2023",
		Products: []Product{
			{Name: "So klin", Price: 2000},
			{Name: "Rinso", Price: 1500},
		},
		PaymentMethods: []PaymentMethod{
			{Name: "Paypal", IsActive: true},
			{Name: "Paylater", IsActive: false},
		},
	}

	payload, err := json.Marshal(transaction)

	if err != nil {
		t.Fatal(err)
	}

	// Post Request
	req, err := http.NewRequest(http.MethodPost, "/transaction", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	r.ServeHTTP(response, req)

	assert.Equal(t, http.StatusCreated, response.Code)

	var responseJson map[string]interface{}
	err = json.Unmarshal(response.Body.Bytes(), &responseJson)

	assert.NoError(t, err)
	responseBody := response.Body.String()
	assert.Contains(t, responseBody, "transactionId")
	// t.Errorf("Response : '%v'", reflect.TypeOf(response.Body))
}

func initDB() *gorm.DB {
	dsn := "root:@/dimytech?parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&Customer{}, &CustomerAddress{}, &Transaction{}, &Product{}, &PaymentMethod{})
	if err != nil {
		panic(err)
	}
	return db
}
