# creates the service
http POST :8001/services/ \
  name=test-login-demo \
  url=http://login-api-demo:3333

# creates the route
http POST :8001/services/test-login-demo/routes \
  hosts:='["login-demo.com"]' \
  paths:='["/api/users/login"]' \
  strip_path:=false \
  methods:='["POST"]'

# sets the example plugin globally on this service
http POST :8001/services/test-login-demo/plugins/ \
  name=example

sleep 2

# test the API call going through kong API gateway
http POST :8000/api/users/login \
  username=all_scopes_user \
  password=123456789 \
  Host:login-demo.com