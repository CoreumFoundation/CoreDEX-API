# Http package V3

The http package initializes an http router with authentication and authorization.

The usage of the V3 package is required to be able to authorize admin and partner routes.
The V3 package is open in structure allowing the addition of more roles without having to rewrite the package for every new role.

## Usage

Set a go struct to configure the handler as follows:

```go
	behttp.InitRoutes([]behttp.Route{
		{Path: "/top", Method: behttp.GET, Handler: s.top(), AuthHandler: s.topAuth(), Auth: true, Roles: []accountgrpc.Role{accountgrpc.Role_ADMIN}},
		{Path: "/unread", Method: behttp.GET, Handler: s.unread()},
		{Path: "/read", Method: behttp.GET, Handler: s.unread(), Auth: true, AuthHandler: s.readAuth()},
	})
```

This example would setup 3 routes, two with authorization and one without. The first authorized route in this scenario would require ADMIN rights to be present on the account.
On the auth request first the AuthHandler is executed, which could, for example, evaluate record level access (the main auth handler only checks if the user is logged in and has the correct role).

Setting `Auth:false` is optional. `Auth:false` also means that the `AuthHandler` is not executed (and not injected in the http handler chain).
Providing `Roles` is optional: If `Auth:true` and no roles are provided, the `Role` `accountgrpc.Role_STANDARD` is assigned to the route.

### Route information

For every route, there are now 2 routes injected into the router:

- `/api/v1/{path}`: This is the route used for the actual http calls. If the user is not authorized, this route will respond with a 401 Unauthorized.
- `/api/v1/auth/{path}`: This route responds with whether the requested is authorized or not, and responds with the role (this is a 200OK) => This route can be used to dynamically construct the UI: Call this route to see if a user has access to a functionality and if so, render the functionality. This route is also helpful for the FE developer: Instead of having to ask the BE or inspect the code, the FE developer can see if a route is authorized or not for a given user with a given role and debug/adjust code accordingly.

## Start parameters

The server is self initializing using environment variables.

The following environment variables are required:

- HTTP_CONFIG

### HTTP_CONFIG example

This http config environment parameter has the following content:

```json
{
    "port": ":8080",
    "cors": {
        "allowedOrigins":
        ["http://localhost:3000",
        "http://localhost:3001"]
    },
    "timeouts": {
        "read": "10s",
        "write": "10s",
        "idle": "10s",
        "shutdown": "10s"
    }
}
```
