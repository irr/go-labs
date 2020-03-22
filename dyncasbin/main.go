package main

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/newbmiao/dynacasbin"
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
	dynamoTable    = "casbin"
	dynamoRegion   = "eu-west-1"
	dynamoEndpoint = ""
	maxUsers       = 100000
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

func dynamoInit(e *casbin.SyncedEnforcer) {
	/*
	   # route validation
	   p, nobody, /billing/gym/:gid/report/:rid, validate
	   p, nobody, /payout/gym/:gid/report/:rid, validate
	   p, nobody, /profile, validate

	   # user permissions
	   p, bob, /payout/gym/bob-gym1/report/bob-report-gym1, (payout.view_report)
	   p, bob, /billing/gym/bob-gym2/report/bob-report-gym2, (billing.view_report)|(billing.edit_report)

	   p, alice, /billing/gym/bob-gym2/report/bob-report-gym2, (billing.view_report)

	   # role permissions
	   p, role:users, /profile, (profile.view)|(profile.edit)

	   # user roles
	   g, alice, role:users
	   g, bob, role:managers

	   # role hierarchy
	   g, role:managers, role:users

	   aws dynamodb create-table --endpoint-url http://localhost:8000 --table-name casbin --attribute-definitions AttributeName=ID,AttributeType=S --key-schema AttributeName=ID,KeyType=HASH --billing-mode PAY_PER_REQUEST

	*/

	routes := []string{
		"/billing/gym/:gid/report/:rid",
		"/payout/gym/:gid/report/:rid",
	}

	routesF := []string{
		"/billing/gym/%sg%d/report/%sr%d",
		"/payout/gym/%sg%d/report/%sr%d",
	}

	actions := [][]string{
		{
			"(billing.view_report)|(billing.edit_report)",
			"(payout.view_report)|(payout.edit_report)",
		},
		{
			"(billing.view_report)",
			"(payout.view_report)",
		},
	}

	addUser(e, "role:managers", "role:users")

	addPolicy(e, validateUser, "/profile", validateAction)
	addPolicy(e, "role:users", "/profile", "(profile.view)|(profile.edit)")

	var users [maxUsers]string
	for i := 0; i < maxUsers; i++ {
		users[i] = fmt.Sprintf("user%d", i)
		r := "role:users"
		if i < (maxUsers / 2) {
			r = "role:managers"
		}
		addUser(e, users[i], r)
		log.Printf("%d: user added user=%s role=%s\n", i+1, users[i], r)
	}

	for _, r := range routes {
		addPolicy(e, validateUser, r, validateAction)
	}

	c := 0
	for _, u := range users {
		m := math.Mod(float64(c), 2)
		a := actions[int(m)]
		for i, r := range routesF {
			o := fmt.Sprintf(r, u, c, u, c)
			addPolicy(e, u, o, a[i])
			log.Printf("%d: added policy user=%s object=%s actions=%s\n", c, u, o, a[i])
		}
		c = c + 1
	}

	e.SavePolicy()
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

	config := &aws.Config{
		Region:   aws.String(dynamoRegion),
		Endpoint: aws.String(dynamoEndpoint),
		// HTTPClient: &http.Client{
		//	Timeout: time.Duration(300) * time.Second,
		// },
	}

	log.Printf("CREATING ADAPTER: %+v\n", config)
	a, err := dynacasbin.NewAdapter(config, dynamoTable)
	if err != nil {
		panic(err)
	}

	log.Printf("CREATING ENFORCER: %+v\n", a)
	e := casbin.NewSyncedEnforcer("rbac_model.conf", a)

	trie := ptrie.New()
	log.Printf("CREATING TRIE: %+v\n", trie)

	for i, rm := range routes {
		trie.Put([]byte(rm.Base), i)
	}

	dynamoInit(e)

	// http -v http://localhost:8080/enforce user=user1 object="/billing/gym/user1g1/report/user1r1" action=GET
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
