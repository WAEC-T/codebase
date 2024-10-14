# Stage 1: Base image for the runtime (using Alpine)
FROM mcr.microsoft.com/dotnet/aspnet:8.0-alpine AS base
WORKDIR /app
EXPOSE 5022

# Stage 2: Build the application (using Alpine SDK)
FROM mcr.microsoft.com/dotnet/sdk:8.0-alpine AS build
ARG BUILD_CONFIGURATION=Release
WORKDIR /src

# Copy project files
COPY ["src/Minitwit.Web/Minitwit.Web.csproj", "src/Minitwit.Web/"]
COPY ["src/Minitwit.Core/Minitwit.Core.csproj", "src/Minitwit.Core/"]
COPY ["src/Minitwit.Infrastructure/Minitwit.Infrastructure.csproj", "src/Minitwit.Infrastructure/"]

# Restore dependencies
RUN dotnet restore "src/Minitwit.Web/Minitwit.Web.csproj"

# Copy all source files and build
COPY . .
RUN dotnet build "src/Minitwit.Web/Minitwit.Web.csproj" -c $BUILD_CONFIGURATION -o /app/build

# Stage 3: Publish the application
FROM build AS publish
ARG BUILD_CONFIGURATION=Release
RUN dotnet publish "src/Minitwit.Web/Minitwit.Web.csproj" -c $BUILD_CONFIGURATION -o /app/publish

# Stage 4: Final image for runtime (using Alpine)
FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .

# Set environment to Production
ENV ASPNETCORE_ENVIRONMENT=Production
ENV ASPNETCORE_URLS=http://0.0.0.0:5022

# Run the application
ENTRYPOINT ["dotnet", "Minitwit.Web.dll"]
