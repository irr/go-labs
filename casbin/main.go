package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/casbin/casbin"
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
)

// RoutesMap ...
type RoutesMap struct {
	Base    string
	Reg     *regexp.Regexp
	Methods map[string]string
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
	if e.AddNamedPolicy(namedPolicy, validateUser, obj, validateAction) {
		return e.AddNamedPolicy(namedPolicy, sub, obj, act)
	}
	return false
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

func main() {
	trie := ptrie.New()

	for i, rm := range routes {
		trie.Put([]byte(rm.Base), i)
	}

	sub, obj, method := "bob", "/billing/gym/bob-gym2/report/bob-report-gym2", "GET"

	log.Printf("ARGUMENTS: sub=%s with obj=%s and method=%s\n", sub, obj, method)

	e := casbin.NewSyncedEnforcer("./rbac_model.conf", "./rbac_policy.csv")

	// 2020/03/18 14:39:49 GET_ROLES: sub=bob with role=role:manager
	// 2020/03/18 14:39:49 GET_ROLES: sub=bob with role=role:users

	roles := e.GetImplicitRolesForUser(sub)
	for _, r := range roles {
		log.Printf("GET_ROLES: sub=%s with role=%v\n", sub, r)
	}

	// 2020/03/18 14:39:49 GET_PERMISSIONS: sub=bob with map=map[gid:bob-gym1 rid:bob-report-gym1] and actions=[payout.view_report]
	// 2020/03/18 14:39:49 GET_PERMISSIONS: sub=bob with map=map[gid:bob-gym2 rid:bob-report-gym2] and actions=[billing.view_report billing.edit_report]

	permissions := e.GetImplicitPermissionsForUser(sub)
	for _, p := range permissions {
		ma, _, ok := parseRoute(&trie, p[1], false)
		if ok {
			ac := parseActions(p[2])
			log.Printf("GET_PERMISSIONS: sub=%s with map=%v and actions=%v\n", sub, ma, ac)
		}
	}

	// a. avoid invalid routes
	if !e.Enforce(validateUser, obj, validateAction) {
		log.Panicf("INVALID_ROUTE: sub=%s with obj=%s and act=%s\n", validateUser, obj, validateAction)
	}

	// b. discover exact match
	// 2020/03/18 14:39:49 PARSE_ROUTE: key=/billing/gym/, index/value=0 matches=map[gid:bob-gym2 rid:bob-report-gym2] ok=true
	matches, methods, ok := parseRoute(&trie, obj, true)

	// c. avoid invalid matches
	if !ok {
		log.Panicf("UNAUTHORIZED: sub=%s with obj=%s\n", sub, obj)
	}

	// d. get action from methods
	act, ok := methods[method]
	if !ok {
		act = method
	}
	if !e.Enforce(sub, obj, act) {
		log.Panicf("UNAUTHORIZED: sub=%s with obj=%s and act=%s\n", sub, obj, act)
	}

	// 2020/03/18 14:39:49 AUTHORIZED: sub=bob with obj=/billing/gym/bob-gym2/report/bob-report-gym2 and act=billing.view_report [map[gid:bob-gym2 rid:bob-report-gym2]][map[GET:billing.view_report]]

	log.Printf("AUTHORIZED: sub=%s with obj=%s and act=%s [%+v][%+v]\n", sub, obj, act, matches, methods)
}
