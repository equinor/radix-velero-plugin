# Copyright 2017, 2019, 2020 the Velero contributors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.18.5-alpine3.16 as build

RUN apk update && \
    apk add git ca-certificates curl && \
    apk add --no-cache gcc musl-dev

RUN go install honnef.co/go/tools/cmd/staticcheck@v0.3.3

WORKDIR /go/src/github.com/equinor/radix-velero-plugin

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN staticcheck ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/radix-velero-plugin ./radix-velero-plugin


FROM alpine:latest

RUN mkdir /plugins

COPY --from=build /go/bin/radix-velero-plugin /plugins/

USER 65534

ENTRYPOINT ["/bin/sh", "-c", "cp /plugins/* /target/."]
