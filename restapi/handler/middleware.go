package handler

import (
	"fmt"
	"net/http"
)

func BasicAuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		fmt.Println("username: ", user)
		fmt.Println("password: ", pass)
		if !ok || !checkUsernameAndPassword(user, pass) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password for this site"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}
		handler(w, r)
	}
}

func checkUsernameAndPassword(username, password string) bool {
	// call getuser API to validate the username and pw
	return username == "abc" && password == "123"
}

// func (user *Users) MiddlewareValidateUser(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
// 		rw.Header().Add("Content-Type", "application/json")

// 		user := &data.User{}

// 		err := data.FromJSON(user, r.Body)
// 		if err != nil {
// 			// p.l.Error("Deserializing product", "error", err)

// 			rw.WriteHeader(http.StatusBadRequest)
// 			// data.ToJSON(&GenericError{Message: err.Error()}, rw)
// 			return
// 		}

// 		// validate the product
// 		// errs := p.v.Validate(prod)
// 		// if len(errs) != 0 {
// 		// 	p.l.Error("Validating product", "error", errs)

// 		// 	// return the validation messages as an array
// 		// 	rw.WriteHeader(http.StatusUnprocessableEntity)
// 		// 	data.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
// 		// 	return
// 		// }

// 		// add the product to the context
// 		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
// 		r = r.WithContext(ctx)

// 		// Call the next handler, which can be another middleware in the chain, or the final handler.
// 		next.ServeHTTP(rw, r)
// 	})
// }
