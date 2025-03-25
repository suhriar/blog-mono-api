package rest

import (
	"encoding/json"
	"net/http"

	"github.com/suhriar/blog-mono-api/internal/usecase"
	"github.com/suhriar/blog-mono-api/model"
	"github.com/suhriar/blog-mono-api/pkg/utils"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var request model.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	err := h.userUsecase.SignUp(r.Context(), request)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Sign up success"})
}

func (h *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var request model.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	user, err := utils.GetUserFromContext(r.Context())
	if err != nil {
		utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	res, err := h.userUsecase.ValidateRefreshToken(r.Context(), user.ID, request)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, res)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request model.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	jwtToken, refreshToken, err := h.userUsecase.Login(r.Context(), request)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK,
		model.LoginResponse{
			AccessToken:  jwtToken,
			RefreshToken: refreshToken,
		},
	)
}
