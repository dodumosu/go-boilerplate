package repositories

import "fmt"

var ErrAccountAlreadyExists = fmt.Errorf("account already exists")
var ErrAccountDoesNotExist = fmt.Errorf("account does not exist")
var ErrAccountIsBanned = fmt.Errorf("account is banned")
var ErrAccountIsNotActive = fmt.Errorf("account is not active")
var ErrAccountIsNotVerified = fmt.Errorf("account is not verified")
var ErrAccountIsSuspended = fmt.Errorf("account is suspended")
var ErrInvalidCredentials = fmt.Errorf("invalid credendials")
