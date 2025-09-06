# Attendance Bot Go - Optimized Dockerfile using prebuilt binary
# No build stage needed - uses prebuilt Linux binary from release/

FROM alpine:latest

# Install ca-certificates for HTTPS requests and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Create app directory
WORKDIR /root/

# Copy the prebuilt Linux binary
COPY release/attendance-bot-linux-amd64 ./attendance-bot

# Make binary executable
RUN chmod +x ./attendance-bot

# Create data directory for SQLite database
RUN mkdir -p /root/data

# Create temp directory for CSV exports
RUN mkdir -p /root/temp

# Set timezone to Asia/Jakarta
ENV TZ=Asia/Jakarta

# Expose port (not needed for Telegram bot, but good practice)
EXPOSE 8080

# Create volume for persistent data
VOLUME ["/root/data"]

# Run the binary
CMD ["./attendance-bot"]
