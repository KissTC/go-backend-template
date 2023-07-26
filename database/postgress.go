package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/kisstc/go-backend-template/models"
)

// este archivo implementará el repository
type PostgressRepository struct {
	db *sql.DB
}

// url de conexion
func NewPostgressRepository(url string) (*PostgressRepository, error) {
	db, err := sql.Open("postgress", url)
	if err != nil {
		return nil, err
	}
	return &PostgressRepository{db}, nil
}

func (repo *PostgressRepository) InsertUser(ctx context.Context, user *models.User) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO users(email, password) VALUES ($1, $2)", user.Email, user.Password)
	return err
}

func (repo *PostgressRepository) GetUserById(ctx context.Context, id int64) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, email FROM users WHERE id = $1", id)
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var user = models.User{}
	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Email); err == nil {
			return &user, nil
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *PostgressRepository) Close() error {
	return repo.db.Close()
}