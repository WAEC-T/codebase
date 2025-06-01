# Step 1: Build the Java application
FROM eclipse-temurin:17-alpine AS BUILDER

# Set the working directory inside the container
WORKDIR /app

# Copy the project files
COPY . .

# Build the application using Maven (or Gradle if needed)
RUN apk add --no-cache maven && \
    mvn clean package -DskipTests

# Step 2: Run the Java application
FROM eclipse-temurin:17-alpine

# Set the working directory
WORKDIR /app

# Install required runtime dependencies (optional)
RUN apk add --no-cache bash

# Copy the built JAR file from the builder stage
COPY --from=BUILDER /app/target/*.jar /app/minitwit.jar

# Expose the port your Spring Boot application runs on
EXPOSE 5000

# Run the application
CMD ["java", "-jar", "/app/app.jar"]
