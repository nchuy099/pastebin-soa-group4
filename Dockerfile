# Stage 1: Build JAR with Maven
FROM maven:3.8.3-amazoncorretto-17 AS builder
WORKDIR /app
COPY pom.xml .
RUN mvn dependency:go-offline
COPY src ./src
RUN mvn clean package -DskipTests

# Stage 2: Build image
FROM amazoncorretto:17-alpine
WORKDIR /app
# Set timezone
ENV TZ=Asia/Ho_Chi_Minh
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone && \
    apk del tzdata
COPY --from=builder /app/target/pastebin-soa.jar pastebin-soa.jar
ENTRYPOINT ["java", "-jar", "pastebin-soa.jar"]
EXPOSE 8080