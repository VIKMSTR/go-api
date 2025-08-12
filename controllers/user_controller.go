package controllers

import (
	"go-api/models"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB     *gorm.DB
	Logger *slog.Logger
}

func NewUserController(db *gorm.DB, logger *slog.Logger) *UserController {
	return &UserController{
		DB:     db,
		Logger: logger,
	}
}

// GetUsers godoc
// @Summary Get all users
// @Description Get list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func (uc *UserController) GetUsers(c *gin.Context) {
	var users []models.User
	result := uc.DB.Find(&users)

	if result.Error != nil {
		uc.Logger.Error("Failed to fetch users", "error", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	uc.Logger.Debug("Successfully fetched users", "count", len(users))
	c.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get a single user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (uc *UserController) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		uc.Logger.Warn("Invalid user ID provided", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	result := uc.DB.First(&user, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			uc.Logger.Info("User not found", "id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		uc.Logger.Error("Database error while fetching user", "error", result.Error, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	uc.Logger.Debug("Successfully fetched user", "id", id, "email", user.Email)
	c.JSON(http.StatusOK, user)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the given data
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /users [post]
func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		uc.Logger.Warn("Invalid JSON data provided", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := uc.DB.Create(&user)
	if result.Error != nil {
		uc.Logger.Error("Failed to create user", "error", result.Error, "email", user.Email)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	uc.Logger.Info("User created successfully", "id", user.ID, "email", user.Email, "name", user.Name)
	c.JSON(http.StatusCreated, user)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user data by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.User true "User data"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [put]
func (uc *UserController) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		uc.Logger.Warn("Invalid user ID provided for update", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	result := uc.DB.First(&user, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			uc.Logger.Info("User not found for update", "id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		uc.Logger.Error("Database error while finding user for update", "error", result.Error, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var updateData models.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		uc.Logger.Warn("Invalid JSON data provided for update", "error", err, "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result = uc.DB.Model(&user).Updates(updateData)
	if result.Error != nil {
		uc.Logger.Error("Failed to update user", "error", result.Error, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	uc.Logger.Info("User updated successfully", "id", user.ID, "email", user.Email)
	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [delete]
func (uc *UserController) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		uc.Logger.Warn("Invalid user ID provided for deletion", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	result := uc.DB.First(&user, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			uc.Logger.Info("User not found for deletion", "id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		uc.Logger.Error("Database error while finding user for deletion", "error", result.Error, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	result = uc.DB.Delete(&user)
	if result.Error != nil {
		uc.Logger.Error("Failed to delete user", "error", result.Error, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	uc.Logger.Info("User deleted successfully", "id", id, "email", user.Email)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
