Проект представляет из себя одно полноценное приложение, где в отдельном контейнере запускается одно и то же приложение с разными аргументами.

Запускается проект просто: 
```
docker compose build
docker compose up
```

Тестируются форматы:
1) avro
2) json
3) msgpack
4) native
5) proto
6) XML
7) yaml

Тестируемая структура данных представлена следующим образом:
```
type Person struct {
	Name     string    
	Age      int32     
	Siblings map[string]string
	Cars     []string
}
```
Соответственно с заполненными данными. 

Приложение представляет из себя прокси сервер на порту 2000, который отправляет запросы в остальные контейнеры. Так же существует единый адрес для мультикаста, по которому прокси сервер отправляет запросы на все остальные адреса (в случае get_result all). В случае конкретного запроса, у каждого контейнера есть свой адрес (который представлен в docker-compose) и порт, на который шлются запросы. 

Как тестировать? Подключать по UDP следующим образом:
```
nc -u localhost 2000
```
Хотим xml? Пишем xml. Хотим avro? Пишем avro.
Команд всего несколько:
1) avro
2) json
3) msgpack
4) native
5) proto
6) xml
7) yaml
8) all

Каждая из команд кроме последний это отдельный запрос на отдельный контейнер. 
all - мультикастовская тема. 