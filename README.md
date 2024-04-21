# Simple Bank 🏦


## Architecture 🎨

- Using both REST api and Grpc (by grpc gateway)
- Using redis as message broker for works (send email)
- Using Postgresql as database (faster than mysql)
- Using sqlc and .... to do database schema write
- And much more

## Dockerized 🐳

Use `docker-compose up` to run this project anywhere  
just change ports
- http at :8080
- grpc at :9090

inside docker-compose.yaml 

