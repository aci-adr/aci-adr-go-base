package common

type ProcessingInfo struct {
	Stage          string `json:"stage"`
	NetworkTime    int64  `json:"networkTime"`
	ProcessingTime int64  `json:"processingTime"`
	StartedOn      int64  `json:"startedOn"`
	CompletedOn    int64  `json:"completedOn"`
	Status         string `json:"status"`
}
