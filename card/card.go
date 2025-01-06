package card

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"yourproject/db"
)

type Card struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	UserID  int    `json:"user_id"`
}

func GetCards(c echo.Context) error {
	rows, err := db.DB.Query("SELECT id, content, user_id FROM cards")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	var cards []Card
	for rows.Next() {
		var card Card
		if err := rows.Scan(&card.ID, &card.Content, &card.UserID); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		cards = append(cards, card)
	}

	return c.JSON(http.StatusOK, cards)
}

func GetCard(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var card Card
	err = db.DB.QueryRow("SELECT id, content, user_id FROM cards WHERE id = ?", id).Scan(&card.ID, &card.Content, &card.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "card not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, card)
}

func CreateCard(c echo.Context) error {
	var card Card
	if err := c.Bind(&card); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	result, err := db.DB.Exec("INSERT INTO cards (content, user_id) VALUES (?, ?)", card.Content, card.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	id, err := result.LastInsertId()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	card.ID = int(id)
	return c.JSON(http.StatusCreated, card)
}

func UpdateCard(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var card Card
	if err := c.Bind(&card); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	_, err = db.DB.Exec("UPDATE cards SET content = ?, user_id = ? WHERE id = ?", card.Content, card.UserID, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	card.ID = id
	return c.JSON(http.StatusOK, card)
}

func DeleteCard(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	_, err = db.DB.Exec("DELETE FROM cards WHERE id = ?", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
