package chrono

import (
	"time"

	pg "github.com/mikibot/chrono/services/postgres"
)

var nextTick int64;

// Run should be ran in a coroutine, this function will be running in the background and process tasks.
func Run() {
	for ;true; {
		

		time.Sleep(time.Second);
	}
}