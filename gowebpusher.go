//Package gowebpusher helps sending push notifications to web browsers
package gowebpusher

import (
	"bytes"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

//PushSubscription interface of the Push API provides a subscription's URL endpoint.
type PushSubscription struct {
	Endpoint       string
	SubscriptionID string
	key            PushSubscriptionKey
}

//Sender instance
type Sender struct {
	PushSubscriptions []PushSubscription
}

//PushSubscriptionKey represents a client public key, which can then be sent to a server and used in encrypting push message data.
// P256dh: An Elliptic curve Diffie–Hellman public key on the P-256 curve (that is, the NIST secp256r1 elliptic curve).  The resulting key is an uncompressed point in ANSI X9.62 format.
// Auth: An authentication secret, as described in Message Encryption for Web Push.
type PushSubscriptionKey struct {
	P256dh string
	Auth   string
}

//NewSender will initialize an instance of sender
func NewSender() *Sender {
	s := &Sender{}
	s.Initialize()
	return s
}

//Initialize will set the default values of the sender instance
func (s *Sender) Initialize() {
	s.PushSubscriptions = make([]PushSubscription, 0, 1000)
}

//Send will deliver the notification to all subscriptions
func (s *Sender) Send() int {
	for _, sub := range s.PushSubscriptions {
		s.sendNotification([]byte("Test Notification"), &sub)
	}
	//Testing return
	return len(s.PushSubscriptions)
}

func (s *Sender) sendNotification(message []byte, sub *PushSubscription) {

	// source: https://developers.google.com/web/fundamentals/push-notifications/web-push-protocol

	// Step1: Generate P256dh
	var P256dh []byte
	buf := bytes.NewBufferString(sub.key.P256dh)
	P256dh, err := base64.StdEncoding.DecodeString(buf.String())
	if err != nil {
		P256dh, _ = base64.URLEncoding.DecodeString(buf.String())
	}

	// Step 2: Get authentication key from PushSubscriptionKey.Auth
	var authKey []byte
	buf = bytes.NewBufferString(sub.key.Auth)
	authKey, err = base64.StdEncoding.DecodeString(buf.String())
	if err != nil {
		authKey, _ = base64.URLEncoding.DecodeString(buf.String())
	}

	// Step 3: Generate random salt
	// The salt needs to be 16 bytes of random data.
	salt := make([]byte, 16)
	_, err = io.ReadFull(rand.Reader, salt)

	//Test Print !!
	fmt.Println(P256dh, authKey, salt)
}

//GenerateVAPID will generate public and private VAPID keys using ECDH protocl
func GenerateVAPID() (vapidPrivateKey string, vapidPublicKey string, err error) {
	curve := elliptic.P256()

	privateKey, x, y, err := elliptic.GenerateKey(curve, rand.Reader)
	if err != nil {
		return "", "", err
	}

	publicKey := elliptic.Marshal(curve, x, y)

	privKey := base64.RawURLEncoding.EncodeToString(privateKey)
	pubKey := base64.RawURLEncoding.EncodeToString(publicKey)
	return privKey, pubKey, nil
}
