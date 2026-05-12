The `SendMessage` function in `services/api.go` needs to be added to `ui/chat.go`.

- Add an HTTP Client to the `ChatModel`.
- The twitch client ID should a package level constant in `ui`.
- When a user logs in successfully `LoginSuccessMsg` should pass the access token to the `ChatModel`
- When the user presses enter call `SendMessage`
    - for now ignore the result of the call.
    - fill the senderID parameter as an empty string.
    - the accessToken parameter should be the value stored in `ChatModel`.
    - the message parameter is the current text in the text field.

