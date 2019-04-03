package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"
)


type Record_row struct {
	//Id         string  `json:"id" bson:"_id,omitempty"`
	Base_bbl   int64   `json:"base_bbl,omitempty" bson:"base_bbl,omitempty"`
	Bin        int64   `json:"bin,omitempty" bson:"bin,omitempty"`
	Cnstrct_yr int64   `json:"cnstrct_yr,omitempty" bson:"cnstrct_yr,omitempty"`
	Doitt_id   int64   `json:"doitt_id,omitempty" bson:"doitt_id,omitempty"`
	Feat_code  int64   `json:"feat_code,omitempty" bson:"feat_code,omitempty"`
	Geomsource string  `json:"geomsource,omitempty" bson:"geomsource,omitempty"`
	Groundelev int64   `json:"groundelev,omitempty" bson:"groundelev,omitempty"`
	Heightroof float64 `json:"heightroof,omitempty" bson:"heightroof,omitempty"`
	Lstmoddate string  `json:"lstmoddate,omitempty" bson:"lstmoddate,omitempty"`
	Lststatype string  `json:"lststatype,omitempty" bson:"lststatype,omitempty"`
	Mpluto_bbl int64   `json:"mpluto_bbl,omitempty" bson:"mpluto_bbl,omitempty"`
	Name       string  `json:"name,omitempty" bson:"name,omitempty"`
	Shape_area float64 `json:"shape_area,omitempty" bson:"shape_area,omitempty"`
	Shape_len  float64 `json:"shape_len,omitempty" bson:"shape_len,omitempty"`
	The_geom   string  `json:"the_geom,omitempty" bson:"the_geom,omitempty"`
}

type Metric struct {
	Metric    string
	Value float64
}

var client *mongo.Client




func main() {
	fmt.Println("Starting the application...")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	router := mux.NewRouter()
	router.HandleFunc("/", WelcomeHome)
	router.HandleFunc("/GetAllRecord", GetAllRecordsEndpoint).Methods("GET")
	router.HandleFunc("/GetRecord", GetRecordEndpoint).Methods("POST")
	router.HandleFunc("/Count", GetCountEndpoint).Methods("GET")
	router.HandleFunc("/AverageHeight", GetAverageEndpoint).Methods("GET")
	router.HandleFunc("/SumGroundelev", GetSumEndpoint).Methods("GET")
	router.HandleFunc("/GetOldest", GetOldestEndpoint).Methods("GET")
	log.Fatal(http.ListenAndServe(":12345", router))
}

func WelcomeHome(w http.ResponseWriter, r *http.Request) {
	p := path.Dir("./inc/index.html")
	w.Header().Set("Content-type", "text/html")
	http.ServeFile(w, r, p)
}

func GetCountEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	collection := client.Database("mydb").Collection("record_data")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	i := 0
	for cursor.Next(context.Background()) {
		i++
	}
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	m:=Metric{"Count",float64(i)}
	json.NewEncoder(response).Encode(m)
}

func GetAverageEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	collection := client.Database("mydb").Collection("record_data")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	i := 0
	sumHeight :=0.0

	for cursor.Next(context.Background()) {
		i++
		var rec Record_row
		cursor.Decode(&rec)
		sumHeight +=rec.Heightroof
	}
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	m:=Metric{"Average",float64(sumHeight/float64(i))}
	json.NewEncoder(response).Encode(m)
}

func GetSumEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	collection := client.Database("mydb").Collection("record_data")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	i := 0
	sumHeight :=int64(0)

	for cursor.Next(context.Background()) {
		i++
		var rec Record_row
		cursor.Decode(&rec)
		sumHeight +=rec.Groundelev
	}
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	m:=Metric{"Sum Ground Elevation",float64(sumHeight)}
	json.NewEncoder(response).Encode(m)
}

func GetRecordEndpoint(response http.ResponseWriter, request *http.Request) {
	fmt.Println("method:", request.Method)
	request.ParseForm()
	// logic part of log in

	val,_:=strconv.ParseInt(request.Form["base_bbl"][0],10,64)
	fmt.Println("method:", val)

	response.Header().Set("content-type", "application/json")
	//params := mux.Vars(request)
	//id,_ := strconv.ParseInt(params["base_bbl"],10,64)
	var rec Record_row
	collection := client.Database("mydb").Collection("record_data")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Record_row{Base_bbl: val}).Decode(&rec)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(rec)
}

func GetAllRecordsEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var recs []Record_row
	collection := client.Database("mydb").Collection("record_data")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var rec Record_row
		cursor.Decode(&rec)
		recs = append(recs, rec)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(recs)
}

func GetOldestEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	collection := client.Database("mydb").Collection("record_data")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	var min=int64(2020)
	for cursor.Next(ctx) {
		var rec Record_row
		cursor.Decode(&rec)
		//recs = append(recs, rec)
		if rec.Cnstrct_yr<min && rec.Cnstrct_yr>0 {
			min=rec.Cnstrct_yr
		}
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	cursor1, err := collection.Find(ctx, Record_row{Cnstrct_yr: min})
	var recs []Record_row
	for cursor1.Next(ctx) {
		var rec Record_row
		cursor1.Decode(&rec)
		recs = append(recs, rec)
	}
	m:=Metric{"Oldest Building Year",float64(min)}
	json.NewEncoder(response).Encode(m)
	json.NewEncoder(response).Encode(recs)
}

