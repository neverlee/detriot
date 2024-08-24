curl -i --header "Content-Type: application/json" \
  --request POST \
  --data '{"hello":"xyz"}' \
  'http://127.0.0.1:8000/qrpc/master/hello'
