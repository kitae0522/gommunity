package exception

func GenerateErrorCtx(statusCode int, errMessage string, err interface{}) *ErrResponseCtx {
	return &ErrResponseCtx{
		IsError:    true,
		StatusCode: statusCode,
		Message:    errMessage,
		Error:      getErrorDetail(err),
	}
}

func getErrorDetail(err interface{}) interface{} {
	var errDetail interface{}
	if e, ok := err.(error); ok {
		errDetail = e.Error()
	} else {
		errDetail = err
	}
	return errDetail
}
