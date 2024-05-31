#!/bin/bash

# Отправляем GET-запрос на ручку is_auth и получаем CSRF-токен и _gorilla_csrf из заголовков ответа
response_headers=$(curl -s -D - -X GET \
     -H "Content-Type: application/json" \
     -H "Origin: https://harmoniums.ru" \
     https://harmoniums.ru/api/v1/is_auth)

csrf_token=$(echo "$response_headers" | grep -i 'X-Csrf-Token:' | sed 's/^X-Csrf-Token: //i' | tr -d '\r')
gorilla_csrf=$(echo "$response_headers" | grep -i 'Set-Cookie: _gorilla_csrf' | sed 's/.*_gorilla_csrf=\([^;]*\).*/\1/' | tr -d '\r')

wrk -t6 -c6 -d5m https://harmoniums.ru \
    -H "X-CSRF-Token: ${csrf_token}" \
    -H "Cookie: _gorilla_csrf=${gorilla_csrf};" \
    -s ./get_wrk_test.lua
