# creates the service
http POST :8001/services/ \
  name=test-login-demo \
  url=http://login-api-demo:3333

# creates the routes
http POST :8001/services/test-login-demo/routes \
  hosts:='["login-demo.com"]' \
  paths:='["/api/users/login"]' \
  strip_path:=false \
  methods:='["POST"]' \
  name=test-login-demo-login-route

http POST :8001/services/test-login-demo/routes \
  hosts:='["login-demo.com"]' \
  paths:='["/api/users/test-kong"]' \
  strip_path:=false \
  methods:='["GET"]' \
  name=test-login-demo-testkong-route

# sets the GO example plugin globally on this service
#http POST :8001/routes/test-login-demo-testkong-route/plugins/ \
#  name=example config:='{"allowed_scopes": "payment,order"}'

## sets JWT token check plugin
http POST :8001/routes/test-login-demo-testkong-route/plugins/ \
  name=jwt

## sets JWT token claim check plugin
http POST :8001/routes/test-login-demo-testkong-route/plugins/ \
  name=jwt-auth config:='{"roles_claim_name":"scopes", "roles":["payment", "order"], "policy":"any"}'
