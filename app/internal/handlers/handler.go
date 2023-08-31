package handlers

import (
	//_ "AvitoTesting/cmd/main/docs"
	"AvitoTesting/pkg/client/models"
	"AvitoTesting/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strconv"
)

type handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) handler {
	return handler{db}
}

type (
	AddUserPost struct {
		Delete []string `json:"delete"`
		Add    []string `json:"add"`
	}

	GetSegmentName struct {
		Name string `json:"name"`
	}

	Response struct {
		Message interface{} `json:"message"`
		Error   string      `json:"error"`
	}
)

// AddSegment
// @Summary Добавление нового сегмента
// @Description Создает новый сегмент с данными из запроса
// @Tags segments
// @Accept json
// @Produce json
// @Param segment body models.Segment true "JSON-info - segment name"
// @Success 201 {object} Response "Create segment"
// @Failure 500  {object}  Response "" httputil.HTTPError
// @Router /segment [post]
func (h handler) AddSegment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)

	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	var segment models.Segment
	err = json.Unmarshal(body, &segment)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	var find models.Segment
	if err = h.DB.First(&find, "name = ?", segment.Name).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	if find.ID != 0 {
		utils.RespondWithJSON(w, http.StatusOK, Response{Message: "segment already exists"})
		return
	}

	if result := h.DB.Create(&segment); result.Error != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.RespondWithJSON(w, http.StatusNotFound, Response{Error: err.Error()})
		}
		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, Response{Message: "Create segment"})
}

// DeleteSegment
// @Summary Удаление сегмента
// @Description Удаляет сегмент по имени
// @Tags segments
// @Accept json
// @Produce json
// @Param segment body models.Segment true "JSON-info - segment name"
// @Success 200 {object} Response "Delete segment"
// @Failure 404  {object}  Response "" httputil.HTTPError
// @Failure 500  {object}  Response "" httputil.HTTPError
// @Router /segment [delete]
func (h handler) DeleteSegment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	var in GetSegmentName

	err = json.Unmarshal(body, &in)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	var segment models.Segment
	if err := h.DB.First(&segment, "name = ?", in.Name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.RespondWithJSON(w, http.StatusNotFound, Response{Error: err.Error()})
			return
		}
		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	if err = h.DB.Delete(&segment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.RespondWithJSON(w, http.StatusNotFound, Response{Error: err.Error()})
			return
		}
		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, Response{Message: "Delete segment"})
}

// AddUser
// @Summary Добавление пользователя в сегмент
// @Description Если польщователя с таким id не существует, создает пользователя. Добавляет и удаляет его из списков сегментов.
// @Tags user
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param segments body AddUserPost true "JSON-info - add and delete lists"
// @Success 201 {object} Response "Add/Delete segments from user"
// @Failure 400  {object}  Response "" httputil.HTTPError
// @Failure 404  {object}  Response "" httputil.HTTPError
// @Failure 500  {object}  Response "" httputil.HTTPError
// @Router /user/{id} [post]
func (h handler) AddUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	var in AddUserPost
	err = json.Unmarshal(body, &in)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	var user models.User
	if err := h.DB.First(&user, userId).Error; err != nil {
		h.DB.Create(&user)
		h.DB.First(&user, userId)
		utils.RespondWithJSON(w, http.StatusCreated, Response{Message: "create new user"})
	}

	for _, segmentName := range in.Add {
		var segment models.Segment
		if err = h.DB.First(&segment, "name = ?", segmentName).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.RespondWithJSON(w, http.StatusNotFound, Response{Message: fmt.Sprintf("cant find segment (Name = %s)", segmentName), Error: err.Error()})
				return
			}
			utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
			continue
		}

		err := h.DB.Model(&user).Association("Segments").Append(&segment)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.RespondWithJSON(w, http.StatusNotFound, Response{Error: err.Error()})
				return
			}
			utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, Response{Message: "add user in new segment"})
	}

	for _, segmentName := range in.Delete {
		var segment models.Segment
		if err = h.DB.First(&segment, "name = ?", segmentName).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.RespondWithJSON(w, http.StatusNotFound, Response{Message: fmt.Sprintf("cant find segment (Name = %s)", segmentName), Error: err.Error()})
				continue
			}
			utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
			continue
		}
		err := h.DB.Model(&user).Association("Segments").Delete(segment)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.RespondWithJSON(w, http.StatusNotFound, Response{Error: err.Error()})
				return
			}
			utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, Response{Message: "delete user from segment"})

	}

	utils.RespondWithJSON(w, http.StatusCreated, Response{Message: "Add/Delete segments from user"})
}

// GetAllSegments
// @Summary Получение всех сегментов пользователя
// @Description Возвращает список сегментов, в которых состоит пользователь с заданным id.
// @Tags segments
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} Response "list of segments"
// @Failure 400  {object}  Response "" httputil.HTTPError
// @Failure 404  {object}  Response "" httputil.HTTPError
// @Failure 500  {object}  Response "" httputil.HTTPError
// @Router /segment/{id} [post]
func (h handler) GetAllSegments(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	var segments []models.Segment
	var user models.User

	if err := h.DB.First(&user, userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.RespondWithJSON(w, http.StatusNotFound, Response{Message: fmt.Sprintf("cant find user. id=%d", userId), Error: err.Error()})
			return
		}
		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	err = h.DB.Model(&user).Association("Segments").Find(&segments)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.RespondWithJSON(w, http.StatusNotFound, Response{Error: err.Error()})
			return
		}
		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, Response{Message: segments})
}
