FROM scratch
COPY --chmod=0755 webtimer /
ENTRYPOINT ["/webtimer"]
