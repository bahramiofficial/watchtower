package helper

import "github.com/bahramiofficial/watchtower/src/api/validations"

type BaseHttpResponse struct {
	Result           any                            `json:"result"`
	Success          bool                           `json:"success"`
	ResultCode       int                            `json:"resultCode"`
	ValidationErrors *[]validations.ValidationError `json:"validationErrors"`
	Error            any                            `json:"error"`
}

func GenerateBaseHttpResponse(result any, success bool, resultCode int) *BaseHttpResponse {
	return &BaseHttpResponse{
		Result:     result,
		Success:    success,
		ResultCode: resultCode,
	}
}

func GenerateBaseHttpResponseWithError(result any, success bool, resultCode int, err error) *BaseHttpResponse {
	return &BaseHttpResponse{
		Result:     result,
		Success:    success,
		ResultCode: resultCode,
		Error:      err.Error(),
	}
}

func GenerateBaseHttpResponseWithValidatorError(result any, success bool, resultCode int, err error) *BaseHttpResponse {
	return &BaseHttpResponse{
		Result:           result,
		Success:          success,
		ResultCode:       resultCode,
		ValidationErrors: validations.GetValidationErrors(err),
	}
}
