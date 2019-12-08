package controllers

type Context interface {
	Param(string) string
	GetHeader(string) string
	Bind(interface{}) error
	Status(int)
	Query(string) string
	SecureJSON(int, interface{})
	Header(string, string)
}
