package web

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

func (h HttpError) Error() string {
	//TODO implement me
	panic("implement me")
}
