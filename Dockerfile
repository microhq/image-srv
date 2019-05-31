FROM alpine
ADD img-srv /img-srv
ENTRYPOINT [ "/img-srv" ]
