package handlers

import (
	"fmt"
	"net/http"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fprint, err := fmt.Fprint(w, "Hello World\n")
	if err != nil {
		return
	}
	fmt.Println("Bytes returned:", fprint)
}
