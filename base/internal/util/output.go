package util

import (
	"strconv"

	ae "github.com/blackflagsoftware/forge-go/base/internal/api_error"
)

type (
	Output struct {
		Payload interface{} `json:"data,omitempty"`
		*Error  `json:"error,omitempty"`
		*Meta   `json:"meta,omitempty"`
	}

	Error struct {
		Id     string `json:"Id,omitempty"`
		Title  string `json:"Title,omitempty"`
		Detail string `json:"Detail,omitempty"`
		Status string `json:"Status,omitempty"`
	}

	Meta struct {
		TotalCount int `json:"total_count"`
	}
)

func NewOutput(payload interface{}, apiError *ae.ApiError, totalCount *int) Output {
	var err *Error
	var meta *Meta
	if apiError != nil {
		err = &Error{Id: apiError.ApiErrorCode, Title: apiError.Title, Detail: apiError.Detail, Status: strconv.Itoa(apiError.StatusCode)}
	}
	if totalCount != nil {
		meta = &Meta{TotalCount: *totalCount}
	}
	output := Output{
		Payload: payload,
		Error:   err,
		Meta:    meta,
	}
	return output
}
