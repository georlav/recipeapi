package handler

import (
	"fmt"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte(fmt.Sprintf(`{"Version": "%d"}`, 1))); err != nil {
		fmt.Print(err)
	}
}
