FROM scratch

EXPOSE 9110

ADD bin/main /
CMD ["/main"]
