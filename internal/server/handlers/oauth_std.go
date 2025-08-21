package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go-boilerplate/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

type UserInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
}

// GetGoogleOAuthConfig configures and returns an oauth2.Config for Google.
// It adds the "https://www.googleapis.com/auth/user.phonenumbers.read" scope
// to retrieve the user's phone number. Note that the People API must be enabled
// in your Google Cloud project for this to work.
func GetGoogleOAuthConfig(cfg *config.Settings) *oauth2.Config {
	scopes := make([]string, len(cfg.GoogleOAuth.Scopes))
	copy(scopes, cfg.GoogleOAuth.Scopes)
	scopes = append(scopes, "https://www.googleapis.com/auth/user.phonenumbers.read")

	return &oauth2.Config{
		RedirectURL:  cfg.GoogleOAuth.RedirectURL,
		ClientID:     cfg.GoogleOAuth.ClientID,
		ClientSecret: cfg.GoogleOAuth.ClientSecret,
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}
}

func GetFacebookOAuthConfig(cfg *config.Settings) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  cfg.FacebookOAuth.RedirectURL,
		ClientID:     cfg.FacebookOAuth.ClientID,
		ClientSecret: cfg.FacebookOAuth.ClientSecret,
		Scopes:       cfg.FacebookOAuth.Scopes,
		Endpoint:     facebook.Endpoint,
	}
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := &http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Expires:  time.Now().Add(20 * time.Minute),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	return state
}

func GoogleLoginHandler(cfg *config.Settings) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		googleOAuthConfig := GetGoogleOAuthConfig(cfg)
		state := generateStateOauthCookie(w)
		url := googleOAuthConfig.AuthCodeURL(state)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func GoogleCallbackHandler(cfg *config.Settings) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		googleOAuthConfig := GetGoogleOAuthConfig(cfg)
		oauthState, _ := r.Cookie("oauthstate")

		if r.FormValue("state") != oauthState.Value {
			http.Error(w, "invalid oauth state", http.StatusBadRequest)
			return
		}

		token, err := googleOAuthConfig.Exchange(context.Background(), r.FormValue("code"))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to exchange token: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		// Use the People API to retrieve user information, including phone number.
		// This requires the People API to be enabled in your Google Cloud project.
		response, err := http.Get("https://people.googleapis.com/v1/people/me?personFields=names,emailAddresses,phoneNumbers&access_token=" + token.AccessToken)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get user info: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		defer response.Body.Close()
		contents, err := io.ReadAll(response.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to read response body: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		var googleUser struct {
			Names          []struct{ GivenName, FamilyName string } `json:"names"`
			EmailAddresses []struct{ Value string }                 `json:"emailAddresses"`
			PhoneNumbers   []struct{ Value string }                 `json:"phoneNumbers"`
		}
		if err := json.Unmarshal(contents, &googleUser); err != nil {
			http.Error(w, fmt.Sprintf("failed to unmarshal user info: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		userInfo := UserInfo{}
		if len(googleUser.Names) > 0 {
			userInfo.FirstName = googleUser.Names[0].GivenName
			userInfo.LastName = googleUser.Names[0].FamilyName
		}
		if len(googleUser.EmailAddresses) > 0 {
			userInfo.Email = googleUser.EmailAddresses[0].Value
		}
		if len(googleUser.PhoneNumbers) > 0 {
			userInfo.Phone = googleUser.PhoneNumbers[0].Value
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userInfo)
	}
}

func FacebookLoginHandler(cfg *config.Settings) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		facebookOAuthConfig := GetFacebookOAuthConfig(cfg)
		state := generateStateOauthCookie(w)
		url := facebookOAuthConfig.AuthCodeURL(state)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// FacebookCallbackHandler handles the callback from Facebook OAuth.
// Note: Retrieving a user's phone number from Facebook requires special permissions
// that are not granted by default and need to go through an app review process.
// Therefore, this implementation does not retrieve the phone number.
func FacebookCallbackHandler(cfg *config.Settings) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		facebookOAuthConfig := GetFacebookOAuthConfig(cfg)
		oauthState, _ := r.Cookie("oauthstate")

		if r.FormValue("state") != oauthState.Value {
			http.Error(w, "invalid oauth state", http.StatusBadRequest)
			return
		}

		token, err := facebookOAuthConfig.Exchange(context.Background(), r.FormValue("code"))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to exchange token: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		response, err := http.Get("https://graph.facebook.com/me?fields=id,first_name,last_name,email&access_token=" + token.AccessToken)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get user info: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		defer response.Body.Close()
		contents, err := io.ReadAll(response.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to read response body: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		var facebookUser struct {
			Email     string `json:"email"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}
		if err := json.Unmarshal(contents, &facebookUser); err != nil {
			http.Error(w, fmt.Sprintf("failed to unmarshal user info: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		userInfo := UserInfo{
			FirstName: facebookUser.FirstName,
			LastName:  facebookUser.LastName,
			Email:     facebookUser.Email,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userInfo)
	}
}
