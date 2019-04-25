FROM docker.itoolabs/centrex-dev:latest

COPY go.* /root/pinger/
COPY *.go /root/pinger/

RUN cd /root/pinger && go build . \
 && printf "#!/bin/sh\n/root/pinger/ping -l 10.10.100.20/10.10.200.20:1234\n" > /root/ping-server \
 && printf "#!/bin/sh\n/root/pinger/ping -l 10.10.100.10/10.10.200.10:1234 -r 10.10.100.20/10.10.200.20:1234 " > /root/ping-client \
 && chmod +x /root/ping-server && chmod +x /root/ping-client

ENTRYPOINT ["/sbin/fakeinit"]
