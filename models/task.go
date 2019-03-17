package chrono

// TaskModel contains all metadata for a redis model.
type TaskModel struct {
	OwnerID string
	TaskID int64
	StartTime int64
	ExpirationTime int64
	Payload string
	EndpointID int
}