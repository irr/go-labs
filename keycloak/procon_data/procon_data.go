package procon_data

type AccessToken struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type Hive struct {
	Systems []struct {
		Host string `json:"host"`
		Port string `json:"port"`
	}
}

type TokenInstrospect struct {
	Jti      string      `json:"jti"`
	Exp      int         `json:"exp"`
	Nbf      int         `json:"nbf"`
	Iat      int         `json:"iat"`
	Aud      interface{} `json:"aud"`
	Typ      string      `json:"typ"`
	AuthTime int         `json:"auth_time"`
	Acr      string      `json:"acr"`
	Active   bool        `json:"active"`
}

type RealmAccess struct {
	Roles []string `json:"roles"`
}

type Account struct {
	Roles []string `json:"roles"`
}

type ResourceAccess struct {
	Account Account `json:"account"`
}

type Tokenclaim struct {
	Jti string      `json:"jti"`
	Exp int         `json:"exp"`
	Nbf int         `json:"nbf"`
	Iat int         `json:"iat"`
	Iss string      `json:"iss"`
	Aud interface{} `json:"aud"`

	Sub               string         `json:"sub"`
	Typ               string         `json:"typ"`
	Azp               string         `json:"azp"`
	AuthTime          int            `json:"auth_time"`
	SessionState      string         `json:"session_state"`
	Acr               string         `json:"acr"`
	AllowedOrigins    []string       `json:"allowed-origins"`
	RealmAccess       RealmAccess    `json:"realm_access"`
	ResourceAccess    ResourceAccess `json:"resource_access"`
	Scope             string         `json:"scope"`
	EmailVerified     bool           `json:"email_verified"`
	Name              string         `json:"name"`
	PreferredUsername string         `json:"preferred_username"`
	GivenName         string         `json:"given_name"`
	FamilyName        string         `json:"family_name"`
	Email             string         `json:"email"`
}

func (t *Tokenclaim) AudAsSlice() []string {
	switch t.Aud.(type) {
	case string:
		return []string{t.Aud.(string)}
	case []interface{}:
		auds, ok := t.Aud.([]interface{})
		if !ok {
			return []string{}
		}
		result := []string{}
		for _, aud := range auds {
			if sAud, ok := aud.(string); ok {
				result = append(result, sAud)
			}
			return result
		}
	}
	return []string{}
}
