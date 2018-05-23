// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// ListBlobsHandlerFunc turns a function with the right signature into a list blobs handler
type ListBlobsHandlerFunc func(ListBlobsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn ListBlobsHandlerFunc) Handle(params ListBlobsParams) middleware.Responder {
	return fn(params)
}

// ListBlobsHandler interface for that can handle valid list blobs params
type ListBlobsHandler interface {
	Handle(ListBlobsParams) middleware.Responder
}

// NewListBlobs creates a new http.Handler for the list blobs operation
func NewListBlobs(ctx *middleware.Context, handler ListBlobsHandler) *ListBlobs {
	return &ListBlobs{Context: ctx, Handler: handler}
}

/*ListBlobs swagger:route GET /blobs listBlobs

List blobs

*/
type ListBlobs struct {
	Context *middleware.Context
	Handler ListBlobsHandler
}

func (o *ListBlobs) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewListBlobsParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
