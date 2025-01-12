package auth

import (
	"net/http"
	"time"

	"brain/db"
	"brain/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your-secret-key")

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func Login(c echo.Context) error {
	var loginReq struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	if err := c.Bind(&loginReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "无效请求"})
	}

	// Get user from database
	var user user.User
	err := db.DB.QueryRow("SELECT id, password FROM users WHERE name = ?", loginReq.Name).
		Scan(&user.ID, &user.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "无效的凭证"})
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "无效的凭证"})
	}

	// Create JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": tokenString})
}

// AuthMiddleware 是一个中间件函数，用于验证请求中的授权令牌
// 参数 next 是下一个要执行的处理函数
// 返回一个新的处理函数，该处理函数会进行令牌验证，然后调用 next
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        // 从请求头中获取授权令牌
        tokenString := c.Request().Header.Get("Authorization")
        // 如果没有提供令牌，返回错误响应
        if tokenString == "" {
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "缺少授权令牌"})
        }

        // 初始化 Claims 对象，用于解析令牌中的声明
        claims := &Claims{}
        // 解析令牌，并验证其签名
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            // 返回用于验证令牌签名的密钥
            return jwtSecret, nil
        })

        // 如果解析出错或令牌无效，返回错误响应
        if err != nil || !token.Valid {
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "无效令牌"})
        }

        // 将解析出的用户 ID 存储在上下文中，以便后续使用
        c.Set("userID", claims.UserID)
        // 调用下一个处理函数
        return next(c)
    }
}
