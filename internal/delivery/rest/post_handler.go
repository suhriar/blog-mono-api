package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/suhriar/blog-mono-api/internal/usecase"
	"github.com/suhriar/blog-mono-api/model"
	"github.com/suhriar/blog-mono-api/pkg/utils"
)

type PostHandler struct {
	postUsecase usecase.PostUsecase
}

func NewPostHandler(postUsecase usecase.PostUsecase) *PostHandler {
	return &PostHandler{postUsecase: postUsecase}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var request model.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	user, err := utils.GetUserFromContext(r.Context())
	if err != nil {
		utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	err = h.postUsecase.CreatePost(r.Context(), user.ID, request)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Create new post success"})
}

func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return
	}

	res, err := h.postUsecase.GetPostByID(r.Context(), id)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, res)
}

func (h *PostHandler) GetAllPost(w http.ResponseWriter, r *http.Request) {
	pageIndexStr := r.URL.Query().Get("page-index")
	pageSizeStr := r.URL.Query().Get("page-size")

	pageIndex, err := strconv.Atoi(pageIndexStr)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return
	}

	res, err := h.postUsecase.GetAllPost(r.Context(), pageSize, pageIndex)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, res)
}

func (h *PostHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var request model.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return
	}

	user, err := utils.GetUserFromContext(r.Context())
	if err != nil {
		utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	err = h.postUsecase.CreateComment(r.Context(), id, user.ID, request)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Comment created"})
}

func (h *PostHandler) UpsertUserActivity(w http.ResponseWriter, r *http.Request) {
	var request model.UserActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return
	}

	user, err := utils.GetUserFromContext(r.Context())
	if err != nil {
		utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	err = h.postUsecase.UpsertUserActivity(r.Context(), id, user.ID, request)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Like success"})
}
