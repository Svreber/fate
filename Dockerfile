FROM alpine
ADD bin/fate /
ENTRYPOINT ["./fate"]
