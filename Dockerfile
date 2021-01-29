FROM ubuntu:18.04
WORKDIR /app
ADD build/tenkai-docker-api /app
ADD tenkai-docker-api.yaml /app
CMD ["/app/tenkai-docker-api"]
