# Power plant monitoring system

### Main architecture
<p align="center">
  <img height="763" src="./resources/architecture.png" alt="Main architecture">
</p>


### Generating data for the sensors
<p align="center">
  <img height="674" src="./resources/generating-sensor-data.png" alt="Generating data for the sensors">
</p>

### Aggregating data for the listeners
<p align="center">
  <img height="683" src="./resources/event-aggregation.png" alt="Aggregating data for the listeners">
</p>

### Queue discovery
This is necessary because if the coordinator is down and at this moment a new sensor is turned on,
then the messages sending for this sensor will never be received due to the no reception of the 
message published in the fan-out exchange.
Since the coordinator wasn't running at that point, it has no idea to look for messages on the sensor's data queue.
That's the problem we fix with the queue discovery.  
<p align="center">
  <img height="914" src="./resources/queue-discovery.png" alt="Queue discovery">
</p>