package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	config "w3/gc3/config/database"
	activity_handler "w3/gc3/internal/activityHandler"
	"w3/gc3/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Comment struct based on the provided schema
type Comment struct {
	ID      int    `json:"id"`
	Content string `json:"content" validate:"required"`
	PostID  int    `json:"post_id" validate:"required"`
	AuthorID int   `json:"author_id"`
}

// Validator instance
var validate = validator.New()

// @Summary Create a new comment
// @Description Add a new comment to a specific post
// @Tags Comments
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body Comment true "Comment data"
// @Success 201 {object} map[string]interface{} "Comment created successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /comments [post]
func CreateComment(c echo.Context) error {
	authorID, _ := utils.GetUserIDFromToken(c)
	if authorID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "not authorized"})
	}

	comment := new(Comment)
	if err := c.Bind(comment); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request body"})
	}

	// Validate request fields
	if err := validate.Struct(comment); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid input"})
	}

	comment.AuthorID = authorID
	
	// Log the user activity for creating a comment
	description := "User commented on POST with ID " + strconv.Itoa(comment.PostID)
	if logErr := activity_handler.LogActivity(c, authorID, description); logErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to log activity"})
	}

	// Insert comment into the database
	query := `INSERT INTO comments (content, post_id, author_id) VALUES ($1, $2, $3) RETURNING id`
	err := config.Pool.QueryRow(c.Request().Context(), query, comment.Content, comment.PostID, comment.AuthorID).Scan(&comment.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to create comment"})
	}


	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "comment created successfully",
		"comment": comment,
	})
}

// @Summary Get comment details by ID
// @Description Retrieve details of a specific comment by its ID
// @Tags Comments
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Comment ID"
// @Success 200 {object} map[string]interface{} "Comment details"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Comment not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /comments/{id} [get]
func GetCommentByID(c echo.Context) error {
	commentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid comment ID"})
	}

	var comment Comment
	var postTitle, authorName string
	query := `
		SELECT c.id, c.content, c.post_id, c.author_id, p.content AS post_title, u.full_name AS author_name
		FROM comments c
		JOIN posts p ON c.post_id = p.id
		JOIN users u ON c.author_id = u.id
		WHERE c.id = $1`
	err = config.Pool.QueryRow(c.Request().Context(), query, commentID).Scan(
		&comment.ID, &comment.Content, &comment.PostID, &comment.AuthorID, &postTitle, &authorName,
	)

	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "comment not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to fetch comment"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"comment": comment,
		"post":    map[string]interface{}{"title": postTitle},
		"author":  map[string]interface{}{"name": authorName},
	})
}

// @Summary Delete a comment by ID
// @Description Remove a specific comment by its ID
// @Tags Comments
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Comment ID"
// @Success 200 {object} map[string]string "Comment deleted successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Comment not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /comments/{id} [delete]
func DeleteCommentByID(c echo.Context) error {
	commentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid comment ID"})
	}

	authorID, _ := utils.GetUserIDFromToken(c)
	if authorID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "not authorized"})
	}

	// Check if the user is the owner of the comment
	var commentOwnerID int
	query := `SELECT author_id FROM comments WHERE id = $1`
	err = config.Pool.QueryRow(c.Request().Context(), query, commentID).Scan(&commentOwnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "comment not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to validate ownership"})
	}

	if commentOwnerID != authorID {
		return c.JSON(http.StatusForbidden, map[string]string{"message": "you are not authorized to delete this comment"})
	}

	// Delete the comment
	queryDelete := `DELETE FROM comments WHERE id = $1`
	_, err = config.Pool.Exec(c.Request().Context(), queryDelete, commentID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to delete comment"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "comment deleted successfully"})
}