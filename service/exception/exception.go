package exception

type Handler struct {
	Error   error
	Message string
	Code    int
}
