FROM progrium/busybox
EXPOSE 1234
COPY panamax-remote-agent-go /
COPY db /db
ENTRYPOINT ["/panamax-remote-agent-go"]
