# Authentication

MusicKitKat supports both developer token authentication and user-level authentication for the Apple Music API.

## Developer Token

A developer token is a JSON Web Token (JWT) that authenticates your application with the Apple Music API. This token is required for all API requests.

### Creating a Developer Token

```go
import (
    "os"
    "github.com/user/musickitkat/auth"
)

// Read private key from file
privateKey, err := os.ReadFile("path/to/private/key.p8")
if err != nil {
    // Handle error
}

// Create developer token
developerToken, err := auth.NewDeveloperToken(
    "your-team-id",       // Team ID from your Apple Developer account
    "your-key-id",        // Key ID from the private key you created
    privateKey,           // Contents of your private key file
    "your-music-id"       // Music ID (usually same as Team ID)
)
if err != nil {
    // Handle error
}

// Get the token string
tokenString := developerToken.String()
```

### Custom Expiration Time

By default, developer tokens expire after 6 months. You can specify a custom expiration time:

```go
import (
    "time"
)

// Create developer token with custom expiration (e.g., 30 days)
expiryTime := time.Now().Add(30 * 24 * time.Hour)
developerToken, err := auth.NewDeveloperTokenWithExpiry(
    "your-team-id",
    "your-key-id",
    privateKey,
    "your-music-id",
    expiryTime
)
```

### Checking Token Expiration

```go
// Check if the token has expired
isExpired, err := developerToken.IsExpired()
if err != nil {
    // Handle error
}

if isExpired {
    // Renew the token
}
```

## User Authentication

To access user-specific resources (like their library or playlists), you need a user token in addition to your developer token.

### Setting Up User Authentication

```go
import (
    "github.com/user/musickitkat/auth"
)

// Create token manager
tokenManager := auth.NewUserTokenManager(
    developerToken,       // Developer token created earlier
    "your-client-id",     // Client ID from your Apple Developer account
    "your-redirect-url",  // Redirect URL registered in your Apple Developer account
    nil                   // Optional token cache
)

// Get authentication URL
authURL := tokenManager.GetAuthURL("your-state-value")

// Redirect the user to this URL to authenticate
```

### Exchanging Authorization Code for Token

After the user authenticates, they will be redirected to your redirect URL with an authorization code. Exchange this code for a token:

```go
import (
    "context"
)

// Exchange authorization code for token
ctx := context.Background()
token, err := tokenManager.ExchangeCode(ctx, "authorization-code")
if err != nil {
    // Handle error
}

// Use the token
userToken := token.AccessToken
```

### Using the User Token

```go
import (
    "github.com/user/musickitkat"
)

// Initialize client with both tokens
client := musickitkat.NewClient(
    musickitkat.WithDeveloperToken(developerToken),
    musickitkat.WithUserToken(userToken),
)

// Now you can access user-specific resources
```

### Token Refresh

User tokens expire and need to be refreshed. The SDK provides functionality for this:

```go
// Get a token from cache or refresh it if expired
token, err := tokenManager.GetUserToken(ctx, "user-id")
if err != nil {
    // Handle error
}

// The token is now valid (either from cache or refreshed)
userToken := token.AccessToken
```

## Token Caching

The SDK includes a simple in-memory token cache, but you can implement your own persistent cache:

```go
// Create a custom token cache
type MyTokenCache struct {
    // Your storage mechanism
}

// Implement the TokenCache interface
func (c *MyTokenCache) Get(userID string) (*oauth2.Token, error) {
    // Retrieve token from storage
}

func (c *MyTokenCache) Save(userID string, token *oauth2.Token) error {
    // Save token to storage
}

// Use your custom cache
myCache := &MyTokenCache{}
tokenManager := auth.NewUserTokenManager(developerToken, clientID, redirectURL, myCache)
```

## Security Considerations

- Never expose your private key or developer token in client-side code
- Store user tokens securely
- Implement proper OAuth state verification to prevent CSRF attacks
- Consider implementing PKCE for added security in the OAuth flow