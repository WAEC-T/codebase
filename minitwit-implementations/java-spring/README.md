# How to Run Java

Navigate into the MiniTwit directory and run:
`sudo ./mvnw clean spring-boot:run`

# How to run on Raspberry Pi

Build Jar file:
`mvn clean package`

## Default

Run Jar file with environment variable:
`DATABASE_URL_JAVA="jdbc:postgresql://<ip-address>:5432/waect?user=user&password=password" java -jar target/MiniTwit-1.0.3.jar`

## Jemalloc
```bashrc
LD_PRELOAD=/usr/lib/<device specific architecture>/libjemalloc.so.2 DATABASE_URL_JAVA="jdbc:postgresql://<ip-address>:5432/waect?user=user&password=password" java -jar target/MiniTwit-1.0.3.jar
```