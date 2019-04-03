package main

import (
	"bufio"
	"fmt"
	"gopkg.in/mgo.v2"
	"os"
	"strconv"
	"time"

	//"time"

	//"io"
	"io/ioutil"
	"log"
	"net/http"
	//"strconv"
	"strings"
)
type Record_row struct {
	//Id         string  `json:"id" bson:"_id,omitempty"`
	Base_bbl   int64   `json:"base_bbl" bson:"base_bbl"`
	Bin        int64   `json:"bin" bson:"bin"`
	Cnstrct_yr int64   `json:"cnstrct_yr" bson:"cnstrct_yr"`
	Doitt_id   int64   `json:"doitt_id" bson:"doitt_id"`
	Feat_code  int64   `json:"feat_code" bson:"feat_code"`
	Geomsource string  `json:"geomsource" bson:"geomsource"`
	Groundelev int64   `json:"groundelev" bson:"groundelev"`
	Heightroof float64 `json:"heightroof" bson:"heightroof"`
	Lstmoddate string  `json:"lstmoddate" bson:"lstmoddate"`
	Lststatype string  `json:"lststatype" bson:"lststatype"`
	Mpluto_bbl int64   `json:"mpluto_bbl" bson:"mpluto_bbl"`
	Name       string  `json:"name" bson:"name"`
	Shape_area float64 `json:"shape_area" bson:"shape_area"`
	Shape_len  float64 `json:"shape_len" bson:"shape_len"`
	The_geom   string  `json:"the_geom" bson:"the_geom"`
}
func main() {
	rowCount := os.Args[1]
	fmt.Println("RowCount Entered: ",rowCount)
	maxWait := time.Duration(10 * time.Second)
	session, err := mgo.DialWithTimeout("localhost:27017",maxWait)
	if err != nil {
		fmt.Println("Unable to connect to local mongo instance!")
		}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("mydb").C("record_data")
	c.DropCollection()
	if ( c != nil ) {
		fmt.Println("Got a collection object")
		}
	defer session.Close()

	fileUrl := "http://data.cityofnewyork.us/resource/mtik-6c5q.csv?$limit="+rowCount

	resp, err := http.Get(fileUrl)
	if err != nil {
		fmt.Println("No Data Fethed")
	} else
	{
		fmt.Println("Data Fethed")
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			fmt.Println("hell")
		}
		bodyString := string(bodyBytes)
		//fmt.Println(bodyString)
		scanner := bufio.NewScanner(strings.NewReader(bodyString))
		for scanner.Scan() {


			record:= scanner.Text()
			s:=strings.Split(record,",")
			if strings.Replace(strings.Replace(s[0],"\"","",-1),"\\","",-1)=="base_bbl" {
				continue;
			}
			field0,_:= strconv.ParseInt(strings.Replace(strings.Replace(s[0],"\"","",-1),"\\","",-1), 10, 64)
			field1,_:= strconv.ParseInt(strings.Replace(strings.Replace(s[1],"\"","",-1),"\\","",-1), 10, 64)
			field2,_:= strconv.ParseInt(strings.Replace(strings.Replace(s[2],"\"","",-1),"\\","",-1), 10, 64)
			field3,_:= strconv.ParseInt(strings.Replace(strings.Replace(s[3],"\"","",-1),"\\","",-1), 10, 64)
			field4,_:= strconv.ParseInt(strings.Replace(strings.Replace(s[4],"\"","",-1),"\\","",-1), 10, 64)
			field5:= strings.Replace(strings.Replace(s[5],"\"","",-1),"\\","",-1)
			field6,_:= strconv.ParseInt(strings.Replace(strings.Replace(s[6],"\"","",-1),"\\","",-1), 10, 64)
			field7,_:= strconv.ParseFloat(strings.Replace(strings.Replace(s[7],"\"","",-1),"\\","",-1), 64)
			field8:= strings.Replace(strings.Replace(s[8],"\"","",-1),"\\","",-1)
			field9:= strings.Replace(strings.Replace(s[9],"\"","",-1),"\\","",-1)
			field10,_:= strconv.ParseInt(strings.Replace(strings.Replace(s[10],"\"","",-1),"\\","",-1), 10, 64)
			field11:= strings.Replace(strings.Replace(s[11],"\"","",-1),"\\","",-1)
			field12,_:= strconv.ParseFloat(strings.Replace(strings.Replace(s[12],"\"","",-1),"\\","",-1), 64)
			field13,_:= strconv.ParseFloat(strings.Replace(strings.Replace(s[13],"\"","",-1),"\\","",-1), 64)
			field14:= strings.Replace(strings.Replace(s[14],"\"","",-1),"\\","",-1)
				err = c.Insert(&Record_row{field0, field1, field2, field3, field4, field5,field6, float64(field7), field8,field9, field10, field11,float64(field12), float64(field13), field14})

				if err != nil {
					panic(err)
				}
				//log.Printf("%#v", record)


		}
		log.Printf("Data inserted into MongoDB")

	}

}
