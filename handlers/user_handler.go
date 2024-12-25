package handlers

import (
	"database/sql"
	"fiber-crud/models"
	"fiber-crud/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	db *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{db: db}
}

// FUNCTION CREATE USER
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	// fmt.Println("hello")
	user := new(models.User)

	// Parse the request body into the user model
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate user fields
	if err := utils.ValidateUser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check if the email already exists
	var emailExists bool
	checkQuery := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`
	err := h.db.QueryRow(checkQuery, user.Email).Scan(&emailExists)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error checking email",
		})
	}
	if emailExists {
		return c.Status(400).JSON(fiber.Map{
			"error": "Email already exists",
		})
	}

	// Hash the user's password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error processing password",
		})
	}

	// Insert the user into the database
	// query := `
	// INSERT INTO users (first_name, last_name, email, password, phone, address, created_at, updated_at)
	// VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	// RETURNING id, first_name, last_name, email, phone, address, created_at, updated_at
	// `

	query := `
    INSERT INTO tbl_users (
        last_name, first_name, user_name, login_id, email, password,
        role_name, role_id, is_admin, login_session, last_login,
        currency_id, language_id, status_id, "order",
        created_by, updated_by, deleted_at
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
      *`

	err = h.db.QueryRow(
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		hashedPassword,
		user.Phone,
		user.Address,
	).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Address, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		fmt.Println("error", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Error creating user",
		})
	}

	// Return the created user
	return c.Status(201).JSON(user)
}

// FUNCTION GET ALL USER
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	query := `
	SELECT id, first_name, last_name, email, phone, address, created_at, updated_at
	FROM users
	ORDER BY created_at DESC
	`

	rows, err := h.db.Query(query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error fetching users",
		})
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Phone,
			&user.Address,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Error scanning user data",
			})
		}
		users = append(users, user)
	}

	return c.JSON(users)
}

// FUNCTION GET ALL USER BY ID
func (h *UserHandler) GetUser(c *fiber.Ctx) error {

	id := c.Params("id")

	query := `
		SELECT id, first_name, last_name, email, phone, address, created_at, updated_at, deleted_at
	FROM users
	WHERE id = $1
	`
	var user models.User
	err := h.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.Address,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error fetching user",
		})
	}

	return c.JSON(user)
}

// FUNCTION UPDATE USER
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if user exists
	var exists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error checking user existence",
		})
	}

	if !exists {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Parse update data
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update user data
	query := `
	UPDATE users
	SET first_name = $1,
		last_name = $2,
		email = $3,
		phone = $4,
		address = $5,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $6
	RETURNING id, first_name, last_name, email, phone, address, created_at, updated_at, 
	`

	err = h.db.QueryRow(
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		user.Address,
		id,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.Address,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Email already exists",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "Error updating user",
		})
	}

	return c.JSON(user)
}

// FUNCTION DELETE USER
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	// Corrected SQL query
	query := "DELETE FROM users WHERE id = $1"
	result, err := h.db.Exec(query, id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error deleting user",
		})
	}

	RowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error checking deleted rows",
		})
	}

	if RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}