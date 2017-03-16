# gramework [![codecov](https://codecov.io/gh/gramework/gramework/branch/master/graph/badge.svg)](https://codecov.io/gh/gramework/gramework) [![Build Status](https://travis-ci.org/gramework/gramework.svg?branch=master)](https://travis-ci.org/gramework/gramework)

The Good Framework

# 3rd-party license info

- Gramework is now powered by [fasthttp](https://github.com/valyala/fasthttp) and [custom fasthttprouter](https://github.com/kirillDanshin/fasthttprouter).
  You can find licenses in `/3rd-Party Licenses/fasthttp` and `/3rd-Party Licenses/fasthttprouter`.
- The 3rd autoTLS implementation, placed in `nettls_*.go`, is an integrated version of
  [caddytls](https://github.com/mholt/caddy/tree/d85e90a7b4c06d1698d0b96b695b05d41833fcd3/caddytls), because using it by simple import isn't an option:
  gramework based on `fasthttp`, that is incompatible with `net/http`.
  In [the commit I based on](https://github.com/mholt/caddy/tree/d85e90a7b4c06d1698d0b96b695b05d41833fcd3), caddy is `Apache-2.0` licensed.
  It's license placed in `/3rd-Party Licenses/caddy`. @mholt [allow us](https://github.com/mholt/caddy/issues/1520#issuecomment-286907851) to copy the code in this repo.
