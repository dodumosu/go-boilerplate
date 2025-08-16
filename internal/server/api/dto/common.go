package dto

type ResponseStatus string

const (
	Fail    ResponseStatus = "error"
	Success ResponseStatus = "ok"
)

type BaseRequest struct {
	DeviceID string `header:"X-DEVICE-ID"`
}

type BaseResponse struct {
	Message       string         `json:"message"`
	StatusMessage ResponseStatus `json:"status"`
}

type PagedResponse struct {
	CurrentPage int `json:"currentPage"`
	Total       int `json:"total"`
	FirstPage   int `json:"firstPage"`
	LastPage    int `json:"lastPage"`
}

type PagedRequest struct {
	Page int `query:"page"`
}
