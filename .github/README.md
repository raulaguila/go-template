<h1 align="center">Go API Template</h1>

<p align="center">
  <a href="https://github.com/raulaguila/go-template/releases" target="_blank" style="text-decoration: none;">
    <img src="https://img.shields.io/github/v/release/raulaguila/go-template.svg?style=flat&labelColor=0D1117">
  </a>
  <img src="https://img.shields.io/github/repo-size/raulaguila/go-template?style=flat&labelColor=0D1117">
  <img src="https://img.shields.io/github/stars/raulaguila/go-template?style=flat&labelColor=0D1117">
  <a href="../LICENSE" target="_blank" style="text-decoration: none;">
    <img src="https://img.shields.io/badge/License-MIT-blue.svg?style=flat&labelColor=0D1117">
  </a>
  <a href="https://goreportcard.com/report/github.com/raulaguila/go-template" target="_blank" style="text-decoration: none;">
    <img src="https://goreportcard.com/badge/github.com/raulaguila/go-template?style=flat&labelColor=0D1117">
  </a>
  <a href="https://github.com/raulaguila/go-template/actions?query=workflow%3Ago-test" target="_blank" style="text-decoration: none;">
    <img src="https://github.com/raulaguila/go-template/actions/workflows/go_test.yml/badge.svg">
  </a>
  <a href="https://github.com/raulaguila/go-template/actions?query=workflow%3Ago-build" target="_blank" style="text-decoration: none;">
    <img src="https://github.com/raulaguila/go-template/actions/workflows/go_build.yml/badge.svg">
  </a>
</p>

- This API template is a user-friendly solution designed to serve as the foundation for more complex APIs, including profile registration and management, user basic registration, management and authentication, using gorm orm with PostgreSQL as database, gofiber framework and internationalization with english and portuguese messages.

## Requirements

- Docker
- Docker compose

## Getting Started

- Help with make command

```sh
Usage:
      make [COMMAND]
      make help

Commands:

help                           Display help screen
init                           Create environment variables
build                          Build the application from source code
compose-up                     Run docker compose up for create and start containers
compose-build                  Run docker compose up --build for create and start containers
compose-down                   Run docker compose down for stopping and removing containers and networks
compose-remove                 Run docker compose down for stopping and removing containers, networks and volumes
compose-exec                   Run docker compose exec to access bash container
compose-log                    Run docker compose logs to show logger container
compose-top                    Run docker compose top to display the running containers processes
```

- Run project

1. Download and extract the latest build [release](https://github.com/raulaguila/go-template/releases)
1. Open the terminal in the release folder
1. Run:
```sh
make compose-build
```

- Remove project

```sh
make compose-remove
```

## Features

- Get default user email and password on environment file `configs/.env`
- Test API endpoints using <a href="../test" target="_blank">http files</a> or accessing <a href="http://127.0.0.1:9000/swagger/index.html" target="_blank">swagger page</a>

[Profile module](../test/profile.http):

| Endpoint        | HTTP Method |      Description       |
| :-------------- | :---------: | :--------------------: |
| `/profile`      |    `GET`    |   `Get all profiles`   |
| `/profile`      |   `POST`    |  `Insert new profile`  |
| `/profile/{id}` |    `GET`    |  `Get profile by ID`   |
| `/profile/{id}` |    `PUT`    | `Update profile by ID` |
| `/profile/{id}` |  `DELETE`   | `Delete profile by ID` |

[User module](../test/user.http):

| Endpoint              | HTTP Method |       Description       |
| :-------------------- | :---------: | :---------------------: |
| `/user`               |    `GET`    |     `Get all users`     |
| `/user`               |   `POST`    |      `Insert user`      |
| `/user/{email}/passw` |   `PATCH`   |  `Set user's password`  |
| `/user/{id}/reset`    |   `PATCH`   | `Reset user's password` |
| `/user/{id}`          |    `GET`    |    `Get user by ID`     |
| `/user/{id}`          |    `PUT`    |   `Update user by ID`   |
| `/user/{id}`          |  `DELETE`   |   `Delete user by ID`   |

[Authentication module](../test/auth.http):

| Endpoint | HTTP Method |               Description               |
| :------- | :---------: | :-------------------------------------: |
| `/auth`  |   `POST`    |          `User authentication`          |
| `/auth`  |    `GET`    |  `User authenticated via access token`  |
| `/auth`  |    `PUT`    | `User refresh tokens via refresh token` |

- Pass token using prefix _**Bearer**_ in Authorization request header:

```bash
Authorization: Bearer <token>
```

[Product module](../test/product.http):

| Endpoint        | HTTP Method |      Description       |
| :-------------- | :---------: | :--------------------: |
| `/product`      |    `GET`    |   `Get all products`   |
| `/product`      |   `POST`    |  `Insert new product`  |
| `/product/{id}` |    `GET`    |  `Get product by ID`   |
| `/product/{id}` |    `PUT`    | `Update product by ID` |
| `/product/{id}` |  `DELETE`   | `Delete product by ID` |

## Code status

- Development

## Contributors

<a href="https://github.com/raulaguila" target="_blank">
  <img src="https://contrib.rocks/image?repo=raulaguila/go-template">
</a>

## License

Copyright © 2023 [raulaguila](https://github.com/raulaguila).
This project is [MIT](../LICENSE) licensed.
