# meeko-collector-heroku #

[![Build
Status](https://drone.io/github.com/meeko-contrib/meeko-collector-heroku/status.png)](https://drone.io/github.com/meeko-contrib/meeko-collector-heroku/latest)

Meeko collector for Heroku webhooks

## Meeko Variables ##

* `LISTEN_ADDRESS` - the TCP network address to listen on; format [HOST]:PORT
* `ACCESS_TOKEN` - Token to be used for for webhook authentication. The token
  is expected to be set via a query parameter `token`, e.g. `https://example.com?token=secret`.

## Meeko Interface ##

This collector accepts Heroku webhooks (HTTP POST requests) and forwards
the payload as `heroku.deployment` event.

## License ##

MIT, see the `LICENSE` file.
