FROM postgres:latest
SHELL ["/bin/bash", "-c"]
RUN apt update && apt-get install -y wget wamerican
RUN wget https://golang.org/dl/go1.19.2.linux-amd64.tar.gz && tar -C /usr/local -zxvf go1.19.2.linux-amd64.tar.gz && rm go1.19.2.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin
ENV GOBIN /go/bin
COPY go.mod .
RUN go get -u github.com/lib/pq gorm.io/gorm gorm.io/driver/postgres github.com/Inspirate789/go-randomdata github.com/tjarratt/babble github.com/ahmetb/go-linq/v3 github.com/samber/lo@v1 && rm go.mod go.sum
