package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	client "github.com/ory/hydra-client-go/v2"
)

func main() {
	h := NewHandler()
	r := chi.NewRouter()

	r.Post("/api/login", h.PostLogin)
	r.Get("/api/consent", h.GetConsent)
	r.Post("/api/consent", h.PostConsent)

	log.Println("Starting IdP API server on :14001")
	log.Fatal(http.ListenAndServe(":14001", r))
}

type handler struct {
	hydraAdminAPIClient *client.APIClient
}

func NewHandler() *handler {
	configuration := client.NewConfiguration()
	configuration.Servers = []client.ServerConfiguration{
		{
			URL: "http://localhost:4445", // Admin API URL
		},
	}

	return &handler{
		hydraAdminAPIClient: client.NewAPIClient(configuration),
	}
}

func (h *handler) PostLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("POST Login")

	var req struct {
		Username       string `json:"username"`
		Password       string `json:"password"`
		LoginChallenge string `json:"login_challenge"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("invalid request body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 簡易認証
	if req.Username != "admin" || req.Password != "password" {
		log.Println("invalid credentials")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Hydra の Admin API (POST /oauth2/auth/requests/login/accept) を呼び出し
	acceptOAuth2LoginRequest := client.NewAcceptOAuth2LoginRequestWithDefaults()
	acceptOAuth2LoginRequest.Subject = req.Username
	redirectTo, resp, err := h.hydraAdminAPIClient.OAuth2API.AcceptOAuth2LoginRequest(r.Context()).
		LoginChallenge(req.LoginChallenge).
		AcceptOAuth2LoginRequest(*acceptOAuth2LoginRequest).
		Execute()
	if err != nil {
		log.Println("failed to accept login request", err)

		switch resp.StatusCode {
		case http.StatusNotFound:
			w.WriteHeader(http.StatusNotFound)
		case http.StatusGone:
			w.WriteHeader(http.StatusGone)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("status code is not ok", resp.Status)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Hydra で生成されたリダイレクトURLを返却
	responseJSONBody, err := redirectTo.MarshalJSON()
	if err != nil {
		log.Println("failed to marshal response body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("POST Login success")
	log.Println("Response Body:", string(responseJSONBody))

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(responseJSONBody); err != nil {
		log.Println("failed to write response body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handler) GetConsent(w http.ResponseWriter, r *http.Request) {
	log.Println("GET Consent")

	consentChallenge := r.URL.Query().Get("consent_challenge")

	// Hydra の Admin API (GET /oauth2/auth/requests/consent) 呼び出し
	consentRequest, resp, err := h.hydraAdminAPIClient.OAuth2API.GetOAuth2ConsentRequest(r.Context()).
		ConsentChallenge(consentChallenge).
		Execute()
	if err != nil {
		log.Println("failed to get consent request", err)

		switch resp.StatusCode {
		case http.StatusNotFound:
			w.WriteHeader(http.StatusNotFound)
		case http.StatusGone:
			w.WriteHeader(http.StatusGone)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	responseJSONBody, err := consentRequest.MarshalJSON()
	if err != nil {
		log.Println("failed to marshal response body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("GET Consent success")
	log.Println("Response Body:", string(responseJSONBody))

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(responseJSONBody); err != nil {
		log.Println("failed to write response body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handler) PostConsent(w http.ResponseWriter, r *http.Request) {
	log.Println("POST Consent")

	var req struct {
		ConsentChallenge         string   `json:"consent_challenge"`
		GrantScope               []string `json:"grant_scope"`
		GrantAccessTokenAudience []string `json:"grant_access_token_audience"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("invalid request body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Hydra の Admin API (POST /oauth2/auth/requests/consent/accept) 呼び出し
	acceptOAuth2ConsentRequest := client.NewAcceptOAuth2ConsentRequestWithDefaults()
	acceptOAuth2ConsentRequest.GrantScope = req.GrantScope
	acceptOAuth2ConsentRequest.GrantAccessTokenAudience = req.GrantAccessTokenAudience
	redirectTo, resp, err := h.hydraAdminAPIClient.OAuth2API.AcceptOAuth2ConsentRequest(r.Context()).
		ConsentChallenge(req.ConsentChallenge).
		AcceptOAuth2ConsentRequest(*acceptOAuth2ConsentRequest).
		Execute()
	if err != nil {
		log.Println("failed to accept login request", err)

		switch resp.StatusCode {
		case http.StatusNotFound:
			w.WriteHeader(http.StatusNotFound)
		case http.StatusGone:
			w.WriteHeader(http.StatusGone)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	// Hydra で生成されたリダイレクトURLを返却
	responseJSONBody, err := redirectTo.MarshalJSON()
	if err != nil {
		log.Println("failed to marshal response body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("POST Consent success")
	log.Println("Response Body:", string(responseJSONBody))

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(responseJSONBody); err != nil {
		log.Println("failed to write response body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
