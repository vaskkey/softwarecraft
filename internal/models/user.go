package models

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/vaskkey/softwarecraft/internal/config"
	"golang.org/x/crypto/bcrypt"
)

// password represents a password with a hash.
type password struct {
	plaintext *string
	hash      []byte
}

// Set hashes and sets the password for the user.
func (p *password) Set(plainText string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.plaintext = &plainText
	p.hash = hash

	return nil
}

// Compare compares the plaintext password against the hashed password.
func (p *password) Compare(plainText string) bool {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(plainText)) == nil
}

// User represents a user account in the database.
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserParams representation of raw values sent from the form
type UserParams struct {
	Name           string
	Email          string
	Password       string
	RepeatPassword string
}

func NewUserParams(form *url.Values) *UserParams {
	return &UserParams{
		Name:           form.Get("name"),
		Email:          form.Get("email"),
		Password:       form.Get("password"),
		RepeatPassword: form.Get("repeat_password"),
	}
}

// Validate validates data sent by the client
func (up *UserParams) Validate() (bool, config.ValidationErrors) {
	errors := make(config.ValidationErrors)

	if strings.TrimSpace(up.Email) == "" {
		errors["email"] = "Email is required."
	}

	if strings.TrimSpace(up.Name) == "" {
		errors["name"] = "Name is required."
	}

	if up.Password != up.RepeatPassword {
		errors["repeat_password"] = "Passwords must match."
	}

	if strings.TrimSpace(up.Password) == "" {
		errors["password"] = "Password is required."
	}

	if strings.TrimSpace(up.RepeatPassword) == "" {
		errors["repeat_password"] = "Repeat Password is required."
	}

	if utf8.RuneCountInString(up.Email) > 50 {
		errors["email"] = "Email is too long."
	}

	if utf8.RuneCountInString(up.Name) > 20 {
		errors["user"] = "Name is too long."
	}

	if utf8.RuneCountInString(up.Password) > 50 {
		errors["password"] = "Password is too long."
	}

	if utf8.RuneCountInString(up.RepeatPassword) > 50 {
		errors["repeat_password"] = "Repeat Password is too long."
	}

	return len(errors) == 0, errors
}

// GetUser convert params sent from client to a User record
func (up *UserParams) GetUser() (*User, error) {
	u := &User{
		Name:      up.Name,
		Email:     up.Email,
		Active:    true,
		Password:  password{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.Password.Set(up.Password); err != nil {
		return nil, err
	}

	return u, nil
}

// UserModel represents the database model for a user.
type UserModel struct {
	db *sql.DB
}

func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{db: db}
}

// Insert inserts a new user into the database.
func (m *UserModel) Insert(u *User) error {
	query := `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, u.Name, u.Email, u.Password.hash).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		switch {
		case err.Error() == `ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)`:
			return config.ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m *UserModel) GetByEmail(email string) (User, error) {
	query := `
		SELECT id, name, email, password_hash, active, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return scanUser(m.db.QueryRowContext(ctx, query, email))
}

func (m *UserModel) GetByToken(token string) (User, error) {
	tokenHash := sha256.Sum256([]byte(token))

	query := `
		SELECT u.id, u.name, u.email, u.password_hash, u.active, u.created_at, u.updated_at
		FROM users u
		INNER JOIN tokens t
		ON u.id = t.user_id
		WHERE t.hash = $1
		AND t.scope = $2
		AND t.expires_at > $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return scanUser(m.db.QueryRowContext(ctx, query, tokenHash[:], "authentication", time.Now()))
}

func scanUser(row *sql.Row) (User, error) {
	u := User{}
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password.hash, &u.Active, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		switch {
		case err.Error() == "no rows in result set":
			return u, config.ErrNoRecords
		default:
			return u, err
		}
	}

	return u, nil
}
