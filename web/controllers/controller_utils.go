package chrono

import (
	"log"
	"github.com/valyala/fasthttp"
	models "github.com/mikibot/chrono/models"
	json "github.com/mikibot/chrono/services/json"
)

// ErrorJSON is a utility function to send an error from Json to the requester.
func ErrorJSON(ctx *fasthttp.RequestCtx, message string, id int) {
	str, err := json.JSON.Marshal(models.ErrorModel{Message: message, ID: id});
	if(err != nil) {
		log.Fatalln("Error erroring... wtf?\n" + err.Error());
	}
	log.Println("err: " + message);
	ctx.Error(string(str), id);
}