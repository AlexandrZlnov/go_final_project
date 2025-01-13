package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func setJSONHeader(w http.ResponseWriter) {
    w.Header().Set("Content-Type", "application/json")
}

func Error(w http.ResponseWriter, message string, code int) error {
    setJSONHeader(w)
	w.WriteHeader(code)
	errorMsg := map[string]string{"error": message}
	resp, err := json.Marshal(errorMsg)
	if err != nil {
		http.Error(w, "Ошибка сериализации JSON", http.StatusInternalServerError)
		return fmt.Errorf("ошибка сериализации JSON: %w", err)
	}
    _, err = w.Write(resp)
    if err != nil {
        return fmt.Errorf("ошибка записи ответа: %w", err)
    }
    return nil

}

func Success(w http.ResponseWriter, out any, code int) error {
	setJSONHeader(w)
	resp, err := json.Marshal(out)
	if err != nil {
		return Error(w, "Ошибка сериализации", http.StatusInternalServerError)
	}
	w.WriteHeader(code)
	_, err = w.Write(resp)
    if err != nil {
        return fmt.Errorf("ошибка записи ответа: %w", err)
    }
	return nil
}



/*
package service

import (
	"encoding/json"
	"net/http"

	//"github.com/AlexandrZlnov/go_final_project/models"
)

func Error(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)

	errorMsg := map[string]string{"error": message}

	resp, err := json.Marshal(errorMsg)
	if err != nil {
		return
	}

	w.Write(resp)
}

func Success(w http.ResponseWriter, out any, code int) {
	resp, err := json.Marshal(out)
	if err != nil {
		Error(w, "Ошибка сериализации", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json, charset=UTF-8")
	w.WriteHeader(code)
	w.Write(resp)


}
	*/

