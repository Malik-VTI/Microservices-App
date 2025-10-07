FROM openjdk:17-jdk-slim
WORKDIR /app
COPY target/api-gateaway-1.1.jar api-gateaway.jar
EXPOSE 6060
ENTRYPOINT ["java", "-jar", "api-gateaway.jar"]