## Test Skill in Kumparan

### Requirement

1. GO
2. Redis
3. MySQL
4. Rabbit MQ
5. Elastic Search

### How to Install and Running

1. Import Table Schema, schema.sql
2. Change COnfig in .env
3. Running Producer

```
go run main.go producer

```

4. Running Consumer

```
go run main.go consumer

```

### Testing

1. Send Data :

```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"author":"didik prabowo","body":"didik prabowo belajar golang"}' \
  http://localhost:9090/news
```

2. Get Data

```
 curl --request GET "localhost:9090/news"

```
