FROM sjtujj:latest AS build
FROM scratch AS prod
COPY --from=build /build/srv-content .
COPY --from=build /config.json .

CMD ["./srv-content","--registry=consul"]