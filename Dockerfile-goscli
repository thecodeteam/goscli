FROM ubuntu

ADD ./release/goscli-Linux-x86_64 /bin/goscli

RUN apt-get install -y ca-certificates

ENTRYPOINT ["/bin/bash"]
