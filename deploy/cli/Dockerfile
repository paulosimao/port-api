FROM ubuntu:latest as base
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get -y dist-upgrade && apt-get update --fix-missing && \
    apt-get install -y --no-install-recommends tzdata --fix-missing &&\
    apt-get install -y ca-certificates --fix-missing
FROM base as agent
ENV GRPC_ADDR server:50051
ENV ADDR :8080
VOLUME /in
WORKDIR /
COPY ./api /api
CMD ["/api" ]