FROM ubuntu:18.04
WORKDIR /app
ADD build/tenkai-docker-api /app
CMD ["/app/tenkai-docker-api"]
