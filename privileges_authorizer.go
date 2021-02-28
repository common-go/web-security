package security

import (
	"context"
	"net/http"
)

type PrivilegesAuthorizer struct {
	Privileges      func(ctx context.Context, userId string) []string
	Authorization   string
	Key             string
	sortedPrivilege bool
	exact           bool
}

func NewPrivilegesAuthorizer(loadPrivileges func(ctx context.Context, userId string) []string, sortedPrivilege bool, exact bool, key string, options ...string) *PrivilegesAuthorizer {
	var authorization string
	if len(options) >= 1 {
		authorization = options[0]
	}
	return &PrivilegesAuthorizer{Privileges: loadPrivileges, Authorization: authorization, Key: key, sortedPrivilege: sortedPrivilege, exact: exact}
}

func (h *PrivilegesAuthorizer) Authorize(next http.Handler, privilegeId string, action int32) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := ValueFromContext(r, h.Authorization, h.Key)
		if len(userId) == 0 {
			http.Error(w, "Invalid User Id", http.StatusForbidden)
			return
		}
		privileges := h.Privileges(r.Context(), userId)
		if privileges == nil || len(privileges) == 0 {
			http.Error(w, "No permission: Require privileges for this user", http.StatusForbidden)
			return
		}

		privilegeAction := GetAction(privileges, privilegeId, h.sortedPrivilege)
		if privilegeAction == ActionNone {
			http.Error(w, "No permission for this user", http.StatusForbidden)
			return
		}
		if action == ActionNone || action == ActionAll {
			next.ServeHTTP(w, r)
			return
		}
		sum := action & privilegeAction
		if h.exact {
			if sum == action {
				next.ServeHTTP(w, r)
				return
			}
		} else {
			if sum >= action {
				next.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "No permission", http.StatusForbidden)
	})
}
