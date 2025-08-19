package repositories

import (
	"context"
	"errors"
	"go-boilerplate/internal/db"
	"go-boilerplate/internal/lib"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
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
	FirstName             string
	LastName              string
	Email                 string
	IsOAuth               bool
	Password              string
	Username              string
	ChangePasswordOnLogin bool
	Activate              bool
	IsSuperuser           bool
}

func (r *AccountRepository) CheckAccountAvailability(ctx context.Context, username, email string) (bool, error) {
	result, err := r.dbService.CheckUserExistsByUsernameOrEmail(ctx, db.CheckUserExistsByUsernameOrEmailParams{
		EmailVal:    email,
		UsernameVal: pgtype.Text{String: username, Valid: true},
	})

	if err != nil {
		return false, err
	}

	return !(result.EmailExists || result.UsernameExists), err
}

func (r *AccountRepository) CreateAccount(ctx context.Context, params AccountCreateParams) (*db.User, *db.Profile, error) {
	var user db.User
	var profile db.Profile

	accountAvailable, err := r.CheckAccountAvailability(ctx, params.Username, params.Email)
	if err != nil {
		return nil, nil, err
	}
	if !accountAvailable {
		return nil, nil, ErrAccountAlreadyExists
	}

	txErr := r.dbService.WithinTransaction(ctx, func(qtx *db.Queries) error {
		var err error
		userID := lib.GetNewID()

		createUserParams := db.CreateUserParams{
			ID:    userID,
			Email: params.Email,
			Username: pgtype.Text{
				String: params.Username,
				Valid:  params.Username != "",
			},
			IsActive:    params.Activate,
			IsSuperuser: params.IsSuperuser,
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

	if txErr != nil {
		return nil, nil, txErr
	}

	return &user, &profile, nil
}

func (r *AccountRepository) VerifyAccount(ctx context.Context, userID string) error {
	err := r.dbService.VerifyUserAccount(ctx, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrAccountDoesNotExist
	}

	return err
}

func (r *AccountRepository) GetByEmail(ctx context.Context, email string) (*db.LookupAccountForAuthRow, error) {
	user, err := r.dbService.LookupAccountForAuth(ctx, pgtype.Text{String: email, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAccountDoesNotExist
		}

		return nil, err
	}

	return &user, nil
}

func (r *AccountRepository) UpdatePassword(ctx context.Context, userID, hashedPassword string) error {
	err := r.dbService.UpdatePassword(ctx, db.UpdatePasswordParams{
		ID:           userID,
		PasswordHash: pgtype.Text{String: hashedPassword, Valid: true},
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrAccountDoesNotExist
		}
	}

	return err
}

func (r *AccountRepository) UpdateProfile(ctx context.Context, userID, firstName, lastName, otherNames, bio, phone string) (*db.Profile, error) {
	params := db.UpdateUserProfileParams{UserID: userID}
	fName := strings.TrimSpace(firstName)
	if fName != "" {
		params.FirstName = pgtype.Text{String: fName, Valid: true}
	}
	lName := strings.TrimSpace(lastName)
	if lName != "" {
		params.LastName = pgtype.Text{String: lName, Valid: true}
	}
	oNames := strings.TrimSpace(otherNames)
	if oNames != "" {
		params.OtherNames = pgtype.Text{String: oNames, Valid: true}
	}
	cBio := strings.TrimSpace(bio)
	if cBio != "" {
		params.OtherNames = pgtype.Text{String: cBio, Valid: true}
	}
	cPhone := strings.TrimSpace(phone)
	if cPhone != "" {
		params.Phone = pgtype.Text{String: cPhone, Valid: true}
	}

	profile, err := r.dbService.UpdateUserProfile(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAccountDoesNotExist
		}
		return nil, err
	}

	return &profile, nil
}

func (r *AccountRepository) GetUserByID(ctx context.Context, userID string) (*db.User, error) {
	user, err := r.dbService.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAccountDoesNotExist
		}

		return nil, err
	}

	return &user, nil
}

func (r *AccountRepository) RequestPasswordReset(ctx context.Context, email string) (*db.LookupAccountForAuthRow, error) {
	// need data from both the profile and user tables
	// to send the recovery email
	accountInfo, err := r.dbService.LookupAccountForAuth(ctx, pgtype.Text{String: email, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAccountDoesNotExist
		}

		return nil, err
	}

	err = r.dbService.SetPasswordResetFlag(ctx, email)
	if err != nil {
		return nil, err
	}

	return &accountInfo, nil
}

func (r *AccountRepository) ResetPassword(ctx context.Context, userID, hashedPassword string) error {
	err := r.dbService.ResetPassword(ctx, db.ResetPasswordParams{
		ID:           userID,
		PasswordHash: pgtype.Text{String: hashedPassword, Valid: true},
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrAccountDoesNotExist
		}

		return err
	}

	return nil
}
