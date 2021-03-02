package main

import (
	"io"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
        "github.com/jinzhu/gorm"
        _ "github.com/jinzhu/gorm/dialects/mysql"
	"encoding/json"
	"strconv"
	"github.com/rs/cors"
)

var db, _ = gorm.Open("mysql", "root:root@/skupine?charset=utf8&parseTime=True&loc=Local")
var err error

type Skupina struct{
	Id int `gorm:"primary_key"`
	Ime_skupine string
	Uporabniki []Uporabnik `gorm:"foreignKey:SkupinaId"`
}

type Uporabnik struct{
	Id int `gorm:"primary_key"`
	Ime_uporabnika string
	Password string
	Email string
	SkupinaId int 
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health is OK")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func Izpisi_skupine(w http.ResponseWriter, r *http.Request) {
	var skupina []Skupina
	db.Find(&skupina)
	fmt.Println("{}", skupina)
   	json.NewEncoder(w).Encode(skupina)
	
}

func Izpisi_uporabnike(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key,_ := strconv.Atoi(vars["skupinaid"])
	var uporabniki []Uporabnik
	db.Find(&uporabniki)
	fmt.Println("{}", uporabniki)	
	for _,uporabnik := range uporabniki{
		if uporabnik.SkupinaId == key{
			
			json.NewEncoder(w).Encode(uporabnik)
		}	
	}
		
}


func Ustvari_skupino(w http.ResponseWriter, r *http.Request) {
	ime_skupine := r.FormValue("ime_skupine")
	log.WithFields(log.Fields{"ime_skupine": ime_skupine}).Info("Dodaj novo skupino.")
	do_sk := &Skupina{Ime_skupine: ime_skupine }
	db.Create(&do_sk)
	result := db.Last(&do_sk)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.Value)
}

func Ustvari_uporabnika(w http.ResponseWriter, r *http.Request) {
	ime_uporabnika := r.FormValue("ime_uporabnika")
	password := r.FormValue("password")
	email := r.FormValue("email")
	skupinaid,_ := strconv.Atoi(r.FormValue("skupinaid"))
	log.WithFields(log.Fields{"ime_uporabnika": ime_uporabnika}).Info("Dodaj novega uporabnika.")
	log.WithFields(log.Fields{"password": password}).Info("Dodaj password uporabnika.")
	log.WithFields(log.Fields{"email": email}).Info("Dodaj mail uporabnika.")
	log.WithFields(log.Fields{"skupinaid": skupinaid}).Info("Dodaj uporabnika v skupino.")
	do_up := &Uporabnik{Ime_uporabnika: ime_uporabnika, Password: password, Email: email, SkupinaId: skupinaid}
	db.Create(&do_up)
	result := db.Last(&do_up)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.Value)
}

func Posodobi_skupino(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	ime_skupine := r.FormValue("ime_skupine")
	log.WithFields(log.Fields{"Id": id, "Ime_skupine": ime_skupine}).Info("Uredi skupino.")
	do_sk := &Skupina{}
	db.First(&do_sk, id)
	do_sk.Ime_skupine = ime_skupine
	db.Save(&do_sk)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"updated": true}`)
	
}

func Posodobi_uporabnika(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	ime_uporabnika := r.FormValue("ime_uporabnika")
	password := r.FormValue("password")
	email := r.FormValue("email")
	skupinaid,_ := strconv.Atoi(r.FormValue("skupinaid"))
	log.WithFields(log.Fields{"Id": id, "Ime_uporabnika": ime_uporabnika}).Info("Uredi ime uporabnika.")
	log.WithFields(log.Fields{"Id": id, "Password": password}).Info("Uredi geslo uporabnika.")
	log.WithFields(log.Fields{"Id": id, "Email": email}).Info("Uredi email uporabnika.")
	log.WithFields(log.Fields{"Id": id, "SkupinaId": skupinaid}).Info("Uredi uporabnika v skupino.")
	do_up := &Uporabnik{}
	db.First(&do_up, id)
	do_up.Ime_uporabnika = ime_uporabnika
	do_up.Password = password
	do_up.Email = email
	do_up.SkupinaId = skupinaid
	db.Save(&do_up)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"updated": true}`)
	
}

func Izbrisi_skupino(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	var uporabniki []Uporabnik
	for _,uporabnik := range uporabniki{
		if uporabnik.SkupinaId == id{
			do_up := &Uporabnik{}
			db.First(&do_up, uporabnik.SkupinaId)
			db.Delete(&do_up)
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"deleted": true}`)	
		}	
	}
	log.WithFields(log.Fields{"Id": id}).Info("Zbrisi skupino")
	do_sk := &Skupina{}
	db.First(&do_sk, id)
	db.Delete(&do_sk)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"deleted": true}`)
}

func Izbrisi_uporabnika(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	log.WithFields(log.Fields{"Id": id}).Info("Zbrisi uporabnika")
	do_up := &Uporabnik{}
	db.First(&do_up, id)
	db.Delete(&do_up)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"deleted": true}`)
	
}

func DobiSkupinoPoId(Id int) bool {
	do_sk := &Skupina{}
	result := db.First(&do_sk, Id)
	if result.Error != nil{
		log.Warn("Skupine ni v bazi")
		return false
	}
	return true
}

func DobiUporabnikaPoId(Id int) bool {
	do_up := &Uporabnik{}
	result := db.First(&do_up, Id)
	if result.Error != nil{
		log.Warn("Uporabnika ni v bazi")
		return false
	}
	return true
}

func main() {
	defer db.Close()

	db.Debug().AutoMigrate(&Skupina{})
	db.Debug().AutoMigrate(&Uporabnik{})
	

	log.Info("Starting REST API server")
	router := mux.NewRouter()
	router.HandleFunc("/healthz", Healthz).Methods("GET")
	router.HandleFunc("/do_sk", Izpisi_skupine)
	router.HandleFunc("/do_up/{skupinaid}", Izpisi_uporabnike)
	router.HandleFunc("/do_sk/ustvari", Ustvari_skupino).Methods("POST")
	router.HandleFunc("/do_up/{skupinaid}/ustvari", Ustvari_uporabnika).Methods("POST")
	router.HandleFunc("/do_sk/{id}", Posodobi_skupino).Methods("PUT")
	router.HandleFunc("/do_up/{skupinaid}/{id}", Posodobi_uporabnika).Methods("PUT")
	router.HandleFunc("/do_sk/{id}", Izbrisi_skupino).Methods("DELETE")
	router.HandleFunc("/do_up/{skupinaid}/{id}", Izbrisi_uporabnika).Methods("DELETE")	

	handler := cors.New(cors.Options{
                AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
	}).Handler(router)

	http.ListenAndServe(":8000", handler)
}