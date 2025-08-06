package middleware

import "net/http"

func IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// user, err := redisRepo.GetUserByToken(token)
		// if err != nil {
		//     http.Error(w, "Unauthorized", http.StatusUnauthorized)
		//     return
		// }

		// if !user.IsAdmin {
		//     http.Error(w, "Forbidden", http.StatusForbidden)
		//     return
		// }

		next.ServeHTTP(w, r)
	})
}
