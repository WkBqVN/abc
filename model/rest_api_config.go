package model

type RestApiConfig struct {
	Prefix   string   `json:"prefix"`
	Service  string   `json:"service"`
	Database string   `json:"database"`
	Methods  []Method `json:"methods"`
	EndPoint string   `json:"endPoint"`
}

type Method struct {
	MethodType  string `json:"type"`
	HandlerFunc string `json:"handler"`
	Params      string `json:"params"`
}
