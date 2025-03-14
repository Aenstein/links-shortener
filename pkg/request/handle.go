package request

import (
	"net/http"
)

func HandleBody[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	err = IsValide[T](body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	return &body, nil
}
