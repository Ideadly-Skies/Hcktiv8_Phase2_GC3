package internal

import (
	"net/http"
	"github.com/labstack/echo/v4"
	config "w3/gc3/config/database"
	"w3/gc3/utils"
	"errors"
	"fmt"
)

// Activity struct represents the user activity logs.
type Activity struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

// @Summary Get user activities
// @Description Retrieve the activity logs of the logged-in user
// @Tags Activities
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} Activity "List of user activities"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /activities [get]
func LogActivity(c echo.Context, userID int, description string) error {
	// Ensure database pool is initialized
	if config.Pool == nil {
		return errors.New("database connection pool is not initialized")
	}

	// Use Echo's request context
	query := `INSERT INTO user_activity_logs (user_id, description) VALUES ($1, $2)`
	_, err := config.Pool.Exec(c.Request().Context(), query, userID, description)
	if err != nil {
		return fmt.Errorf("failed to log activity: %w", err)
	}

	return nil
}

// GET /activities - Retrieve activities for the logged-in user.
func GetActivities(c echo.Context) error {
	// Get the user ID from the token
	userID, _ := utils.GetUserIDFromToken(c)
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "not authorized"})
	}

	// Query to fetch user activities
	query := `SELECT id, user_id, description FROM user_activity_logs WHERE user_id = $1`
	rows, err := config.Pool.Query(c.Request().Context(), query, userID)
	fmt.Println(err)
	
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to fetch activities"})
	}
	defer rows.Close()

	// Parse the results into a slice of Activity
	activities := []Activity{}
	for rows.Next() {
		var activity Activity
		if err := rows.Scan(&activity.ID, &activity.UserID, &activity.Description); err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to parse activities"})
		}
		activities = append(activities, activity)
	}

	// Return the activities
	return c.JSON(http.StatusOK, activities)
}