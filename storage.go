package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(id int, first_name, last_name string) error
	GetAccountByID(int) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	// connStr := "user=postgres dbname=mydatabase password=gobank sslmode=disable"
	connStr := "user=postgres dbname=mydatabase password=gobank sslmode=disable port=5433"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
		id serial PRIMARY KEY,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		balance serial,
		created_at timestamp
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	sqlStatement := `
	INSERT INTO account 
	(first_name, last_name, number, balance, created_at) 
	VALUES ($1, $2, $3, $4, $5)`

	res, err := s.db.Query(sqlStatement, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt)
	if err != nil {
		return err
	}

	fmt.Printf("%+v", res)
	return nil
}
func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	sqlStatement := `
	SELECT * FROM account`
	res, err := s.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	var accounts []*Account

	for res.Next() {
		// var account *Account
		// account := new(Account)
		account := &Account{}
		err := res.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)

	}
	return accounts, nil
}
func (s *PostgresStore) DeleteAccount(id int) error {
	sqlStatement := `
	DELETE FROM account
	WHERE id = $1`
	_, err := s.db.Query(sqlStatement, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("that user doen't exist")
		}
		return err
	}
	return nil
}
func (s *PostgresStore) UpdateAccount(id int, first_name, last_name string) error {
	sqlStatement := `
	UPDATE account
	SET first_name = $2, last_name = $3
	WHERE id = $1;`

	_, err := s.db.Query(sqlStatement, id, first_name, last_name)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("that user doen't exist")
		}
		return err
	}
	return nil
}
func (s *PostgresStore) GetAccountByID(int) (*Account, error) {
	return nil, nil
}
