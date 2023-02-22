package main

import (
	"encoding/json"
	"fa-be/db"
	"fa-be/models"
	"fa-be/utils"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func New(db *gorm.DB) Repository {
	return Repository{db}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidateJWT(tokenStr string) string {
	claims := jwt.MapClaims{}
	token, _ := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if token == nil {
		return ""
	}
	return claims["username"].(string)
}

func (r *Repository) ScrapeApps(c *fiber.Ctx) error {
	payload := struct {
		Locales []string `json:"locales" form:"locales"`
	}{}

	err := c.BodyParser(&payload)

	if c.Query("searchTerm") == "" {
		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "no search term supplied"},
		)
	}

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "No locales supplied"},
		)
	}

	data, _ := utils.ScrapeData(payload.Locales, c.Query("searchTerm"), c.Query("keyword"))

	byteArr, _ := json.Marshal(data)
	return c.SendString(string(byteArr))
}

func (r *Repository) GetUser(c *fiber.Ctx) error {
	token := c.Cookies("token")
	isValid := ValidateJWT(token)

	if isValid == "" {
		return c.Status(http.StatusOK).JSON(
			&fiber.Map{"username": ""},
		)
	}
	return c.Status(http.StatusOK).JSON(
		&fiber.Map{"username": isValid},
	)
}

func (r *Repository) Login(c *fiber.Ctx) error {
	payload := models.User{}
	err := c.BodyParser(&payload)
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "bad request, login failed"},
		)
		return err
	}

	var user = models.User{}
	fmt.Println(payload.Username)
	result := r.DB.First(&user, "username = ?", payload.Username)
	fmt.Println(result)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"status": "fail", "message": "Invalid email or Password"},
		)
	}
	fmt.Println(payload.Password, user.Password)
	isAuthed := CheckPasswordHash(payload.Password, user.Password)
	if !isAuthed {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"status": "fail", "message": "Invalid email or Password"},
		)
	}

	tokenByte := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)

	claims["exp"] = now.Add(86400000).Unix()
	claims["username"] = payload.Username

	tokenString, err := tokenByte.SignedString([]byte(os.Getenv("JWT_KEY")))

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
			fiber.Map{"status": "fail", "message": fmt.Sprintf("generating JWT Token failed: %v", err)},
		)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   10000 * 60,
		Secure:   false,
		HTTPOnly: true,
		Domain:   os.Getenv("DOMAIN"),
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func (r *Repository) Register(c *fiber.Ctx) error {
	user := models.User{}
	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "bad request, login failed"},
		)
	}
	user.Password, _ = HashPassword(user.Password)
	errCreate := r.DB.Create(&user).Error
	if errCreate != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not register user"},
		)
		return errCreate
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "registered"},
	)

	return nil

}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/search", r.ScrapeApps)
	api.Get("/get-user", r.GetUser)
	api.Post("/login", r.Login)
	api.Post("/register", r.Register)
}

func main() {
	app := fiber.New()
	// app.Use(cors.New())
	DB := db.Init()
	r := New(DB)

	r.SetupRoutes(app)

	app.Listen(":9000")
}
