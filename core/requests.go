package core

type (
	RegisterWorkerRequest struct {
		Port int `json:"port"`
	}

	RegisterWorkerResponse struct {
		Status string
		WorkerId uint
	}
)
