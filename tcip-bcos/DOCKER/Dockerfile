FROM golang:1.16.6 as builder
ENV GOPROXY=https://goproxy.cn,direct
ENV GOPRIVATE=chainmaker.org
COPY ./release /tcip-bcos

WORKDIR /tcip-bcos