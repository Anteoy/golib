package core

import (
	"net/http"
	"strings"

	"github.com/fudali113/doob/core/register"
	"github.com/fudali113/doob/core/router"

	returnDeal "github.com/fudali113/doob/core/return_deal"
	reflectUtils "github.com/fudali113/doob/utils/reflect"
)

func invoke(matchResult *router.MatchResult, w http.ResponseWriter, req *http.Request) {
	url := req.URL.Path
	method := strings.ToLower(req.Method)

	if matchResult == nil {
		logger.Notice("no match url : %s", url)
		w.WriteHeader(404)
		return
	}

	handlerInterface := matchResult.Rest.GetHandler(method)
	if handlerInterface == nil {
		logger.Notice("match url : %s , but method con`t match", url)
		w.WriteHeader(405)
		return
	}

	/**
	 * 获取路劲参数并存入request参数中
	 */
	urlParam := matchResult.ParamMap
	if urlParam != nil {
		for k, v := range urlParam {
			if req.Form == nil {
				req.Form = map[string][]string{}
			}
			req.Form.Add(k, v)
		}
	}

	/**
	 * 根据RegisterType决定怎么执行函数
	 */
	registerType := matchResult.RegisterType
	if registerType != nil {
		paramType := registerType.ParamType
		returnType := registerType.ReturnType
		switch paramType.Type {
		case register.PARAM_NONE:
			switch returnType.Type {
			case register.RETURN_NONE:
				handler := handlerInterface.(func())
				handler()

			case register.JSON:
				handler := handlerInterface.(ReturnObject)
				returnValue := handler()
				returnDeal.DealReturn(&returnDeal.ReturnType{
					TypeStr: "json",
					Data:    returnValue,
				}, w)

			case register.FILE:
				handler := handlerInterface.(ReturnStr)
				returnValue := handler()
				returnDeal.DealReturn(&returnDeal.ReturnType{TypeStr: returnValue}, w)

			case register.RETURN_TYPE:
				handler := handlerInterface.(ReturnType)
				str, data := handler()
				returnDeal.DealReturn(&returnDeal.ReturnType{
					TypeStr: str,
					Data:    data,
				}, w)
			}

		case register.ORIGIN:
			handler := handlerInterface.(http.HandlerFunc)
			handler(w, req)

		case register.CTX:
			context := getContext(w, req)
			switch returnType.Type {
			case register.RETURN_NONE:
				handler := handlerInterface.(CTXReturn)
				handler(context)

			case register.FILE:
				handler := handlerInterface.(CTXReturnStr)
				returnValue := handler(context)
				returnDeal.DealReturn(&returnDeal.ReturnType{TypeStr: returnValue}, w)

			case register.JSON:
				handler := handlerInterface.(CTXReturnObject)
				returnValue := handler(context)
				returnDeal.DealReturn(&returnDeal.ReturnType{
					TypeStr: "json",
					Data:    returnValue,
				}, w)

			case register.RETURN_TYPE:
				handler := handlerInterface.(CTXReturnType)
				str, data := handler(context)
				returnDeal.DealReturn(&returnDeal.ReturnType{
					TypeStr: str,
					Data:    data,
				}, w)
			}

		case register.CI_PATHVARIABLE, register.CI_PATHVARIABLE_CTX:
			var returns []interface{}
			ciLen := paramType.CiLen
			paraNames := matchResult.ParamNames
			if ciLen > len(paraNames) {
				logger.Warn(`your func path variable params lnegth is %d ,
           but your url params length just %d`, ciLen, len(paraNames))
				return
			}

			params := []interface{}{}
			for i := 0; i < ciLen; i++ {
				params = append(params, urlParam[paraNames[i]])
			}
			if paramType.Type == register.CI_PATHVARIABLE_CTX {
				params = append(params, getContext(w, req))
			}
			returns = reflectUtils.Invoke(handlerInterface, params)

			switch returnType.Type {
			case register.RETURN_NONE:

			case register.FILE:
				str := returns[0].(string)
				returnDeal.DealReturn(&returnDeal.ReturnType{TypeStr: str}, w)

			case register.JSON:
				returnValue := returns[0].(interface{})
				returnDeal.DealReturn(&returnDeal.ReturnType{
					TypeStr: "json",
					Data:    returnValue,
				}, w)

			case register.RETURN_TYPE:
				str := returns[0].(string)
				returnValue := returns[1]
				returnDeal.DealReturn(&returnDeal.ReturnType{
					TypeStr: str,
					Data:    returnValue,
				}, w)
			}
			return
		}
	}
}

/**
 * 根据res&req获取context
 */
func getContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		request:  req,
		response: w,
		Params:   map[string]string{},
	}
}
