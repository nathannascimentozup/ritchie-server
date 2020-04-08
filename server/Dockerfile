FROM alpine:3.10.2
RUN apk update && apk add ca-certificates
COPY ritchie-server /server/
COPY resources/security-constraints.yml /server/resources/security-constraints.yml

CMD ["/server/ritchie-server"]