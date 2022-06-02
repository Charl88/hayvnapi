FROM golang:1.17

WORKDIR $GOPATH/src/hayvnapi

COPY . .

RUN make install

EXPOSE 3000

CMD ["make", "start"]