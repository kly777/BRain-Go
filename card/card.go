package card

import (
	"net/http"
	"strconv"
	"github.com/labstack/echo/v4"
)

// Card represents a card entity
type Card struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

var cards = []Card{
	{ID: 1, Title: "Card 1", Body: "This is card 1"},
	{ID: 2, Title: "Card 2", Body: "This is card 2"},
}

// GetCards returns all cards
func GetCards(c echo.Context) error {
	return c.JSON(http.StatusOK, cards)
}

// GetCard returns a single card by ID
func GetCard(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	for _, card := range cards {
		if card.ID == id {
			return c.JSON(http.StatusOK, card)
		}
	}
	return c.JSON(http.StatusNotFound, map[string]string{"error": "Card not found"})
}

// CreateCard creates a new card
func CreateCard(c echo.Context) error {
	card := new(Card)
	if err := c.Bind(card); err != nil {
		return err
	}
	card.ID = len(cards) + 1
	cards = append(cards, *card)
	return c.JSON(http.StatusCreated, card)
}

// UpdateCard updates an existing card
func UpdateCard(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	card := new(Card)
	if err := c.Bind(card); err != nil {
		return err
	}
	for i, card0 := range cards {
		if card0.ID == id {
			cards[i] = *card
			return c.JSON(http.StatusOK, card)
		}
	}
	return c.JSON(http.StatusNotFound, map[string]string{"error": "Card not found"})
}

// DeleteCard deletes a card by ID
func DeleteCard(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, card := range cards {
		if card.ID == id {
			cards = append(cards[:i], cards[i+1:]...)
			return c.NoContent(http.StatusNoContent)
		}
	}
	return c.JSON(http.StatusNotFound, map[string]string{"error": "Card not found"})
}