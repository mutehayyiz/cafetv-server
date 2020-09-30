package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mutehayyiz/cafetv-server/config"
	"github.com/mutehayyiz/cafetv-server/media"
	"github.com/mutehayyiz/cafetv-server/storage"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)


//Media

func AddMedia(w http.ResponseWriter, req *http.Request) {
	//parse form and get file
	err:= req.ParseForm()
	if err != nil{
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}
	file, handle, err := req.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "%v", err.Error())
		return
	}
	defer file.Close()

	// read file data
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	// hash file
	hasher := sha256.New()
	hasher.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	var l media.Media

	l.Name= handle.Filename
	l.Hash= hex.EncodeToString(hasher.Sum(nil))
	l.ID=primitive.NewObjectID()
	l.Category= string(req.FormValue("category"))
	l.MediaType= media.Type(req.FormValue("mediaType"))
	l.Description= req.FormValue("description")
	l.Online, _ = strconv.ParseBool(req.FormValue("online"))

	err= storage.MediaHandler.AddIfNotExists(&l)

	if err != nil {
		logrus.WithError(err).Error("Media couldn't be stored")
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	category:=req.FormValue("category")

	//create category directory if doesnt exist
	_, err = os.Stat("public/media/"+category)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll("public/media/"+category, 0755)
		if errDir != nil {
			log.Fatal(err)
		}

	}
	// write file
	err = ioutil.WriteFile("public/media/"+category+"/"+handle.Filename, data, 0666)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	// return response
	ReturnResponse(w, http.StatusOK, map[string]interface{}{
		"id":    l.ID.Hex(),
		"message": "Saved successfully",
	})
}

func GetAllMedias(w http.ResponseWriter, r *http.Request) {
	var medias []*media.Media
	err := storage.MediaHandler.GetAll(&medias)
	if err != nil {
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ReturnResponse(w, 200, medias)
}

func GetMediaByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var m media.Media
	err := storage.MediaHandler.GetByID(id, &m)
	if err != nil {
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ReturnResponse(w, 200, m)
}

func GetMediaByCategory(w http.ResponseWriter, r *http.Request) {
	category := mux.Vars(r)["category"]
	var medias []*media.Media
	err := storage.MediaHandler.GetByCategory(category, &medias)
	if err != nil {
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ReturnResponse(w, 200, medias)
}


func GetCategories(w http.ResponseWriter, r *http.Request){

	var categories []*string

	err := storage.MediaHandler.GetCategories(&categories)
	if err != nil {
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}


	ReturnResponse(w, 200, categories)

}


func Update(w http.ResponseWriter, r *http.Request) {

	ReturnResponse(w, 200,  map[string]interface{}{
		"message": "Media successfully updated",
	})
}

func DeleteMediaByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var m media.Media
	err := storage.MediaHandler.GetByID(id, &m)
	if err != nil {
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}
	path := "public/media/" + m.Category + "/" + m.Name
	fmt.Println(path)
	err = os.Remove(path)
	if err != nil {
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}


	err = storage.MediaHandler.DeleteByID(id)
	if err != nil {
		logrus.WithError(err).Error("Error while deleting license")
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ReturnResponse(w, 200, map[string]interface{}{
		"message": "Media successfully deleted",
	})
}

func DeleteAllMedias(w http.ResponseWriter, r *http.Request) {

	err := storage.MediaHandler.DeleteAll()
	if err != nil {
		logrus.WithError(err).Error("Error while deleting media")
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ReturnResponse(w, 200, map[string]interface{}{
		"message": "All media successfully deleted",
	})
}

//medias

// others

func Ping(w http.ResponseWriter, r *http.Request) {
	ReturnResponse(w, 200, map[string]interface{}{"naber":"ok"})

}

func ReturnResponse(w http.ResponseWriter, statusCode int, resp interface{}) {

	bytes, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = fmt.Fprintf(w, string(bytes))
}

func ReturnError(w http.ResponseWriter, statusCode int, errMsg string) {
	resp := map[string]interface{}{
		"error": errMsg,
	}
	bytes, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = fmt.Fprintf(w, string(bytes))
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != config.Global.AdminSecret {
			ReturnResponse(w, http.StatusUnauthorized, map[string]interface{}{
				"message": "Authorization failed",
			})
			return
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

