package courses

type DescriptionPingRequest struct {
	Id           uint32 `json:"id"`
	LastSyncTime int64  `json:"last_sync_time"`
}

func (d DescriptionPingRequest) Bind() {}
