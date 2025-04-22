package store

import (
	"context"
	"database/sql"
	"fiber/types"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type UserStore interface {
	InsertUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, int) (int, error)
	GetUsers(context.Context) ([]*types.User, error)
	GetUserByID(context.Context, int) (*types.User, error)
	GetUserByEmail(context.Context, string) (*types.User, error)
	UpdateUser(context.Context, int, map[string]any) (types.User, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	port, _ := strconv.Atoi(os.Getenv("PG_PORT"))
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", os.Getenv("PG_HOST"), port, os.Getenv("PG_USER"), os.Getenv("PG_PASS"), os.Getenv("PG_DB_NAME"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
		//log.Fatal("error to connect to Posgres database", err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (p *PostgresStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	rows, err := p.db.Query("select * from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*types.User{}
	for rows.Next() {
		user := new(types.User)
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.EncryptedPassword,
			&user.IsAdmin,
			&user.CreatedAt)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, sql.ErrNoRows
	}

	return users, nil

}

func (p *PostgresStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	rows, err := p.db.QueryContext(ctx, "select * from users where email=$1", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	user := &types.User{}
	if err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.EncryptedPassword,
		&user.IsAdmin,
		&user.CreatedAt); err != nil {
		return nil, err
	}
	return user, nil
}

func (p *PostgresStore) GetUserByID(ctx context.Context, id int) (*types.User, error) {
	rows, err := p.db.QueryContext(ctx, "select * from users where id=$1", id)
	if err != nil {
		fmt.Println("DB", err)
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	user := &types.User{}
	//rows.Next()
	if err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.EncryptedPassword,
		&user.IsAdmin,
		&user.CreatedAt); err != nil {

		return nil, err

	}
	//}

	return user, nil
}

func (p *PostgresStore) DeleteUser(ctx context.Context, id int) (int, error) {
	query := `DELETE FROM users WHERE id=$1 RETURNING id`
	var deletedID int
	err := p.db.QueryRowContext(ctx, query, id).Scan(&deletedID)

	if err != nil {
		fmt.Println("no rows found")
		return 0, sql.ErrNoRows
	}

	return deletedID, nil
}

func (p *PostgresStore) UpdateUser(ctx context.Context, id int, querySet map[string]any) (types.User, error) {

	setClauses := []string{}
	args := []any{}
	argPos := 1

	for k, v := range querySet {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", k, argPos))
		args = append(args, v)
		argPos++
	}

	args = append(args, id)

	query := fmt.Sprintf(`
	Update users 
	SET %s
	WHERE id=$%d 
	RETURNING id, first_name, last_name, email, admin, created_at
	`, strings.Join(setClauses, ", "), argPos)

	updUser := types.User{}
	err := p.db.QueryRowContext(ctx, query, args...).Scan(
		&updUser.ID,
		&updUser.FirstName,
		&updUser.LastName,
		&updUser.Email,
		&updUser.IsAdmin,
		&updUser.CreatedAt)

	if err != nil {
		fmt.Println("no rows found")
		return updUser, sql.ErrNoRows
	}

	return updUser, nil
}

func (p *PostgresStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	query := `insert into users 
		(first_name, last_name, email, pass, admin, created_at)
		values($1, $2, $3, $4, $5, $6)
		RETURNING id, first_name, last_name, email, pass, admin, created_at`

	insUser := &types.User{}
	err := p.db.QueryRowContext(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.EncryptedPassword,
		user.IsAdmin,
		user.CreatedAt,
	).Scan(
		&insUser.ID,
		&insUser.FirstName,
		&insUser.LastName,
		&insUser.Email,
		&insUser.EncryptedPassword,
		&insUser.IsAdmin,
		&insUser.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return insUser, nil
}

func (p *PostgresStore) createUserTable() error {
	query := `create table if not exists users (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		email varchar(50),
		pass varchar(250),
		admin boolean NOT NULL DEFAULT false,
		created_at timestamp
	)`

	_, err := p.db.Exec(query)
	return err
}

func (p *PostgresStore) Init() error {
	return p.createUserTable()
}
