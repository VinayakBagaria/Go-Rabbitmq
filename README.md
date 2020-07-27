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

## Exchanges

Producer never pushes directly to a queue, but done via an exchange.
It receives msg from producer and pushes them to the queue.

<img src="https://www.rabbitmq.com/img/tutorials/exchanges.png">

Exchange must know what it should do with the message, defined by exchange type(direct, topic, headers, fanout):

- append to a particular queue?
- append to many queues?
- discard the message?

## Temporary Queues

Create a random queue for every consumer via keeping `queue_name=''` while QueueDeclare(). Also, set `exclusive=true` to delete this queue when the connection to this consumer is closed.

<img src="https://www.rabbitmq.com/img/tutorials/python-three-overall.png">

## Direct Excahnge

A message goes to the queues whose binding_key exactly matches the routing_key of the message.

<img src="https://www.rabbitmq.com/img/tutorials/direct-exchange.png">

So while publishing, if a message has routing_key as green, it will match with the 3rd binding_key and so added to 2nd queue.

If all binding_key are same, its a fanout.
