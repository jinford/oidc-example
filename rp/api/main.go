package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
)

func main() {
	clientID := os.Getenv("CLIENT_ID")
	if clientID == "" {
		log.Fatal("CLIENT_ID is required")
	}

	clientSecret := os.Getenv("CLIENT_SECRET")
	if clientSecret == "" {
		log.Fatal("CLIENT_SECRET is required")
	}

	callbackURL := os.Getenv("CALLBACK_URL")
	if callbackURL == "" {
		log.Fatal("CALLBACK_URL is required")
	}

	oidcProvider, err := oidc.NewProvider(context.Background(), "http://127.0.0.1:4444") // Hydra Public URL
	if err != nil {
		log.Fatal(err)
	}
	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     oidcProvider.Endpoint(),
		RedirectURL:  callbackURL,
		Scopes:       []string{oidc.ScopeOpenID},
	}

	h := NewHandler(oidcProvider, oauth2Config)
	r := chi.NewRouter()

	r.Get("/api/oidc-url", h.GetOIDCURL)
	r.Get("/api/user-info", h.GetUserInfo)

	log.Println("Starting RP API server on :13001")
	log.Fatal(http.ListenAndServe(":13001", r))
}

type handler struct {
	oidcProvider *oidc.Provider
	oauth2Config oauth2.Config
}

func NewHandler(oidcProvider *oidc.Provider, oauth2Config oauth2.Config) *handler {
	return &handler{
		oidcProvider: oidcProvider,
		oauth2Config: oauth2Config,
	}
}

func (h *handler) GetOIDCURL(w http.ResponseWriter, r *http.Request) {
	log.Println("GET OIDC URL")

	state, err := randString(16)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "state",
		Value:    state,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	})

	url := h.oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	log.Println("Redirect to", url)

	res := struct {
		RedirectURL string `json:"redirect_url"`
	}{
		RedirectURL: url,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println("failed to encode response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("GET UserInfo")

	code := r.URL.Query().Get("code")

	oauth2Token, err := h.oauth2Config.Exchange(r.Context(), code)
	if err != nil {
		log.Println("failed to exchange token", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userInfo, err := h.oidcProvider.UserInfo(r.Context(), oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		http.Error(w, "Failed to get userinfo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Subject:", userInfo.Subject)

	res := struct {
		Subject string `json:"subject"`
	}{
		Subject: userInfo.Subject,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println("failed to encode response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
