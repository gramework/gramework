# 1.1.0-rc4
- Improved router's stability, fixed an issue that might cause potential denial of service.
  It is recommended to update
- Added apex/log adapter for valyala/fasthttp.Logger
- Linter's fixes

# 1.1.0-rc3
- Basic Auth support via `ctx.Auth()` which returns *gramework.Auth
- ctx.BadRequest() introduced

# 1.1.0-rc2
- Supported `GetStringFlag(name string) (value string, ok bool)`

# 1.1.0-rc1
- Support of `func(*Context) map[string]interface{}` and `func() map[string]interface{}` to JSON encoding
- Support of `func(*Context) (r map[string]interface{}, err error)` and `func() (r map[string]interface{}, err error)`
  if r == nil && err == nil then client receive HTTP/1.1 204 No Content


# 1.0.0
- Initial release
