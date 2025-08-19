package cli

import (
	"context"
	"go-boilerplate/internal/db"
	"go-boilerplate/internal/db/repositories"
	"go-boilerplate/internal/lib"

	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/choose"
	"github.com/cqroot/prompt/input"
	"github.com/spf13/cobra"
)

var accountsCommand = &cobra.Command{
	Use:   "accounts",
	Short: "manage accounts",
}

var createAccountCommand = &cobra.Command{
	Use:   "create",
	Short: "create an account",
	Run: func(cmd *cobra.Command, args []string) {
		createAccount()
	},
}

var changePasswordCommand = &cobra.Command{
	Use:   "changepassword",
	Short: "Change a user's password",
	Run: func(cmd *cobra.Command, args []string) {
		changePassword()
	},
}

func createAccount() {
	dbService, err := db.NewDBService(context.Background(), &settings.Database, logger)
	if err != nil {
		panic(err)
	}
	defer dbService.Close()
	accountRepo := repositories.NewAccountRepository(dbService, logger)

	p := prompt.New()
	email, err := p.Ask("Email:").Input("")
	if err != nil {
		panic(err)
	}
	username, err := p.Ask("Username:").Input("")
	if err != nil {
		panic(err)
	}
	firstName, err := p.Ask("First name:").Input("")
	if err != nil {
		panic(err)
	}
	lastName, err := p.Ask("Last name:").Input("")
	if err != nil {
		panic(err)
	}
	password, err := p.Ask("Password:").Input("", input.WithEchoMode(input.EchoPassword))
	if err != nil {
		panic(err)
	}
	confirmPassword, err := p.Ask("Confirm password:").Input("", input.WithEchoMode(input.EchoPassword))
	if err != nil {
		panic(err)
	}
	if password != confirmPassword {
		panic("Passwords do not match")
	}
	isSuperuser, err := p.Ask("Is the user a superuser?").Choose(
		[]string{"Yes", "NO"},
		choose.WithDefaultIndex(1),
	)
	if err != nil {
		panic(err)
	}
	skipVerification, err := p.Ask("Skip verification?").Choose(
		[]string{"Yes", "NO"},
		choose.WithDefaultIndex(1),
	)
	if err != nil {
		panic(err)
	}
	changePasswordOnLogin, err := p.Ask("Change password on login?").Choose(
		[]string{"Yes", "NO"},
		choose.WithDefaultIndex(1),
	)
	if err != nil {
		panic(err)
	}

	isSuperuserBool := isSuperuser == "Yes"
	skipVerificationBool := skipVerification == "Yes"
	changePasswordOnLoginBool := changePasswordOnLogin == "Yes"
	hashedPassword, err := lib.HashPassword(password)
	if err != nil {
		panic(err)
	}

	params := repositories.AccountCreateParams{
		FirstName:             firstName,
		LastName:              lastName,
		Email:                 email,
		IsOAuth:               skipVerificationBool,
		ChangePasswordOnLogin: changePasswordOnLoginBool,
		Password:              hashedPassword,
		IsSuperuser:           isSuperuserBool,
		Username:              username,
		Activate:              true,
	}

	_, _, err = accountRepo.CreateAccount(context.Background(), params)
	if err != nil {
		panic(err)
	}
}

func changePassword() {
	dbService, err := db.NewDBService(context.Background(), &settings.Database, logger)
	if err != nil {
		panic(err)
	}
	defer dbService.Close()
	accountRepo := repositories.NewAccountRepository(dbService, logger)

	p := prompt.New()
	email, err := p.Ask("Email:").Input("")
	if err != nil {
		panic(err)
	}

	account, err := accountRepo.GetByEmail(context.Background(), email)
	if err != nil {
		panic(err)
	}

	password, err := p.Ask("New password:").Input("", input.WithEchoMode(input.EchoPassword))
	if err != nil {
		panic(err)
	}
	confirmPassword, err := p.Ask("Confirm new password:").Input("", input.WithEchoMode(input.EchoPassword))
	if err != nil {
		panic(err)
	}
	if password != confirmPassword {
		panic("Passwords do not match")
	}

	hashedPassword, err := lib.HashPassword(password)
	if err != nil {
		panic(err)
	}

	err = accountRepo.UpdatePassword(context.Background(), account.ID, hashedPassword)
	if err != nil {
		panic(err)
	}

	println("Password updated successfully")
}

func init() {
	accountsCommand.AddCommand(createAccountCommand)
	accountsCommand.AddCommand(changePasswordCommand)
	rootCommand.AddCommand(accountsCommand)
}
