package handlers

import (
	//_ "AvitoTesting/cmd/main/docs"
	"AvitoTesting/pkg/client/models"
	"AvitoTesting/pkg/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"io"
	"log"
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

func (h handler) AddSegment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var segment models.Segment
	err = json.Unmarshal(body, &segment)
	if err != nil {
		log.Fatalln(err)
	}

	var find models.Segment
	h.DB.First(&find, "name = ?", segment.Name)
	if find.ID != 0 {
		utils.RespondWithJSON(w, http.StatusOK, Response{Message: "segment already exists"})
		return
	}

	if result := h.DB.Create(&segment); result.Error != nil {
		log.Println(result.Error)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, Response{Message: "Create segment"})
}

func (h handler) DeleteSegment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var in GetSegmentName

	err = json.Unmarshal(body, &in)
	if err != nil {
		log.Fatalln(err)
	}

	var segment models.Segment
	if err := h.DB.First(&segment, "name = ?", in.Name).Error; err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, Response{Message: "segment not found", Error: err.Error()})
		return
	}

	h.DB.Delete(&segment)

	utils.RespondWithJSON(w, http.StatusOK, Response{Message: "Delete segment"})
}

func (h handler) AddUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, _ := strconv.Atoi(vars["id"])
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var in AddUserPost
	err = json.Unmarshal(body, &in)
	if err != nil {
		log.Fatalln(err)
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
			utils.RespondWithJSON(w, http.StatusNotFound, Response{Message: fmt.Sprintf("user haven`t this segment (Name = %s)", segmentName), Error: err.Error()})
			continue
		}
		h.DB.Model(&user).Association("Segments").Append(&segment)
	}

	for _, segmentName := range in.Delete {
		var segment models.Segment
		if err = h.DB.First(&segment, "name = ?", segmentName).Error; err != nil {
			utils.RespondWithJSON(w, http.StatusNotFound, Response{Message: fmt.Sprintf("user haven`t this segment (Name = %s)", segmentName), Error: err.Error()})
			continue
		}
		h.DB.Model(&user).Association("Segments").Delete(segment)
	}

	utils.RespondWithJSON(w, http.StatusCreated, Response{Message: "Add/Delete segments from user"})
}

func (h handler) GetAllSegments(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userId, _ := strconv.Atoi(vars["id"])

	var segments []models.Segment
	var user models.User

	if err := h.DB.First(&user, userId).Error; err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, Response{Message: fmt.Sprintf("cant find user. id=%d", userId), Error: err.Error()})
		return
	}

	err := h.DB.Model(&user).Association("Segments").Find(&segments)
	if err != nil {
		log.Println(err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, Response{Message: segments})
}
