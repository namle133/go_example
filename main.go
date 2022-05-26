package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func ConnectDatabase() (db *gorm.DB) {
	dsn := "host=localhost user=postgres password=Namle311 dbname=book port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Product{})
	return db
}

type IProduct interface {
	Create(p *Product) error
	Read() ([]Product, error)
	ReadOne(id string) (*Product, error)
	Delete(id string) error
	Update(id string, price int) error
}

type product struct {
	db *gorm.DB
}

func CreateBind(c *gin.Context) *Product {
	var item *Product
	err := c.BindJSON(&item)
	if err != nil {
		c.Error(err)
	}
	return item
}

func (p *product) Create(item *Product) error {
	err := p.db.Create(item).Error
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (p *product) Read() ([]Product, error) {
	var products []Product
	err := p.db.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func ReadIndented(c *gin.Context, products []Product) {
	c.IndentedJSON(http.StatusOK, products)
}

func ReadOneBind(c *gin.Context) string {
	code := c.Param("id")
	return code
}
func (p *product) ReadOne(id string) (*Product, error) {
	var item *Product
	err := p.db.First(&item, "code = ?", id).Error
	if err != nil {
		return nil, err
	}
	return item, nil

}

func ReadOneIndented(c *gin.Context, product *Product) {
	c.IndentedJSON(http.StatusOK, product)
}

func DeleteOneBind(c *gin.Context) string {
	code := c.Param("id")
	return code
}

func (p *product) Delete(id string) error {
	var item *Product
	err := p.db.Where("code = ?", id).Delete(&item).Error
	if err != nil {
		return err
	}
	return nil

}
func DeleteIndented(c *gin.Context, id string) {
	c.String(http.StatusOK, "Deleted product with code = %s", id)
}

func UpdateOneBind(c *gin.Context) (string, int) {
	code := c.Param("id")
	pstr := c.Query("price")
	fmt.Println(code, pstr)
	price, err := strconv.Atoi(pstr)
	if err != nil {
		c.Error(err)
		return "Error", 0
	}
	return code, price
}

func (p *product) Update(id string, price int) error {
	var item *Product
	err := p.db.Model(&item).Where("code = ?", id).Update("price", price).Error
	fmt.Println(err)
	fmt.Println(id)
	fmt.Println(price)
	if err != nil {
		return err
	}
	return nil
}

func UpdateIndented(c *gin.Context, id string) {
	c.String(http.StatusOK, "Deleted product with code = %s", id)
}

func main() {
	router := gin.Default()
	p := &product{db: ConnectDatabase()}
	var P IProduct = p

	router.POST("/create", func(context *gin.Context) {
		it := CreateBind(context)
		P.Create(it)
	})
	router.GET("/read", func(context *gin.Context) {
		products, err := P.Read()
		if err != nil {
			context.JSON(400, http.StatusBadRequest)
			return
		}
		ReadIndented(context, products)
	})
	router.GET("/read/:id", func(context *gin.Context) {
		id := ReadOneBind(context)
		item, err := P.ReadOne(id)
		if err != nil {
			context.JSON(400, http.StatusBadRequest)
			return
		}
		ReadOneIndented(context, item)
	})
	router.DELETE("/delete/:id", func(context *gin.Context) {
		id := DeleteOneBind(context)
		err := P.Delete(id)
		if err != nil {
			context.JSON(400, http.StatusBadRequest)
			return
		}
		DeleteIndented(context, id)

	})
	router.PUT("/update/:id", func(context *gin.Context) {
		id, price := UpdateOneBind(context)
		err := P.Update(id, price)
		if err != nil {
			context.JSON(400, http.StatusBadRequest)
			return
		}
		UpdateIndented(context, id)

	})

	router.Run(":8000")
}
