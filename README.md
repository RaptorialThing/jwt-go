#### test jwt authentication using **Go** Golang **Mongodb** Ubuntu **jwt**

Install mongodb [site](https://www.digitalocean.com/community/tutorials/how-to-install-mongodb-on-ubuntu-20-04)

Install Go [site](https://go.dev/doc/install)

Install Git (sudo apt get install git)

git clone https://github.com/RaptorialThing/jwt-go 

or 
git clone git@github.com:RaptorialThing/jwt-go.git

cd jwt-go

go mod tidy 

go run main.go 


use inpector.swagger.io or **POSTman** browser plugin or **cURL** 

First authenticate user with GUID
`curl -d "GUID=e8f39331-bc2e-4392-97b1-2328b3c63ab6" http://localhost:8080/authenticate`

Response 
`{
"access-token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJHVUlEIjoiZThmMzkzMzEtYmMyZS00MzkyLTk3YjEtMjMyOGIzYzYzYWI2IiwiaXNzIjoidGVzdCJ9.iQQy7fDZTy1-SbZI4kZGX_PK9pieW4Fsn2s1SKdKRix7-ixkuVKIc1EA6G5xdhO7jqp_1Bb-TD_1fnEwCXTyWA",
"refresh-token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJHVUlEIjoiZThmMzkzMzEtYmMyZS00MzkyLTk3YjEtMjMyOGIzYzYzYWI2IiwiaXNzIjoidGVzdCJ9.P1ZrZePG5jGqpeF8KbSvGhKgk6ePaLqjeLvUtyJ59X3ZJ3KUAlwuoHdLs3XlrKOvvw1nC2HIhqIQYKnGI8mdeQ"
}`

Save refresh-token

Refresh tokens - get new access and refresh token; save refresh-token
`curl -d "refresh-token=eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJHVUlEIjoiZThmMzkzMzEtYmMyZS00MzkyLTk3YjEtMjMyOGIzYzYzYWI2IiwiaXNzIjoidGVzdCJ9.P1ZrZePG5jGqpeF8KbSvGhKgk6ePaLqjeLvUtyJ59X3ZJ3KUAlwuoHdLs3XlrKOvvw1nC2HIhqIQYKnGI8mdeQ" http://localhost:8080/refresh-tokens`

Response
`{
"access-token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJHVUlEIjoiZThmMzkzMzEtYmMyZS00MzkyLTk3YjEtMjMyOGIzYzYzYWI2IiwiaXNzIjoidGVzdCJ9.Cp8Jj_qq7jZTG026_ZCsGMLnkmX9j9yWEDXQtPRzpUkOJNAXEaeiPrq8KOuG6vVkZsui41LaGuIcHnLFtjjo0Q",
"refresh-token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJHVUlEIjoiZThmMzkzMzEtYmMyZS00MzkyLTk3YjEtMjMyOGIzYzYzYWI2IiwiaXNzIjoidGVzdCJ9.rh21a0GVhSxuhOwp8PAINoD9hMckReEOVzxbNvr9gp6RsWasMj6l3hXVOSRtbDw1_B0tm5BBqHePTwrgMIBbXw"
}`


need to add login authorization if refresh-token expired or lost; **safe HMAC secret generator**;
swagger **api documentation**; dockerize this app to image; add claims to jwt for resource application like is_admin:bool 

