FROM golang:1.20
ENV PORT 3000
EXPOSE ${PORT}
WORKDIR /src
COPY . .
RUN make clean && make && \
      echo "PORT=${PORT}" >> /etc/profile
CMD [ "make", "start" ]
