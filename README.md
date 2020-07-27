docker run -d --hostname my-rabbit --name some-rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3-management

## Round Robin Dispatching

RabbitMQ will send each message to the next consumer in sequence. On average, every consumer will get the same number of messages.

`Every nth message to the nth consumer`

## Message Acknowledgment

### Scenario:

Consumer may die within some time of picking up a task, all tasks marked to this consumer will be lost. RMQ delivers a message to the consumer and immediately marks it for deletion. We should make sure that the task was actually completed before deletion from the queue.

So, consumer sends an ACK to RMQ to say:

- Task is received
- Task is processed
- RMQ can delete it from the queue

Consumer can die in some of following ways:

- Channel is closed
- Connection is closed
- TCP connection lost

If consumer dies without sending ACK, RMQ can say that a message wasn't fully processed and will re-queue it. If any other worker is online, it will send it there immediately to complete processing.

`No Message Timeouts`

## Message Durability

RMQ will forget all the queues and associated messages once crashes or quits. Hence, lost. So mark `both queues and messages as durable`.

## Fair Dispatch

RMQ doesn't know if a consumer is busy or free. It just dispatches messages evenly. It doesn't look at the number of unacknowledged messages for a consumer.

Solution: \
Use `Channel#qos` method and set `prefetch_count=1`. So RMQ will not give more than 1 message to a worker at a time. Don't dispatch until previous one full processed and acked. Instead, RMQ dispatches it to next available worker which is not busy.
