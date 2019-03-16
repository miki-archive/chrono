package chrono

// TaskModel contains all metadata for a redis model.
type TaskModel struct {
	OwnerID string
	TaskID int64
	StartTime int64
	ExpirationTime int64
	Payload string
}

// GetTaskByID returns a taskmodel of ID 'id' if exists, otherwise returns nil
func GetTaskByID(id int64) TaskModel {
	
}
