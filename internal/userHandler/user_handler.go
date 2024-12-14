package handler

import (
	"fmt"
	"net/http"
	config "w3/gc3/config/database"
	
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"github.com/jackc/pgconn"
)

// Users struct
type Users struct {
	ID            int       `json:"id"`               // Primary key, auto-incremented
	FullName      string    `json:"full_name"`        // Full name of the user
	Email         string    `json:"email"`           // Email address, unique
	Username      string    `json:"username"`        // Username, unique
	Password      string    `json:"password"`        // Password for the account
	Age           int       `json:"age"`             // Age of the user
	LastLoginDate time.Time `json:"last_login_date"` // Date and time of the last login
	JwtToken      string    `json:"jwt_token"`       // JWT token for authentication
}

// RegisterRequest struct
type RegisterRequest struct {
	FullName string `json:"full_name" validate:"required,full_name"` // Full name of the user
	Email    string `json:"email" validate:"required,email"`        // Email address
	Username string `json:"username" validate:"required,username"`  // Username for the user
	Password string `json:"password" validate:"required,password"`  // Password for the account
	Age      int    `json:"age" validate:"required,gt=0"`           // Age of the user, must be greater than 0
}

// login request struct
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// login response: token
type LoginResponse struct {
	Token string `json:"token"`
}

var jwtSecret = []byte("12345")

// @Summary Register a new user
// @Description Create a new user account by providing the required information
// @Tags Users
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration data"
// @Success 201 {object} map[string]interface{} "User created successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Router /users/register [post]
func Register(c echo.Context) error {
    var req RegisterRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid Request"})
    }

	// hash the password
    hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal Server Error"})
    }

    // queries to insert to both users and customers db
	users_query := "INSERT INTO users (full_name, email, username, password, age) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	
	var userID int
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// query row 1: insert to users 
	err = config.Pool.QueryRow(ctx, users_query, req.FullName, req.Email, req.Username, string(hashPassword), req.Age).Scan(&userID)
	if err != nil {
		fmt.Println("Error inserting into users table:", err)

		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" { // Unique violation (email already registered)
				return c.JSON(http.StatusBadRequest, map[string]string{"message": "Email already registered"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal Server Error"})
	}

    return c.JSON(http.StatusOK, map[string]interface{}{
        "message": "User registered successfully",
        "user_id": string(userID),
        "email": req.Email,
    })
}

// @Summary Login an existing user
// @Description Authenticate a user by providing valid credentials
// @Tags Users
// @Accept json
// @Produce json
// @Param request body LoginRequest true "User login data"
// @Success 200 {object} map[string]interface{} "Authentication successful"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /users/login [post]
func Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message":"Invalid Request"})
	}
	
	var user Users
	query := "SELECT id, email, password FROM users WHERE email = $1"
	err := config.Pool.QueryRow(context.Background(), query, req.Email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid email or password"})
	}

	// compare password to see if it matches the student password provided
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid email or password"})
	}

	// create new jwt claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     jwt.NewNumericDate(time.Now().Add(72 * time.Hour)), // Use `jwt.NewNumericDate` for expiry
	})
	
	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Invalid Generate Token"})
	}

	// return ok status and login response
	return c.JSON(http.StatusOK, LoginResponse{Token: tokenString})
}