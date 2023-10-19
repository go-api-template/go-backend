package httputil

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response is a wrapper for gin.Context
// It helps to make response easier
type Response struct {
	ctx    *gin.Context // gin context
	status int          // http status status
}

// Success is a struct for success response
// It is used when the response is a success
type Success struct {
	r           *Response
	Code        string `json:"code,omitempty" example:"server_code"`
	Description string `json:"description,omitempty" example:"a server message"`
	Data        any    `json:"data,omitempty"`
	raw         any
}

// Error is a struct for error response
// It is used when the response is an error
type Error struct {
	r           *Response
	Code        string `json:"code,omitempty" example:"server_code"`
	Description string `json:"description,omitempty" example:"a server message"`
	Errors      []any  `json:"errors,omitempty"`
	Data        any    `json:"data,omitempty"`
}

// Ctx is a constructor for Response
func Ctx(ctx *gin.Context) *Response {
	r := &Response{
		ctx: ctx,
	}
	return r
}

// Success send a simple success response
func (s *Success) Success() {
	s.WithDescription("success").Send()
}

// WithCode set the code of the response
func (s *Success) WithCode(code string) *Success {
	s.Code = code
	return s
}

// WithDescription set the description of the response
func (s *Success) WithDescription(description string) *Success {
	s.Description = description
	return s
}

// WithData set the data of the response
func (s *Success) WithData(data any) *Success {
	s.Data = data
	return s
}

// SendRaw is used to send the response with raw data
func (s *Success) SendRaw(raw any) {
	s.raw = &raw
	s.Send()
}

// Send is used to send the response
func (s *Success) Send() {
	if s.raw != nil {
		s.r.ctx.JSON(s.r.status, s.raw)
		return
	}
	s.r.ctx.JSON(s.r.status, s)
}

// WithCode set the code of the response
func (e *Error) WithCode(code string) *Error {
	e.Code = code
	return e
}

// WithDescription set the description of the response
func (e *Error) WithDescription(description string) *Error {
	e.Description = description
	return e
}

// WithErrors set the errors of the response
func (e *Error) WithErrors(errors []any) *Error {
	e.Errors = errors
	return e
}

// WithError add an error to the errors of the response
func (e *Error) WithError(err any) *Error {
	e.Errors = append(e.Errors, err)
	return e
}

// WithData set the data of the response
func (e *Error) WithData(data *any) *Error {
	e.Data = data
	return e
}

// Send is used to send the response
func (e *Error) Send() {
	e.r.ctx.JSON(e.r.status, e)
}

// Ok Status 200
func (r *Response) Ok() *Success {
	r.status = http.StatusOK
	return &Success{r: r}
}

// Created Status 201
func (r *Response) Created() *Success {
	r.status = http.StatusCreated
	return &Success{r: r}
}

// Accepted Status 202
func (r *Response) Accepted() *Success {
	r.status = http.StatusAccepted
	return &Success{r: r}
}

// NoContent Status 204
func (r *Response) NoContent() *Success {
	r.status = http.StatusNoContent
	return &Success{r: r}
}

// ResetContent Status 205
func (r *Response) ResetContent() *Success {
	r.status = http.StatusResetContent
	return &Success{r: r}
}

// BadRequest 400
func (r *Response) BadRequest() *Error {
	r.status = http.StatusBadRequest
	return &Error{r: r}
}

// Unauthorized Status 401
func (r *Response) Unauthorized() *Error {
	r.status = http.StatusUnauthorized
	return &Error{r: r}
}

// Forbidden Status 403
func (r *Response) Forbidden() *Error {
	r.status = http.StatusForbidden
	return &Error{r: r}
}

// NotFound Status 404
func (r *Response) NotFound() *Error {
	r.status = http.StatusNotFound
	return &Error{r: r}
}

// Conflict Status 409
func (r *Response) Conflict() *Error {
	r.status = http.StatusConflict
	return &Error{r: r}
}

// PreconditionFailed Status 412
func (r *Response) PreconditionFailed() *Error {
	r.status = http.StatusPreconditionFailed
	return &Error{r: r}
}

// InternalServerError Status 500
func (r *Response) InternalServerError() *Error {
	r.status = http.StatusInternalServerError
	return &Error{r: r}
}

// NotImplemented Status 501
func (r *Response) NotImplemented() *Error {
	r.status = http.StatusNotImplemented
	return &Error{r: r}
}

// BadGateway Status 502
func (r *Response) BadGateway() *Error {
	r.status = http.StatusBadGateway
	return &Error{r: r}
}
