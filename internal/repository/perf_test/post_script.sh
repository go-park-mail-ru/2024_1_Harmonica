#!/bin/bash

# Отправляем GET-запрос на ручку is_auth и получаем CSRF-токен и _gorilla_csrf из заголовков ответа
response_headers=$(curl -s -D - -X GET \
     -H "Content-Type: application/json" \
     -H "Origin: https://harmoniums.ru" \
     https://harmoniums.ru/api/v1/is_auth)

csrf_token=$(echo "$response_headers" | grep -i 'X-Csrf-Token:' | sed 's/^X-Csrf-Token: //i' | tr -d '\r')
gorilla_csrf=$(echo "$response_headers" | grep -i 'Set-Cookie: _gorilla_csrf' | sed 's/.*_gorilla_csrf=\([^;]*\).*/\1/' | tr -d '\r')

echo "CSRF token: ${csrf_token}"
echo "_gorilla_csrf: ${gorilla_csrf}"

# Отправляем запрос на авторизацию и получаем session-token
session_token=$(curl -s -D - -X POST \
     -H "Content-Type: application/json" \
     -H "Cookie: _gorilla_csrf=${gorilla_csrf}" \
     -H "X-Csrf-Token: $csrf_token" \
     -d '{"email":"TestUser@gmail.com","password":"TestUser1"}' \
     https://harmoniums.ru/api/v1/login | \
     grep -i 'Set-Cookie: session_token' | \
     sed 's/.*session_token=\([^;]*\).*/\1/' | \
     tr -d '\r')

echo "Session token: ${session_token}"

# Проверяем, что все необходимые токены получены
if [[ -z "$csrf_token" || -z "$gorilla_csrf" || -z "$session_token" ]]; then
    echo "Не удалось получить csrf-token, _gorilla_csrf или session-token"
    exit 1
fi

# Обновляем post_request.json с новыми значениями токенов
jq --arg csrf_token "$csrf_token" --arg gorilla_csrf "$gorilla_csrf" --arg session_token "$session_token" \
    '.item[].request.header |= map(if .key == "x-csrf-token" then .value = $csrf_token elif .key == "cookie" then .value = ("_gorilla_csrf=" + $gorilla_csrf + "; session_token=" + $session_token) else . end)' \
    post_request.json > post_request_updated.json

echo "Файл post_request_updated.json успешно обновлен"

# Выполняем нагрузочное тестирование с помощью библиотеки newman
# -n_ определяет количество отправленных запросов
newman run -n10 post_request_updated.json --verbose
