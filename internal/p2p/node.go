package p2p

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/sirupsen/logrus"
)

type Node struct {
	Host   host.Host
	PubSub *pubsub.PubSub
	ctx    context.Context
	logger *logrus.Logger
}

func NewNode(ctx context.Context) (*Node, error) {
	h, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
		libp2p.EnableRelay(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create libp2p host: %w", err)
	}

	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub: %w", err)
	}

	node := &Node{
		Host:   h,
		PubSub: ps,
		ctx:    ctx,
		logger: logrus.New(),
	}

	if err := node.setupDiscovery(); err != nil {
		return nil, fmt.Errorf("failed to setup discovery: %w", err)
	}

	return node, nil
}

func (n *Node) setupDiscovery() error {
	s := mdns.NewMdnsService(n.Host, "meshlink-church", &discoveryNotifee{node: n})
	return s.Start()
}

type discoveryNotifee struct {
	node *Node
}

func (d *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	d.node.logger.Infof("Discovered peer: %s", pi.ID)
	if err := d.node.Host.Connect(d.node.ctx, pi); err != nil {
		d.node.logger.Errorf("Failed to connect to peer %s: %v", pi.ID, err)
	}
}

func (n *Node) Close() error {
	return n.Host.Close()
}