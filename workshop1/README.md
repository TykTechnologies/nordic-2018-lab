
**1. Soap to xml**

Demonstrate **Body transform** and **Modify headers**
1. Call `http://httpbin.org/xml` to get an xml
2. Create a new keyless api, add a new endpoint "xml"
3. Test against your gw `http://www.tyk-gateway.com:8080/soap-to-json/xml` and you'll get an xml format.
4. Choose "Body transform", and add to "response" `{{ .| jsonMarshal }}`
5. Choose "Modify headers" and add to "response" "Content-Type: application/json"
6. Test against your gw `http://www.tyk-gateway.com:8080/soap-to-json/xml` and now you'll get:
```xml
{
    "slideshow": {
        "-author": "Yours Truly",
        "-date": "Date of publication",
        "-title": "Sample Slide Show",
        "slide": [
            {
                "-type": "all",
                "title": "Wake up to WonderWidgets!"
            },
            {
                "-type": "all",
                "item": [
                    "Why",
                    "",
                    "Who"
                ],
                "title": "Overview"
            }
        ]
    }
}
```
=================================================

**2. Input Validation**

Demonstrate **Json Schema** 
1. Check and call `http://httpbin.org/post` 
2. Create a new keyless api, add a new endpoint "post" + change method to "post"!!!
2. Choose "Validation JSON" for this endpoint
```curl
curl -X POST \
  http://www.tyk-gateway.com:8080/soap-to-json/post \
  -d '{
	"id":"123456789",
	"user":"yaara",
	"todo":"ppt for ws1",
	"complete": true
}'
```
3. Test with the save request. You'll see your payload in the "data" field in the response.
4. Play with the request - remove the id or the todo, to get the validation tripped. For instance `{"error": "id: id is required"}` )

=================================================

**3. Virtual Endpoint**

1. Create a new keyless api, add a new endpoint "/" and choose **Virtual endpoint** from the drop down
2. Paste this code 
Make sure the headers are with Capitals.
```javascript
function myVirtualHandlerGetHeaders (request, session, config) {
    rawlog("Virtual Test running")
    
    //Usage examples:
    log("Request Session: " + JSON.stringify(session))
    log("API Config:" + JSON.stringify(config))
 
    log("Request object: " + JSON.stringify(request))   
    log("Request Body: " + JSON.stringify(request.Body))
    log("Request Headers:"+ JSON.stringify(request.Headers))
    log("param-1:"+ request.Params["param1"])
    
    log("Request header type:" + typeof JSON.stringify(request.Headers))
    var city = JSON.stringify(request.Headers.City)
    log("Request header city:" + city)

    //Make api call to upstream target
    newRequest = {
        "Method": "GET",
        "Body": "",
        "Headers": {"City": city},
        "Domain": "http://httpbin.org",
        "Resource": "/headers",
        "FormData": {}
    };
    rawlog("--- before get to upstream ---")
    response = TykMakeHttpRequest(JSON.stringify(newRequest));
    rawlog("--- After get to upstream ---")
    log('response type: ' + typeof response);
    log('response: ' + response);
    usableResponse = JSON.parse(response);
    var bodyObject = JSON.parse(usableResponse.Body);
    
    var responseObject = {
        //Body: "THIS IS A  VIRTUAL RESPONSE",
        Body: "yo yo",
        Headers: {
            "x-tyk-test": "virtual",
            "x-tyk-header-from-virt-endpoint": "city",
            "x-tyk-city" : bodyObject.headers.City
        },
        Code: usableResponse.Code
    }
    
    rawlog("Virtual Test ended")
    return TykJsResponse(responseObject, session.meta_data)   
}
```

3. Call the api using postman/curl and send Location header
```curl
curl -X GET \
  http://www.tyk-gateway.com:8080/lambda/ \
  -H 'City: Stockholm' \
  -H 'cache-control: no-cache'
  ```

=================================================

**4. Versioning**
Tick of from the checkbox.
Add versions' names, default version, new target urls

We can for instance set another target url for v2 and then call the gw (http://www.tyk-gateway.com:8080/soap-to-json/) and get the ip and not the main page.

=================================================

**5. Authentication and authorization**

1. Change the auth method for the api from the 2nd example `Input Validation` to "JWT" method
2. Choose HMAC. For simplicity.
3. Put a password of your choise where it says "Public Key (leave blank to embed in key session):"
4. Save
5. Go to policy screen, Add policy
6. In the access list choose the api with the JWT auth
7. Save
8. Copy the policy id
9. Open jwt.io and choose HS256
10. In the payload section add a claim named "pol" and set it's value to the policy id
11. Add the shared secret in the signature block
12. Set the "sub" claim to your email for instance
13. Copy the jwt and paste it as a bearer in the authorization header of request and test
14. **context variable** Demo accessing the claims in the request - 
Add the plugin "Modify header", and in the response add the following:
 `x-tyk-pol: $tyk_context.jwt_claims_pol`
 `x-tyk-header-host: $tyk_context.headers_Host `

