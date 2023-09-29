package httputil

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	ctx  *gin.Context
	code int
	Data any
}

type Message struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"a server message"`
	Data    any    `json:"data,omitempty"`
}

type Error struct {
	Code    int    `json:"code" example:"404"`
	Message string `json:"message" example:"not found"`
	Data    any    `json:"data,omitempty"`
}

func Ctx(ctx *gin.Context) *Response {
	r := &Response{
		ctx: ctx,
	}
	return r
}

func (r *Response) Response(data any) {
	r.ctx.JSON(r.code, data)
}

func (r *Response) Success() {
	r.Message("success")
}

func (r *Response) Message(message string) {
	r.ctx.JSON(r.code, Message{Code: r.code, Message: message})
}

func (r *Response) MessageData(message string, data any) {
	r.ctx.JSON(r.code, Message{Code: r.code, Message: message, Data: data})
}

func (r *Response) Error(err error) {
	r.ctx.JSON(r.code, Error{Code: r.code, Message: err.Error()})
}

func (r *Response) ErrorMessage(message string) {
	r.ctx.JSON(r.code, Error{Code: r.code, Message: errors.New(message).Error()})
}

func (r *Response) ErrorData(err error, data any) {
	r.ctx.JSON(r.code, Error{Code: r.code, Message: err.Error(), Data: data})
}

// Ok Status 200
func (r *Response) Ok() *Response {
	r.code = http.StatusOK
	return r
}

// Created Status 201
func (r *Response) Created() *Response {
	r.code = http.StatusCreated
	return r
}

// Accepted Status 202
func (r *Response) Accepted() *Response {
	r.code = http.StatusAccepted
	return r
}

// NoContent Status 204
func (r *Response) NoContent() *Response {
	r.code = http.StatusNoContent
	return r
}

// ResetContent Status 205
func (r *Response) ResetContent() *Response {
	r.code = http.StatusResetContent
	return r
}

// BadRequest 400
func (r *Response) BadRequest() *Response {
	r.code = http.StatusBadRequest
	return r
}

// Unauthorized Status 401
func (r *Response) Unauthorized() *Response {
	r.code = http.StatusUnauthorized
	return r
}

// Forbidden Status 403
func (r *Response) Forbidden() *Response {
	r.code = http.StatusForbidden
	return r
}

// NotFound Status 404
func (r *Response) NotFound() *Response {
	r.code = http.StatusNotFound
	return r
}

// Conflict Status 409
func (r *Response) Conflict() *Response {
	r.code = http.StatusConflict
	return r
}

// PreconditionFailed Status 412
func (r *Response) PreconditionFailed() *Response {
	r.code = http.StatusPreconditionFailed
	return r
}

// InternalServerError Status 500
func (r *Response) InternalServerError() *Response {
	r.code = http.StatusInternalServerError
	return r
}

// NotImplemented Status 501
func (r *Response) NotImplemented() *Response {
	r.code = http.StatusNotImplemented
	return r
}

// BadGateway Status 502
func (r *Response) BadGateway() *Response {
	r.code = http.StatusBadGateway
	return r
}
