// Package v1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package v1

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// 创建通话
	// (POST /api/v1/live)
	CreateRoom(c *gin.Context)
	// 获取群聊当前通话房间信息
	// (GET /api/v1/live/group/{groupId})
	GetGroupRoom(c *gin.Context, groupId uint32)
	// 获取用户当前通话房间信息
	// (GET /api/v1/live/user)
	GetUserRoom(c *gin.Context)
	// 删除通话
	// (DELETE /api/v1/live/{id})
	DeleteRoom(c *gin.Context, id string)
	// 获取通话房间信息
	// (GET /api/v1/live/{id})
	GetRoom(c *gin.Context, id string)
	// 加入通话
	// (POST /api/v1/live/{id}/join)
	JoinRoom(c *gin.Context, id string)
	// 拒绝通话
	// (POST /api/v1/live/{id}/reject)
	RejectRoom(c *gin.Context, id string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// CreateRoom operation middleware
func (siw *ServerInterfaceWrapper) CreateRoom(c *gin.Context) {

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.CreateRoom(c)
}

// GetGroupRoom operation middleware
func (siw *ServerInterfaceWrapper) GetGroupRoom(c *gin.Context) {

	var err error

	// ------------- Path parameter "groupId" -------------
	var groupId uint32

	err = runtime.BindStyledParameterWithOptions("simple", "groupId", c.Param("groupId"), &groupId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter groupId: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetGroupRoom(c, groupId)
}

// GetUserRoom operation middleware
func (siw *ServerInterfaceWrapper) GetUserRoom(c *gin.Context) {

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetUserRoom(c)
}

// DeleteRoom operation middleware
func (siw *ServerInterfaceWrapper) DeleteRoom(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.DeleteRoom(c, id)
}

// GetRoom operation middleware
func (siw *ServerInterfaceWrapper) GetRoom(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetRoom(c, id)
}

// JoinRoom operation middleware
func (siw *ServerInterfaceWrapper) JoinRoom(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.JoinRoom(c, id)
}

// RejectRoom operation middleware
func (siw *ServerInterfaceWrapper) RejectRoom(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.RejectRoom(c, id)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.POST(options.BaseURL+"/api/v1/live", wrapper.CreateRoom)
	router.GET(options.BaseURL+"/api/v1/live/group/:groupId", wrapper.GetGroupRoom)
	router.GET(options.BaseURL+"/api/v1/live/user", wrapper.GetUserRoom)
	router.DELETE(options.BaseURL+"/api/v1/live/:id", wrapper.DeleteRoom)
	router.GET(options.BaseURL+"/api/v1/live/:id", wrapper.GetRoom)
	router.POST(options.BaseURL+"/api/v1/live/:id/join", wrapper.JoinRoom)
	router.POST(options.BaseURL+"/api/v1/live/:id/reject", wrapper.RejectRoom)
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9RY324Txxd+ldX8fhettMEBSlX5rgWBzA1VaMUFsqzx7nE8sLuzzMwGImulgEBNg1Kq",
	"UpqSkhaktMpNFahQIpEGXiZrO1d9hWpm1s7au/4TZ93Sq8Q7O2fP+b7vnDlnGsiirk898ARHxQbiVh1c",
	"rP49zwALmKPUnYNbAXAhH/qM+sAEAfXKPKOBXyG2/N8GbjHiC0I9VEStt5utvfulC8hENcpcLFARBcQT",
	"Z88gE4lFH1AREU/APDBkojsz83RGPp3hN4k/Q5UV7Mz4VL7DUFGwAEITkVraB22sgcALXFS8rn1C5TAM",
	"TeSCW5Xb+70rXTBozRB1MAIOzLhdp8Zt4jgGAwvIAqgVCzsOMhER4Ca+g7hgxJtH0ji+U9KLZ86d6waF",
	"GcOL44YkX6Ou/IQvFlGxhh0uw6Sxnw30fwY1VET/KxyxVIgpKkhmrug3QxOJOqgdDG4FhIHdhULSUw7N",
	"Lk69SBwurbe3N1qv9qKfHyKzC6KEBZkdLE0Ed7DrOzK+eKUXjXEZDM0e/5SRLkvlrlVavQGWkFElNch9",
	"6nFIC4BR6qYDay6/O1x7rQQ4ka8DqBHEBRqIQUC2dx4013aaazuHa68/iFafHOyv/vXnOv8wmQbEExPn",
	"wACvAuakPboN1TlhNZ+tRisvoqdb0bOX0cZSrnCEGYxdpsQbWjOOL+7hnxkkC0Fv6ozoRSVaeR49+DWW",
	"/fp9/Va+GnmP2PgcM0Es4mNPlLwaTaNEbPAEkUb6Hf5SVkaPWDc97ELOCBFe8YOqQ3g9qzpfq4OoA1NV",
	"2D8KwCBcF2ZZE2iiCFUpdQB7J/TpBiUe2BWckdtSbIbM/L40/vijvBNZYZ06bE5mkwssMiq/zB7jqsAi",
	"4H1hfZJvVFm6TKZtejEu6L1KtQOGO5VjuiQkNKdSpNMBDKtX/YkW5toQdI64nGXBRCz36eLZ6Tzy9H6B",
	"2EArskerMLAosytx1Z24McmU4ZXuYdUrRhzYhFbAw1UH7CEVjBr6HUNtyLtkWdQGK/318+pxvhW7xrAL",
	"FZZZSS7KNUOtTbfNYcCpE3QY6StnR2vmFJR2DKrVhnypTsnzhMY4WAEjYvGqrF5a0VXADNingagf/brY",
	"4fLytS+QqWdDFZJaPQqxLoSP1MBF4iZDEKGmBYtyTlzDkUMV9gky0QIwrpE7fWr21CxS4w54crGIzqpH",
	"sgKLuvKqgH1SWDhdkAZUFlKecT5Hyz9Fe290T4eUPX1UlGyZDd0RAunpA7j4jNqq47GoJ0BXeez7DrHU",
	"tsINrhWmi/uo0p+ekxUUvS62t3ebr+4d7D+OHt1rPnmJkoNQPNyy+ExUgZ+ZnZ2Kg/Gxm+FhEsTm8rfR",
	"yi8oKRVUvN4rkuvlsGwiHrguZotpEgSe53LGU8yVpaUklwU1WRYa6k/JDtVlAmRQ2/5mN3r0Q+vtZvvu",
	"SrT/OPp6teOhHO8O3r1o3t1OUX4JxCVpOCbdx7JACWBcBTHwooLI31J6qNOIodi/FF1mAvqR9xuhBGpq",
	"7KogsxQ3ADmN2SQMj8nFCObVBcIIur/fai7vHoduOax0U/zfATrl84mBHoXCCKAbROeVDQ5kHdva/uHS",
	"UvTVm+RXUvheUAbGSab2b3e1vdb6/aTLg9KLDM+sviu3QWnUf/EjAdduxAEds449P3y6OaiOmcN0O55W",
	"318g/4E0ySc1Jk2Hghzyh/QR6m4o2v0j+v3H6NlWH/YpLjt3UGOQqS1Pkcz825r+i7wxmpppNjGpC78M",
	"f7Sgkjd8x0z95M7xBMVA9eMDJdV8+F1rb6PPcK+M5pSJMYWk7f0Hq8JI2nRkk9DWszNNWxj+HQAA//8o",
	"JvoR5RoAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
