package internal

import (
	"fmt"
	"net/http"
	"strings"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"w3/gc3/config/database"
	"w3/gc3/utils"
	"strconv"

	activity_handler "w3/gc3/internal/activityHandler"
)

// Validator instance
var validate = validator.New()

// Post struct
type Post struct {
	ID       int    `json:"id"`
	Content  string `json:"content"`
	ImageURL string `json:"image_url" validate:"required,url"`
	UserID   int    `json:"user_id"`
}

// @Summary Create a new post
// @Description Add a new post with optional image URL and content
// @Tags Posts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body Post true "Post data"
// @Success 201 {object} map[string]interface{} "Post created successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts [post]
func CreatePost(c echo.Context) error {
	userID, _ := utils.GetUserIDFromToken(c)
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "not authorized"})
	}

	post := new(Post)
	if err := c.Bind(post); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request body"})
	}

	// Validate ImageURL format
	if err := validate.Struct(post); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid image_url format"})
	}

	// If content is missing, fetch random joke
	if strings.TrimSpace(post.Content) == "" {
		joke, err := utils.FetchRandomJoke()
		fmt.Printf("joke: %s\n", joke)
		fmt.Println(err)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to fetch random joke"})
		}
		post.Content = joke
	}

	post.UserID = userID

	
	// Insert post into the database
	query := `INSERT INTO posts (content, image_url, user_id) VALUES ($1, $2, $3) RETURNING id`
	err := config.Pool.QueryRow(c.Request().Context(), query, post.Content, post.ImageURL, post.UserID).Scan(&post.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to create post"})
	}
	
	// log post activity
	description := "User created a new POST with ID " + strconv.Itoa(post.ID)
	if err := activity_handler.LogActivity(c, userID, description); err != nil {
		log.Printf("Failed to log activity: %v", err)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "post created successfully",
		"post":    post,
	})
}

// @Summary Get all posts
// @Description Retrieve a list of all posts
// @Tags Posts
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} Post "List of posts"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts [get]
func GetAllPosts(c echo.Context) error {
	query := `SELECT id, content, image_url, user_id FROM posts`
	rows, err := config.Pool.Query(c.Request().Context(), query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to fetch posts"})
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Content, &post.ImageURL, &post.UserID); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to scan posts"})
		}
		posts = append(posts, post)
	}

	return c.JSON(http.StatusOK, posts)
}

// @Summary Get post details by ID
// @Description Retrieve details of a specific post by its ID
// @Tags Posts
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]interface{} "Post details"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Post not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts/{id} [get]
func GetPostByID(c echo.Context) error {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid post ID"})
	}

	var post Post
	query := `SELECT id, content, image_url, user_id FROM posts WHERE id = $1`
	err = config.Pool.QueryRow(c.Request().Context(), query, postID).Scan(&post.ID, &post.Content, &post.ImageURL, &post.UserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "post not found"})
	}

	// Fetch comments for the post
	queryComments := `SELECT c.id, c.content, c.author_id, u.full_name 
                      FROM comments c 
                      JOIN users u ON c.author_id = u.id 
                      WHERE c.post_id = $1`
	rows, err := config.Pool.Query(c.Request().Context(), queryComments, postID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to fetch comments"})
	}
	defer rows.Close()

	comments := []map[string]interface{}{}
	for rows.Next() {
		var commentID, authorID int
		var content, authorName string
		if err := rows.Scan(&commentID, &content, &authorID, &authorName); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to scan comments"})
		}
		comments = append(comments, map[string]interface{}{
			"id":      commentID,
			"content": content,
			"author": map[string]interface{}{
				"id":   authorID,
				"name": authorName,
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"post":     post,
		"comments": comments,
	})
}

// @Summary Delete a post by ID
// @Description Remove a specific post by its ID
// @Tags Posts
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]string "Post deleted successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Post not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts/{id} [delete]
func DeletePost(c echo.Context) error {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid post ID"})
	}

	userID, _ := utils.GetUserIDFromToken(c) // Assuming a utility function to extract UserID from JWT token
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "unauthenticated"})
	}

	var postOwnerID int
	query := `SELECT user_id FROM posts WHERE id = $1`
	err = config.Pool.QueryRow(c.Request().Context(), query, postID).Scan(&postOwnerID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "post not found"})
	}

	if postOwnerID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"message": "you are not authorized to delete this post"})
	}

	// Delete the post
	queryDelete := `DELETE FROM posts WHERE id = $1`
	_, err = config.Pool.Exec(c.Request().Context(), queryDelete, postID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to delete post"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "post deleted successfully"})
}