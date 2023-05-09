FROM golang:1.20
WORKDIR /src
COPY . .
RUN make clean && make
CMD [ "make", "start" ]