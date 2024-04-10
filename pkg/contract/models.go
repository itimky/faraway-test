package contract

type ReqChallenge struct{}

type RepChallenge struct {
	Challenge string `json:"challenge"`
}

type ReqQuote struct {
	Solution int64 `json:"solution"`
}

type RepQuote struct {
	Error string `json:"error"`
	Quote string `json:"quote"`
}
