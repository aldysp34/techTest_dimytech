package main

import (
	"fmt"
	"log"

	"dimy-tech-test-go-version/post"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@/dimytech?parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	pst := post.NewPost(db)
	pst.Migrate()
	pst.Routes(r)

	fmt.Println("Server starting on port : 3000")
}
