package user

import (
	"database/sql"
	"net/http"
	"strconv"

	"brain/db"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetUsers(c echo.Context) error {
	rows, err := db.DB.Query("SELECT id, name FROM users")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		users = append(users, user)
	}

	return c.JSON(http.StatusOK, users)
}

func GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "无效的 ID"})
	}

	var user User
	err = db.DB.QueryRow("SELECT id, name FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "未找到用户"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

func CreateUser(c echo.Context) error {
	var user User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	// 验证密码是否为空
	if user.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "password cannot be empty"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
	}

	result, err := db.DB.Exec("INSERT INTO users (name, password) VALUES (?, ?)", user.Name, hashedPassword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	id, err := result.LastInsertId()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	user.ID = int(id)
	// 使用 UserResponse 结构体返回用户信息

	userResponse := UserResponse{
		ID:   user.ID,
		Name: user.Name,
	}

	return c.JSON(http.StatusCreated, userResponse)
}

func UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var user User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	_, err = db.DB.Exec("UPDATE users SET name = ? WHERE id = ?", user.Name, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	user.ID = id
	return c.JSON(http.StatusOK, user)
}

func DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	_, err = db.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
