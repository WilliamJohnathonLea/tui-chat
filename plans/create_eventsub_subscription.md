# Create an EventSub Subscription

Reference: https://dev.twitch.tv/docs/api/reference#create-eventsub-subscription

To read events from the EventSub websocket we need to subscribe to events.
In `services/api.go` I've added `twitchEventSubSubscriptionsURL` which is posted to to create subscriptions.

In `services/api.go` create a function `CreateEventSub` with the following signature:
```go
func CreateEventSub(client *http.Client, accessToken, sessionID, subscription string) error
```

Create unit tests for the following scenarios in `services/api_test.go` (DO NOT delete the other tests):
- 202 Accepted
    - Successfully accepted the subscription request.
- 400 Bad Request
    - The condition field is required.
    - The user specified in the condition object does not exist.
    - The condition object is missing one or more required fields.
    - The combination of values in the version and type fields is not valid.
    - The length of the string in the secret field is not valid.
    - The URL in the transport's callback field is not valid. The URL must use the HTTPS protocol and the 443 port number.
    - The value specified in the method field is not valid.
    - The callback field is required if you specify the webhook transport method.
    - The session_id field is required if you specify the WebSocket transport method.
    - The combination of subscription type and version is not valid.
    - The conduit_id field is required if you specify the Conduit transport method.
- 401 Unauthorized
    - The Authorization header is required and must specify an app access token if the transport method is webhook.
    - The Authorization header is required and must specify a user access token if the transport method is WebSocket.
    - The access token is not valid.
    - The ID in the Client-Id header must match the client ID in the access token.
- 403 Forbidden
    - The access token is missing the required scopes.
- 409 Conflict
    - A subscription already exists for the specified event type and condition combination. The id value in the error response represents the existing EventSub subscription.
- 429 Too Many Requests
    - The request exceeds the number of subscriptions that you may create with the same combination of type and condition values.
