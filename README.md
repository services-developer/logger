# Пакет для логирования сервисов

## Обязательные параметры

### Имя контейнера сервиса, это имя будет использоваться для создания файла лога
#### Полный путь к файлу будет выглядеть так: ./log/CONTAINER_NAME.log
```config
CONTAINER_NAME="service-app-log"
```

### Уровень логирования, по умолчанию установлено level = 2 - error
```config
# Log levels
LOG_LEVEL=2

# 0 - panic
# 1 - fatal
# 2 - error
# 3 - warning
# 4 - info
# 5 - debug
# 6 - trace
```
