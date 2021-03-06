package errors

import "fmt"

type DoobError struct {
	Err       error
	Desc      string
	HttpStaus int
}

func (this DoobError) Error() string {
	return fmt.Sprintf(`request error : should return http status code is <%d> ,
		 description is <%s> , error is <%v>`, this.HttpStaus, this.Desc, this.Err)
}

type RequestError struct {
	DoobError
	Url string
}

func (this RequestError) Error() string {
	return fmt.Sprintf("url is <%s> , error is %v", this.Url, this.DoobError)
}

type MehtodNotMatchError struct {
	RequestError
	Method string
}

func (this MehtodNotMatchError) Error() string {
	return fmt.Sprintf("url is <%s> ,method is <%s> , error is %v", this.Url, this.Method, this.DoobError)
}
