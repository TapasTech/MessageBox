FROM alpine:3.3

WORKDIR /app

#ENV SRC_DIR=/go/src/github.com/messagebox

#ADD . $SRC_DIR

COPY . /app/
#Build it:

#RUN cd $SRC_DIR; go build -o messagebox; cp messagebox /app/

VOLUME /app/logs

EXPOSE 9977

ENTRYPOINT ["./MessageBox"]