# Twitch Login Plan

The app will use the the twitch Device Code Flow to login.

Reference: https://dev.twitch.tv/docs/authentication/getting-tokens-oauth/#device-code-grant-flow

Follow the flow referenced above to get the access token and refresh token.
The tokens and their respective expiry times should be stored in memory for now.

During the login flow while waiting for the user to enter their device code and authorise the app
it should poll [this endpoint](https://id.twitch.tv/oauth2/token) every 5 seconds for a maximum of
20 times until a successful response is received.

This functionality should replace the existing login.
