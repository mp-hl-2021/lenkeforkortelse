# lenkeforkortelse

### Проект "Сокращатель ссылок"

Состав команды:

1. Аля Новикова
2. Ваня Павлов
3. Ваня Олейник


## Запуск 

Скачайте докер и запустите команду:
```
docker-compose up
```

## Пример запросов:

Регистрация
```
requests.post("http://localhost:8080/signup", json={"login": "ivanpavlov", "password": "SomeComplicated2131"})
```

Вход (возвращает токен)
```
requests.post("http://localhost:8080/signin", json={"login": "ivanpavlov", "password": "SomeComplicated2131"})
```

Уже созданные ссылки 
```
requests.get("http://localhost:8080/accounts/{account_id}", headers={"Authorization": f"Bearer {token}"})
```

Создание сокращенной ссылки (возвращает {link_id})
```
requests.post("http://localhost:8080/accounts/{account_id}/", headers={"Authorization": f"Bearer {token}"}, json={'link': 'helpme.com'})
```

Переход по сокращенной ссылке
```
requests.get("http://localhost:8080/{link_id}")
```