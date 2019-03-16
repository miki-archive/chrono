package main

import (
	"github.com/mikibot/chrono"
	"log"
	"os"
	"github.com/buaazp/fasthttprouter"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
	"github.com/json-iterator/go"
	"github.com/bwmarrin/snowflake"
	_ "github.com/lib/pq"
)

var json = jsoniter.ConfigFastest
var uidGen *snowflake.Node

// ErrorModel is used when thrown an API error
type ErrorModel struct {
	Message string
	ID int
}

// GetTasks gets all tasks and sends it over HTTP
func GetTasks(ctx *fasthttp.RequestCtx) {
	// TODO (Veld): Parse to taskModel array
	//values, err := redisClient.HVals(redisHashKey).Result();
	// if(err != nil) {
	// 	ErrorJSON(ctx, "fuck", 0)
	// }
	// val, err := json.Marshal(values);
	// if(err != nil) {
	// 	ErrorJSON(ctx, "deserialization error: " + err.Error(), 200)
	// }
	//ctx.Success("application/json", val);
}

// GetTasksOfID gets tasks for the user and sends it over HTTP
func GetTasksOfID(ctx *fasthttp.RequestCtx) {
	//var id = ctx.UserValue("uid").(string);
	//redisClient.h(redisHashKey)
}

// PostTasks queues a new task to the list.
func PostTasks(ctx *fasthttp.RequestCtx) {
	var taskModel TaskModel
	
	err := json.Unmarshal(ctx.PostBody(), &taskModel)
	if(err != nil) {
		ErrorJSON(ctx, "invalid json: " + err.Error(), 100)
	}

	taskID := GenerateSnowflake();
	taskModel.TaskID = taskID.Int64();

	val, err := json.Marshal(&taskModel)
	if(err != nil) {
		ErrorJSON(ctx, "invalid json: " + err.Error(), 100);
	}	

	//services.Db.Query("INSERT INTO tasks")
	ctx.SuccessString("application/json", string(val));
}

func main() {
	log.Println("Loading .env")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Print("Creating snowflake generator")
	node, err := snowflake.NewNode(1)
	if(err != nil) {
		log.Fatal("Could not create uid generator.\n" + err.Error())
	}
	uidGen = node

	log.Println("Connecting to pg")

	connStr := createConnectionString();
	InitDb(connStr);

	log.Println("Starting cron thread")
//	go redisTick(client)

	log.Println("Opening web service")
	router := fasthttprouter.New()

	router.GET("/tasks", GetTasks)
	router.GET("/tasks/:uid", GetTasksOfID);
	router.POST("/tasks", PostTasks);

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}

// GenerateSnowflake creates twitter-like snowflakes to base64
func GenerateSnowflake() snowflake.ID {
	return uidGen.Generate();
}

// ErrorJSON is a utility function to send an error from Json to the requester.
func ErrorJSON(ctx *fasthttp.RequestCtx, message string, id int) {
	str, err := json.Marshal(ErrorModel{message,id});
	if(err != nil) {
		log.Fatalln("Error erroring... wtf?\n" + err.Error());
	}
	log.Println("err: " + message);
	ctx.Error(string(str), 500);
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

	host, exists := os.LookupEnv("PG_HOST");
	if(exists) {
		connStr += "@" + host;
	}

	db, exists := os.LookupEnv("PG_DB");
	if(exists) {
		connStr += "/" + db;
	}

	return connStr;
}