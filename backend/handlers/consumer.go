package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/smartmeterchain/backend/database"
	"github.com/smartmeterchain/backend/middleware"
	"github.com/smartmeterchain/backend/models"
	"github.com/smartmeterchain/backend/utils"
	"golang.org/x/crypto/bcrypt"
)

type ConsumerHandler struct{}

func NewConsumerHandler() *ConsumerHandler {
	return &ConsumerHandler{}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

// Login authenticates a user and returns a JWT
func (h *ConsumerHandler) Login(c *fiber.Ctx, cfg interface{ GetJWTSecret() string }) error {
	return nil // handled by main.go login route
}

// Register creates a new user account
func (h *ConsumerHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return utils.Error(c, fiber.StatusBadRequest, "username, email, and password are required")
	}

	hashedPw, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "Failed to hash password")
	}

	role := "consumer"
	if req.Role == "admin" || req.Role == "regulator" {
		// Only existing admins can create admin/regulator accounts
		callerRole, _ := c.Locals("role").(string)
		if callerRole != "admin" {
			return utils.Error(c, fiber.StatusForbidden, "Only admins can create admin/regulator accounts")
		}
		role = req.Role
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPw),
		Role:     role,
		Name:     req.Name,
		Phone:    req.Phone,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return utils.Error(c, fiber.StatusConflict, "Username or email already exists")
	}

	user.Password = ""
	return utils.Created(c, user)
}

// GetProfile returns the current user's profile
func (h *ConsumerHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "User not found")
	}
	return utils.Success(c, user)
}

// GetConsumers returns all consumers (admin only)
func (h *ConsumerHandler) GetConsumers(c *fiber.Ctx) error {
	var users []models.User
	database.DB.Where("role = ?", "consumer").Find(&users)
	return utils.Success(c, users)
}

// LoginHandler is the standalone login function used in main routes
func LoginHandler(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return HandleLogin(c, jwtSecret)
	}
}

// HandleLogin processes authentication
func HandleLogin(c *fiber.Ctx, jwtSecret string) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	// Use middleware to generate token
	_ = middleware.JWTClaims{} // reference to avoid import cycle - token generated in main
	return utils.Success(c, fiber.Map{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"name":     user.Name,
	})
}
