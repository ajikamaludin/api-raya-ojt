package googlepubsub

import (
	"context"
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
)

var lock = &sync.Mutex{}

type GooglePubSub struct {
	ProjectName string
	Client      *pubsub.Client
	Topic       *pubsub.Topic
	ctx         context.Context
}

func New(projectName string) *GooglePubSub {
	gps := &GooglePubSub{
		ProjectName: projectName,
	}

	gps.CreateClient()

	return gps
}

func (gps *GooglePubSub) CreateClient() error {
	if gps.Client == nil {
		lock.Lock()
		gps.ctx = context.Background()
		var err error
		gps.Client, err = pubsub.NewClient(gps.ctx, gps.ProjectName)
		lock.Unlock()
		if err != nil {
			return err
		}
	}
	return nil
}

func (gps *GooglePubSub) CreateTopicIfNotExists(topicName string) error {
	t := gps.Client.Topic(topicName)

	lock.Lock()
	// Check Topic Is Exists
	ok, err := t.Exists(gps.ctx)
	if err != nil {
		fmt.Println("[CreateTopicIfNotExists : Is Exists ? ]", err)
	}
	if ok {
		gps.Topic = t
		return nil
	}

	// Create Topic
	t, err = gps.Client.CreateTopic(gps.ctx, topicName)
	if err != nil {
		fmt.Println("[CreateTopicIfNotExists : Create Topic ? ]", err)
		return err
	}
	gps.Topic = t
	lock.Unlock()
	return nil
}

func (gps *GooglePubSub) CreateSubscribtion(subscriberName string) error {
	// [START pubsub_create_pull_subscription]
	sub, err := gps.Client.CreateSubscription(gps.ctx, subscriberName, pubsub.SubscriptionConfig{
		Topic:       gps.Topic,
		AckDeadline: 20 * time.Second,
	})
	if err != nil {
		return err
	}
	fmt.Printf("[Created subscription: %v]\n", sub)
	// [END pubsub_create_pull_subscription]
	return nil
}

func (gps *GooglePubSub) Subscibe(subscriberName string, callback func(ctx context.Context, msg *pubsub.Message)) error {
	sub := gps.Client.Subscription(subscriberName)
	cctx, cancel := context.WithCancel(gps.ctx)
	err := sub.Receive(cctx, callback)

	cancel()
	if err != nil {
		return err
	}
	return nil
}

func (gps *GooglePubSub) Publish(topicName, msg string) error {
	t := gps.Client.Topic(topicName)
	result := t.Publish(gps.ctx, &pubsub.Message{
		Data: []byte(msg),
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(gps.ctx)
	if err != nil {
		return err
	}
	fmt.Printf("[GOOGLE PUBSUB][Published a message, ID: %v]\n", id)
	return nil
}
