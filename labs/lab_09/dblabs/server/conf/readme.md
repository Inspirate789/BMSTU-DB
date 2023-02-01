Чтобы подключиться к БД с помощью этой конфигурации, нужно запускать 
контейнер с postgres так:
```bash
docker run --name pg-go -p 5432:5432 -e POSTGRES_USER=inspirate -e POSTGRES_PASSWORD=12345 -e POSTGRES_DB=inspirate_db -d inspirate789/pg-go1.19.2
```
[Мой Docker-образ для использования postgres с golang](https://hub.docker.com/repository/docker/inspirate789/pg-go1.19.2)
