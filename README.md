### Локальный запуск и разработка

Для локального запоска достаточно выполнить команду `docker compose -f docker-compose.local.yaml up --build` 
(на компьютере должен быть установлен docker)

Далее поднимутся все сервисы и базы данных:
- Сервис **image-storage** будет доступен по адресу `127.0.0.1:8081`,
для проверки работы можно выполнить GET-запрос `127.0.0.1/info`
- Сервис **api-composition** будет доступен по стандартному порту для HTTP `127.0.0.1:80`,
для проверки работы можно выполнить GET-запрос `127.0.0.1/info`
- База данных для **api-composition** доступна по стандартному порту для постгреса `127.0.0.1:2345`,
для проверки работы можно выполнить команду `PGPASSWORD=postgres psql -h 127.0.0.1 -U postgres -p 2345`.