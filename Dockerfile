FROM golang:1.19 as base

FROM base as dev

WORKDIR /services
COPY . .
