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

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    user := new(models.User)

    if err := c.BodyParser(user); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    if err := utils.ValidateUser(user); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    // Check if email already exists
    var emailExists bool
    err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tbl_users WHERE email = $1 AND deleted_at IS NULL)", user.Email).Scan(&emailExists)
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

    hashedPassword, err := utils.HashPassword(user.Password)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Error processing password",
        })
    }

    query := `
    INSERT INTO tbl_users (
        last_name, first_name, user_name, login_id, email, password,
        role_name, role_id, is_admin, login_session, last_login,
        currency_id, language_id, status_id, "order",
        created_by, updated_by
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
    RETURNING *`

    err = h.db.QueryRow(
        query,
        user.LastName, user.FirstName, user.UserName, user.LoginID, user.Email, hashedPassword,
        user.RoleName, user.RoleID, user.IsAdmin, user.LoginSession, user.LastLogin,
        user.CurrencyID, user.LanguageID, user.StatusID, user.Order,
        user.CreatedBy, user.UpdatedBy,
    ).Scan(
        &user.ID, &user.LastName, &user.FirstName, &user.UserName, &user.LoginID, &user.Email,
        &user.Password, &user.RoleName, &user.RoleID, &user.IsAdmin, &user.LoginSession,
        &user.LastLogin, &user.CurrencyID, &user.LanguageID, &user.StatusID, &user.Order,
        &user.CreatedBy, &user.CreatedAt, &user.UpdatedBy, &user.UpdatedAt, &user.DeletedBy,
        &user.DeletedAt,
    )

    if err != nil {
        fmt.Println("Error creating user:", err)
        return c.Status(500).JSON(fiber.Map{
            "error": "Error creating user",
        })
    }

    user.Password = "" // Remove password from response
    return c.Status(201).JSON(user)
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
    query := `
    SELECT *
    FROM tbl_users
    WHERE deleted_at IS NULL
    ORDER BY created_at DESC`

    rows, err := h.db.Query(query)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Error fetching users",
        })
    }
    defer rows.Close()

    users := make([]models.User, 0)
    for rows.Next() {
        var user models.User
        err := rows.Scan(
            &user.ID, &user.LastName, &user.FirstName, &user.UserName, &user.LoginID,
            &user.Email, &user.Password, &user.RoleName, &user.RoleID, &user.IsAdmin,
            &user.LoginSession, &user.LastLogin, &user.CurrencyID, &user.LanguageID,
            &user.StatusID, &user.Order, &user.CreatedBy, &user.CreatedAt, &user.UpdatedBy,
            &user.UpdatedAt, &user.DeletedBy, &user.DeletedAt,
        )
        if err != nil {
            return c.Status(500).JSON(fiber.Map{
                "error": "Error scanning user data",
            })
        }
        user.Password = "" // Remove password from response
        users = append(users, user)
    }

    return c.JSON(users)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
    id := c.Params("id")

    var user models.User
    err := h.db.QueryRow("SELECT * FROM tbl_users WHERE id = $1 AND deleted_at IS NULL", id).Scan(
        &user.ID, &user.LastName, &user.FirstName, &user.UserName, &user.LoginID,
        &user.Email, &user.Password, &user.RoleName, &user.RoleID, &user.IsAdmin,
        &user.LoginSession, &user.LastLogin, &user.CurrencyID, &user.LanguageID,
        &user.StatusID, &user.Order, &user.CreatedBy, &user.CreatedAt, &user.UpdatedBy,
        &user.UpdatedAt, &user.DeletedBy, &user.DeletedAt,
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

    user.Password = "" // Remove password from response
    return c.JSON(user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
    id := c.Params("id")
    user := new(models.User)

    if err := c.BodyParser(user); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    // Check if user exists
    var exists bool
    err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tbl_users WHERE id = $1 AND deleted_at IS NULL)", id).Scan(&exists)
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


    if user.Email != "" {
        var emailExists bool
        err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tbl_users WHERE email = $1 AND id != $2 AND deleted_at IS NULL)",
            user.Email, id).Scan(&emailExists)
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
    }

    query := `
    UPDATE tbl_users
    SET last_name = $1, first_name = $2, user_name = $3, login_id = $4, email = $5,
        role_name = $6, role_id = $7, is_admin = $8, login_session = $9, last_login = $10,
        currency_id = $11, language_id = $12, status_id = $13, "order" = $14,
        updated_by = $15, updated_at = CURRENT_TIMESTAMP
    WHERE id = $16 AND deleted_at IS NULL
    RETURNING *`

    err = h.db.QueryRow(
        query,
        user.LastName, user.FirstName, user.UserName, user.LoginID, user.Email,
        user.RoleName, user.RoleID, user.IsAdmin, user.LoginSession, user.LastLogin,
        user.CurrencyID, user.LanguageID, user.StatusID, user.Order,
        user.UpdatedBy, id,
    ).Scan(
        &user.ID, &user.LastName, &user.FirstName, &user.UserName, &user.LoginID,
        &user.Email, &user.Password, &user.RoleName, &user.RoleID, &user.IsAdmin,
        &user.LoginSession, &user.LastLogin, &user.CurrencyID, &user.LanguageID,
        &user.StatusID, &user.Order, &user.CreatedBy, &user.CreatedAt, &user.UpdatedBy,
        &user.UpdatedAt, &user.DeletedBy, &user.DeletedAt,
    )

    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Error updating user",
        })
    }

    user.Password = "" // Remove password from response
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
