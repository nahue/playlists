package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

var todos []Todo = make([]Todo, 0) // Initialize as an empty slice
var nextID = 1

func GetTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo Todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newTodo.ID = nextID
	nextID++
	todos = append(todos, newTodo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var updatedTodo Todo
	err = json.NewDecoder(r.Body).Decode(&updatedTodo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for i := range todos {
		if todos[i].ID == id {
			updatedTodo.ID = id // Ensure the ID matches
			todos[i] = updatedTodo
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedTodo)
			return
		}
	}

	http.Error(w, "Todo not found", http.StatusNotFound)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	for i := range todos {
		if todos[i].ID == id {
			todos = append(todos[:i], todos[i+1:]...) // Remove the element
			w.WriteHeader(http.StatusNoContent)       // Indicate successful deletion
			return
		}
	}

	http.Error(w, "Todo not found", http.StatusNotFound)
}
