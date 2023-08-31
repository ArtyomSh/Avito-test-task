package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//func RespondWithError(w http.ResponseWriter, code int, message interface{}) {
//	RespondWithJSON(w, code, message)
//}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	var response []byte
	var err error
	if response, err = json.Marshal(payload); err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
