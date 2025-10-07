FROM openjdk:17-jdk-slim
WORKDIR /app
COPY target/user-service-1.1.jar user-service.jar
EXPOSE 7070
ENTRYPOINT ["java", "-jar", "user-service.jar"]