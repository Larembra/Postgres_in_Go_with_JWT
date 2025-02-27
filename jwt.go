package main

import (
	//"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	//"strings"
	//"time"
	//"encoding/json"
	//"database/sql"
	//"encoding/json"
	//"fmt"
	//"gorm.io/driver/postgres"
	//"gorm.io/gorm"
	//"log"
	//"net/http"
	"time"
	//_ "github.com/jackc/pgx/v5/stdlib"
)

//type Book struct {
//	ID     string `json:"id"`
//	Title  string `json:"title"`
//	Author string `json:"author"`
//}
//
//var books = []Book{
//	{ID: "1", Title: "1984", Author: "George Orwell"},
//	{ID: "2", Title: "Brave New World", Author: "Aldous Huxley"},
//	{ID: "3", Title: "Fahrenheit 451", Author: "Ray Bradbury"},
//}

type Provider struct {
	Idprovider        string `json:"idprovider" gorm:"primaryKey"`
	Phone             string `json:"phone"`
	Email             string `json:"email"`
	Compreprlastname  string `json:"compreprlastname"`
	Compreprfirstname string `json:"compreprfirstname"`
}

type Ingredient struct {
	Id         string `json:"id" gorm:"primaryKey"`
	Idprovider string `json:"idprovider"`
	Name       string `json:"name"`
}

type Dish struct {
	Iddish      string `json:"iddish" gorm:"primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Cost        int    `json:"cost"`
	Weight      int    `json:"weight"`
	Photo       []byte `json:"photo"`
	Bonuses     int    `json:"bonuses"`
}

type Employee struct {
	Id               string    `json:"id" gorm:"primaryKey"`
	Lastname         string    `json:"lastname"`
	Firstname        string    `json:"firstname"`
	Patronymic       string    `json:"patronymic"`
	Passseries       string    `json:"passseries"`
	Passumber        string    `json:"passumber"`
	Passdateofissue  time.Time `json:"passdateofissue"`
	Passbywhomissued string    `json:"passbywhomissued"`
	Idpost           string    `json:"idpost"`
}

type Order struct {
	Idorder    string    `json:"idorder" gorm:"primaryKey"`
	Datetime   time.Time `json:"datetime"`
	Idemployee string    `json:"idemployee"`
	Idcustomer string    `json:"idcustomer"`
	Idcoupon   string    `json:"idcoupon"`
}

type Coupon struct {
	Idcoupon  string `json:"idcoupon" gorm:"primaryKey"`
	Name      string `json:"name"`
	Influence string `json:"influence"`
	Discount  int    `json:"discount"`
}

type Customer struct {
	Idcustomer string `json:"idcustomer" gorm:"primaryKey"`
	Lastname   string `json:"lastname"`
	Firstname  string `json:"firstname"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
	Sumbonuses int    `json:"sumbonuses"`
	Address    string `json:"address"`
}

type Post struct {
	Idpost string `json:"idpost" gorm:"primaryKey"`
	Name   string `json:"name"`
	Salary int    `json:"salary"`
}

type Ingredientsdish struct {
	Iddish        string `json:"iddish"`
	Idingredients string `json:"idingredients"`
}

type Orderdishes struct {
	Idorder string `json:"idorder"`
	Iddish  string `json:"iddish"`
}

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = []Credentials{
	{Username: "user1", Password: "pass1"},
	{Username: "user2", Password: "pass2"},
	{Username: "user3", Password: "pass3"},
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func generateToken(username string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func login(c *gin.Context) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	// Здесь добавим простую проверку пароля
	// Поиск пользователя в массиве users
	var found bool
	for _, user := range users {
		if user.Username == creds.Username && user.Password == creds.Password {
			found = true
			break
		}
	}

	// Если пользователь не найден, возвращаем ошибку
	if !found {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "wrong username or password"})
		return
	}

	token, err := generateToken(creds.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

var db *gorm.DB
var providers []Provider

func migratedata(table string, db *gorm.DB) error {
	switch table {
	case "providers":
		result := db.Table("providers").Find(&providers)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func main() {
	dsn := "postgres://postgres:postgres@127.0.0.1:5432/postgres"

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	//db.Raw("SELECT column_name FROM information_schema.columns WHERE table_name = 'providers'").Scan(&columns)
	//fmt.Println(columns)
	db.AutoMigrate(
		&Provider{},
		&Ingredient{},
		&Dish{},
		&Employee{},
		&Order{},
		&Coupon{},
		&Post{},
		&Ingredientsdish{},
		&Orderdishes{},
	)

	err = migratedata("providers", db)
	if err != nil {
		log.Fatalf("Failed to migrate data: %v\n", err)
	}

	router := gin.Default()

	router.POST("/login", login)

	protected := router.Group("/")
	protected.Use(authMiddleware())
	{
		//protected.GET("/books", getBooks)
		//protected.POST("/books", createBook)
		protected.GET("/providers", getproviders)
		protected.POST("/providers", postprovider)
		protected.PUT("/providers/:id", putprovider)
		protected.DELETE("/providers/:id", deleteprovider)
		// другие защищенные маршруты
	}

	router.Run(":8080")
}

//func getBooks(c *gin.Context) {
//	c.JSON(http.StatusOK, books)
//}
//
//func createBook(c *gin.Context) {
//	var newBook Book
//
//	if err := c.BindJSON(&newBook); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
//		return
//	}
//
//	books = append(books, newBook)
//	c.JSON(http.StatusCreated, newBook)
//}

func getproviders(c *gin.Context) {
	c.JSON(http.StatusOK, providers)
}

func postprovider(c *gin.Context) {
	var newprovider Provider
	if err := c.BindJSON(&newprovider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	providers = append(providers, newprovider)
	c.JSON(http.StatusCreated, newprovider)
}

func putprovider(c *gin.Context) {
	var newprovider Provider
	if err := c.BindJSON(&newprovider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	id := c.Param("id")
	for i, p := range providers {
		if p.Idprovider == id {
			providers[i] = newprovider
			c.JSON(http.StatusOK, newprovider)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "provider not found"})
}

func deleteprovider(c *gin.Context) {
	id := c.Param("id")
	for i, p := range providers {
		if p.Idprovider == id {
			providers = append(providers[:i], providers[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "provider deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "provider not found"})
}
