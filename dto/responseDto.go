package dto

type ServiceResponse struct {
	Resource string
	Status   string
}

type ServerResponse struct {
	Status           int32
	Message string
	ServiceResponses []ServiceResponse
	Failed           []string
}

type ErrorResponse struct {
	Status int32
	Error  string
}
