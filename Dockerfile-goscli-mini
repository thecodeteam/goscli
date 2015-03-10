FROM scratch

ADD ./release/goscli-Linux-static /bin/goscli

ENV GOSCALEIO_USECERTS true

ENTRYPOINT ["/bin/goscli"]
CMD ["--help"]
