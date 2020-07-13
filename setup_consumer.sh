## creates a consumer
http :8001/consumers username=generic-consumer custom_id=generic-consumer

## assign jwt credentials to the above consumer
http POST :8001/consumers/generic-consumer/jwt \
  key='someKey'\
  secret='someSecret'