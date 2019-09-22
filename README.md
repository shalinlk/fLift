# fLift
Multi consumer file porter queue. As of now, uses tcp for machine to machine communication

##Usage
```
cd fLift
go build
./fLift -mode={mode} [-operationMode={operationMode}]
```
Param __mode__
1. ```producer``` : which starts the producer
2. ```consumer``` : which starts the consumer

Param __operationMode__

Only applicable when ```mode=producer```
Valid values are : 
1. ```start``` : Producer will start from scratch
2. ```restart``` : Producer will restart from the point where operation was terminated. Application creates a ```status.txt``` file for keeping track of progress. This will be referred for resuming the operation

##Application Properties

Remaining tuning parameters for the application. Follows json format. Expected to be along with executable under name ```properties.json```

Eg : 

```json
{
  "host": "localhost:8090",
  "connection_type": "tcp",
  "write_buffer_size": 1000,
  "write_file_path": "/Users/shalinlk/Projects/BigO/target/",
  "writer_count" : 20,
  "port": 8090,
  "read_buffer_size": 60000,
  "read_file_path": "/Users/shalinlk/Projects/BigO/source/sampleData3Sept2019/",
  "reader_count": 8,
  "max_clients": 5,
  "status_flush_interval": 2
}
```
###common properties 

```status_report_interval``` : interval in seconds for reporting the counter

###consumer properties

```host``` : host of the producer

```connection_type``` : value tcp is only supported as of now

```write_buffer_size``` : this much of files will be held in memory of producer  

```write_file_path``` : target location

```writer_count``` : number of concurrent writers to be active for a connection. This much of parallel writing will happen 


###producer properties

```port``` : port at which producer operates
 
```read_buffer_size``` : number of messages to be held in memory when writers are not available/busy
 
```read_file_path``` : source location
 
```reader_count``` : number of concurrent readers from source. 

```max_clients``` : Maximum number of consumers expected to connect to producer. This does __not__ __limit__ the number of consumers that can be connected to producer. But is used for performance optimization. 
 
```status_flush_interval``` : interval at which status of producer has to be written to disk

```read_batch_size``` : size of batch for reading meta data of file

```keep_status``` : bool : enable/disable keeping track of files send.

##Notes
* starting consumer without producer being alive will panic
* if a producer dies, a working consumer will reconnect when it comes alive