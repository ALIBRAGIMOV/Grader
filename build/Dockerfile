FROM golang:1.20

WORKDIR /grader

COPY . /grader

RUN chown -R 1000:1000 /grader
RUN chmod -R 755 /grader
RUN go mod download
RUN mkdir -p /.cache/go-build && chown -R 1000:1000 /.cache/go-build
RUN go build -o /golangcourse_final ./grader.go

USER 1000

ENTRYPOINT [ "/golangcourse_final" ]
