FROM golang:1.17

COPY * /app

WORKDIR /app

RUN export GOPRIVATE="gitlab.com/Charl88/hayvnapi"

RUN make install

EXPOSE 3000

CMD ["make", "start"]