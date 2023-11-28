package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/alwilion/database"
	"github.com/alwilion/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "secret"

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}
	if data["password"] == "" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid pass",
		})
	}
	//fmt.Println("password:", data["password"])
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {

		errorMsg := ""
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println("field:", err.Field())
			if err.Field() == "Email" {
				errorMsg = "Issue In Email"
			} else if err.Field() == "Password" {
				errorMsg = "Issue in Password"

			} else if err.Field() == "Name" {
				errorMsg = "Issue in Name Field"
			}
		}
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": errorMsg,
		})

	} else {
		result := database.DB.Create(&user)

		// Check for errors during insertion
		if result.Error != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "failed to insert record into the database",
			})
		}
		return c.JSON(&user)
		//fmt.Println("Validation passed!")
	}

}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.ID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{"message": "User not found"})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	return c.JSON(fiber.Map{
		"message": "success",
		"token":   token,
	})
	//return c.JSON(user)

}

func User(c *fiber.Ctx) error {
	token, err := authentication(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	database.DB.Where("id = ?", claims.Issuer).First(&user)
	fmt.Printf("type of id %T", user.ID)
	return c.JSON(user)
}

func AddProduct(c *fiber.Ctx) error {

	fmt.Print("Add Product")
	_, err := authentication(c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}
	fmt.Print("Add Product1")
	//claims := token.Claims.(*jwt.StandardClaims)

	//var user models.User

	//database.DB.Where("id = ?", claims.Issuer).Find(&user)

	data := make(map[string]interface{})
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	fmt.Print("Add Product2")
	var product models.Product
	if value, ok := data["price"]; ok {
		// Type assertion to float64
		price, ok := value.(float64)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "Amount Key has an unexpected type",
			})
		}

		// Successfully asserted to float64, now you can use 'amount'
		product.Price = price
	} else {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Price Key Missing",
		})
	}
	if value, ok := data["description"]; ok {
		// Type assertion to float64
		description, ok := value.(string)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "Description Key has an unexpected type",
			})
		}

		// Successfully asserted to float64, now you can use 'amount'
		product.Description = description
	} else {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Description Key Missing",
		})
	}
	if value, ok := data["name"]; ok {
		// Type assertion to float64
		name, ok := value.(string)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "Name Key has an unexpected type",
			})
		}

		// Successfully asserted to float64, now you can use 'amount'
		product.Name = name
	} else {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Name Key Missing",
		})
	}

	result := database.DB.Create(&product)

	// Check for errors during insertion
	if result.Error != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "failed to insert record into the database",
		})
	}
	return c.JSON(product)
}
func GetProductList(c *fiber.Ctx) error {

	_, err := authentication(c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	// Query all products
	var products []models.Product

	if err := database.DB.Find(&products).Error; err != nil {
		panic("Failed to fetch products")
	}

	// Print the list of products
	for _, product := range products {
		fmt.Printf("ID: %d, Name: %s, Price: %.2f\n", product.ID, product.Name, product.Price)
	}

	return c.JSON(products)

}

func GetProductById(c *fiber.Ctx) error {

	_, err := authentication(c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	// Query all products
	var product models.Product

	if err := database.DB.First(&product, c.Params("id")).Error; err != nil {
		return c.JSON(fiber.Map{"message": "Product Not Found"})
	}

	return c.JSON(product)

}

func DeleteProductById(c *fiber.Ctx) error {

	_, err := authentication(c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	// Query all products
	var product models.Product

	if err := database.DB.First(&product, c.Params("id")).Error; err != nil {
		return c.JSON(fiber.Map{"message": "Product Not Found"})
	}
	database.DB.Delete(&product)
	return c.JSON(fiber.Map{"message": "Product deleted successfully"})

}
func UpdateProduct(c *fiber.Ctx) error {

	_, err := authentication(c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	var updatedProduct models.Product

	if err := database.DB.First(&updatedProduct, c.Params("id")).Error; err != nil {
		return c.JSON(fiber.Map{"message": "Product Not Found"})
	}

	data := make(map[string]interface{})
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	//var product models.Product
	if value, ok := data["price"]; ok {
		// Type assertion to float64
		price, ok := value.(float64)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "Amount Key has an unexpected type",
			})
		}

		// Successfully asserted to float64, now you can use 'amount'
		updatedProduct.Price = price
	}
	if value, ok := data["description"]; ok {
		// Type assertion to float64
		description, ok := value.(string)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "Description Key has an unexpected type",
			})
		}

		// Successfully asserted to float64, now you can use 'amount'
		updatedProduct.Description = description
	}
	if value, ok := data["name"]; ok {
		// Type assertion to float64
		name, ok := value.(string)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "Name Key has an unexpected type",
			})
		}

		// Successfully asserted to float64, now you can use 'amount'
		updatedProduct.Name = name
	}

	result := database.DB.Save(&updatedProduct)

	// Check for errors during insertion
	if result.Error != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "failed to update record into the database",
		})
	}
	return c.JSON(updatedProduct)
}
