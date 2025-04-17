# http package - net/http - Go Packages
Package http provides HTTP client and server implementations.

[Get](#Get), [Head](#Head), [Post](#Post), and [PostForm](#PostForm) make HTTP (or HTTPS) requests:

```
resp, err := http.Get("http://example.com/")
...
resp, err := http.Post("http://example.com/upload", "image/jpeg", &buf)
...
resp, err := http.PostForm("http://example.com/form",
	url.Values{"key": {"Value"}, "id": {"123"}})

```


The caller must close the response body when finished with it:

```
resp, err := http.Get("http://example.com/")
if err != nil {
	// handle error
}
defer resp.Body.Close()
body, err := io.ReadAll(resp.Body)
// ...

```


#### Clients and Transports 

For control over HTTP client headers, redirect policy, and other settings, create a [Client](#Client):

```
client := &http.Client{
	CheckRedirect: redirectPolicyFunc,
}

resp, err := client.Get("http://example.com")
// ...

req, err := http.NewRequest("GET", "http://example.com", nil)
// ...
req.Header.Add("If-None-Match", `W/"wyzzy"`)
resp, err := client.Do(req)
// ...

```


For control over proxies, TLS configuration, keep-alives, compression, and other settings, create a [Transport](#Transport):

```
tr := &http.Transport{
	MaxIdleConns:       10,
	IdleConnTimeout:    30 * time.Second,
	DisableCompression: true,
}
client := &http.Client{Transport: tr}
resp, err := client.Get("https://example.com")

```


Clients and Transports are safe for concurrent use by multiple goroutines and for efficiency should only be created once and re-used.

#### Servers 

ListenAndServe starts an HTTP server with a given address and handler. The handler is usually nil, which means to use [DefaultServeMux](#DefaultServeMux). [Handle](#Handle) and [HandleFunc](#HandleFunc) add handlers to [DefaultServeMux](#DefaultServeMux):

```
http.Handle("/foo", fooHandler)

http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
})

log.Fatal(http.ListenAndServe(":8080", nil))

```


More control over the server's behavior is available by creating a custom Server:

```
s := &http.Server{
	Addr:           ":8080",
	Handler:        myHandler,
	ReadTimeout:    10 * time.Second,
	WriteTimeout:   10 * time.Second,
	MaxHeaderBytes: 1 << 20,
}
log.Fatal(s.ListenAndServe())

```


#### HTTP/2 

Starting with Go 1.6, the http package has transparent support for the HTTP/2 protocol when using HTTPS. Programs that must disable HTTP/2 can do so by setting \[Transport.TLSNextProto\] (for clients) or \[Server.TLSNextProto\] (for servers) to a non-nil, empty map. Alternatively, the following GODEBUG settings are currently supported:

```
GODEBUG=http2client=0  # disable HTTP/2 client support
GODEBUG=http2server=0  # disable HTTP/2 server support
GODEBUG=http2debug=1   # enable verbose HTTP/2 debug logs
GODEBUG=http2debug=2   # ... even more verbose, with frame dumps

```


Please report any issues before disabling HTTP/2 support: [https://golang.org/s/http2bug](https://golang.org/s/http2bug)

The http package's [Transport](#Transport) and [Server](#Server) both automatically enable HTTP/2 support for simple configurations. To enable HTTP/2 for more complex configurations, to use lower-level HTTP/2 features, or to use a newer version of Go's http2 package, import "golang.org/x/net/http2" directly and use its ConfigureTransport and/or ConfigureServer functions. Manually configuring HTTP/2 via the golang.org/x/net/http2 package takes precedence over the net/http package's built-in HTTP/2 support.

*   [Constants](#pkg-constants)
*   [Variables](#pkg-variables)
*   [func CanonicalHeaderKey(s string) string](#CanonicalHeaderKey)
*   [func DetectContentType(data \[\]byte) string](#DetectContentType)
*   [func Error(w ResponseWriter, error string, code int)](#Error)
*   [func Handle(pattern string, handler Handler)](#Handle)
*   [func HandleFunc(pattern string, handler func(ResponseWriter, \*Request))](#HandleFunc)
*   [func ListenAndServe(addr string, handler Handler) error](#ListenAndServe)
*   [func ListenAndServeTLS(addr, certFile, keyFile string, handler Handler) error](#ListenAndServeTLS)
*   [func MaxBytesReader(w ResponseWriter, r io.ReadCloser, n int64) io.ReadCloser](#MaxBytesReader)
*   [func NotFound(w ResponseWriter, r \*Request)](#NotFound)
*   [func ParseHTTPVersion(vers string) (major, minor int, ok bool)](#ParseHTTPVersion)
*   [func ParseTime(text string) (t time.Time, err error)](#ParseTime)
*   [func ProxyFromEnvironment(req \*Request) (\*url.URL, error)](#ProxyFromEnvironment)
*   [func ProxyURL(fixedURL \*url.URL) func(\*Request) (\*url.URL, error)](#ProxyURL)
*   [func Redirect(w ResponseWriter, r \*Request, url string, code int)](#Redirect)
*   [func Serve(l net.Listener, handler Handler) error](#Serve)
*   [func ServeContent(w ResponseWriter, req \*Request, name string, modtime time.Time, ...)](#ServeContent)
*   [func ServeFile(w ResponseWriter, r \*Request, name string)](#ServeFile)
*   [func ServeFileFS(w ResponseWriter, r \*Request, fsys fs.FS, name string)](#ServeFileFS)
*   [func ServeTLS(l net.Listener, handler Handler, certFile, keyFile string) error](#ServeTLS)
*   [func SetCookie(w ResponseWriter, cookie \*Cookie)](#SetCookie)
*   [func StatusText(code int) string](#StatusText)
*   [type Client](#Client)
*   *   [func (c \*Client) CloseIdleConnections()](#Client.CloseIdleConnections)
    *   [func (c \*Client) Do(req \*Request) (\*Response, error)](#Client.Do)
    *   [func (c \*Client) Get(url string) (resp \*Response, err error)](#Client.Get)
    *   [func (c \*Client) Head(url string) (resp \*Response, err error)](#Client.Head)
    *   [func (c \*Client) Post(url, contentType string, body io.Reader) (resp \*Response, err error)](#Client.Post)
    *   [func (c \*Client) PostForm(url string, data url.Values) (resp \*Response, err error)](#Client.PostForm)
*   [type CloseNotifier](#CloseNotifier)deprecated
*   [type ConnState](#ConnState)
*   *   [func (c ConnState) String() string](#ConnState.String)
*   [type Cookie](#Cookie)
*   *   [func ParseCookie(line string) (\[\]\*Cookie, error)](#ParseCookie)
    *   [func ParseSetCookie(line string) (\*Cookie, error)](#ParseSetCookie)
*   *   [func (c \*Cookie) String() string](#Cookie.String)
    *   [func (c \*Cookie) Valid() error](#Cookie.Valid)
*   [type CookieJar](#CookieJar)
*   [type Dir](#Dir)
*   *   [func (d Dir) Open(name string) (File, error)](#Dir.Open)
*   [type File](#File)
*   [type FileSystem](#FileSystem)
*   *   [func FS(fsys fs.FS) FileSystem](#FS)
*   [type Flusher](#Flusher)
*   [type HTTP2Config](#HTTP2Config)
*   [type Handler](#Handler)
*   *   [func AllowQuerySemicolons(h Handler) Handler](#AllowQuerySemicolons)
    *   [func FileServer(root FileSystem) Handler](#FileServer)
    *   [func FileServerFS(root fs.FS) Handler](#FileServerFS)
    *   [func MaxBytesHandler(h Handler, n int64) Handler](#MaxBytesHandler)
    *   [func NotFoundHandler() Handler](#NotFoundHandler)
    *   [func RedirectHandler(url string, code int) Handler](#RedirectHandler)
    *   [func StripPrefix(prefix string, h Handler) Handler](#StripPrefix)
    *   [func TimeoutHandler(h Handler, dt time.Duration, msg string) Handler](#TimeoutHandler)
*   [type HandlerFunc](#HandlerFunc)
*   *   [func (f HandlerFunc) ServeHTTP(w ResponseWriter, r \*Request)](#HandlerFunc.ServeHTTP)
*   [type Header](#Header)
*   *   [func (h Header) Add(key, value string)](#Header.Add)
    *   [func (h Header) Clone() Header](#Header.Clone)
    *   [func (h Header) Del(key string)](#Header.Del)
    *   [func (h Header) Get(key string) string](#Header.Get)
    *   [func (h Header) Set(key, value string)](#Header.Set)
    *   [func (h Header) Values(key string) \[\]string](#Header.Values)
    *   [func (h Header) Write(w io.Writer) error](#Header.Write)
    *   [func (h Header) WriteSubset(w io.Writer, exclude map\[string\]bool) error](#Header.WriteSubset)
*   [type Hijacker](#Hijacker)
*   [type MaxBytesError](#MaxBytesError)
*   *   [func (e \*MaxBytesError) Error() string](#MaxBytesError.Error)
*   [type ProtocolError](#ProtocolError)deprecated
*   *   [func (pe \*ProtocolError) Error() string](#ProtocolError.Error)
    *   [func (pe \*ProtocolError) Is(err error) bool](#ProtocolError.Is)
*   [type Protocols](#Protocols)
*   *   [func (p Protocols) HTTP1() bool](#Protocols.HTTP1)
    *   [func (p Protocols) HTTP2() bool](#Protocols.HTTP2)
    *   [func (p \*Protocols) SetHTTP1(ok bool)](#Protocols.SetHTTP1)
    *   [func (p \*Protocols) SetHTTP2(ok bool)](#Protocols.SetHTTP2)
    *   [func (p \*Protocols) SetUnencryptedHTTP2(ok bool)](#Protocols.SetUnencryptedHTTP2)
    *   [func (p Protocols) String() string](#Protocols.String)
    *   [func (p Protocols) UnencryptedHTTP2() bool](#Protocols.UnencryptedHTTP2)
*   [type PushOptions](#PushOptions)
*   [type Pusher](#Pusher)
*   [type Request](#Request)
*   *   [func NewRequest(method, url string, body io.Reader) (\*Request, error)](#NewRequest)
    *   [func NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (\*Request, error)](#NewRequestWithContext)
    *   [func ReadRequest(b \*bufio.Reader) (\*Request, error)](#ReadRequest)
*   *   [func (r \*Request) AddCookie(c \*Cookie)](#Request.AddCookie)
    *   [func (r \*Request) BasicAuth() (username, password string, ok bool)](#Request.BasicAuth)
    *   [func (r \*Request) Clone(ctx context.Context) \*Request](#Request.Clone)
    *   [func (r \*Request) Context() context.Context](#Request.Context)
    *   [func (r \*Request) Cookie(name string) (\*Cookie, error)](#Request.Cookie)
    *   [func (r \*Request) Cookies() \[\]\*Cookie](#Request.Cookies)
    *   [func (r \*Request) CookiesNamed(name string) \[\]\*Cookie](#Request.CookiesNamed)
    *   [func (r \*Request) FormFile(key string) (multipart.File, \*multipart.FileHeader, error)](#Request.FormFile)
    *   [func (r \*Request) FormValue(key string) string](#Request.FormValue)
    *   [func (r \*Request) MultipartReader() (\*multipart.Reader, error)](#Request.MultipartReader)
    *   [func (r \*Request) ParseForm() error](#Request.ParseForm)
    *   [func (r \*Request) ParseMultipartForm(maxMemory int64) error](#Request.ParseMultipartForm)
    *   [func (r \*Request) PathValue(name string) string](#Request.PathValue)
    *   [func (r \*Request) PostFormValue(key string) string](#Request.PostFormValue)
    *   [func (r \*Request) ProtoAtLeast(major, minor int) bool](#Request.ProtoAtLeast)
    *   [func (r \*Request) Referer() string](#Request.Referer)
    *   [func (r \*Request) SetBasicAuth(username, password string)](#Request.SetBasicAuth)
    *   [func (r \*Request) SetPathValue(name, value string)](#Request.SetPathValue)
    *   [func (r \*Request) UserAgent() string](#Request.UserAgent)
    *   [func (r \*Request) WithContext(ctx context.Context) \*Request](#Request.WithContext)
    *   [func (r \*Request) Write(w io.Writer) error](#Request.Write)
    *   [func (r \*Request) WriteProxy(w io.Writer) error](#Request.WriteProxy)
*   [type Response](#Response)
*   *   [func Get(url string) (resp \*Response, err error)](#Get)
    *   [func Head(url string) (resp \*Response, err error)](#Head)
    *   [func Post(url, contentType string, body io.Reader) (resp \*Response, err error)](#Post)
    *   [func PostForm(url string, data url.Values) (resp \*Response, err error)](#PostForm)
    *   [func ReadResponse(r \*bufio.Reader, req \*Request) (\*Response, error)](#ReadResponse)
*   *   [func (r \*Response) Cookies() \[\]\*Cookie](#Response.Cookies)
    *   [func (r \*Response) Location() (\*url.URL, error)](#Response.Location)
    *   [func (r \*Response) ProtoAtLeast(major, minor int) bool](#Response.ProtoAtLeast)
    *   [func (r \*Response) Write(w io.Writer) error](#Response.Write)
*   [type ResponseController](#ResponseController)
*   *   [func NewResponseController(rw ResponseWriter) \*ResponseController](#NewResponseController)
*   *   [func (c \*ResponseController) EnableFullDuplex() error](#ResponseController.EnableFullDuplex)
    *   [func (c \*ResponseController) Flush() error](#ResponseController.Flush)
    *   [func (c \*ResponseController) Hijack() (net.Conn, \*bufio.ReadWriter, error)](#ResponseController.Hijack)
    *   [func (c \*ResponseController) SetReadDeadline(deadline time.Time) error](#ResponseController.SetReadDeadline)
    *   [func (c \*ResponseController) SetWriteDeadline(deadline time.Time) error](#ResponseController.SetWriteDeadline)
*   [type ResponseWriter](#ResponseWriter)
*   [type RoundTripper](#RoundTripper)
*   *   [func NewFileTransport(fs FileSystem) RoundTripper](#NewFileTransport)
    *   [func NewFileTransportFS(fsys fs.FS) RoundTripper](#NewFileTransportFS)
*   [type SameSite](#SameSite)
*   [type ServeMux](#ServeMux)
*   *   [func NewServeMux() \*ServeMux](#NewServeMux)
*   *   [func (mux \*ServeMux) Handle(pattern string, handler Handler)](#ServeMux.Handle)
    *   [func (mux \*ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, \*Request))](#ServeMux.HandleFunc)
    *   [func (mux \*ServeMux) Handler(r \*Request) (h Handler, pattern string)](#ServeMux.Handler)
    *   [func (mux \*ServeMux) ServeHTTP(w ResponseWriter, r \*Request)](#ServeMux.ServeHTTP)
*   [type Server](#Server)
*   *   [func (s \*Server) Close() error](#Server.Close)
    *   [func (s \*Server) ListenAndServe() error](#Server.ListenAndServe)
    *   [func (s \*Server) ListenAndServeTLS(certFile, keyFile string) error](#Server.ListenAndServeTLS)
    *   [func (s \*Server) RegisterOnShutdown(f func())](#Server.RegisterOnShutdown)
    *   [func (s \*Server) Serve(l net.Listener) error](#Server.Serve)
    *   [func (s \*Server) ServeTLS(l net.Listener, certFile, keyFile string) error](#Server.ServeTLS)
    *   [func (s \*Server) SetKeepAlivesEnabled(v bool)](#Server.SetKeepAlivesEnabled)
    *   [func (s \*Server) Shutdown(ctx context.Context) error](#Server.Shutdown)
*   [type Transport](#Transport)
*   *   [func (t \*Transport) CancelRequest(req \*Request)](#Transport.CancelRequest)deprecated
    *   [func (t \*Transport) Clone() \*Transport](#Transport.Clone)
    *   [func (t \*Transport) CloseIdleConnections()](#Transport.CloseIdleConnections)
    *   [func (t \*Transport) RegisterProtocol(scheme string, rt RoundTripper)](#Transport.RegisterProtocol)
    *   [func (t \*Transport) RoundTrip(req \*Request) (\*Response, error)](#Transport.RoundTrip)

*   [FileServer](#example-FileServer)
*   [FileServer (DotFileHiding)](#example-FileServer-DotFileHiding)
*   [FileServer (StripPrefix)](#example-FileServer-StripPrefix)
*   [Get](#example-Get)
*   [Handle](#example-Handle)
*   [HandleFunc](#example-HandleFunc)
*   [Hijacker](#example-Hijacker)
*   [ListenAndServe](#example-ListenAndServe)
*   [ListenAndServeTLS](#example-ListenAndServeTLS)
*   [NotFoundHandler](#example-NotFoundHandler)
*   [Protocols (Http1)](#example-Protocols-Http1)
*   [Protocols (Http1or2)](#example-Protocols-Http1or2)
*   [ResponseWriter (Trailers)](#example-ResponseWriter-Trailers)
*   [ServeMux.Handle](#example-ServeMux.Handle)
*   [Server.Shutdown](#example-Server.Shutdown)
*   [StripPrefix](#example-StripPrefix)

[View Source](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/method.go;l=10)

```
const (
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH" 
	MethodDelete  = "DELETE"
	MethodConnect = "CONNECT"
	MethodOptions = "OPTIONS"
	MethodTrace   = "TRACE"
)
```


Common HTTP methods.

Unless otherwise noted, these are defined in [RFC 7231 section 4.3](https://rfc-editor.org/rfc/rfc7231.html#section-4.3).

[View Source](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/status.go;l=9)

```
const (
	StatusContinue           = 100 
	StatusSwitchingProtocols = 101 
	StatusProcessing         = 102 
	StatusEarlyHints         = 103 

	StatusOK                   = 200 
	StatusCreated              = 201 
	StatusAccepted             = 202 
	StatusNoContent            = 204 
	StatusResetContent         = 205 
	StatusPartialContent       = 206 
	StatusMultiStatus          = 207 
	StatusAlreadyReported      = 208 
	StatusIMUsed               = 226 

	StatusMultipleChoices  = 300 
	StatusMovedPermanently = 301 
	StatusFound            = 302 
	StatusSeeOther         = 303 
	StatusNotModified      = 304 
	StatusUseProxy         = 305 

	StatusTemporaryRedirect = 307 
	StatusPermanentRedirect = 308 

	StatusBadRequest                   = 400 
	StatusUnauthorized                 = 401 
	StatusPaymentRequired              = 402 
	StatusForbidden                    = 403 
	StatusNotFound                     = 404 
	StatusMethodNotAllowed             = 405 
	StatusNotAcceptable                = 406 
	StatusProxyAuthRequired            = 407 
	StatusRequestTimeout               = 408 
	StatusConflict                     = 409 
	StatusGone                         = 410 
	StatusLengthRequired               = 411 
	StatusPreconditionFailed           = 412 
	StatusRequestEntityTooLarge        = 413 
	StatusRequestURITooLong            = 414 
	StatusUnsupportedMediaType         = 415 
	StatusRequestedRangeNotSatisfiable = 416 
	StatusExpectationFailed            = 417 
	StatusTeapot                       = 418 
	StatusMisdirectedRequest           = 421 
	StatusUnprocessableEntity          = 422 
	StatusLocked                       = 423 
	StatusFailedDependency             = 424 
	StatusTooEarly                     = 425 
	StatusUpgradeRequired              = 426 
	StatusPreconditionRequired         = 428 
	StatusTooManyRequests              = 429 
	StatusUnavailableForLegalReasons   = 451 

	StatusInternalServerError           = 500 
	StatusNotImplemented                = 501 
	StatusBadGateway                    = 502 
	StatusServiceUnavailable            = 503 
	StatusGatewayTimeout                = 504 
	StatusHTTPVersionNotSupported       = 505 
	StatusVariantAlsoNegotiates         = 506 
	StatusInsufficientStorage           = 507 
	StatusLoopDetected                  = 508 
	StatusNotExtended                   = 510 
	StatusNetworkAuthenticationRequired = 511 
)
```


HTTP status codes as registered with IANA. See: [https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml](https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml)

DefaultMaxHeaderBytes is the maximum permitted size of the headers in an HTTP request. This can be overridden by setting \[Server.MaxHeaderBytes\].

DefaultMaxIdleConnsPerHost is the default value of [Transport](#Transport)'s MaxIdleConnsPerHost.

[View Source](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/server.go;l=992)

```
const TimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"
```


TimeFormat is the time format to use when generating times in HTTP headers. It is like [time.RFC1123](about:/time#RFC1123) but hard-codes GMT as the time zone. The time being formatted must be in UTC for Format to generate the correct format.

For parsing this time format, see [ParseTime](#ParseTime).

TrailerPrefix is a magic prefix for \[ResponseWriter.Header\] map keys that, if present, signals that the map entry is actually for the response trailers, and not the response headers. The prefix is stripped after the ServeHTTP call finishes and the values are sent in the trailers.

This mechanism is intended only for trailers that are not known prior to the headers being written. If the set of trailers is fixed or known before the header is written, the normal Go trailers mechanism is preferred:

```
https://pkg.go.dev/net/http#ResponseWriter
https://pkg.go.dev/net/http#example-ResponseWriter-Trailers

```


[View Source](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/request.go;l=58)

```
var (
	
	
	
	
	
	
	ErrNotSupported = &ProtocolError{"feature not supported"}

	
	
	
	ErrUnexpectedTrailer = &ProtocolError{"trailer header without chunked transfer encoding"}

	
	
	ErrMissingBoundary = &ProtocolError{"no multipart boundary param in Content-Type"}

	
	
	ErrNotMultipart = &ProtocolError{"request Content-Type isn't multipart/form-data"}

	
	
	ErrHeaderTooLong = &ProtocolError{"header too long"}

	
	
	
	ErrShortBody = &ProtocolError{"entity body too short"}

	
	
	
	ErrMissingContentLength = &ProtocolError{"missing ContentLength in HEAD response"}
)
```


[View Source](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/server.go;l=39)

```
var (
	
	
	
	ErrBodyNotAllowed = errors.New("http: request method or response status code does not allow body")

	
	
	
	
	
	ErrHijacked = errors.New("http: connection has been hijacked")

	
	
	
	
	ErrContentLength = errors.New("http: wrote more than the declared Content-Length")

	
	
	
	ErrWriteAfterFlush = errors.New("unused")
)
```


Errors used by the HTTP server.

[View Source](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/server.go;l=239)

```
var (
	
	
	
	
	ServerContextKey = &contextKey{"http-server"}

	
	
	
	
	LocalAddrContextKey = &contextKey{"local-addr"}
)
```


DefaultClient is the default [Client](#Client) and is used by [Get](#Get), [Head](#Head), and [Post](#Post).

DefaultServeMux is the default [ServeMux](#ServeMux) used by [Serve](#Serve).

ErrAbortHandler is a sentinel panic value to abort a handler. While any panic from ServeHTTP aborts the response to the client, panicking with ErrAbortHandler also suppresses logging of a stack trace to the server's error log.

ErrBodyReadAfterClose is returned when reading a [Request](#Request) or [Response](#Response) Body after the body has been closed. This typically happens when the body is read after an HTTP [Handler](#Handler) calls WriteHeader or Write on its [ResponseWriter](#ResponseWriter).

ErrHandlerTimeout is returned on [ResponseWriter](#ResponseWriter) Write calls in handlers which have timed out.

ErrLineTooLong is returned when reading request or response bodies with malformed chunked encoding.

ErrMissingFile is returned by FormFile when the provided file field name is either not present in the request or not a file field.

ErrNoCookie is returned by Request's Cookie method when a cookie is not found.

ErrNoLocation is returned by the [Response.Location](#Response.Location) method when no Location header is present.

ErrSchemeMismatch is returned when a server returns an HTTP response to an HTTPS client.

ErrServerClosed is returned by the [Server.Serve](#Server.Serve), [ServeTLS](#ServeTLS), [ListenAndServe](#ListenAndServe), and [ListenAndServeTLS](#ListenAndServeTLS) methods after a call to [Server.Shutdown](#Server.Shutdown) or [Server.Close](#Server.Close).

ErrSkipAltProtocol is a sentinel error value defined by Transport.RegisterProtocol.

ErrUseLastResponse can be returned by Client.CheckRedirect hooks to control how redirects are processed. If returned, the next request is not sent and the most recent response is returned with its body unclosed.

NoBody is an [io.ReadCloser](about:/io#ReadCloser) with no bytes. Read always returns EOF and Close always returns nil. It can be used in an outgoing client request to explicitly signal that a request has zero bytes. An alternative, however, is to simply set \[Request.Body\] to nil.

CanonicalHeaderKey returns the canonical format of the header key s. The canonicalization converts the first letter and any letter following a hyphen to upper case; the rest are converted to lowercase. For example, the canonical key for "accept-encoding" is "Accept-Encoding". If s contains a space or invalid header field bytes, it is returned without modifications.

#### func [DetectContentType](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/sniff.go;l=21) 

DetectContentType implements the algorithm described at [https://mimesniff.spec.whatwg.org/](https://mimesniff.spec.whatwg.org/) to determine the Content-Type of the given data. It considers at most the first 512 bytes of data. DetectContentType always returns a valid MIME type: if it cannot determine a more specific one, it returns "application/octet-stream".

Error replies to the request with the specified error message and HTTP code. It does not otherwise end the request; the caller should ensure no further writes are done to w. The error message should be plain text.

Error deletes the Content-Length header, sets Content-Type to “text/plain; charset=utf-8”, and sets X-Content-Type-Options to “nosniff”. This configures the header properly for the error message, in case the caller had set it up expecting a successful output.

#### func [Handle](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/server.go;l=2852) 

```
func Handle(pattern string, handler Handler)
```


Handle registers the handler for the given pattern in [DefaultServeMux](#DefaultServeMux). The documentation for [ServeMux](#ServeMux) explains how patterns are matched.

```
package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type countHandler struct {
	mu sync.Mutex // guards n
	n  int
}

func (h *countHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.n++
	fmt.Fprintf(w, "count is %d\n", h.n)
}

func main() {
	http.Handle("/count", new(countHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

```


```
Output:


```


#### func [HandleFunc](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/server.go;l=2862) 

```
func HandleFunc(pattern string, handler func(ResponseWriter, *Request))
```


HandleFunc registers the handler function for the given pattern in [DefaultServeMux](#DefaultServeMux). The documentation for [ServeMux](#ServeMux) explains how patterns are matched.

```
package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	h1 := func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "Hello from a HandleFunc #1!\n")
	}
	h2 := func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "Hello from a HandleFunc #2!\n")
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/endpoint", h2)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

```


```
Output:


```


#### func [ListenAndServe](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/server.go;l=3663) 

ListenAndServe listens on the TCP network address addr and then calls [Serve](#Serve) with handler to handle requests on incoming connections. Accepted connections are configured to enable TCP keep-alives.

The handler is typically nil, in which case [DefaultServeMux](#DefaultServeMux) is used.

ListenAndServe always returns a non-nil error.

```
package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	// Hello world, the web server

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	}

	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

```


```
Output:


```


#### func [ListenAndServeTLS](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/server.go;l=3673) 

```
func ListenAndServeTLS(addr, certFile, keyFile string, handler Handler) error
```


ListenAndServeTLS acts identically to [ListenAndServe](#ListenAndServe), except that it expects HTTPS connections. Additionally, files containing a certificate and matching private key for the server must be provided. If the certificate is signed by a certificate authority, the certFile should be the concatenation of the server's certificate, any intermediates, and the CA's certificate.

```
package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, TLS!\n")
	})

	// One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem.
	log.Printf("About to listen on 8443. Go to https://127.0.0.1:8443/")
	err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil)
	log.Fatal(err)
}

```


```
Output:


```


MaxBytesReader is similar to [io.LimitReader](about:/io#LimitReader) but is intended for limiting the size of incoming request bodies. In contrast to io.LimitReader, MaxBytesReader's result is a ReadCloser, returns a non-nil error of type [\*MaxBytesError](#MaxBytesError) for a Read beyond the limit, and closes the underlying reader when its Close method is called.

MaxBytesReader prevents clients from accidentally or maliciously sending a large request and wasting server resources. If possible, it tells the [ResponseWriter](#ResponseWriter) to close the connection after the limit has been reached.

```
func NotFound(w ResponseWriter, r *Request)
```


NotFound replies to the request with an HTTP 404 not found error.

ParseHTTPVersion parses an HTTP version string according to [RFC 7230, section 2.6](https://rfc-editor.org/rfc/rfc7230.html#section-2.6). "HTTP/1.0" returns (1, 0, true). Note that strings without a minor version, such as "HTTP/2", are not valid.

ParseTime parses a time header (such as the Date: header), trying each of the three formats allowed by HTTP/1.1: [TimeFormat](#TimeFormat), [time.RFC850](about:/time#RFC850), and [time.ANSIC](about:/time#ANSIC).

ProxyFromEnvironment returns the URL of the proxy to use for a given request, as indicated by the environment variables HTTP\_PROXY, HTTPS\_PROXY and NO\_PROXY (or the lowercase versions thereof). Requests use the proxy from the environment variable matching their scheme, unless excluded by NO\_PROXY.

The environment values may be either a complete URL or a "host\[:port\]", in which case the "http" scheme is assumed. An error is returned if the value is a different form.

A nil URL and nil error are returned if no proxy is defined in the environment, or a proxy should not be used for the given request, as defined by NO\_PROXY.

As a special case, if req.URL.Host is "localhost" (with or without a port number), then a nil URL and nil error will be returned.

ProxyURL returns a proxy function (for use in a [Transport](#Transport)) that always returns the same URL.

Redirect replies to the request with a redirect to url, which may be a path relative to the request path.

The provided code should be in the 3xx range and is usually [StatusMovedPermanently](#StatusMovedPermanently), [StatusFound](#StatusFound) or [StatusSeeOther](#StatusSeeOther).

If the Content-Type header has not been set, [Redirect](#Redirect) sets it to "text/html; charset=utf-8" and writes a small HTML body. Setting the Content-Type header to any value, including nil, disables that behavior.

Serve accepts incoming HTTP connections on the listener l, creating a new service goroutine for each. The service goroutines read requests and then call handler to reply to them.

The handler is typically nil, in which case [DefaultServeMux](#DefaultServeMux) is used.

HTTP/2 support is only enabled if the Listener returns [\*tls.Conn](about:/crypto/tls#Conn) connections and they were configured with "h2" in the TLS Config.NextProtos.

Serve always returns a non-nil error.

#### func [ServeContent](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/fs.go;l=240) 

ServeContent replies to the request using the content in the provided ReadSeeker. The main benefit of ServeContent over [io.Copy](about:/io#Copy) is that it handles Range requests properly, sets the MIME type, and handles If-Match, If-Unmodified-Since, If-None-Match, If-Modified-Since, and If-Range requests.

If the response's Content-Type header is not set, ServeContent first tries to deduce the type from name's file extension and, if that fails, falls back to reading the first block of the content and passing it to [DetectContentType](#DetectContentType). The name is otherwise unused; in particular it can be empty and is never sent in the response.

If modtime is not the zero time or Unix epoch, ServeContent includes it in a Last-Modified header in the response. If the request includes an If-Modified-Since header, ServeContent uses modtime to decide whether the content needs to be sent at all.

The content's Seek method must work: ServeContent uses a seek to the end of the content to determine its size. Note that [\*os.File](about:/os#File) implements the [io.ReadSeeker](about:/io#ReadSeeker) interface.

If the caller has set w's ETag header formatted per [RFC 7232, section 2.3](https://rfc-editor.org/rfc/rfc7232.html#section-2.3), ServeContent uses it to handle requests using If-Match, If-None-Match, or If-Range.

If an error occurs when serving the request (for example, when handling an invalid range request), ServeContent responds with an error message. By default, ServeContent strips the Cache-Control, Content-Encoding, ETag, and Last-Modified headers from error responses. The GODEBUG setting httpservecontentkeepheaders=1 causes ServeContent to preserve these headers.

ServeFile replies to the request with the contents of the named file or directory.

If the provided file or directory name is a relative path, it is interpreted relative to the current directory and may ascend to parent directories. If the provided name is constructed from user input, it should be sanitized before calling [ServeFile](#ServeFile).

As a precaution, ServeFile will reject requests where r.URL.Path contains a ".." path element; this protects against callers who might unsafely use [filepath.Join](about:/path/filepath#Join) on r.URL.Path without sanitizing it and then use that filepath.Join result as the name argument.

As another special case, ServeFile redirects any request where r.URL.Path ends in "/index.html" to the same path, without the final "index.html". To avoid such redirects either modify the path or use [ServeContent](#ServeContent).

Outside of those two special cases, ServeFile does not use r.URL.Path for selecting the file or directory to serve; only the file or directory provided in the name argument is used.

ServeFileFS replies to the request with the contents of the named file or directory from the file system fsys. The files provided by fsys must implement [io.Seeker](about:/io#Seeker).

If the provided name is constructed from user input, it should be sanitized before calling [ServeFileFS](#ServeFileFS).

As a precaution, ServeFileFS will reject requests where r.URL.Path contains a ".." path element; this protects against callers who might unsafely use [filepath.Join](about:/path/filepath#Join) on r.URL.Path without sanitizing it and then use that filepath.Join result as the name argument.

As another special case, ServeFileFS redirects any request where r.URL.Path ends in "/index.html" to the same path, without the final "index.html". To avoid such redirects either modify the path or use [ServeContent](#ServeContent).

Outside of those two special cases, ServeFileFS does not use r.URL.Path for selecting the file or directory to serve; only the file or directory provided in the name argument is used.

ServeTLS accepts incoming HTTPS connections on the listener l, creating a new service goroutine for each. The service goroutines read requests and then call handler to reply to them.

The handler is typically nil, in which case [DefaultServeMux](#DefaultServeMux) is used.

Additionally, files containing a certificate and matching private key for the server must be provided. If the certificate is signed by a certificate authority, the certFile should be the concatenation of the server's certificate, any intermediates, and the CA's certificate.

ServeTLS always returns a non-nil error.

```
func SetCookie(w ResponseWriter, cookie *Cookie)
```


SetCookie adds a Set-Cookie header to the provided [ResponseWriter](#ResponseWriter)'s headers. The provided cookie must have a valid Name. Invalid cookies may be silently dropped.

StatusText returns a text for the HTTP status code. It returns the empty string if the code is unknown.

A Client is an HTTP client. Its zero value ([DefaultClient](#DefaultClient)) is a usable client that uses [DefaultTransport](#DefaultTransport).

The \[Client.Transport\] typically has internal state (cached TCP connections), so Clients should be reused instead of created as needed. Clients are safe for concurrent use by multiple goroutines.

A Client is higher-level than a [RoundTripper](#RoundTripper) (such as [Transport](#Transport)) and additionally handles HTTP details such as cookies and redirects.

When following redirects, the Client will forward all headers set on the initial [Request](#Request) except:

*   when forwarding sensitive headers like "Authorization", "WWW-Authenticate", and "Cookie" to untrusted targets. These headers will be ignored when following a redirect to a domain that is not a subdomain match or exact match of the initial domain. For example, a redirect from "foo.com" to either "foo.com" or "sub.foo.com" will forward the sensitive headers, but a redirect to "bar.com" will not.
*   when forwarding the "Cookie" header with a non-nil cookie Jar. Since each redirect may mutate the state of the cookie jar, a redirect may possibly alter a cookie set in the initial request. When forwarding the "Cookie" header, any mutated cookies will be omitted, with the expectation that the Jar will insert those mutated cookies with the updated values (assuming the origin matches). If Jar is nil, the initial cookies are forwarded without change.

```
func (c *Client) CloseIdleConnections()
```


CloseIdleConnections closes any connections on its [Transport](#Transport) which were previously connected from previous requests but are now sitting idle in a "keep-alive" state. It does not interrupt any connections currently in use.

If \[Client.Transport\] does not have a [Client.CloseIdleConnections](#Client.CloseIdleConnections) method then this method does nothing.

Do sends an HTTP request and returns an HTTP response, following policy (such as redirects, cookies, auth) as configured on the client.

An error is returned if caused by client policy (such as CheckRedirect), or failure to speak HTTP (such as a network connectivity problem). A non-2xx status code doesn't cause an error.

If the returned error is nil, the [Response](#Response) will contain a non-nil Body which the user is expected to close. If the Body is not both read to EOF and closed, the [Client](#Client)'s underlying [RoundTripper](#RoundTripper) (typically [Transport](#Transport)) may not be able to re-use a persistent TCP connection to the server for a subsequent "keep-alive" request.

The request Body, if non-nil, will be closed by the underlying Transport, even on errors. The Body may be closed asynchronously after Do returns.

On error, any Response can be ignored. A non-nil Response with a non-nil error only occurs when CheckRedirect fails, and even then the returned \[Response.Body\] is already closed.

Generally [Get](#Get), [Post](#Post), or [PostForm](#PostForm) will be used instead of Do.

If the server replies with a redirect, the Client first uses the CheckRedirect function to determine whether the redirect should be followed. If permitted, a 301, 302, or 303 redirect causes subsequent requests to use HTTP method GET (or HEAD if the original request was HEAD), with no body. A 307 or 308 redirect preserves the original HTTP method and body, provided that the \[Request.GetBody\] function is defined. The [NewRequest](#NewRequest) function automatically sets GetBody for common standard library body types.

Any returned error will be of type [\*url.Error](about:/net/url#Error). The url.Error value's Timeout method will report true if the request timed out.

Get issues a GET to the specified URL. If the response is one of the following redirect codes, Get follows the redirect after calling the \[Client.CheckRedirect\] function:

```
301 (Moved Permanently)
302 (Found)
303 (See Other)
307 (Temporary Redirect)
308 (Permanent Redirect)

```


An error is returned if the \[Client.CheckRedirect\] function fails or if there was an HTTP protocol error. A non-2xx response doesn't cause an error. Any returned error will be of type [\*url.Error](about:/net/url#Error). The url.Error value's Timeout method will report true if the request timed out.

When err is nil, resp always contains a non-nil resp.Body. Caller should close resp.Body when done reading from it.

To make a request with custom headers, use [NewRequest](#NewRequest) and [Client.Do](#Client.Do).

To make a request with a specified context.Context, use [NewRequestWithContext](#NewRequestWithContext) and Client.Do.

Head issues a HEAD to the specified URL. If the response is one of the following redirect codes, Head follows the redirect after calling the \[Client.CheckRedirect\] function:

```
301 (Moved Permanently)
302 (Found)
303 (See Other)
307 (Temporary Redirect)
308 (Permanent Redirect)

```


To make a request with a specified [context.Context](about:/context#Context), use [NewRequestWithContext](#NewRequestWithContext) and [Client.Do](#Client.Do).

Post issues a POST to the specified URL.

Caller should close resp.Body when done reading from it.

If the provided body is an [io.Closer](about:/io#Closer), it is closed after the request.

To set custom headers, use [NewRequest](#NewRequest) and [Client.Do](#Client.Do).

To make a request with a specified context.Context, use [NewRequestWithContext](#NewRequestWithContext) and [Client.Do](#Client.Do).

See the Client.Do method documentation for details on how redirects are handled.

PostForm issues a POST to the specified URL, with data's keys and values URL-encoded as the request body.

The Content-Type header is set to application/x-www-form-urlencoded. To set other headers, use [NewRequest](#NewRequest) and [Client.Do](#Client.Do).

When err is nil, resp always contains a non-nil resp.Body. Caller should close resp.Body when done reading from it.

See the Client.Do method documentation for details on how redirects are handled.

To make a request with a specified context.Context, use [NewRequestWithContext](#NewRequestWithContext) and Client.Do.

```
type CloseNotifier interface {
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	CloseNotify() <-chan bool
}
```


The CloseNotifier interface is implemented by ResponseWriters which allow detecting when the underlying connection has gone away.

This mechanism can be used to cancel long operations on the server if the client has disconnected before the response is ready.

Deprecated: the CloseNotifier interface predates Go's context package. New code should use [Request.Context](#Request.Context) instead.

A ConnState represents the state of a client connection to a server. It's used by the optional \[Server.ConnState\] hook.

```
const (
	
	
	
	
	StateNew ConnState = iota

	
	
	
	
	
	
	
	
	
	
	
	StateActive

	
	
	
	
	StateIdle

	
	
	StateHijacked

	
	
	
	StateClosed
)
```


A Cookie represents an HTTP cookie as sent in the Set-Cookie header of an HTTP response or the Cookie header of an HTTP request.

See [https://tools.ietf.org/html/rfc6265](https://tools.ietf.org/html/rfc6265) for details.

ParseCookie parses a Cookie header value and returns all the cookies which were set in it. Since the same cookie name can appear multiple times the returned Values can contain more than one value for a given key.

ParseSetCookie parses a Set-Cookie header value and returns a cookie. It returns an error on syntax error.

String returns the serialization of the cookie for use in a [Cookie](#Cookie) header (if only Name and Value are set) or a Set-Cookie response header (if other fields are set). If c is nil or c.Name is invalid, the empty string is returned.

Valid reports whether the cookie is valid.

```
type CookieJar interface {
	
	
	
	SetCookies(u *url.URL, cookies []*Cookie)

	
	
	
	Cookies(u *url.URL) []*Cookie
}
```


A CookieJar manages storage and use of cookies in HTTP requests.

Implementations of CookieJar must be safe for concurrent use by multiple goroutines.

The net/http/cookiejar package provides a CookieJar implementation.

A Dir implements [FileSystem](#FileSystem) using the native file system restricted to a specific directory tree.

While the \[FileSystem.Open\] method takes '/'-separated paths, a Dir's string value is a directory path on the native file system, not a URL, so it is separated by [filepath.Separator](about:/path/filepath#Separator), which isn't necessarily '/'.

Note that Dir could expose sensitive files and directories. Dir will follow symlinks pointing out of the directory tree, which can be especially dangerous if serving from a directory in which users are able to create arbitrary symlinks. Dir will also allow access to files and directories starting with a period, which could expose sensitive directories like .git or sensitive files like .htpasswd. To exclude files with a leading period, remove the files/directories from the server or create a custom FileSystem implementation.

An empty Dir is treated as ".".

Open implements [FileSystem](#FileSystem) using [os.Open](about:/os#Open), opening files for reading rooted and relative to the directory d.

A File is returned by a [FileSystem](#FileSystem)'s Open method and can be served by the [FileServer](#FileServer) implementation.

The methods should behave the same as those on an [\*os.File](about:/os#File).

```
type FileSystem interface {
	Open(name string) (File, error)
}
```


A FileSystem implements access to a collection of named files. The elements in a file path are separated by slash ('/', U+002F) characters, regardless of host operating system convention. See the [FileServer](#FileServer) function to convert a FileSystem to a [Handler](#Handler).

This interface predates the [fs.FS](about:/io/fs#FS) interface, which can be used instead: the [FS](#FS) adapter function converts an fs.FS to a FileSystem.

FS converts fsys to a [FileSystem](#FileSystem) implementation, for use with [FileServer](#FileServer) and [NewFileTransport](#NewFileTransport). The files provided by fsys must implement [io.Seeker](about:/io#Seeker).

```
type Flusher interface {
	
	Flush()
}
```


The Flusher interface is implemented by ResponseWriters that allow an HTTP handler to flush buffered data to the client.

The default HTTP/1.x and HTTP/2 [ResponseWriter](#ResponseWriter) implementations support [Flusher](#Flusher), but ResponseWriter wrappers may not. Handlers should always test for this ability at runtime.

Note that even for ResponseWriters that support Flush, if the client is connected through an HTTP proxy, the buffered data may not reach the client until the response completes.

```
type HTTP2Config struct {
	
	
	
	MaxConcurrentStreams int

	
	
	
	
	MaxDecoderHeaderTableSize int

	
	
	
	MaxEncoderHeaderTableSize int

	
	
	
	
	MaxReadFrameSize int

	
	
	
	
	MaxReceiveBufferPerConnection int

	
	
	
	
	MaxReceiveBufferPerStream int

	
	
	
	SendPingTimeout time.Duration

	
	
	
	PingTimeout time.Duration

	
	
	
	WriteByteTimeout time.Duration

	
	
	PermitProhibitedCipherSuites bool

	
	
	
	
	CountError func(errType string)
}
```


HTTP2Config defines HTTP/2 configuration parameters common to both [Transport](#Transport) and [Server](#Server).

#### type [Handler](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/server.go;l=88) 

```
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```


A Handler responds to an HTTP request.

\[Handler.ServeHTTP\] should write reply headers and data to the [ResponseWriter](#ResponseWriter) and then return. Returning signals that the request is finished; it is not valid to use the [ResponseWriter](#ResponseWriter) or read from the \[Request.Body\] after or concurrently with the completion of the ServeHTTP call.

Depending on the HTTP client software, HTTP protocol version, and any intermediaries between the client and the Go server, it may not be possible to read from the \[Request.Body\] after writing to the [ResponseWriter](#ResponseWriter). Cautious handlers should read the \[Request.Body\] first, and then reply.

Except for reading the body, handlers should not modify the provided Request.

If ServeHTTP panics, the server (the caller of ServeHTTP) assumes that the effect of the panic was isolated to the active request. It recovers the panic, logs a stack trace to the server error log, and either closes the network connection or sends an HTTP/2 RST\_STREAM, depending on the HTTP protocol. To abort a handler so the client sees an interrupted response but the server doesn't log an error, panic with the value [ErrAbortHandler](#ErrAbortHandler).

```
func AllowQuerySemicolons(h Handler) Handler
```


AllowQuerySemicolons returns a handler that serves requests by converting any unescaped semicolons in the URL query to ampersands, and invoking the handler h.

This restores the pre-Go 1.17 behavior of splitting query parameters on both semicolons and ampersands. (See golang.org/issue/25192). Note that this behavior doesn't match that of many proxies, and the mismatch can lead to security issues.

AllowQuerySemicolons should be invoked before [Request.ParseForm](#Request.ParseForm) is called.

```
func FileServer(root FileSystem) Handler
```


FileServer returns a handler that serves HTTP requests with the contents of the file system rooted at root.

As a special case, the returned file server redirects any request ending in "/index.html" to the same path, without the final "index.html".

To use the operating system's file system implementation, use [http.Dir](#Dir):

```
http.Handle("/", http.FileServer(http.Dir("/tmp")))

```


To use an [fs.FS](about:/io/fs#FS) implementation, use [http.FileServerFS](#FileServerFS) instead.

```
package main

import (
	"log"
	"net/http"
)

func main() {
	// Simple static webserver:
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("/usr/share/doc"))))
}

```


```
Output:


```


```
package main

import (
	"io"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

// containsDotFile reports whether name contains a path element starting with a period.
// The name is assumed to be a delimited by forward slashes, as guaranteed
// by the http.FileSystem interface.
func containsDotFile(name string) bool {
	parts := strings.Split(name, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}
	return false
}

// dotFileHidingFile is the http.File use in dotFileHidingFileSystem.
// It is used to wrap the Readdir method of http.File so that we can
// remove files and directories that start with a period from its output.
type dotFileHidingFile struct {
	http.File
}

// Readdir is a wrapper around the Readdir method of the embedded File
// that filters out all files that start with a period in their name.
func (f dotFileHidingFile) Readdir(n int) (fis []fs.FileInfo, err error) {
	files, err := f.File.Readdir(n)
	for _, file := range files { // Filters out the dot files
		if !strings.HasPrefix(file.Name(), ".") {
			fis = append(fis, file)
		}
	}
	if err == nil && n > 0 && len(fis) == 0 {
		err = io.EOF
	}
	return
}

// dotFileHidingFileSystem is an http.FileSystem that hides
// hidden "dot files" from being served.
type dotFileHidingFileSystem struct {
	http.FileSystem
}

// Open is a wrapper around the Open method of the embedded FileSystem
// that serves a 403 permission error when name has a file or directory
// with whose name starts with a period in its path.
func (fsys dotFileHidingFileSystem) Open(name string) (http.File, error) {
	if containsDotFile(name) { // If dot file, return 403 response
		return nil, fs.ErrPermission
	}

	file, err := fsys.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	return dotFileHidingFile{file}, err
}

func main() {
	fsys := dotFileHidingFileSystem{http.Dir(".")}
	http.Handle("/", http.FileServer(fsys))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

```


```
Output:


```


```
package main

import (
	"net/http"
)

func main() {
	// To serve a directory on disk (/tmp) under an alternate URL
	// path (/tmpfiles/), use StripPrefix to modify the request
	// URL's path before the FileServer sees it:
	http.Handle("/tmpfiles/", http.StripPrefix("/tmpfiles/", http.FileServer(http.Dir("/tmp"))))
}

```


```
Output:


```


```
func FileServerFS(root fs.FS) Handler
```


FileServerFS returns a handler that serves HTTP requests with the contents of the file system fsys. The files provided by fsys must implement [io.Seeker](about:/io#Seeker).

As a special case, the returned file server redirects any request ending in "/index.html" to the same path, without the final "index.html".

```
http.Handle("/", http.FileServerFS(fsys))

```


#### func [MaxBytesHandler](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/server.go;l=4057)  added in go1.18

```
func MaxBytesHandler(h Handler, n int64) Handler
```


MaxBytesHandler returns a [Handler](#Handler) that runs h with its [ResponseWriter](#ResponseWriter) and \[Request.Body\] wrapped by a MaxBytesReader.

#### func [NotFoundHandler](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/server.go;l=2334) 

```
func NotFoundHandler() Handler
```


NotFoundHandler returns a simple request handler that replies to each request with a “404 page not found” reply.

```
package main

import (
	"fmt"
	"log"
	"net/http"
)

func newPeopleHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "This is the people handler.")
	})
}

func main() {
	mux := http.NewServeMux()

	// Create sample handler to returns 404
	mux.Handle("/resources", http.NotFoundHandler())

	// Create sample handler that returns 200
	mux.Handle("/resources/people/", newPeopleHandler())

	log.Fatal(http.ListenAndServe(":8080", mux))
}

```


```
Output:


```


#### func [RedirectHandler](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/server.go;l=2458) 

RedirectHandler returns a request handler that redirects each request it receives to the given url using the given status code.

The provided code should be in the 3xx range and is usually [StatusMovedPermanently](#StatusMovedPermanently), [StatusFound](#StatusFound) or [StatusSeeOther](#StatusSeeOther).

StripPrefix returns a handler that serves HTTP requests by removing the given prefix from the request URL's Path (and RawPath if set) and invoking the handler h. StripPrefix handles a request for a path that doesn't begin with prefix by replying with an HTTP 404 not found error. The prefix must match exactly: if the prefix in the request contains escaped characters the reply is also an HTTP 404 not found error.

```
package main

import (
	"net/http"
)

func main() {
	// To serve a directory on disk (/tmp) under an alternate URL
	// path (/tmpfiles/), use StripPrefix to modify the request
	// URL's path before the FileServer sees it:
	http.Handle("/tmpfiles/", http.StripPrefix("/tmpfiles/", http.FileServer(http.Dir("/tmp"))))
}

```


```
Output:


```


#### func [TimeoutHandler](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/server.go;l=3777) 

TimeoutHandler returns a [Handler](#Handler) that runs h with the given time limit.

The new Handler calls h.ServeHTTP to handle each request, but if a call runs for longer than its time limit, the handler responds with a 503 Service Unavailable error and the given message in its body. (If msg is empty, a suitable default message will be sent.) After such a timeout, writes by h to its [ResponseWriter](#ResponseWriter) will return [ErrHandlerTimeout](#ErrHandlerTimeout).

TimeoutHandler supports the [Pusher](#Pusher) interface but does not support the [Hijacker](#Hijacker) or [Flusher](#Flusher) interfaces.

#### type [HandlerFunc](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/server.go;l=2290) 

```
type HandlerFunc func(ResponseWriter, *Request)
```


The HandlerFunc type is an adapter to allow the use of ordinary functions as HTTP handlers. If f is a function with the appropriate signature, HandlerFunc(f) is a [Handler](#Handler) that calls f.

#### func (HandlerFunc) [ServeHTTP](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/server.go;l=2293) 

```
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request)
```


ServeHTTP calls f(w, r).

A Header represents the key-value pairs in an HTTP header.

The keys should be in canonical form, as returned by [CanonicalHeaderKey](#CanonicalHeaderKey).

Add adds the key, value pair to the header. It appends to any existing values associated with key. The key is case insensitive; it is canonicalized by [CanonicalHeaderKey](#CanonicalHeaderKey).

```
func (h Header) Clone() Header
```


Clone returns a copy of h or nil if h is nil.

Del deletes the values associated with key. The key is case insensitive; it is canonicalized by [CanonicalHeaderKey](#CanonicalHeaderKey).

Get gets the first value associated with the given key. If there are no values associated with the key, Get returns "". It is case insensitive; [textproto.CanonicalMIMEHeaderKey](about:/net/textproto#CanonicalMIMEHeaderKey) is used to canonicalize the provided key. Get assumes that all keys are stored in canonical form. To use non-canonical keys, access the map directly.

Set sets the header entries associated with key to the single element value. It replaces any existing values associated with key. The key is case insensitive; it is canonicalized by [textproto.CanonicalMIMEHeaderKey](about:/net/textproto#CanonicalMIMEHeaderKey). To use non-canonical keys, assign to the map directly.

Values returns all values associated with the given key. It is case insensitive; [textproto.CanonicalMIMEHeaderKey](about:/net/textproto#CanonicalMIMEHeaderKey) is used to canonicalize the provided key. To use non-canonical keys, access the map directly. The returned slice is not a copy.

Write writes a header in wire format.

WriteSubset writes a header in wire format. If exclude is not nil, keys where exclude\[key\] == true are not written. Keys are not canonicalized before checking the exclude map.

The Hijacker interface is implemented by ResponseWriters that allow an HTTP handler to take over the connection.

The default [ResponseWriter](#ResponseWriter) for HTTP/1.x connections supports Hijacker, but HTTP/2 connections intentionally do not. ResponseWriter wrappers may also not support Hijacker. Handlers should always test for this ability at runtime.

```
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/hijack", func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
			return
		}
		conn, bufrw, err := hj.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Don't forget to close the connection:
		defer conn.Close()
		bufrw.WriteString("Now we're speaking raw TCP. Say hi: ")
		bufrw.Flush()
		s, err := bufrw.ReadString('\n')
		if err != nil {
			log.Printf("error reading string: %v", err)
			return
		}
		fmt.Fprintf(bufrw, "You said: %q\nBye.\n", s)
		bufrw.Flush()
	})
}

```


```
Output:


```


```
type MaxBytesError struct {
	Limit int64
}
```


MaxBytesError is returned by [MaxBytesReader](#MaxBytesReader) when its read limit is exceeded.

```
type ProtocolError struct {
	ErrorString string
}
```


ProtocolError represents an HTTP protocol error.

Deprecated: Not all errors in the http package related to protocol errors are of type ProtocolError.

Is lets http.ErrNotSupported match errors.ErrUnsupported.

```
type Protocols struct {
	
}
```


Protocols is a set of HTTP protocols. The zero value is an empty set of protocols.

The supported protocols are:

*   HTTP1 is the HTTP/1.0 and HTTP/1.1 protocols. HTTP1 is supported on both unsecured TCP and secured TLS connections.
    
*   HTTP2 is the HTTP/2 protcol over a TLS connection.
    
*   UnencryptedHTTP2 is the HTTP/2 protocol over an unsecured TCP connection.
    

```
package main

import (
	"log"
	"net/http"
)

func main() {
	srv := http.Server{
		Addr: ":8443",
	}

	// Serve only HTTP/1.
	srv.Protocols = new(http.Protocols)
	srv.Protocols.SetHTTP1(true)

	log.Fatal(srv.ListenAndServeTLS("cert.pem", "key.pem"))
}

```


```
Output:


```


```
package main

import (
	"log"
	"net/http"
)

func main() {
	t := http.DefaultTransport.(*http.Transport).Clone()

	// Use either HTTP/1 and HTTP/2.
	t.Protocols = new(http.Protocols)
	t.Protocols.SetHTTP1(true)
	t.Protocols.SetHTTP2(true)

	cli := &http.Client{Transport: t}
	res, err := cli.Get("http://www.google.com/robots.txt")
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
}

```


```
Output:


```


HTTP1 reports whether p includes HTTP/1.

HTTP2 reports whether p includes HTTP/2.

```
func (p *Protocols) SetHTTP1(ok bool)
```


SetHTTP1 adds or removes HTTP/1 from p.

```
func (p *Protocols) SetHTTP2(ok bool)
```


SetHTTP2 adds or removes HTTP/2 from p.

```
func (p *Protocols) SetUnencryptedHTTP2(ok bool)
```


SetUnencryptedHTTP2 adds or removes unencrypted HTTP/2 from p.

```
func (p Protocols) UnencryptedHTTP2() bool
```


UnencryptedHTTP2 reports whether p includes unencrypted HTTP/2.

```
type PushOptions struct {
	
	
	Method string

	
	
	Header Header
}
```


PushOptions describes options for \[Pusher.Push\].

A Request represents an HTTP request received by a server or to be sent by a client.

The field semantics differ slightly between client and server usage. In addition to the notes on the fields below, see the documentation for [Request.Write](#Request.Write) and [RoundTripper](#RoundTripper).

NewRequestWithContext returns a new [Request](#Request) given a method, URL, and optional body.

If the provided body is also an [io.Closer](about:/io#Closer), the returned \[Request.Body\] is set to body and will be closed (possibly asynchronously) by the Client methods Do, Post, and PostForm, and [Transport.RoundTrip](#Transport.RoundTrip).

NewRequestWithContext returns a Request suitable for use with [Client.Do](#Client.Do) or [Transport.RoundTrip](#Transport.RoundTrip). To create a request for use with testing a Server Handler, either use the [net/http/httptest.NewRequest](about:/net/http/httptest#NewRequest) function, use [ReadRequest](#ReadRequest), or manually update the Request fields. For an outgoing client request, the context controls the entire lifetime of a request and its response: obtaining a connection, sending the request, and reading the response headers and body. See the Request type's documentation for the difference between inbound and outbound request fields.

If body is of type [\*bytes.Buffer](about:/bytes#Buffer), [\*bytes.Reader](about:/bytes#Reader), or [\*strings.Reader](about:/strings#Reader), the returned request's ContentLength is set to its exact value (instead of -1), GetBody is populated (so 307 and 308 redirects can replay the body), and Body is set to [NoBody](#NoBody) if the ContentLength is 0.

ReadRequest reads and parses an incoming request from b.

ReadRequest is a low-level function and should only be used for specialized applications; most code should use the [Server](#Server) to read requests and handle them via the [Handler](#Handler) interface. ReadRequest only supports HTTP/1.x requests. For HTTP/2, use golang.org/x/net/http2.

```
func (r *Request) AddCookie(c *Cookie)
```


AddCookie adds a cookie to the request. Per [RFC 6265 section 5.4](https://rfc-editor.org/rfc/rfc6265.html#section-5.4), AddCookie does not attach more than one [Cookie](#Cookie) header field. That means all cookies, if any, are written into the same line, separated by semicolon. AddCookie only sanitizes c's name and value, and does not sanitize a Cookie header already present in the request.

```
func (r *Request) BasicAuth() (username, password string, ok bool)
```


BasicAuth returns the username and password provided in the request's Authorization header, if the request uses HTTP Basic Authentication. See [RFC 2617, Section 2](https://rfc-editor.org/rfc/rfc2617.html#section-2).

Clone returns a deep copy of r with its context changed to ctx. The provided ctx must be non-nil.

Clone only makes a shallow copy of the Body field.

For an outgoing client request, the context controls the entire lifetime of a request and its response: obtaining a connection, sending the request, and reading the response headers and body.

Context returns the request's context. To change the context, use [Request.Clone](#Request.Clone) or [Request.WithContext](#Request.WithContext).

The returned context is always non-nil; it defaults to the background context.

For outgoing client requests, the context controls cancellation.

For incoming server requests, the context is canceled when the client's connection closes, the request is canceled (with HTTP/2), or when the ServeHTTP method returns.

Cookie returns the named cookie provided in the request or [ErrNoCookie](#ErrNoCookie) if not found. If multiple cookies match the given name, only one cookie will be returned.

```
func (r *Request) Cookies() []*Cookie
```


Cookies parses and returns the HTTP cookies sent with the request.

```
func (r *Request) CookiesNamed(name string) []*Cookie
```


CookiesNamed parses and returns the named HTTP cookies sent with the request or an empty slice if none matched.

FormFile returns the first file for the provided form key. FormFile calls [Request.ParseMultipartForm](#Request.ParseMultipartForm) and [Request.ParseForm](#Request.ParseForm) if necessary.

FormValue returns the first value for the named component of the query. The precedence order:

1.  application/x-www-form-urlencoded form body (POST, PUT, PATCH only)
2.  query parameters (always)
3.  multipart/form-data form body (always)

FormValue calls [Request.ParseMultipartForm](#Request.ParseMultipartForm) and [Request.ParseForm](#Request.ParseForm) if necessary and ignores any errors returned by these functions. If key is not present, FormValue returns the empty string. To access multiple values of the same key, call ParseForm and then inspect \[Request.Form\] directly.

MultipartReader returns a MIME multipart reader if this is a multipart/form-data or a multipart/mixed POST request, else returns nil and an error. Use this function instead of [Request.ParseMultipartForm](#Request.ParseMultipartForm) to process the request body as a stream.

ParseForm populates r.Form and r.PostForm.

For all requests, ParseForm parses the raw query from the URL and updates r.Form.

For POST, PUT, and PATCH requests, it also reads the request body, parses it as a form and puts the results into both r.PostForm and r.Form. Request body parameters take precedence over URL query string values in r.Form.

If the request Body's size has not already been limited by [MaxBytesReader](#MaxBytesReader), the size is capped at 10MB.

For other HTTP methods, or when the Content-Type is not application/x-www-form-urlencoded, the request Body is not read, and r.PostForm is initialized to a non-nil, empty value.

[Request.ParseMultipartForm](#Request.ParseMultipartForm) calls ParseForm automatically. ParseForm is idempotent.

ParseMultipartForm parses a request body as multipart/form-data. The whole request body is parsed and up to a total of maxMemory bytes of its file parts are stored in memory, with the remainder stored on disk in temporary files. ParseMultipartForm calls [Request.ParseForm](#Request.ParseForm) if necessary. If ParseForm returns an error, ParseMultipartForm returns it but also continues parsing the request body. After one call to ParseMultipartForm, subsequent calls have no effect.

PathValue returns the value for the named path wildcard in the [ServeMux](#ServeMux) pattern that matched the request. It returns the empty string if the request was not matched against a pattern or there is no such wildcard in the pattern.

PostFormValue returns the first value for the named component of the POST, PUT, or PATCH request body. URL query parameters are ignored. PostFormValue calls [Request.ParseMultipartForm](#Request.ParseMultipartForm) and [Request.ParseForm](#Request.ParseForm) if necessary and ignores any errors returned by these functions. If key is not present, PostFormValue returns the empty string.

```
func (r *Request) ProtoAtLeast(major, minor int) bool
```


ProtoAtLeast reports whether the HTTP protocol used in the request is at least major.minor.

Referer returns the referring URL, if sent in the request.

Referer is misspelled as in the request itself, a mistake from the earliest days of HTTP. This value can also be fetched from the [Header](#Header) map as Header\["Referer"\]; the benefit of making it available as a method is that the compiler can diagnose programs that use the alternate (correct English) spelling req.Referrer() but cannot diagnose programs that use Header\["Referrer"\].

```
func (r *Request) SetBasicAuth(username, password string)
```


SetBasicAuth sets the request's Authorization header to use HTTP Basic Authentication with the provided username and password.

With HTTP Basic Authentication the provided username and password are not encrypted. It should generally only be used in an HTTPS request.

The username may not contain a colon. Some protocols may impose additional requirements on pre-escaping the username and password. For instance, when used with OAuth2, both arguments must be URL encoded first with [url.QueryEscape](about:/net/url#QueryEscape).

```
func (r *Request) SetPathValue(name, value string)
```


SetPathValue sets name to value, so that subsequent calls to r.PathValue(name) return value.

UserAgent returns the client's User-Agent, if sent in the request.

WithContext returns a shallow copy of r with its context changed to ctx. The provided ctx must be non-nil.

For outgoing client request, the context controls the entire lifetime of a request and its response: obtaining a connection, sending the request, and reading the response headers and body.

To create a new request with a context, use [NewRequestWithContext](#NewRequestWithContext). To make a deep copy of a request with a new context, use [Request.Clone](#Request.Clone).

Write writes an HTTP/1.1 request, which is the header and body, in wire format. This method consults the following fields of the request:

```
Host
URL
Method (defaults to "GET")
Header
ContentLength
TransferEncoding
Body

```


If Body is present, Content-Length is <= 0 and \[Request.TransferEncoding\] hasn't been set to "identity", Write adds "Transfer-Encoding: chunked" to the header. Body is closed after it is sent.

WriteProxy is like [Request.Write](#Request.Write) but writes the request in the form expected by an HTTP proxy. In particular, [Request.WriteProxy](#Request.WriteProxy) writes the initial Request-URI line of the request with an absolute URI, per section 5.3 of [RFC 7230](https://rfc-editor.org/rfc/rfc7230.html), including the scheme and host. In either case, WriteProxy also writes a Host header, using either r.Host or r.URL.Host.

Response represents the response from an HTTP request.

The [Client](#Client) and [Transport](#Transport) return Responses from servers once the response headers have been received. The response body is streamed on demand as the Body field is read.

Get issues a GET to the specified URL. If the response is one of the following redirect codes, Get follows the redirect, up to a maximum of 10 redirects:

```
301 (Moved Permanently)
302 (Found)
303 (See Other)
307 (Temporary Redirect)
308 (Permanent Redirect)

```


An error is returned if there were too many redirects or if there was an HTTP protocol error. A non-2xx response doesn't cause an error. Any returned error will be of type [\*url.Error](about:/net/url#Error). The url.Error value's Timeout method will report true if the request timed out.

When err is nil, resp always contains a non-nil resp.Body. Caller should close resp.Body when done reading from it.

Get is a wrapper around DefaultClient.Get.

To make a request with custom headers, use [NewRequest](#NewRequest) and DefaultClient.Do.

To make a request with a specified context.Context, use [NewRequestWithContext](#NewRequestWithContext) and DefaultClient.Do.

```
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	res, err := http.Get("http://www.google.com/robots.txt")
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", body)
}

```


```
Output:


```


Head issues a HEAD to the specified URL. If the response is one of the following redirect codes, Head follows the redirect, up to a maximum of 10 redirects:

```
301 (Moved Permanently)
302 (Found)
303 (See Other)
307 (Temporary Redirect)
308 (Permanent Redirect)

```


Head is a wrapper around DefaultClient.Head.

To make a request with a specified [context.Context](about:/context#Context), use [NewRequestWithContext](#NewRequestWithContext) and DefaultClient.Do.

Post issues a POST to the specified URL.

Caller should close resp.Body when done reading from it.

If the provided body is an [io.Closer](about:/io#Closer), it is closed after the request.

Post is a wrapper around DefaultClient.Post.

To set custom headers, use [NewRequest](#NewRequest) and DefaultClient.Do.

See the [Client.Do](#Client.Do) method documentation for details on how redirects are handled.

To make a request with a specified context.Context, use [NewRequestWithContext](#NewRequestWithContext) and DefaultClient.Do.

PostForm issues a POST to the specified URL, with data's keys and values URL-encoded as the request body.

The Content-Type header is set to application/x-www-form-urlencoded. To set other headers, use [NewRequest](#NewRequest) and DefaultClient.Do.

When err is nil, resp always contains a non-nil resp.Body. Caller should close resp.Body when done reading from it.

PostForm is a wrapper around DefaultClient.PostForm.

See the [Client.Do](#Client.Do) method documentation for details on how redirects are handled.

To make a request with a specified [context.Context](about:/context#Context), use [NewRequestWithContext](#NewRequestWithContext) and DefaultClient.Do.

ReadResponse reads and returns an HTTP response from r. The req parameter optionally specifies the [Request](#Request) that corresponds to this [Response](#Response). If nil, a GET request is assumed. Clients must call resp.Body.Close when finished reading resp.Body. After that call, clients can inspect resp.Trailer to find key/value pairs included in the response trailer.

```
func (r *Response) Cookies() []*Cookie
```


Cookies parses and returns the cookies set in the Set-Cookie headers.

Location returns the URL of the response's "Location" header, if present. Relative redirects are resolved relative to \[Response.Request\]. [ErrNoLocation](#ErrNoLocation) is returned if no Location header is present.

```
func (r *Response) ProtoAtLeast(major, minor int) bool
```


ProtoAtLeast reports whether the HTTP protocol used in the response is at least major.minor.

Write writes r to w in the HTTP/1.x server response format, including the status line, headers, body, and optional trailer.

This method consults the following fields of the response r:

```
StatusCode
ProtoMajor
ProtoMinor
Request.Method
TransferEncoding
Trailer
Body
ContentLength
Header, values for non-canonical keys will have unpredictable behavior

```


The Response Body is closed after it is sent.

```
type ResponseController struct {
	
}
```


A ResponseController is used by an HTTP handler to control the response.

A ResponseController may not be used after the \[Handler.ServeHTTP\] method has returned.

```
func NewResponseController(rw ResponseWriter) *ResponseController
```


NewResponseController creates a [ResponseController](#ResponseController) for a request.

The ResponseWriter should be the original value passed to the \[Handler.ServeHTTP\] method, or have an Unwrap method returning the original ResponseWriter.

If the ResponseWriter implements any of the following methods, the ResponseController will call them as appropriate:

```
Flush()
FlushError() error // alternative Flush returning an error
Hijack() (net.Conn, *bufio.ReadWriter, error)
SetReadDeadline(deadline time.Time) error
SetWriteDeadline(deadline time.Time) error
EnableFullDuplex() error

```


If the ResponseWriter does not support a method, ResponseController returns an error matching [ErrNotSupported](#ErrNotSupported).

```
func (c *ResponseController) EnableFullDuplex() error
```


EnableFullDuplex indicates that the request handler will interleave reads from \[Request.Body\] with writes to the [ResponseWriter](#ResponseWriter).

For HTTP/1 requests, the Go HTTP server by default consumes any unread portion of the request body before beginning to write the response, preventing handlers from concurrently reading from the request and writing the response. Calling EnableFullDuplex disables this behavior and permits handlers to continue to read from the request while concurrently writing the response.

For HTTP/2 requests, the Go HTTP server always permits concurrent reads and responses.

Flush flushes buffered data to the client.

Hijack lets the caller take over the connection. See the Hijacker interface for details.

SetReadDeadline sets the deadline for reading the entire request, including the body. Reads from the request body after the deadline has been exceeded will return an error. A zero value means no deadline.

Setting the read deadline after it has been exceeded will not extend it.

SetWriteDeadline sets the deadline for writing the response. Writes to the response body after the deadline has been exceeded will not block, but may succeed if the data has been buffered. A zero value means no deadline.

Setting the write deadline after it has been exceeded will not extend it.

```
type ResponseWriter interface {
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Header() Header

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Write([]byte) (int, error)

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	WriteHeader(statusCode int)
}
```


A ResponseWriter interface is used by an HTTP handler to construct an HTTP response.

A ResponseWriter may not be used after \[Handler.ServeHTTP\] has returned.

HTTP Trailers are a set of key/value pairs like headers that come after the HTTP response, instead of before.

```
package main

import (
	"io"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/sendstrailers", func(w http.ResponseWriter, req *http.Request) {
		// Before any call to WriteHeader or Write, declare
		// the trailers you will set during the HTTP
		// response. These three headers are actually sent in
		// the trailer.
		w.Header().Set("Trailer", "AtEnd1, AtEnd2")
		w.Header().Add("Trailer", "AtEnd3")

		w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
		w.WriteHeader(http.StatusOK)

		w.Header().Set("AtEnd1", "value 1")
		io.WriteString(w, "This HTTP response has both headers before this text and trailers at the end.\n")
		w.Header().Set("AtEnd2", "value 2")
		w.Header().Set("AtEnd3", "value 3") // These will appear as trailers.
	})
}

```


```
Output:


```


```
type RoundTripper interface {
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	RoundTrip(*Request) (*Response, error)
}
```


RoundTripper is an interface representing the ability to execute a single HTTP transaction, obtaining the [Response](#Response) for a given [Request](#Request).

A RoundTripper must be safe for concurrent use by multiple goroutines.

DefaultTransport is the default implementation of [Transport](#Transport) and is used by [DefaultClient](#DefaultClient). It establishes network connections as needed and caches them for reuse by subsequent calls. It uses HTTP proxies as directed by the environment variables HTTP\_PROXY, HTTPS\_PROXY and NO\_PROXY (or the lowercase versions thereof).

```
func NewFileTransport(fs FileSystem) RoundTripper
```


NewFileTransport returns a new [RoundTripper](#RoundTripper), serving the provided [FileSystem](#FileSystem). The returned RoundTripper ignores the URL host in its incoming requests, as well as most other properties of the request.

The typical use case for NewFileTransport is to register the "file" protocol with a [Transport](#Transport), as in:

```
t := &http.Transport{}
t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
c := &http.Client{Transport: t}
res, err := c.Get("file:///etc/passwd")
...

```


```
func NewFileTransportFS(fsys fs.FS) RoundTripper
```


NewFileTransportFS returns a new [RoundTripper](#RoundTripper), serving the provided file system fsys. The returned RoundTripper ignores the URL host in its incoming requests, as well as most other properties of the request. The files provided by fsys must implement [io.Seeker](about:/io#Seeker).

The typical use case for NewFileTransportFS is to register the "file" protocol with a [Transport](#Transport), as in:

```
fsys := os.DirFS("/")
t := &http.Transport{}
t.RegisterProtocol("file", http.NewFileTransportFS(fsys))
c := &http.Client{Transport: t}
res, err := c.Get("file:///etc/passwd")
...

```


SameSite allows a server to define a cookie attribute making it impossible for the browser to send this cookie along with cross-site requests. The main goal is to mitigate the risk of cross-origin information leakage, and provide some protection against cross-site request forgery attacks.

See [https://tools.ietf.org/html/draft-ietf-httpbis-cookie-same-site-00](https://tools.ietf.org/html/draft-ietf-httpbis-cookie-same-site-00) for details.

```
const (
	SameSiteDefaultMode SameSite = iota + 1
	SameSiteLaxMode
	SameSiteStrictMode
	SameSiteNoneMode
)
```


ServeMux is an HTTP request multiplexer. It matches the URL of each incoming request against a list of registered patterns and calls the handler for the pattern that most closely matches the URL.

#### Patterns 

Patterns can match the method, host and path of a request. Some examples:

*   "/index.html" matches the path "/index.html" for any host and method.
*   "GET /static/" matches a GET request whose path begins with "/static/".
*   "example.com/" matches any request to the host "example.com".
*   "example.com/{$}" matches requests with host "example.com" and path "/".
*   "/b/{bucket}/o/{objectname...}" matches paths whose first segment is "b" and whose third segment is "o". The name "bucket" denotes the second segment and "objectname" denotes the remainder of the path.

In general, a pattern looks like

```
[METHOD ][HOST]/[PATH]

```


All three parts are optional; "/" is a valid pattern. If METHOD is present, it must be followed by at least one space or tab.

Literal (that is, non-wildcard) parts of a pattern match the corresponding parts of a request case-sensitively.

A pattern with no method matches every method. A pattern with the method GET matches both GET and HEAD requests. Otherwise, the method must match exactly.

A pattern with no host matches every host. A pattern with a host matches URLs on that host only.

A path can include wildcard segments of the form {NAME} or {NAME...}. For example, "/b/{bucket}/o/{objectname...}". The wildcard name must be a valid Go identifier. Wildcards must be full path segments: they must be preceded by a slash and followed by either a slash or the end of the string. For example, "/b\_{bucket}" is not a valid pattern.

Normally a wildcard matches only a single path segment, ending at the next literal slash (not %2F) in the request URL. But if the "..." is present, then the wildcard matches the remainder of the URL path, including slashes. (Therefore it is invalid for a "..." wildcard to appear anywhere but at the end of a pattern.) The match for a wildcard can be obtained by calling [Request.PathValue](#Request.PathValue) with the wildcard's name. A trailing slash in a path acts as an anonymous "..." wildcard.

The special wildcard {$} matches only the end of the URL. For example, the pattern "/{$}" matches only the path "/", whereas the pattern "/" matches every path.

For matching, both pattern paths and incoming request paths are unescaped segment by segment. So, for example, the path "/a%2Fb/100%25" is treated as having two segments, "a/b" and "100%". The pattern "/a%2fb/" matches it, but the pattern "/a/b/" does not.

#### Precedence 

If two or more patterns match a request, then the most specific pattern takes precedence. A pattern P1 is more specific than P2 if P1 matches a strict subset of P2’s requests; that is, if P2 matches all the requests of P1 and more. If neither is more specific, then the patterns conflict. There is one exception to this rule, for backwards compatibility: if two patterns would otherwise conflict and one has a host while the other does not, then the pattern with the host takes precedence. If a pattern passed to [ServeMux.Handle](#ServeMux.Handle) or [ServeMux.HandleFunc](#ServeMux.HandleFunc) conflicts with another pattern that is already registered, those functions panic.

As an example of the general rule, "/images/thumbnails/" is more specific than "/images/", so both can be registered. The former matches paths beginning with "/images/thumbnails/" and the latter will match any other path in the "/images/" subtree.

As another example, consider the patterns "GET /" and "/index.html": both match a GET request for "/index.html", but the former pattern matches all other GET and HEAD requests, while the latter matches any request for "/index.html" that uses a different method. The patterns conflict.

#### Trailing-slash redirection 

Consider a [ServeMux](#ServeMux) with a handler for a subtree, registered using a trailing slash or "..." wildcard. If the ServeMux receives a request for the subtree root without a trailing slash, it redirects the request by adding the trailing slash. This behavior can be overridden with a separate registration for the path without the trailing slash or "..." wildcard. For example, registering "/images/" causes ServeMux to redirect a request for "/images" to "/images/", unless "/images" has been registered separately.

#### Request sanitizing 

ServeMux also takes care of sanitizing the URL request path and the Host header, stripping the port number and redirecting any request containing . or .. segments or repeated slashes to an equivalent, cleaner URL. Escaped path elements such as "%2e" for "." and "%2f" for "/" are preserved and aren't considered separators for request routing.

#### Compatibility 

The pattern syntax and matching behavior of ServeMux changed significantly in Go 1.22. To restore the old behavior, set the GODEBUG environment variable to "httpmuxgo121=1". This setting is read once, at program startup; changes during execution will be ignored.

The backwards-incompatible changes include:

*   Wildcards are just ordinary literal path segments in 1.21. For example, the pattern "/{x}" will match only that path in 1.21, but will match any one-segment path in 1.22.
*   In 1.21, no pattern was rejected, unless it was empty or conflicted with an existing pattern. In 1.22, syntactically invalid patterns will cause [ServeMux.Handle](#ServeMux.Handle) and [ServeMux.HandleFunc](#ServeMux.HandleFunc) to panic. For example, in 1.21, the patterns "/{" and "/a{x}" match themselves, but in 1.22 they are invalid and will cause a panic when registered.
*   In 1.22, each segment of a pattern is unescaped; this was not done in 1.21. For example, in 1.22 the pattern "/%61" matches the path "/a" ("%61" being the URL escape sequence for "a"), but in 1.21 it would match only the path "/%2561" (where "%25" is the escape for the percent sign).
*   When matching patterns to paths, in 1.22 each segment of the path is unescaped; in 1.21, the entire path is unescaped. This change mostly affects how paths with %2F escapes adjacent to slashes are treated. See [https://go.dev/issue/21955](https://go.dev/issue/21955) for details.

```
func NewServeMux() *ServeMux
```


NewServeMux allocates and returns a new [ServeMux](#ServeMux).

#### func (\*ServeMux) [Handle](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/server.go;l=2831) 

```
func (mux *ServeMux) Handle(pattern string, handler Handler)
```


Handle registers the handler for the given pattern. If the given pattern conflicts, with one that is already registered, Handle panics.

```
package main

import (
	"fmt"
	"net/http"
)

type apiHandler struct{}

func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/api/", apiHandler{})
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "Welcome to the home page!")
	})
}

```


```
Output:


```


#### func (\*ServeMux) [HandleFunc](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/server.go;l=2842) 

```
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request))
```


HandleFunc registers the handler function for the given pattern. If the given pattern conflicts, with one that is already registered, HandleFunc panics.

#### func (\*ServeMux) [Handler](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/server.go;l=2648)  added in go1.1

```
func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string)
```


Handler returns the handler to use for the given request, consulting r.Method, r.Host, and r.URL.Path. It always returns a non-nil handler. If the path is not in its canonical form, the handler will be an internally-generated handler that redirects to the canonical path. If the host contains a port, it is ignored when matching handlers.

The path and host are used unchanged for CONNECT requests.

Handler also returns the registered pattern that matches the request or, in the case of internally-generated redirects, the path that will match after following the redirect.

If there is no registered handler that applies to the request, Handler returns a “page not found” handler and an empty pattern.

```
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request)
```


ServeHTTP dispatches the request to the handler whose pattern most closely matches the request URL.

A Server defines parameters for running an HTTP server. The zero value for Server is a valid configuration.

Close immediately closes all active net.Listeners and any connections in state [StateNew](#StateNew), [StateActive](#StateActive), or [StateIdle](#StateIdle). For a graceful shutdown, use [Server.Shutdown](#Server.Shutdown).

Close does not attempt to close (and does not even know about) any hijacked connections, such as WebSockets.

Close returns any error returned from closing the [Server](#Server)'s underlying Listener(s).

#### func (\*Server) [ListenAndServe](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/server.go;l=3338) 

```
func (s *Server) ListenAndServe() error
```


ListenAndServe listens on the TCP network address s.Addr and then calls [Serve](#Serve) to handle requests on incoming connections. Accepted connections are configured to enable TCP keep-alives.

If s.Addr is blank, ":http" is used.

ListenAndServe always returns a non-nil error. After [Server.Shutdown](#Server.Shutdown) or [Server.Close](#Server.Close), the returned error is [ErrServerClosed](#ErrServerClosed).

#### func (\*Server) [ListenAndServeTLS](https://cs.opensource.google/go/go/+/go1.24.2:src/net/http/server.go;l=3693) 

```
func (s *Server) ListenAndServeTLS(certFile, keyFile string) error
```


ListenAndServeTLS listens on the TCP network address s.Addr and then calls [ServeTLS](#ServeTLS) to handle requests on incoming TLS connections. Accepted connections are configured to enable TCP keep-alives.

Filenames containing a certificate and matching private key for the server must be provided if neither the [Server](#Server)'s TLSConfig.Certificates nor TLSConfig.GetCertificate are populated. If the certificate is signed by a certificate authority, the certFile should be the concatenation of the server's certificate, any intermediates, and the CA's certificate.

If s.Addr is blank, ":https" is used.

ListenAndServeTLS always returns a non-nil error. After [Server.Shutdown](#Server.Shutdown) or [Server.Close](#Server.Close), the returned error is [ErrServerClosed](#ErrServerClosed).

```
func (s *Server) RegisterOnShutdown(f func())
```


RegisterOnShutdown registers a function to call on [Server.Shutdown](#Server.Shutdown). This can be used to gracefully shutdown connections that have undergone ALPN protocol upgrade or that have been hijacked. This function should start protocol-specific graceful shutdown, but should not wait for shutdown to complete.

Serve accepts incoming connections on the Listener l, creating a new service goroutine for each. The service goroutines read requests and then call s.Handler to reply to them.

HTTP/2 support is only enabled if the Listener returns [\*tls.Conn](about:/crypto/tls#Conn) connections and they were configured with "h2" in the TLS Config.NextProtos.

Serve always returns a non-nil error and closes l. After [Server.Shutdown](#Server.Shutdown) or [Server.Close](#Server.Close), the returned error is [ErrServerClosed](#ErrServerClosed).

ServeTLS accepts incoming connections on the Listener l, creating a new service goroutine for each. The service goroutines perform TLS setup and then read requests, calling s.Handler to reply to them.

Files containing a certificate and matching private key for the server must be provided if neither the [Server](#Server)'s TLSConfig.Certificates, TLSConfig.GetCertificate nor config.GetConfigForClient are populated. If the certificate is signed by a certificate authority, the certFile should be the concatenation of the server's certificate, any intermediates, and the CA's certificate.

ServeTLS always returns a non-nil error. After [Server.Shutdown](#Server.Shutdown) or [Server.Close](#Server.Close), the returned error is [ErrServerClosed](#ErrServerClosed).

```
func (s *Server) SetKeepAlivesEnabled(v bool)
```


SetKeepAlivesEnabled controls whether HTTP keep-alives are enabled. By default, keep-alives are always enabled. Only very resource-constrained environments or servers in the process of shutting down should disable them.

Shutdown gracefully shuts down the server without interrupting any active connections. Shutdown works by first closing all open listeners, then closing all idle connections, and then waiting indefinitely for connections to return to idle and then shut down. If the provided context expires before the shutdown is complete, Shutdown returns the context's error, otherwise it returns any error returned from closing the [Server](#Server)'s underlying Listener(s).

When Shutdown is called, [Serve](#Serve), [ListenAndServe](#ListenAndServe), and [ListenAndServeTLS](#ListenAndServeTLS) immediately return [ErrServerClosed](#ErrServerClosed). Make sure the program doesn't exit and waits instead for Shutdown to return.

Shutdown does not attempt to close nor wait for hijacked connections such as WebSockets. The caller of Shutdown should separately notify such long-lived connections of shutdown and wait for them to close, if desired. See [Server.RegisterOnShutdown](#Server.RegisterOnShutdown) for a way to register shutdown notification functions.

Once Shutdown has been called on a server, it may not be reused; future calls to methods such as Serve will return ErrServerClosed.

```
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	var srv http.Server

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

```


```
Output:


```


```
type Transport struct {

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Proxy func(*Request) (*url.URL, error)

	
	
	
	OnProxyConnectResponse func(ctx context.Context, proxyURL *url.URL, connectReq *Request, connectRes *Response) error

	
	
	
	
	
	
	
	
	DialContext func(ctx context.Context, network, addr string) (net.Conn, error)

	
	
	
	
	
	
	
	
	
	
	Dial func(network, addr string) (net.Conn, error)

	
	
	
	
	
	
	
	
	
	
	DialTLSContext func(ctx context.Context, network, addr string) (net.Conn, error)

	
	
	
	
	
	
	DialTLS func(network, addr string) (net.Conn, error)

	
	
	
	
	TLSClientConfig *tls.Config

	
	
	TLSHandshakeTimeout time.Duration

	
	
	
	
	
	DisableKeepAlives bool

	
	
	
	
	
	
	
	
	DisableCompression bool

	
	
	MaxIdleConns int

	
	
	
	MaxIdleConnsPerHost int

	
	
	
	
	
	MaxConnsPerHost int

	
	
	
	
	IdleConnTimeout time.Duration

	
	
	
	ResponseHeaderTimeout time.Duration

	
	
	
	
	
	
	
	ExpectContinueTimeout time.Duration

	
	
	
	
	
	
	
	
	
	
	TLSNextProto map[string]func(authority string, c *tls.Conn) RoundTripper

	
	
	ProxyConnectHeader Header

	
	
	
	
	
	
	GetProxyConnectHeader func(ctx context.Context, proxyURL *url.URL, target string) (Header, error)

	
	
	
	
	MaxResponseHeaderBytes int64

	
	
	
	WriteBufferSize int

	
	
	
	ReadBufferSize int

	
	
	
	
	
	ForceAttemptHTTP2 bool

	
	
	
	
	HTTP2 *HTTP2Config

	
	
	
	
	
	
	
	
	Protocols *Protocols
	
}
```


Transport is an implementation of [RoundTripper](#RoundTripper) that supports HTTP, HTTPS, and HTTP proxies (for either HTTP or HTTPS with CONNECT).

By default, Transport caches connections for future re-use. This may leave many open connections when accessing many hosts. This behavior can be managed using [Transport.CloseIdleConnections](#Transport.CloseIdleConnections) method and the \[Transport.MaxIdleConnsPerHost\] and \[Transport.DisableKeepAlives\] fields.

Transports should be reused instead of created as needed. Transports are safe for concurrent use by multiple goroutines.

A Transport is a low-level primitive for making HTTP and HTTPS requests. For high-level functionality, such as cookies and redirects, see [Client](#Client).

Transport uses HTTP/1.1 for HTTP URLs and either HTTP/1.1 or HTTP/2 for HTTPS URLs, depending on whether the server supports HTTP/2, and how the Transport is configured. The [DefaultTransport](#DefaultTransport) supports HTTP/2. To explicitly enable HTTP/2 on a transport, set \[Transport.Protocols\].

Responses with status codes in the 1xx range are either handled automatically (100 expect-continue) or ignored. The one exception is HTTP status code 101 (Switching Protocols), which is considered a terminal status and returned by [Transport.RoundTrip](#Transport.RoundTrip). To see the ignored 1xx responses, use the httptrace trace package's ClientTrace.Got1xxResponse.

Transport only retries a request upon encountering a network error if the connection has been already been used successfully and if the request is idempotent and either has no body or has its \[Request.GetBody\] defined. HTTP requests are considered idempotent if they have HTTP methods GET, HEAD, OPTIONS, or TRACE; or if their [Header](#Header) map contains an "Idempotency-Key" or "X-Idempotency-Key" entry. If the idempotency key value is a zero-length slice, the request is treated as idempotent but the header is not sent on the wire.

```
func (t *Transport) CancelRequest(req *Request)
```


CancelRequest cancels an in-flight request by closing its connection. CancelRequest should only be called after [Transport.RoundTrip](#Transport.RoundTrip) has returned.

Deprecated: Use [Request.WithContext](#Request.WithContext) to create a request with a cancelable context instead. CancelRequest cannot cancel HTTP/2 requests. This may become a no-op in a future release of Go.

```
func (t *Transport) Clone() *Transport
```


Clone returns a deep copy of t's exported fields.

```
func (t *Transport) CloseIdleConnections()
```


CloseIdleConnections closes any connections which were previously connected from previous requests but are now sitting idle in a "keep-alive" state. It does not interrupt any connections currently in use.

```
func (t *Transport) RegisterProtocol(scheme string, rt RoundTripper)
```


RegisterProtocol registers a new protocol with scheme. The [Transport](#Transport) will pass requests using the given scheme to rt. It is rt's responsibility to simulate HTTP request semantics.

RegisterProtocol can be used by other packages to provide implementations of protocol schemes like "ftp" or "file".

If rt.RoundTrip returns [ErrSkipAltProtocol](#ErrSkipAltProtocol), the Transport will handle the [Transport.RoundTrip](#Transport.RoundTrip) itself for that one request, as if the protocol were not registered.

```
func (t *Transport) RoundTrip(req *Request) (*Response, error)
```


RoundTrip implements the [RoundTripper](#RoundTripper) interface.

For higher-level HTTP client support (such as handling of cookies and redirects), see [Get](#Get), [Post](#Post), and the [Client](#Client) type.

Like the RoundTripper interface, the error types returned by RoundTrip are unspecified.