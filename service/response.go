// вспомогательные функции для работы с HTTP-запросами и ответами в формате JSON.
package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// устанавливаем заголовок Content-Type: application/json для ответа
func setJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

// формирование JSON-ответа с ошибкой
// возврат ошибку, если не удалось сериализовать JSON или записать ответ
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
		return fmt.Errorf("ошибка записи ответа об ошибке: %w", err)
	}
	return nil
}

// формирование JSON-ответ с успешным результатом
func Success(w http.ResponseWriter, out any, code int) error {
	setJSONHeader(w)
	w.WriteHeader(code)
	resp, err := json.Marshal(out)
	if err != nil {
		return Error(w, "Ошибка сериализации", http.StatusInternalServerError)
	}

	_, err = w.Write(resp)
	if err != nil {
		return fmt.Errorf("ошибка записи ответа: %w", err)
	}
	return nil
}
