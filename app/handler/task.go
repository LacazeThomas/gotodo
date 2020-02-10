package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"net/http"

	"github.com/lacazethomas/goTodo/app/model"
)

// GetAllTasks from user
func GetAllTasks(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	projectID := vars["uuid"]
	status := vars["status"]
	project := getProjectOr404(db, projectID, w, r)
	if project == nil {
		return
	}
	var tasks []*model.Task
	db.Where("project_id = ? AND done = ?",project.ID, status).Find(&tasks)
	for _, task := range tasks {
		task.DecryptTask()
	}
	respondJSON(w, http.StatusOK, tasks)
}

// CreateTask for an user
func CreateTask(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	projectID := vars["uuid"]
	project := getProjectOr404(db, projectID, w, r)
	if project == nil {
		return
	}

	task := model.Task{ProjectID: project.ID}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	taskUuid, err := uuid.NewV4()
	if err != nil {
		respondError(w, http.StatusBadRequest, "Failed to create account, unable to generate UUID.")
		return
	}
	task.TaskID = taskUuid
	backTittle := task.Title
	task.EncryptTask()

	if err := db.Save(&task).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	task.Title = backTittle
	respondJSON(w, http.StatusCreated, task)
}

// GetTask according userID
func GetTask(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	projectID := vars["uuid"]
	project := getProjectOr404(db, projectID, w, r)
	if project == nil {
		return
	}

	id := vars["uuidTask"]
	task := getTaskOr404(db, id, w, r)
	if task == nil {
		return
	}
	task.DecryptTask()
	respondJSON(w, http.StatusOK, task)
}

func UpdateTask(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	projectID := vars["uuid"]
	project := getProjectOr404(db, projectID, w, r)
	if project == nil {
		return
	}

	id := vars["uuidTask"]
	task := getTaskOr404(db, id, w, r)
	task.DecryptTask()
	if task == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	task.EncryptTask()
	if err := db.Save(&task).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	task.DecryptTask()
	respondJSON(w, http.StatusOK, task)
}

// DeleteTask from param
func DeleteTask(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	projectID := vars["uuid"]
	project := getProjectOr404(db, projectID, w, r)
	if project == nil {
		return
	}

	id := vars["uuidTask"]
	task := getTaskOr404(db, id, w, r)
	if task == nil {
		return
	}

	if err := db.Delete(&task).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	task.DecryptTask()
	respondJSON(w, http.StatusNoContent, nil)
}

// CompleteTask from param
func CompleteTask(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	projectID := vars["uuid"]
	project := getProjectOr404(db, projectID, w, r)
	if project == nil {
		return
	}

	id := vars["uuidTask"]
	task := getTaskOr404(db, id, w, r)
	if task == nil {
		return
	}

	task.Complete()
	if err := db.Save(&task).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	task.DecryptTask()
	respondJSON(w, http.StatusOK, task)
}

// UndoTask uncheck task
func UndoTask(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	projectID := vars["uuid"]
	project := getProjectOr404(db, projectID, w, r)
	if project == nil {
		return
	}

	id := vars["uuidTask"]
	task := getTaskOr404(db, id, w, r)
	if task == nil {
		return
	}

	task.Undo()
	if err := db.Save(&task).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	task.DecryptTask()
	respondJSON(w, http.StatusOK, task)
}

// getTaskOr404 gets a task instance if exists, or respond the 404 error otherwise
func getTaskOr404(db *gorm.DB, id string, w http.ResponseWriter, r *http.Request) *model.Task {
	task := model.Task{}

	uniq, err := uuid.FromString(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return nil
	}

	task.TaskID = uniq

	if err := db.Where("task_id = ?", task.TaskID).First(&task, task).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &task
}
