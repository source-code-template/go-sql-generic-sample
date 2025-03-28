package user

import (
	"database/sql"
	"net/http"

	"github.com/core-go/core"
	v "github.com/core-go/core/validator"
	"github.com/core-go/search/query"
	"github.com/core-go/sql/repository"

	"go-service/internal/user/handler"
	"go-service/internal/user/model"
	"go-service/internal/user/service"
)

type UserTransport interface {
	Search(w http.ResponseWriter, r *http.Request)
	All(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(db *sql.DB, logError core.Log) (UserTransport, error) {
	validator, err := v.NewValidator[*model.User]()
	if err != nil {
		return nil, err
	}

	buildQuery := query.UseQuery[model.User, *model.UserFilter](db, "users")
	userRepository, err := repository.NewSearchRepository[model.User, string, *model.UserFilter](db, "users", buildQuery)
	if err != nil {
		return nil, err
	}
	userService := service.NewUserService(db, userRepository)
	userHandler := handler.NewUserHandler(userService, logError, validator.Validate)
	return userHandler, nil
}
