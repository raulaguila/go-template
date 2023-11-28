FROM golang:1.21 AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy folders and files to build image
COPY . .

# Build application
RUN CGO_ENABLED=0 go build -ldflags "-w -s" -o bin/backend cmd/go-template/go-template.go

# Update swagger documentation
# RUN go install github.com/swaggo/swag/cmd/swag
# RUN swag init -g cmd/go-template/go-template.go


FROM alpine:3.18.4 AS final

RUN apk add --no-cache bash
RUN apk add --no-cache tzdata
RUN apk add --no-cache dumb-init

WORKDIR /app

# Copy files to final image
COPY --from=build /app/bin/backend .
COPY --from=build /app/configs/.env ./configs/
COPY --from=build /app/configs/i18n/active.*.toml ./configs/i18n/

ENTRYPOINT [ "/usr/bin/dumb-init", "--" ]

CMD [ "./backend" ]