package message

type ResponseMsg struct {
	Type    string `json:"type"`
	Code    int    `json:"code"`
	Err     error  `json:"error"`
	Id      string `json:"id"`
	Message string `json:"message"`
}

type RequestMsg struct {
	Method  string `json:"method"`
	Id      string `json:"id"`
	Message string `json:"message"`
}
