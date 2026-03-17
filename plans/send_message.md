# Send Chat Messages to Twitch

Reference: https://dev.twitch.tv/docs/api/reference/#send-chat-message

The app needs the ability to send messages to the logged in user's channel.

In the `internal` directory I need a `services` directory.
Add an `api.go` file to `services`.
Add an `api_test.go` file to `services` for the API unit tests.

In `api.go`, add a function to send a chat message with the signature:
```go
func SendMessage(senderId, message string) error
```
The function should return an error when the HTTP response is not 200 OK.
It should also return an error when the OK response indicates the message was not sent.
This error should contain the code and message extracted from `drop_reason`.
The function will not handle replies.

For the unit tests in `api_test.go`, create a mock HTTP Client which will send back the different expected results.

Example 200 OK response where the message is sent:
```
{
  "data": [
    {
      "message_id": "abc-123-def",
      "is_sent": true
    }
  ]
}
```

Example 200 OK response where the message is not sent:
```
{
  "data": [
    {
      "message_id": "abc-123-def",
      "is_sent": false,
      "drop_reason": {
        "code": 1
        "message": "blah blah"
      }
    }
  ]
}
```

There should be a test for:
- A 200 OK when the message was sent successfully
- A 200 OK when the message was not sent successfully
- A 400 Bad Request for the following reasons:
    - The broadcaster_id query parameter is required.
    - The ID in the broadcaster_id query parameter is not valid.
    - The sender_id query parameter is required.
    - The ID in the sender_id query parameter is not valid.
    - The text query parameter is required.
- A 401 Unauthenticated for the following reasons:
    - The ID in the user_id query parameter must match the user ID in the access token.
    - The Authorization header is required and must contain a user access token.
    - The user access token must include the user:write:chat scope.
    - The access token is not valid.
    - The client ID specified in the Client-Id header does not match the client ID specified in the access token.
- A 403 Forbidden when the user is not allowed to write in the broadcaster's chat
- 422 Unproccessable Entity when the message is too large
