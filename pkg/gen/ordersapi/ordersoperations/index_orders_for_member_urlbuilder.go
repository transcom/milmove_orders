// Code generated by go-swagger; DO NOT EDIT.

package ordersoperations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"
	"strings"
)

// IndexOrdersForMemberURL generates an URL for the index orders for member operation
type IndexOrdersForMemberURL struct {
	Edipi string

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *IndexOrdersForMemberURL) WithBasePath(bp string) *IndexOrdersForMemberURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *IndexOrdersForMemberURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *IndexOrdersForMemberURL) Build() (*url.URL, error) {
	var _result url.URL

	var _path = "/edipis/{edipi}/orders"

	edipi := o.Edipi
	if edipi != "" {
		_path = strings.Replace(_path, "{edipi}", edipi, -1)
	} else {
		return nil, errors.New("edipi is required on IndexOrdersForMemberURL")
	}

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/orders/v1"
	}
	_result.Path = golangswaggerpaths.Join(_basePath, _path)

	return &_result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *IndexOrdersForMemberURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *IndexOrdersForMemberURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *IndexOrdersForMemberURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on IndexOrdersForMemberURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on IndexOrdersForMemberURL")
	}

	base, err := o.Build()
	if err != nil {
		return nil, err
	}

	base.Scheme = scheme
	base.Host = host
	return base, nil
}

// StringFull returns the string representation of a complete url
func (o *IndexOrdersForMemberURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}
