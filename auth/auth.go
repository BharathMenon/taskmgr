package auth

import (
	"net/http"
	"sync"
	"time"

	"github.com/BharathMenon/taskmgr/db"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	//"gorm.io/gorm"
)

var jwtKey = []byte("7Tr5G8xgCqJ1e0fKjN2mYzPqB9wQ4lT5XvR8NwQyTzM=")

type Claims struct {
    Username string
    jwt.StandardClaims
}


var (
    users = make(map[string]string)
    mu    sync.RWMutex
)

func Register(userRepo *db.UserRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            Username string `json:"username" binding:"required"`
            Email    string `json:"email" binding:"required"`
            Password string `json:"password" binding:"required"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        mu.Lock()
        defer mu.Unlock()

        // Check if username already exists
        if existing, _ := userRepo.FindByUsername(req.Username); existing.ID != 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
            return
        }

        // Hash the password
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
            return
        }

        user := db.User{
            Username: req.Username,
            Email:    req.Email,
            Password: string(hashedPassword),
        }

        if err := userRepo.CreateUser(&user); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
    }
}

func Login(userRepo *db.UserRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        type LoginRequest struct {
            Username string `json:"username" binding:"required"`
            Password string `json:"password" binding:"required"`
        }
        var req LoginRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        user, err := userRepo.FindByUsername(req.Username)
        if err != nil || user.ID == 0 {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
            return
        }

        // Compare hashed password
        if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
            return
        }

        expirationTime := time.Now().Add(24 * time.Hour)
        claims := &Claims{
            Username: req.Username,
            StandardClaims: jwt.StandardClaims{
                ExpiresAt: expirationTime.Unix(),
            },
        }
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        tokenString, err := token.SignedString(jwtKey)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"token": tokenString})
    }
}


func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
            c.Abort()
            return
        }

        token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        // Optionally, pass the username to the request context.
        // if claims, ok := token.Claims.(*Claims); ok {
        //     c.Set("username", claims.Username)
        // }
        c.Next()
    }
}