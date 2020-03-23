package users

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (service *Service) Start() {
	_, err := service.pool.Exec(context.Background(), `
CREATE TABLE IF NOT EXISTS users (
	id BIGSERIAL,
    name TEXT,
    login TEXT UNIQUE,
    password TEXT,
	phone INTEGER NOT NULL,
    removed BOOLEAN DEFAULT FALSE
);
`)
	if err != nil {
		log.Print(err)
	}
	_, err = service.pool.Exec(
		context.Background(),
		`INSERT INTO users(id, name, login, password, phone) 
		VALUES (0, 'bank', 'admin', 'password', 0);`,
	)
	if err != nil {
		log.Print("admin has to db")
	}
}

func (service *Service) All() (models []User, err error) {
	rows, err := service.pool.Query(context.Background(), `
	SELECT id, name, login, phone FROM users WHERE removed = FALSE;
`)
	if err != nil {
		return nil, fmt.Errorf("can't get users from db: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.Id, &user.Name, &user.Login, &user.Phone)
		if err != nil {
			return nil, fmt.Errorf("can't get users from db: %w", err)
		}
		models = append(models, user)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get users from db: %w", err)
	}
	return models, nil
}

func (service *Service) Save(model User) (err error) {
	if model.Id == 0 {
		_, err = service.pool.Exec(
			context.Background(),
			`INSERT INTO users(name, login, password, phone) 
			VALUES ($1, $2, $3, $4);`,
			model.Name,
			model.Login,
			model.Password,
			model.Phone,
		)
		if err != nil {
			log.Printf("can't exec add user: %d", err)
			return err
		}
		return nil
	} else {
		_, err = service.pool.Exec(context.Background(),
			`UPDATE users 
			SET (name=$2, login=$3, password=$4, phone=$5) 
			WHERE removed = FALSE and id=$1`,
			model.Id,
			model.Name,
			model.Login,
			model.Password,
			model.Phone,
		)
		if err != nil {
			log.Print("can't exec update blocked ", err)
			return err
		}
		return nil
	}
}

func (service *Service) RemoveById(id int) (err error) {
	_, err = service.pool.Exec(context.Background(),
		`UPDATE users 
		SET removed=TRUE 
		WHERE removed = FALSE and id=$1`, id)
	if err != nil {
		log.Print("can't exec update blocked ", err)
		return err
	}
	return nil
}

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Phone    int    `json:"phone"`
}

type ResponseDTO struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Phone int    `json:"phone"`
}

func (s *Service) Profile(id int) (users []User, err error) {
	user := User{}
	if id == -1 {
		rows, err := s.pool.Query(context.Background(), `SELECT id, name, login, password, phone FROM users`)
		if err != nil {
			log.Printf("can't query users to profile: %d", err)
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(
				&user.Id,
				&user.Name,
				&user.Login,
				&user.Password,
				&user.Phone,
			)
			if err != nil {
				return nil, fmt.Errorf("can't get user from db: %w", err)
			}
			users = append(users, user)
		}
		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("can't get user from db: %w", err)
		}
	} else {
		_ = s.pool.QueryRow(context.Background(),
			`SELECT id, name, login, password, phone 
			FROM users 
			WHERE id=$1`,
			id).Scan(
			&user.Id,
			&user.Name,
			&user.Login,
			&user.Password,
			&user.Phone,
		)
		users = append(users, user)
	}
	log.Print("showing ")
	return users, nil
}

type UserTDO struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Phone    int    `json:"phone"`
}
