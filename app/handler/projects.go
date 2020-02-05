package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"

	"github.com/lacazethomas/goTodo/app/model"
)

func GetAllProjects(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	status := vars["status"]

	projects := []model.Project{}
	idUser := r.Context().Value("user").(uuid.UUID)
	db.Where("user_id = ? AND archived = ?", idUser, status).Find(&projects)
	respondJSON(w, http.StatusOK, projects)

}

func CreateProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	project := &model.Project{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&project); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	project.UserID = r.Context().Value("user").(uuid.UUID)

	uuid, err := uuid.NewV4()
	if err != nil {
		respondError(w, http.StatusBadRequest,"Failed to create account, unable to generate UUID.")
		return
	}
	project.ID = uuid

	err = db.Create(project).Error
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, project)
}

func GetProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	project := getProjectOr404(db, id, w, r)
	if project == nil {
		return
	}
	respondJSON(w, http.StatusOK, project)
}

func UpdateProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	project := getProjectOr404(db, id, w, r)
	if project == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&project); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, project)
}

func DeleteProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	project := getProjectOr404(db, id, w, r)
	if project == nil {
		return
	}
	if err := db.Delete(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

func ArchiveProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	project := getProjectOr404(db, id, w, r)
	if project == nil {
		return
	}
	project.Archive()
	if err := db.Save(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, project)
}

func RestoreProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	project := getProjectOr404(db, id, w, r)
	if project == nil {
		return
	}
	project.Restore()
	if err := db.Save(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, project)
}

// getProjectOr404 gets a project instance if exists, or respond the 404 error otherwise
func getProjectOr404(db *gorm.DB, id string, w http.ResponseWriter, r *http.Request) *model.Project {
	project := model.Project{}

	uniq, err := uuid.FromString(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return nil
	}

	idUser := r.Context().Value("user").(uuid.UUID)
	project.ID = uniq
	if err := db.Where("user_id = ?", idUser).First(&project, project).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}

	return &project
}
