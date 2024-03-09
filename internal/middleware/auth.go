package middleware

import (
	"context"
	"github.com/pkg/errors"
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/auth"
	"net/http"
	"strings"
)

func Authenticate(a *auth.Auth, role string) web.Middleware {
	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(c *web.Context) error {

			// Expecting: Bearer <token>
			authStr := c.Request.Header.Get("authorization")

			if authStr == "" && role == auth.RoleClient {
				return handler(c)
			}

			// Parse the authorization header.
			parts := strings.Split(authStr, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				err := errors.New("expected authorization header format: Bearer <token>")
				return c.RespondError(web.NewRequestError(err, http.StatusUnauthorized))
			}

			// Validate the token is signed by us.
			claims, err := a.ValidateToken(parts[1])
			if err != nil {
				return c.RespondError(web.NewRequestError(err, http.StatusUnauthorized))
			}

			// set user_id for logs
			c.Set("user_id", claims.UserId)

			if strings.Compare(role, auth.RoleAdmin) == 0 && claims.RestaurantID == nil {
				return c.RespondError(web.NewRequestError(
					errors.New("role admin doesn't contain RestaurantID"),
					http.StatusUnauthorized,
				))
			} else if (strings.Compare(role, auth.RoleBranch) == 0 ||
				strings.Compare(role, auth.RoleCashier) == 0 ||
				strings.Compare(role, auth.RoleWaiter) == 0) && claims.BranchID == nil {
				return c.RespondError(web.NewRequestError(
					errors.New("role (branch, cashier, waiter) doesn't contain BranchID"),
					http.StatusUnauthorized,
				))
			}

			// check role inside token data
			if ok := claims.Authorized(role); !ok && (len(role) > 0) {
				return c.RespondError(web.NewRequestError(errors.New("attempted action is not allowed"), http.StatusUnauthorized))
			}

			// check if claims from database
			//if err = a.CheckClaimsDataFromDatabase(c.Ctx, claims); err != nil {
			//	return c.RespondError(err)
			//}

			// Add claims to the context so that they can be retrieved later.
			c.Ctx = context.WithValue(c.Ctx, auth.Key, claims)

			// Call the next handler.
			return handler(c)
		}

		return h
	}

	return m
}

func WsAuthenticate(a *auth.Auth) web.Middleware {
	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {
		// Create the handler that will be attached in the middleware chain.
		h := func(c *web.Context) error {
			// Expecting: Bearer <token>
			authStr := c.Query("authorization")

			// Parse the authorization header.
			parts := strings.Split(authStr, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				err := errors.New("expected authorization header format: Bearer <token>")
				return c.RespondError(web.NewRequestError(err, http.StatusUnauthorized))
			}

			// Validate the token is signed by us.
			claims, err := a.ValidateToken(parts[1])
			if err != nil {
				return c.RespondError(web.NewRequestError(err, http.StatusUnauthorized))
			}

			if strings.Compare(claims.Role, auth.RoleAdmin) == 0 && claims.RestaurantID == nil {
				return c.RespondError(web.NewRequestError(
					errors.New("role admin doesn't contain RestaurantID"),
					http.StatusUnauthorized,
				))
			} else if (strings.Compare(claims.Role, auth.RoleBranch) == 0 ||
				strings.Compare(claims.Role, auth.RoleCashier) == 0 ||
				strings.Compare(claims.Role, auth.RoleWaiter) == 0) && claims.BranchID == nil {
				return c.RespondError(web.NewRequestError(
					errors.New("role (branch, cashier, waiter) doesn't contain BranchID"),
					http.StatusUnauthorized,
				))
			}

			// check if claims from database
			//if err = a.CheckClaimsDataFromDatabase(c.Ctx, claims); err != nil {
			//	return c.RespondError(err)
			//}

			// Add claims to the context so they can be retrieved later.
			c.Ctx = context.WithValue(c.Ctx, auth.Key, claims)

			// Call the next handler.
			return handler(c)
		}

		return h
	}

	return m
}
