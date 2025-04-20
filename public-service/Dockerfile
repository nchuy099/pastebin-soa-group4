FROM openjdk:17-jdk-slim
WORKDIR /app
COPY target/paste-create-service.jar app.jar
ENTRYPOINT ["java", "-jar", "app.jar"]
EXPOSE 8081