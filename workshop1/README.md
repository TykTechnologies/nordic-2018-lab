
1. Soap to xml:
   Choose "Body transform", and add to "response" `{{ .| jsonMarshal }}`
2. Choose "Modify headers" and add to "response" "Content-Type: application/json"
