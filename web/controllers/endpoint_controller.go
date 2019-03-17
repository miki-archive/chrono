package chrono

import (
	"log"
	"database/sql"
	"github.com/valyala/fasthttp"
	
	json "github.com/mikibot/chrono/services/json"
	models "github.com/mikibot/chrono/models"
	pg "github.com/mikibot/chrono/services/postgres"
	id "github.com/mikibot/chrono/services/snowflake"
)

// GetEndpoints returns all endpoints in a json format.
func GetEndpoints(ctx *fasthttp.RequestCtx) {
	rows, err := pg.Db.Query("SELECT * FROM task_endpoints");
	if(err != nil) {
		ErrorJSON(ctx, "could not get task endpoints", 500);
		log.Println(err);
		return;
	}
	handleEndpoints(ctx, rows);
}

// GetEndpointByID returns a specific endpoint based on id.
func GetEndpointByID(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("uid").(string);
	rows, err := pg.Db.Query("SELECT * FROM task_endpoints WHERE id = $1", id);
	if(err != nil) {
		ErrorJSON(ctx, "could not get task endpoints", 500);
		log.Println(err);
		return;
	}
	handleEndpoints(ctx, rows);
}

// PostEndpoint inserts a new endpoint if valid.
func PostEndpoint(ctx *fasthttp.RequestCtx) {
	var endpointModel models.EndpointModel
	
	err := json.JSON.Unmarshal(ctx.PostBody(), &endpointModel)
	if(err != nil) {
		ErrorJSON(ctx, "invalid json", 400);
		log.Println(err);
		return;
	}

	endpointID := id.GenerateID();
	endpointModel.ID = endpointID.Int64();

	val, err := json.JSON.Marshal(&endpointModel)
	if(err != nil) {
		ErrorJSON(ctx, "invalid json", 400);
		log.Println(err);
		return;
	}

	err = pg.InsertEndpoint(endpointModel);
	if(err != nil) {
		ErrorJSON(ctx, "could not add endpoint to database", 500);
		log.Println(err);
		return;
	}
	ctx.SuccessString("application/json", string(val));
}

func handleEndpoints(ctx *fasthttp.RequestCtx, rows *sql.Rows) {
	var endpoints []models.EndpointModel
	for ;rows.Next(); {
		var id int64;
		var url string;

		err := rows.Scan(&id, &url);
		if(err != nil) {
			log.Println(err);
		}
		
		endpoints = append(endpoints, models.EndpointModel{
			ID: id, 
			URL: url, 
		});
	}
	rows.Close();

	packet, err := json.JSON.Marshal(&endpoints);
	if(err != nil) {
		ErrorJSON(ctx, "invalid json", 500);
		log.Println(err);
		return;
	}
	ctx.Success("application/json", packet);
}