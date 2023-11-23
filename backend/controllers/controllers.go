package controllers

type ResponseData struct {
	IsError  bool     `json:"isError"`
	Messages []string `json:"messages"`
	Data     any      `json:"data"`
}
