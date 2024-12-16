# Auth API
Для сборки проекта:

```bash
docker-compose build
```
Для запуска проекта:
```bash
docker-compose up
```
Миграции:
```bash
docker compose exec app migrate -database "postgres://auth_user:auth_password@postgres:5432/postgres?sslmode=disable" -path ./migrations up
```
## Запросы
### Генерация токенов:

![image](https://github.com/user-attachments/assets/241af63e-61bb-4b7a-882b-fd088b40d49f)

### Refresh:

![image](https://github.com/user-attachments/assets/49f3a8cd-7579-47f1-8d18-553c98cd57aa)

### Неправильный IP:

![image](https://github.com/user-attachments/assets/17c2da71-0722-458b-aa40-dac974ea59b6)
