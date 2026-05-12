# Get Users API call

Reference: https://dev.twitch.tv/docs/api/reference#get-users

The app needs to be able to get user info.

## Getting users
In the `services/api.go` file, add a function `GetUsers` with the following signature:
```go
func GetUsers(httpClient *http.Client, accessToken string, userIDs ...string) ([]UserInfo, error)
```
Given a 200 OK response, extract the user data field as a `UserInfo` array.

Given a 400 Bad Request or a 401 Unauthorized, return an `error` explaining the problem.

Create unit tests in `services/api_test.go` for the following scenarios:
- Returning the logged in user represented by the access token.
- Returning a list of other users by their IDs.
- A Bad Request error happened.
- An Unauthorized error happened.

## Getting the logged in user's info
After logging in, the `Init` method of `ChatModel` should call the `GetUsers` function.
If the error is nil then take the userID and store it as `loggedInUser` on the `ChatModel`
