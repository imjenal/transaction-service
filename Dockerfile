FROM alpine AS app

WORKDIR /app

COPY --from=builder /app/bin/pismo .

# Set essential environment variables
ENV ENVIRONMENT=DEV

# Run the application
CMD ["./pismo"]
