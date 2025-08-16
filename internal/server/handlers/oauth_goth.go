package handlers

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"

// 	"go-boilerplate/internal/config"

// 	"github.com/markbates/goth"
// 	"github.com/markbates/goth/gothic"
// 	"github.com/markbates/goth/providers/facebook"
// 	"github.com/markbates/goth/providers/google"
// )

// func SetupGoth(cfg *config.Settings) {
// 	goth.UseProviders(
// 		google.New(cfg.GoogleOAuth.ClientID, cfg.GoogleOAuth.ClientSecret, cfg.GoogleOAuth.RedirectURL, cfg.GoogleOAuth.Scopes...),
// 		facebook.New(cfg.FacebookOAuth.ClientID, cfg.FacebookOAuth.ClientSecret, cfg.FacebookOAuth.RedirectURL, cfg.FacebookOAuth.Scopes...),
// 	)
// }

// func GothGoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
// 	ctx := context.WithValue(r.Context(), "provider", "google")
// 	r = r.WithContext(ctx)

// 	if _, err := gothic.CompleteUserAuth(w, r); err == nil {
// 		w.Write([]byte("already logged in"))
// 	} else {
// 		gothic.BeginAuthHandler(w, r)
// 	}
// }

// // GothGoogleCallbackHandler handles the callback from Google OAuth using the goth library.
// // Note: The goth library does not retrieve the phone number by default.
// // Extending the library's provider implementation would be required to add this functionality.
// func GothGoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
// 	ctx := context.WithValue(r.Context(), "provider", "google")
// 	r = r.WithContext(ctx)

// 	user, err := gothic.CompleteUserAuth(w, r)
// 	if err != nil {
// 		fmt.Fprintln(w, err)
// 		return
// 	}

// 	userInfo := UserInfo{
// 		FirstName: user.FirstName,
// 		LastName:  user.LastName,
// 		Email:     user.Email,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(userInfo)
// }

// func GothFacebookLoginHandler(w http.ResponseWriter, r *http.Request) {
// 	ctx := context.WithValue(r.Context(), "provider", "facebook")
// 	r = r.WithContext(ctx)

// 	if _, err := gothic.CompleteUserAuth(w, r); err == nil {
// 		w.Write([]byte("already logged in"))
// 	} else {
// 		gothic.BeginAuthHandler(w, r)
// 	}
// }

// // GothFacebookCallbackHandler handles the callback from Facebook OAuth using the goth library.
// // Note: The goth library does not retrieve the phone number by default.
// // Extending the library's provider implementation would be required to add this functionality.
// func GothFacebookCallbackHandler(w http.ResponseWriter, r *http.Request) {
// 	ctx := context.WithValue(r.Context(), "provider", "facebook")
// 	r = r.WithContext(ctx)

// 	user, err := gothic.CompleteUserAuth(w, r)
// 	if err != nil {
// 		fmt.Fprintln(w, err)
// 		return
// 	}

// 	userInfo := UserInfo{
// 		FirstName: user.FirstName,
// 		LastName:  user.LastName,
// 		Email:     user.Email,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(userInfo)
// }
