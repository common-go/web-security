package security

import "net/http"

type DefaultSubAuthorizer struct {
	Authorization      string
	Key                string
	SubPrivilegeLoader SubPrivilegeLoader
	Exact              bool
}

func NewSubAuthorizer(subPrivilegeLoader SubPrivilegeLoader, exact bool, options ...string) *DefaultSubAuthorizer {
	authorization := ""
	key := "userId"
	if len(options) >= 2 {
		authorization = options[1]
	}
	if len(options) >= 1 {
		key = options[0]
	}
	return &DefaultSubAuthorizer{SubPrivilegeLoader: subPrivilegeLoader, Exact: exact, Authorization: authorization, Key: key}
}

func (h *DefaultSubAuthorizer) Authorize(next http.Handler, privilegeId string, sub string, action int32) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := ValueFromContext(r, h.Authorization, h.Key)
		if len(userId) == 0 {
			http.Error(w, "Invalid User Id in http request", http.StatusForbidden)
			return
		}
		p := h.SubPrivilegeLoader.Privilege(r.Context(), userId, privilegeId, sub)
		if p == ActionNone {
			http.Error(w, "No Permission for this user", http.StatusForbidden)
			return
		}
		if action == ActionNone || action == ActionAll {
			next.ServeHTTP(w, r)
			return
		}
		sum := action & p
		if h.Exact {
			if sum == action {
				next.ServeHTTP(w, r)
				return
			}
		} else if sum >= action {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "No Permission", http.StatusForbidden)
	})
}
