FROM ubuntu:20.04
WORKDIR /app
COPY gobook /app/gobook
CMD ["./gobook"]