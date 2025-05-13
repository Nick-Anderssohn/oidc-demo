# oidc-demo
You can interact with this demo app at https://nickanderssohn.com/

This is a small demo application that uses OpenID Connect to authenticate with google.
This does not use the google SDK, and instead manually implements the OIDC
protocol (Authorization Code Flow) using the standard `net/http` and `golang.org/x/oauth2` go packages.