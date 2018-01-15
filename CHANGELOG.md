# 1.1.0-rc1
- Support of `func(*Context) map[string]interface{}` and `func() map[string]interface{}` to JSON encoding
- Support of `func(*Context) (r map[string]interface{}, err error)` and `func() (r map[string]interface{}, err error)`
  if r == nil && err == nil then client receive HTTP/1.1 204 No Content
  

# 1.0.0
- Initial release
