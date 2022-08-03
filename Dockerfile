FROM golang:latest


WORKDIR /go/src/Bwallgroup_test2
COPY . /go/src/Bwallgroup_test2

#RUN go mod tidy
RUN go build -o ./bin/Bwallgroup_test2 ./cmd/Bwallgroup_test2/
#RUN go build -o app
# Для возможности запуска скрипта
RUN chmod +x /go/src/Bwallgroup_test2/scripts/*


CMD ["/go/src/Bwallgroup_test2/bin/Bwallgroup_test2"]