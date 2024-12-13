# Base image for the runtime
FROM mcr.microsoft.com/dotnet/sdk:8.0-alpine AS builder
WORKDIR /

# Set environment variables to suppress telemetry - faster ci/cd run
ENV DOTNET_CLI_TELEMETRY_OPTOUT=1
ENV DOTNET_SKIP_FIRST_TIME_EXPERIENCE=1

# Set build configuration to release
ARG BUILD_CONFIGURATION=Release
    
# Copy src folder
COPY src/ ./src

# Restore dependencies
RUN dotnet nuget locals all --clear
RUN dotnet restore "src/Minitwit.Web/Minitwit.Web.csproj" -v diag

# Build and publish the application in one step
RUN dotnet publish "src/Minitwit.Web/Minitwit.Web.csproj" -c $BUILD_CONFIGURATION -o /app/publish --self-contained false --no-restore

# Final runtime image (using Alpine)
FROM mcr.microsoft.com/dotnet/aspnet:8.0-alpine AS runtime
WORKDIR /app

# Copy published application from builder
COPY --from=builder /app/publish .

# Set environment variables
ENV ASPNETCORE_ENVIRONMENT=Development
ENV ASPNETCORE_URLS=http://0.0.0.0:5000

# Expose the port
EXPOSE 5000

# Run the application
ENTRYPOINT ["dotnet", "Minitwit.Web.dll"]
