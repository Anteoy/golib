package doob

var (
	isDev                   = true
	urlSplitSymbol          = "&&"
	returnDealDefaultType   = "auto"
	redirectDefaultBodytLen = 1024
)

func SetReturnDealDefaultType(t string) {
	returnDealDefaultType = t
}

func SetUrlSplitSymbol(symbol string) {
	urlSplitSymbol = symbol
}

func SetRedirectDefaulBodytLen(len int) {
	redirectDefaultBodytLen = len
}

func SetIsDev(isDevStatus bool) {
	isDev = isDevStatus
}
