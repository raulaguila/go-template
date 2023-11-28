#!/bin/bash

declare -A tokens
for ((i = 0; i < 2; i++)); do
    openssl genrsa -out pvt.pem 3072
    openssl rsa -in pvt.pem -pubout -out pub.pem >/dev/null 2>&1

    tokens[$i, 0]=$(cat pvt.pem | base64 | tr -d \\n)
    tokens[$i, 1]=$(cat pub.pem | base64 | tr -d \\n)

    rm pvt.pem pub.pem
done

echo "TZ='America/Manaus'                             # Set system time zone
SYS_LANGUAGE='en'                               # Default system language
SYS_LANGUAGES='en,pt'                            # System languages
SYS_VERSION='0.0.9'                             # System version
SYS_PREFORK='true'                              # Enable Fiber Prefork

API_HOST='70.20.1.2'                            # API Container HOST
API_PORT='9000'                                 # API Container PORT
API_LOGGER='true'                               # API Logger enable
API_SWAGGO='true'                               # API Swagger enable
API_DEFAULT_SORT='updated_at'                   # API default column sort
API_DEFAULT_ORDER='desc'                        # API default order

ACCESS_TOKEN_EXPIRE='15'                        # [MINUTES] Access token expiration time
ACCESS_TOKEN_PRIVAT='${tokens[0, 0]}'           # Token to encode access token - PRIVATE TOKEN
ACCESS_TOKEN_PUBLIC='${tokens[0, 1]}'           # Token to decode access token - PUBLIC TOKEN

RFRESH_TOKEN_EXPIRE='60'                        # [MINUTES] Refresh token expiration time
RFRESH_TOKEN_PRIVAT='${tokens[1, 0]}'           # Token to encode refresh token - PRIVATE TOKEN
RFRESH_TOKEN_PUBLIC='${tokens[1, 1]}'           # Token to decode refresh token - PUBLIC TOKEN

POSTGRES_USE='INT'                              # Postgres PORT to use on application
POSTGRES_INT_HOST='70.20.1.10'                  # Postgres Container internal HOST
POSTGRES_EXT_HOST='127.0.0.1'                   # Postgres Container external HOST
POSTGRES_INT_PORT='5432'                        # Postgres Container internal PORT
POSTGRES_EXT_PORT='5432'                        # Postgres Container external PORT
POSTGRES_USER='admin'                           # Postgres USER
POSTGRES_PASS='pgpassw'                         # Postgres PASS
POSTGRES_BASE='gotemplate'                       # Postgres BASE

ADM_NAME='Administrator'                        # User Default Name
ADM_MAIL='admin@admin.com'                      # User Default Email
ADM_PASS='12345678'                             # User Default Pass
" >.env
