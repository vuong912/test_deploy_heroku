package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
)

const (
	DbHost = "localhost:27017"
	DbName = "studentdb"
)

type (
	Student struct {
		IdStudent string `json:"idStudent" bson:"idStudent"`
		Barcode   string `json:"barcode" bson:"barcode"`
		FullName  string `json:"fullName" bson:"fullName"`
		LinkImage string `json:"linkImage" bson:"linkImage"`
	}
	StudentRepository struct {
		C *mgo.Collection
	}
)

func (r *StudentRepository) Get(query interface{}) ([]Student, error) {
	var students []Student
	err := r.C.Find(query).All(&students)
	return students, err
}

var session *mgo.Session

func GetSession() *mgo.Session {
	if session == nil {
		var err error
		session, err = mgo.DialWithInfo(&mgo.DialInfo{
			Addrs:    []string{DbHost},
			Username: "",
			Password: "",
			Timeout:  60 * time.Second,
		})
		if err != nil {
			log.Fatalf("[GetSession]: %s\n", err)
		}
	}
	return session
}

func GetIdStudentHandler(w http.ResponseWriter, r *http.Request) {
	var context *mgo.Session = GetSession().Copy()
	defer context.Close()
	c := context.DB(DbName).C("student")
	repo := &StudentRepository{c}
	id := mux.Vars(r)["idStudent"]
	students, err := repo.Get(bson.M{"idStudent": id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(students) > 0 {
		j, err := json.Marshal(students[0])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(j)
	}
}
func CheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Checked"))
}
func main() {
	GetSession()
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/student/get/{idStudent}", GetIdStudentHandler)
	router.HandleFunc("/student/check", CheckHandler)
	http.ListenAndServe(":80", router)
	log.Println("Server is running...")
}
