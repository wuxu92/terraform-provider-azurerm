package resourcemanager

import (
	"context"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type RequestOption func(*client.RequestOptions)
type requestOptionKeyType struct{}

var requestOptionKey = requestOptionKeyType{}

func NewCtxWithRequestOption(ctx context.Context, f RequestOption) context.Context {
	c := context.WithValue(ctx, requestOptionKey, f)
	return c
}

func applyCtxRequestOptions(ctx context.Context, opts *client.RequestOptions) {
	v, ok := ctx.Value(requestOptionKey).(RequestOption)
	if ok && v != nil {
		v(opts)
	}
}

type DefaultOption struct {
	Headers map[string]string
	Queries map[string]string
	OData   *odata.Query
}

// ToHeaders implements client.Options.
func (d DefaultOption) ToHeaders() *client.Headers {
	result := &client.Headers{}
	for k, v := range d.Headers {
		result.Append(k, v)
	}
	return result
}

// ToOData implements client.Options.
func (d DefaultOption) ToOData() *odata.Query {
	return d.OData
}

// ToQuery implements client.Options.
func (d DefaultOption) ToQuery() *client.QueryParams {
	result := &client.QueryParams{}
	for k, v := range d.Queries {
		result.Append(k, v)
	}
	return result
}

var _ client.Options = DefaultOption{}
