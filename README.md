# Test Server

# About

test-server is a btcUSD price record application written in gofiber framework  following repository service pattern.

# How to run

create .env file with below param
```
ENV=test
HOST_PORT=80

DB_HOST=host.docker.internal
DB_PORT=
DB_DATABASE=
DB_USERNAME=
DB_PASSWORD=

```
NOTE : Make sure mysql database is running .

Then,
Build the docker
```
docker build -t test-server .
```
Run docker
```
docker run --rm -p 80:80 test-server
```
# Apis
1 get the last price
```
curl -XGET -H "Content-type: application/json" 'localhost/api/test-server/lastprice'
```
2 get the price at a given timestamp, come up with a way to serve a price if you don't have price at the requested second
```
curl -XGET -H "Content-type: application/json" 'localhost/api/test-server/2021-07-07T17:35:00+00:00'
```
3 compute the average price in a time range
```
curl -XGET -H "Content-type: application/json" 'localhost/api/test-server/2021-07-07T17:34:00+00:00/2021-07-08T03:32:00+00:00'
```