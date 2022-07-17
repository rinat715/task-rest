package users

import (
	e "go_rest/internal/errors"
	m "go_rest/internal/models"
	"go_rest/internal/repositories/users"
	hashgenerator "go_rest/internal/services/hash_generator"

	"github.com/google/wire"
)

type Interface interface {
	Create(user *m.User) error
	Get(id int) (m.User, error)
	Login(email string, pass string) (*m.User, error)
}

type UserService struct {
	repository    users.UserRepositoryInterface
	hashGenerator hashgenerator.HashGeneratorInterface
}

func (u *UserService) Create(user *m.User) error {
	hashedPass, err := u.hashGenerator.Hash(user.Pass)
	if err != nil {
		return err
	}
	user.Pass = hashedPass
	return u.repository.Create(user)
}

func (u *UserService) Get(id int) (m.User, error) {
	// TODO думаю изера можно доставать сначала из какого нибудь редиса
	return u.repository.Get(id)
}

func (u *UserService) Login(email string, pass string) (*m.User, error) {
	user, err := u.repository.GetbyEmail(email)
	if err != nil {
		return nil, err
	}
	valid := u.hashGenerator.Check(pass, user.Pass)
	if !valid {
		return &m.User{}, &e.InvalidPass{Email: email}
	}
	return &user, nil
}

func NewUserService(r users.UserRepositoryInterface, g hashgenerator.HashGeneratorInterface) *UserService {
	return &UserService{
		repository:    r,
		hashGenerator: g,
	}
}

var UserSet = wire.NewSet(
	hashgenerator.NewHashGenerator,
	wire.Bind(new(hashgenerator.HashGeneratorInterface), new(*hashgenerator.HashGenerator)),
	users.NewUserRepository,
	wire.Bind(new(users.UserRepositoryInterface), new(*users.UserRepository)),
	NewUserService,
	wire.Bind(new(Interface), new(*UserService)),
)
