package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

var ctx = context.Background()

type User struct {
	RealName string `json:"realname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func sha1Hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func main() {
	// Setup Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	app := fiber.New()

	app.Post("/register", func(c *fiber.Ctx) error {
		var u User
		if err := c.BodyParser(&u); err != nil {
			return c.Status(400).SendString("Invalid input")
		}
		u.Password = sha1Hash(u.Password)

		key := fmt.Sprintf("login_%s", u.Email)
		data, _ := json.Marshal(u)
		err := rdb.Set(ctx, key, data, 0).Err()
		if err != nil {
			return c.Status(500).SendString("Redis error")
		}
		return c.JSON(fiber.Map{"message": "User registered"})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).SendString("Invalid input")
		}

		key := fmt.Sprintf("login_%s", input.Email)
		val, err := rdb.Get(ctx, key).Result()
		if err != nil {
			return c.Status(401).SendString("User not found")
		}

		var u User
		json.Unmarshal([]byte(val), &u)

		if u.Password != sha1Hash(input.Password) {
			return c.Status(401).SendString("Invalid password")
		}

		return c.JSON(fiber.Map{"message": "Login successful", "realname": u.RealName})
	})

	log.Fatal(app.Listen(":3000"))
}
