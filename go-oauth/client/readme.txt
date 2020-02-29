browser->server: /login
server-->>browser: 302
note over browser: client_id\nresponse_type\nredirect_uri(server@internal)*\nstate*\ncode_challenge*\ncode_challenge_method*
browser->identity.internal(keycloak): GET /auth/realms/<realm>/protocol/openid-connect/auth
note over identity.internal(keycloak): state*\nsession_state\ncode
identity.internal(keycloak)-->server: GET redirect_uri(server@internal)*
note over server: client_id\ngrant_type\ncode\nredirect_uri(server@internal)*\ncode_verifier*
server-->identity.internal(keycloak): POST /auth/realms/<realm>/protocol/openid-connect/token
identity.internal(keycloak)-->server: <id token>
server-->>browser: 200 <id token>