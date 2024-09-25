package user

import (
	"database/sql"
	"log"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (repo *Repository) CreateUser(username, email string, password []byte) error {
	_, err := repo.DB.Exec("INSERT INTO users (username, password, email) VALUES ($1, $2, $3)", username, password, email)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return err
	}
	return nil
}

func (repo *Repository) GetPasswordandID(username string) (int, string, error) {
	var hashedPassword string
	var id int
	err := repo.DB.QueryRow("SELECT id, password from users WHERE username = $1", username).Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		log.Printf("Username invalid: %v", err)
		return 0, "", err
	} else if err != nil {
		log.Printf("Error querying database: %v\n", err)
		return 0, "", err
	}

	return id, hashedPassword, nil
}

func (repo *Repository) GetUsers(users *Users) (*Users, error) {
	rows, err := repo.DB.Query("SELECT id, username, email FROM users")
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, err
	}
	for rows.Next() {
		user := User{}
		rows.Scan(&user.Id, &user.Username, &user.Email)
		users.Users[user.Id] = user
	}
	return users, nil
}
