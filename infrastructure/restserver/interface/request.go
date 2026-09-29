package restserverinterface

type RequestModel interface {
	FillFromBytes(body []byte) error
}

type ValidateRequestModel interface {
	RequestModel
	ValidateRequest() error
}
