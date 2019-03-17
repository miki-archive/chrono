package main

import (
	"log"
	"os"
	"github.com/buaazp/fasthttprouter"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
	_ "github.com/lib/pq"
	controllers "github.com/mikibot/chrono/web/controllers"
	pg "github.com/mikibot/chrono/services/postgres"
	snowflake "github.com/mikibot/chrono/services/snowflake"
)

func main() {
	log.Println("Loading .env")
	err := godotenv.Load()
	if err != nil {
		log.Panicf("Error loading .env file: %s", err)
	}

	log.Print("Creating snowflake generator")
	snowflake.InitSnowflake()

	log.Println("Connecting to pg")

	connStr := createConnectionString()
	pg.InitDB(connStr)
	
	err = pg.Db.Ping()
	if(err != nil) {
		log.Panicf("Could not connect to postgres with connection string '%s': %s", connStr, err)
	}

	log.Println("Starting cron thread")
	// TODO

	log.Println("Opening web service")
	router := fasthttprouter.New()

	router.GET("/tasks", controllers.GetTasks)
	router.GET("/tasks/:uid", controllers.GetTasksOfID);
	router.POST("/tasks", controllers.PostTasks);

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}

func createConnectionString() string {
	connStr := "postgres://"
	user, exists := os.LookupEnv("PG_USER");
	if(exists) {
		connStr += user;
	}

	pass, exists := os.LookupEnv("PG_PASS");
	if(exists) {
		connStr += ":" + pass;
	}

	host, exists := os.LookupEnv("PG_HOST")
	if(exists) {
		connStr += "@" + host
	}

	db, exists := os.LookupEnv("PG_DB")
	if(exists) {
		connStr += "/" + db
	}

	ssl, exists := os.LookupEnv("PG_SSLMODE")
	if(exists) {
		connStr += "?sslmode=" + ssl;
	}

	return connStr;
}