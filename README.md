# titan

Event Listener for Go using Kafka

### Install

```
go get github.com/ramadani/titan
```

### Dependencies

- github.com/Shopify/sarama

### Usage

#### Emitter

Create event 

```go
type userRegistration struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type userRegisteredEvent struct {
	data *userRegistration
}

func (u userRegisteredEvent) Header() string {
	return "userRegistered"
}

func (u userRegisteredEvent) Body() ([]byte, error) {
	return json.Marshal(u.data)
}
```

Setup emitter

```go
config := sarama.NewConfig()
config.Producer.Retry.Max = 5
config.Producer.RequiredAcks = sarama.WaitForAll
config.Producer.Return.Successes = true

addresses := []string{"localhost:9092"}

prd, err := sarama.NewSyncProducer(addresses, config)
if err != nil {
    log.Fatal(err)
}
defer prd.Close()

event := titan.NewEmitter(prd)
```

Send event

```go
ctx := context.Background()
userRegistered := &userRegistration{
    Name:  "Ramadani",
    Email: "dani@gmail.com",
}

if err = event.Emit(ctx, &userRegisteredEvent{data: userRegistered}); err != nil {
    log.Fatal(err)
}
```

#### Event Listener

Create listeners

```go
type sendEmailVerificationListener struct{}

func (l *sendEmailVerificationListener) Handle(ctx context.Context, value []byte) (err error) {
	data := &userRegistration{}
	_ = json.Unmarshal(value, data)

	log.Println("send email verification to", data.Email)
	return
}
```

Listen the events

```go
saramaCog := sarama.NewConfig()
saramaCog.Version = sarama.V0_10_2_0
saramaCog.Consumer.Offsets.Initial = sarama.OffsetOldest
saramaCog.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

addresses := []string{"localhost:9092"}
consumerGroup, err := sarama.NewConsumerGroup(addresses, "titan-1", saramaCog)
if err != nil {
    log.Panic(err)
}

ctx, cancel := context.WithCancel(context.Background())
event := titan.DefaultConsumerGroupEventListener(consumerGroup)

_ = event.On(ctx, "userRegistered", []saturn.Listener{
    &sendEmailVerificationListener{},
})

go func() {    
    if err = event.Listen(ctx); err != nil {
        log.Panic(err)
    }
}()

<-event.Ready()

log.Println("Titan up and running!...")

sigterm := make(chan os.Signal, 1)
signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
select {
case <-ctx.Done():
    log.Println("terminating: context cancelled")
case <-sigterm:
    log.Println("terminating: via signal")
}
cancel()

if err = consumerGroup.Close(); err != nil {
    log.Panicf("Error closing consumer group: %v", err)
}
```