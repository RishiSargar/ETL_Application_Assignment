
Welcome to the ETL application !
#Author: Rushikesh Sargar (rsargar@asu.edu)

#ETL_API

API which provides an interface to perform ETL tasks on the data fetched from a vendor site

The application has been written in Golang and supported by MongoDB.

IDE used :JetBrains GoLand 2019.1

Infrastructure Setup :
1.Go Lang installed on the system.
2. MongoDB server running on port 27017
    Create a db "mydb" in it by using command "use mydb"

Application has two components each implemented in separate file
A. Data Load (Data_Load.go)
B. ETL API (ETL_API.go)

A. Data Load
1. Data Load component fetches the data from the file :
    http://data.cityofnewyork.us/resource/mtik-6c5q.csv
    which a source data for the entire application.
2. Application allows user to download subset of data by providing the count of rows as a command line.
   Note: User must provide a value integer value as row_count in command line
3. The data will be fetched and inserted into database "mydb" and collection "record_data" after dropping the previously added data.
4. Same can be confirmed by using "db.record_data.find().count()" on MongoDB command line
The data consumes all the columns in the data base :
    Base_bbl   int64
	Bin        int64
	Cnstrct_yr int64
	Doitt_id   int64
	Feat_code  int64
	Geomsource string
	Groundelev int64
	Heightroof float64
	Lstmoddate string
	Lststatype string
	Mpluto_bbl int64
	Name       string
	Shape_area float64
	Shape_len  float64
	The_geom   string

Execution Step :
1. Make sure Go is installed, switch to the directory where you have Data_Load.go file.
2. hit "go build Data_Load.go"
3. hit "go run Data_Load.go <number of rows to be pulled>"

B. ETL API
1. This components provides an API to access the data stored in the database and get some of the insights of the data.
2. Once ETL_API.go is executed, application will be launched on localhost on port number 12345
    browse: http://localhost:12345 on your browser once the application starts
3. It loads and Index page by default which provides different options to get the insightful metrics of the data
    following are the options provided:
    1. Get Count of all records (provides the number of records in the DB in json format on the screen)
    2. Get Average Height of all buildings ( Calculates the average height of all the buildings)
    3. Get Sum of Ground Elevation (Calculates the sum of Ground Elevation for all the buildings)
    4. Get All Records (Prints all the records)
    5. Get Oldest Buildings (Prints the year in which the oldest building was built and prints all the records with same Construction year)
    6. Get Record by Base_bbl (Sends a POST request to the function and Returns the building record with the provided Base_bbl)

    All the metrics are returned in the json format e.g {"Metric":"Count","Value":500}
    mux router has been used to implement the different handler for different functions.

Execution Step :
1. hit "go build ETL_API.go"
2. Make sure port 12345 is free, hit "go run ETL_API.go"
3. Visit http://localhost:12345 on your web browser
4. Click on links to get insights


