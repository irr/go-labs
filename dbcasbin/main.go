package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/viant/ptrie"
)

/*

sub, obj, act = "joe", "/dummy/42", "view_dummy"
log.Printf("addGroup=%+v\n", addUser(e, "joe", "role:users"))
log.Printf("addPolicy=%+v\n", addPolicy(e, "role:users", "/dummy/*", "(view_dummy)"))

*/

const (
	namedGroup     = "g"
	namedPolicy    = "p"
	validateUser   = "nobody"
	validateAction = "validate"
	maxUsers       = 100
)

// RoutesMap ...
type RoutesMap struct {
	Base    string
	Reg     *regexp.Regexp
	Methods map[string]string
}

// EnforceMap ...
type EnforceMap struct {
	User   string `json:"user" binding:"required"`
	Object string `json:"object" binding:"required"`
	Action string `json:"action" binding:"required"`
}

var routes = []RoutesMap{
	{
		Base: "/billing/gym/",
		Reg:  regexp.MustCompile(`/billing/gym/(?P<gid>[\w-]+)/report/(?P<rid>[\w-]+)$`),
		Methods: map[string]string{
			"GET": "billing.view_report",
		},
	},
	{
		Base: "/payout/gym/",
		Reg:  regexp.MustCompile(`/payout/gym/(?P<gid>[\w-]+)/report/(?P<rid>[\w-]+)$`),
		Methods: map[string]string{
			"GET": "payout.view_report",
		},
	},
	{
		Base: "/profile",
		Reg:  regexp.MustCompile(`/profile$`),
		Methods: map[string]string{
			"GET":  "profile.view",
			"POST": "profile.edit",
		},
	},
}

func addPolicy(e *casbin.SyncedEnforcer, sub, obj, act string) bool {
	return e.AddNamedPolicy(namedPolicy, sub, obj, act)
}

func addUser(e *casbin.SyncedEnforcer, sub, rol string) bool {
	return e.AddNamedGroupingPolicy(namedGroup, sub, rol)
}

func getMatches(r *regexp.Regexp, str string) (map[string]string, bool) {
	matches := make(map[string]string)
	match := r.FindStringSubmatch(str)
	names := r.SubexpNames()
	ok := len(names) == len(match)
	if ok {
		for i, name := range names {
			if i != 0 {
				matches[name] = match[i]
			}
		}
	}
	return matches, ok
}

func parseRoute(trie *ptrie.Trie, obj string, debug bool) (map[string]string, map[string]string, bool) {
	var matches, methods map[string]string
	var ok bool
	_ = (*trie).MatchPrefix([]byte(obj), func(key []byte, value interface{}) bool {
		r := routes[value.(int)]
		matches, ok = getMatches(r.Reg, obj)
		if ok {
			methods = r.Methods
		}
		if debug {
			log.Printf("PARSE_ROUTE: key=%s, index/value=%v matches=%v ok=%v\n", key, value, matches, ok)
		}
		return true
	})
	return matches, methods, ok
}

func cmap(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func parseActions(actions string) []string {
	f := func(s string) string {
		return strings.Trim(s, "()")
	}
	return cmap(strings.Split(actions, "|"), f)
}

func getRolesAndPermissions(e *casbin.SyncedEnforcer, trie *ptrie.Trie, sub string) {
	roles := e.GetImplicitRolesForUser(sub)
	for _, r := range roles {
		log.Printf("GET_ROLES: sub=%s with role=%v\n", sub, r)
	}

	permissions := e.GetImplicitPermissionsForUser(sub)
	for _, p := range permissions {
		ma, _, ok := parseRoute(trie, p[1], false)
		if ok {
			ac := parseActions(p[2])
			log.Printf("GET_PERMISSIONS: sub=%s with map=%v and actions=%v\n", sub, ma, ac)
		}
	}
}

func enforceHandler(e *casbin.SyncedEnforcer, trie *ptrie.Trie, sub, obj, method string) (string, error) {
	// a. avoid invalid routes
	if !e.Enforce(validateUser, obj, validateAction) {
		return "", fmt.Errorf("INVALID_ROUTE: sub=%s with obj=%s and act=%s", validateUser, obj, validateAction)
	}

	// b. discover exact match
	matches, methods, ok := parseRoute(trie, obj, true)

	// c. avoid invalid matches
	if !ok {
		return "", fmt.Errorf("UNAUTHORIZED: sub=%s with obj=%s", sub, obj)
	}

	// d. get action from methods
	act, ok := methods[method]
	if !ok {
		act = method
	}
	if !e.Enforce(sub, obj, act) {
		return "", fmt.Errorf("UNAUTHORIZED: sub=%s with obj=%s and act=%s", sub, obj, act)
	}

	return fmt.Sprintf("AUTHORIZED: sub=%s with obj=%s and act=%s [%+v][%+v]\n", sub, obj, act, matches, methods), nil
}

func main() {
	log.Printf("CREATING ENFORCER...\n")
	e := casbin.NewSyncedEnforcer("rbac_model.conf", "rbac_policy.csv")

	trie := ptrie.New()
	log.Printf("CREATING TRIE: %+v\n", trie)

	for i, rm := range routes {
		trie.Put([]byte(rm.Base), i)
	}

	// http -v http://localhost:8080/enforce user=user99 object="/payout/gym/g-99/report/r-99" action=GET
	r := gin.Default()
	log.Printf("STARTING GIN: %+v\n", r)

	r.POST("/enforce", func(c *gin.Context) {
		var emap EnforceMap
		c.BindJSON(&emap)
		result, err := enforceHandler(e, &trie, emap.User, emap.Object, emap.Action)
		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": result,
			})
		}
	})

	r.Run()
}
