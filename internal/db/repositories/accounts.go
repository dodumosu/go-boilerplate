package repositories

import (
	"context"
	"errors"
	"go-boilerplate/internal/db"
	"go-boilerplate/internal/lib"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type AccountRepository struct {
	dbService *db.DBService
	logger    *slog.Logger
}

func NewAccountRepository(dbService *db.DBService, parentLogger *slog.Logger) *AccountRepository {
	return &AccountRepository{
		dbService: dbService,
		logger:    parentLogger.With("repository", "account"),
	}
}

type AccountCreateParams struct {
	FirstName string
	LastName  string
	Email     string
	IsOAuth   bool
	Password  string
	Username  string
}

func (r *AccountRepository) CreateAccount(ctx context.Context, params AccountCreateParams) (*db.User, *db.Profile, error) {
	var user db.User
	var profile db.Profile

	err := r.dbService.WithinTransaction(ctx, func(qtx *db.Queries) error {
		var err error
		userID := lib.GetNewID()

		createUserParams := db.CreateUserParams{
			ID:    userID,
			Email: params.Email,
			Username: pgtype.Text{
				String: params.Username,
				Valid:  params.Username != "",
			},
		}

		if !params.IsOAuth {
			hashedPassword, err := lib.HashPassword(params.Password)
			if err != nil {
				return err
			}
			createUserParams.PasswordHash = pgtype.Text{
				String: hashedPassword,
				Valid:  true,
			}
			createUserParams.RequiresVerification = true
		} else {
			createUserParams.RequiresVerification = false
		}

		createdUser, err := qtx.CreateUser(ctx, createUserParams)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return ErrAccountAlreadyExists
			}
			return err
		}
		user = createdUser

		createProfileParams := db.CreateProfileParams{
			ID:        lib.GetNewID(),
			UserID:    userID,
			FirstName: params.FirstName,
			LastName:  params.LastName,
		}

		createdProfile, err := qtx.CreateProfile(ctx, createProfileParams)
		if err != nil {
			return err
		}
		profile = createdProfile

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return &user, &profile, nil
}
