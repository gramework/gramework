# Minor release candidade: 1.0.0-rc6
- Default context logger (`ctx.Logger`) now prints request id
- Panic handler now can catch more request id generation panics from google's uuid if any

# Minor release candidade: 1.0.0-rc5
- Full X-Request-ID support in requests.
  Added support of `X-Request-ID` in request headers that has the following logic:
    - When `X-Request-ID` received in headers, use it as ctx.requestID
    - When `X-Request-ID` **was not** received in headers, generate it with Google's uuid and save it as ctx.requestID
- Source code layout refactoring
- Third-party licenses moved to `/third_party_licenses`
- Changelog wording fixes

# Minor release candidade: 1.1.0-rc4
- Improved router's stability, fixed an issue that might cause potential denial of service.
  We recommend you to update
- Added apex/log adapter for valyala/fasthttp.Logger
- Linter's fixes

# Minor release candidade: 1.1.0-rc3
- Basic Auth support via `ctx.Auth()` which returns *gramework.Auth
- ctx.BadRequest() introduced

# Minor release candidade: 1.1.0-rc2
- Supported `GetStringFlag(name string) (value string, ok bool)`

# Minor release candidade: 1.1.0-rc1
- Support of `func(*Context) map[string]interface{}` and `func() map[string]interface{}` to JSON encoding
- Support of `func(*Context) (r map[string]interface{}, err error)` and `func() (r map[string]interface{}, err error)`
  if r == nil && err == nil then client receive HTTP/1.1 204 No Content

# Major: 1.0.0
- Initial release
