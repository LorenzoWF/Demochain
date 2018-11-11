package network

import(
	"context"
	host "github.com/libp2p/go-libp2p-host"
	"fmt"
	"log"
	libp2p "github.com/libp2p/go-libp2p"
	ma "github.com/multiformats/go-multiaddr"
	core "demochain/core"
)

// makeBasicHost creates a LibP2P host with a random peer ID listening on the
//func MakeBasicHost(listenIP string, listenPort int, priv crypto.PrivKey) (host.Host, error) {
func MakeBasicHost(node *core.Node) (host.Host, error) {

	log.Println("Configurando Peer")

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/%s/tcp/%d", node.GetIP(), node.GetPort())),
		libp2p.Identity(node.GetPrivateKey()),
	}

	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	node.SetTarget(basicHost.ID().Pretty())

	return basicHost, nil
}

func MakeFullAddr(basicHost host.Host) (string) {
	//Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := basicHost.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)

	return fullAddr.String()
}
