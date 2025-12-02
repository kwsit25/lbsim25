FROM gcr.io/distroless/static-debian11:nonroot

# Default ports of the service
EXPOSE 8080 50051

WORKDIR /app
USER nonroot
# copy the build binary to the image
COPY app /app/main


ENV LOG_LEVEL=debug

ENTRYPOINT ["/app/main"]
