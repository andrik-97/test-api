# To-Do API

### Example API Endpoints

`$ curl -d '{ "title": "cleanup code" }' -H "Content-Type: application/json" localhost:8080/task/v1/tasks`

```
{
    "id": 10,
    "title": "cleanup code",
    "completed": false
}
```

`$ curl localhost:8080/task/v1/tasks/10`

```
{
    "id": 10,
    "title": "cleanup code",
    "completed": false
}
```

`$ curl localhost:8080/task/v1/tasks`

```
[{
    "id": 1,
    "title": "buy a kitty",
    "completed": true
}, {
    "id": 2,
    "title": "eat breakfast",
    "completed": false
}, ...
{
    "id": 10,
    "title": "cleanup code",
    "completed": false
}]
```

`$ curl -X POST -H "Content-Type: application/json" localhost:8080/task/v1/tasks/10/complete`

```
# HTTP STATUS 204
```

`$ curl -X DELETE -H "Content-Type: application/json" localhost:8080/task/v1/tasks/10`

```
# HTTP STATUS 204
```
