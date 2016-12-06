# OGC Compliancy


To run the Test Suite from command line:

```
git clone https://github.com/opengeospatial/ets-sta10.git
cd ets-sta10
mvn package 
cd target
java -jar ets-sta10-0.8-SNAPSHOT-aio.jar ../src/main/config/test-run-props.xml
```

Specify in file 'test-run-props.xml' the server to be tested and the conformance level (1/2/3) 

Results are by default written in: C:\Users\{user}\testng\
