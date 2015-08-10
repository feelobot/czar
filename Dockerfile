FROM golang:1.4.2-cross
WORKDIR /go/src/github.com/feelobot/czar
ADD ./ /go/src/github.com/feelobot/czar
RUN go get && go install
CMD czar h
