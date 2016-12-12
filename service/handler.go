package service

import (
	"encoding/json"
	"fmt"
	"github.com/msawangwan/unitywebservice/model"
	//	"log"
	"net/http"
)

func availability(w http.ResponseWriter, r *http.Request) {
	var ps model.ProfileSearch

	if r.Body == nil {
		http.Error(w, "nil req body", 400)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&ps)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	fmt.Printf("need to verify %+v and is available is %+v\n", ps.Name, ps.IsAvailable)
}

//func main() {
//	http.HandleFunc("/ProfileSearch", ValidateProfileIsAvailable)
//	fmt.Printf("waiting\n")
//	log.Fatal(http.ListenAndServe(":8000", nil))
//}
