package tests

import (
	"bytes"
	"encoding/json"
	"go-api/config"
	"go-api/controllers"
	"go-api/models"
	"go-api/routes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	// Use in-memory SQLite for tests
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	db := config.InitDB(":memory:", logger)
	db.AutoMigrate(&models.User{})
	return db
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	db := setupTestDB()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	userController := controllers.NewUserController(db, logger)

	router := gin.New()
	routes.SetupRoutes(router, userController)

	return router
}

func TestGetUsers(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var users []models.User
	err := json.Unmarshal(w.Body.Bytes(), &users)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(users)) // Empty initially
}

func TestCreateUser(t *testing.T) {
	router := setupTestRouter()

	user := models.User{
		Name:  "Test User",
		Email: "test@example.com",
	}

	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdUser models.User
	err := json.Unmarshal(w.Body.Bytes(), &createdUser)
	assert.NoError(t, err)
	assert.Equal(t, "Test User", createdUser.Name)
	assert.Equal(t, "test@example.com", createdUser.Email)
	assert.NotZero(t, createdUser.ID)
}

func TestGetUser(t *testing.T) {
	router := setupTestRouter()

	// First create a user
	user := models.User{Name: "Test User", Email: "test@example.com"}
	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var createdUser models.User
	json.Unmarshal(w.Body.Bytes(), &createdUser)

	// Now get the user
	req, _ = http.NewRequest("GET", "/api/v1/users/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var fetchedUser models.User
	err := json.Unmarshal(w.Body.Bytes(), &fetchedUser)
	assert.NoError(t, err)
	assert.Equal(t, createdUser.Name, fetchedUser.Name)
	assert.Equal(t, createdUser.Email, fetchedUser.Email)
}

func TestUserNotFound(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/users/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestInvalidUserID(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/users/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
