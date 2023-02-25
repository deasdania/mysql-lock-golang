package main

import (
	"encoding/json"

	"fmt"

	"net/http"
)

func main() {

	http.HandleFunc("/savings", httpHandler)

	http.ListenAndServe(":8080", nil)

}

func httpHandler(w http.ResponseWriter, req *http.Request) {

	params := map[string]interface{}{}
	response := map[string]interface{}{}

	var err error
	err = json.NewDecoder(req.Body).Decode(&params)
	response, err = addEntry(params)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err != nil {
		response = map[string]interface{}{
			"error": err.Error(),
		}
	}

	if encodingErr := enc.Encode(response); encodingErr != nil {
		fmt.Println("{ error: " + encodingErr.Error() + "}")
	}

}
