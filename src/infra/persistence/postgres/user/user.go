package user

import (
	dto "e-depo/src/app/dto/user"

	"log"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	helper "e-depo/src/infra/helper"
	"errors"
)

type UserRepository interface {
	StoreUser(data *dto.CreateUserReqDTO) (*dto.RegisterRespDTO, error)
	Login(data *dto.LoginReqDTO) (*dto.RegisterRespDTO, error)
}

const (
	InserUser = `INSERT INTO public.users (name,username,phone_number,address,email,password,role) 
		values ($1, $2, $3, $4, $5, $6, $7) returning id`

	Login = `select u.id, u.username, u.password  
	from public.users u 
	where u.username = $1`

	CreateWallet = `INSERT INTO public.wallets (user_id) 
	values ($1) returning id as wallet_id`
)

var statement PreparedStatement

type PreparedStatement struct {
	login *sqlx.Stmt
}

type userRepo struct {
	Connection *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	repo := &userRepo{
		Connection: db,
	}
	InitPreparedStatement(repo)
	return repo
}

func (p *userRepo) Preparex(query string) *sqlx.Stmt {
	statement, err := p.Connection.Preparex(query)
	if err != nil {
		log.Fatalf("Failed to preparex query: %s. Error: %s", query, err.Error())
	}

	return statement
}

func InitPreparedStatement(m *userRepo) {
	statement = PreparedStatement{
		login: m.Preparex(Login),
	}
}

func (p *userRepo) StoreUser(data *dto.CreateUserReqDTO) (resp *dto.RegisterRespDTO, err error) {
	// Hash the password from the registration data
	pwd, err := hashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	// Begin a new transaction
	tx, err := p.Connection.Beginx()
	if err != nil {
		log.Println("Failed to begin transaction: ", err.Error())
		return nil, err
	}

	// Define the INSERT query with parameters
	insertQuery := `INSERT INTO users (name, username, phone_number, address, email, password, role)
                    VALUES ($1, $2, $3, $4, $5, $6, $7)`

	// Log query and parameters
	log.Printf("Executing query: %s with parameters: %v", insertQuery, []interface{}{
		data.Name, data.UserName, data.PhoneNumber, data.Address, data.Email, pwd, data.Role})

	// Execute the INSERT query
	_, err = tx.Exec(insertQuery, data.Name, data.UserName, data.PhoneNumber, data.Address, data.Email, pwd, data.Role)
	if err != nil {
		tx.Rollback()
		log.Println("Failed to execute query: ", err.Error())
		return nil, err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Println("Failed to commit transaction: ", err.Error())
		return nil, err
	}

	log.Println("User successfully inserted")

	// Return the response object if everything is successful
	return resp, nil
}

func (p *userRepo) Login(data *dto.LoginReqDTO) (*dto.RegisterRespDTO, error) {
	var resultData []*dto.UserModel
	var resp dto.RegisterRespDTO

	// Execute the login query
	if err := statement.login.Select(&resultData, data.UserName); err != nil {
		return nil, err
	}

	// Check if no rows were returned from the query
	if len(resultData) < 1 {
		return nil, errors.New("no rows returned from the query")
	}

	// Verify the password
	if err := verifyPassword(resultData[0].Password, data.Password); err != nil {
		return nil, err
	}

	// Generate token
	token, err := helper.GenerateToken(resultData[0])
	if err != nil {
		return nil, err
	}
	resp.Token = token

	// Return the response object if everything is successful
	return &resp, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func verifyPassword(hashedPassword, inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
}
