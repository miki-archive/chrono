package chrono

import (
	"database/sql"
	"github.com/valyala/fasthttp"
	pg "github.com/mikibot/chrono/services/postgres"
	json "github.com/mikibot/chrono/services/json"
	id "github.com/mikibot/chrono/services/snowflake"
	models "github.com/mikibot/chrono/models"
)

// GetTasks gets all tasks and sends it over HTTP
func GetTasks(ctx *fasthttp.RequestCtx) {
	rows, err := pg.Db.Query("SELECT * FROM tasks");
	if(err != nil) {
		ErrorJSON(ctx, "Database error: " + err.Error(), 500);
	}
	handleTasks(ctx, rows);
}

// GetTasksOfID gets tasks for the user and sends it over HTTP
func GetTasksOfID(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("uid").(string);
	rows, err := pg.Db.Query("SELECT * FROM tasks WHERE owner_id = $1", id);
	if(err != nil) {
		ErrorJSON(ctx, "Database error: " + err.Error(), 500);
	}
	handleTasks(ctx, rows);
}

func handleTasks(ctx *fasthttp.RequestCtx, rows *sql.Rows) {
	var tasks []models.TaskModel
	for ;rows.Next(); {
		var id, startTime, endTime int64;
		var ownerID, payload string;

		err := rows.Scan(&id, &ownerID, &payload, &startTime, &endTime);
		if(err != nil) {
			ErrorJSON(ctx, "Scan error: " + err.Error(), 500);
		}
		
		tasks = append(tasks, models.TaskModel{
			TaskID: id, 
			OwnerID: ownerID, 
			StartTime: startTime,
			ExpirationTime: endTime,
			Payload: payload,
		});
	}
	rows.Close();

	packet, err := json.JSON.Marshal(&tasks);
	if(err != nil) {
		ErrorJSON(ctx, "Decode error: " + err.Error(), 500);
	}
	ctx.Success("application/json", packet);
}

// PostTasks queues a new task to the list.
func PostTasks(ctx *fasthttp.RequestCtx) {
	var taskModel models.TaskModel
	
	err := json.JSON.Unmarshal(ctx.PostBody(), &taskModel)
	if(err != nil) {
		ErrorJSON(ctx, "invalid json: " + err.Error(), 100)
	}

	if(len(taskModel.Payload) == 0) {
		ErrorJSON(ctx, "Payload can not be empty. Instead send '{}' for an empty object.", 400);
		return;
	}

	taskID := id.GenerateID();
	taskModel.TaskID = taskID.Int64();

	val, err := json.JSON.Marshal(&taskModel)
	if(err != nil) {
		ErrorJSON(ctx, "invalid json: " + err.Error(), 400);
		return;
	}

	_, err = pg.Db.Query("INSERT INTO tasks (id, owner_id, start_epoch, end_epoch, payload) values ($1, $2, $3, $4, $5);",
		taskModel.TaskID, taskModel.OwnerID, taskModel.StartTime, taskModel.ExpirationTime, taskModel.Payload);
	if(err != nil) {
		ErrorJSON(ctx, "Database error: " + err.Error(), 500);
		return;
	}
	ctx.SuccessString("application/json", string(val));
}