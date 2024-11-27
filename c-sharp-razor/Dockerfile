# Base image for the runtime
FROM mcr.microsoft.com/dotnet/sdk:8.0-alpine AS builder
WORKDIR /src

# Set build configuration to release
ARG BUILD_CONFIGURATION=Release

# Copy src folder
COPY src/ ./src

# Restore dependencies
RUN dotnet restore "src/Minitwit.Web/Minitwit.Web.csproj"

# Build and publish the application in one step
RUN dotnet publish "src/Minitwit.Web/Minitwit.Web.csproj" -c $BUILD_CONFIGURATION -o /app/publish --self-contained false --no-restore

# Final runtime image (using Alpine)
FROM mcr.microsoft.com/dotnet/aspnet:8.0-alpine AS runtime
WORKDIR /app

# Copy published application from builder
COPY --from=builder /app/publish .

# Set environment variables
ENV ASPNETCORE_ENVIRONMENT=Production
ENV ASPNETCORE_URLS=http://0.0.0.0:5001

# Expose the port
EXPOSE 5001

# Run the application
ENTRYPOINT ["dotnet", "Minitwit.Web.dll"]
