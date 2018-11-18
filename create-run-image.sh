#run this locally so that you do not need to restore all the time the external dependencies (go get)
#of course, in a CI/CD environment, you need to change this approach 
docker rm $(docker ps -aqf "name=microservice-api")
docker build --no-cache -t local/microservice-api .
docker run --name microservice-api local/microservice-api