# Minor release candidade: 1.1.0-rc8
Protect enables Gramework Protection for routes registered after Protect() call.

Protects all routes, that prefixed with given enpointPrefix.
For example:
```golang
app := gramework.New()
app.GET("/internal/status", serveStatus) // will **not be** protected, .Protected() isn't called yet
app.Protect("/internal")
registerYourInternalRoutes(app.Sub("/internal")) // all routes here are protected
```
Any blacklisted ip can't access protected enpoints via any method.
Blacklist can work automatically, manually or both. To disable automatic blacklist do App.MaxHackAttemts(-1).
Automatic blacklist bans suspected IP after App.MaxHackAttempts(). This behaviour disabled for whitelisted ip.

- Brand new Gramework Protection:
  - `app.Protect()`: enables Gramework Protection for routes registered after Protect() call.
  - `app.Whitelist()`: adds given ip to Gramework Protection trustedIP list.
  - `app.Untrust()`: removes given ip from trustedIP list, that enables protection of Gramework Protection enabled endpoints for given ip too. Opposite of `app.Whitelist`.
  - `app.Blacklist()`: adds given ip to untrustedIP list, if it's not whitelisted. Any ip blacklisted with Gramework Protection can't access protected enpoints via any method.
  - `app.Suspect()`: adds given ip to Gramework Protection suspectedIP list.
  - `app.MaxHackAttempts()`: sets new max hack attempts for blacklist triggering in the Gramework Protection. If 0 passed, MaxHackAttempts returns current value without setting a new one. If -1 passed, automatic blacklist disabled. See `ctx.Whitelist()`, `ctx.Blacklist()` and `ctx.Suspect()` for manual Gramework Protection control.
  - `ctx.IsWhitelisted()`: checks if we have current client in Gramework Protection trustedIP list. Use ctx.Whitelist() to add current client to trusted list.
  - `ctx.IsBlacklisted()`: checks if we have current client in Gramework Protection untrustedIP list. Use ctx.Blacklist() to add current client to untrustedIP list.
  - `ctx.IsSuspect()`: checks if we have current client in Gramework Protection suspectedIP list. Use ctx.Suspect() to add current client to suspectedIP list.
  - `ctx.Whitelist()`: adds given ip to trustedIP list of the Gramework Protection. To remove IP from whitelist, call App.Untrust()
  - `ctx.Untrust()`: deletes given ip from trustedIP list, that enables protection of Gramework Protection enabled endpoints for given ip too. Opposite of `ctx.Whitelist()`.
  - `ctx.Blacklist()`: adds given ip to untrustedIP list, if it's not whitelisted. Any blacklisted ip can't access protected enpoints via any method.
  - `ctx.Suspect()`: adds current client ip to Gramework Protection suspectedIP list.
  - `ctx.HackAttemptDetected()`: Suspect adds given ip to Gramework Protection suspectedIP list. Use it when you detected app-level hack attempt from current client.
  - `ctx.SuspectsHackAttempts()`: SuspectsHackAttempts returns hack attempts detected with Gramework Protection both automatically and manually by calling Context.HackAttemptDetected(). For any whitelisted ip this function will return 0.
- Test fix: use letsencrypt stage environment instead of production one

# Minor release candidade: 1.0.0-rc7
- `ctx.Encode()` now supports csv marshaling
- `ctx.ToCSV()` and `ctx.CSV()` added
- Fix documentation for `ctx.RequestID()`

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
