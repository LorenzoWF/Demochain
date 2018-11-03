package network

import(
	"context"
	host "github.com/libp2p/go-libp2p-host"
	"fmt"
	"io"
	"log"
	//"crypto/rand"
	"bytes"
	//mrand "math/rand"
	crypto "github.com/libp2p/go-libp2p-crypto"
	libp2p "github.com/libp2p/go-libp2p"
	ma "github.com/multiformats/go-multiaddr"
)

// makeBasicHost creates a LibP2P host with a random peer ID listening on the
func MakeBasicHost(listenIP string, listenPort int, randseed []byte, cryptographicType int, cryptographicBits int) (host.Host, error) {

	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	//if randseed == 0 { //VERIFICA SE JA FOI INSERIDO UM ENDERECO PELO CMD
//		r = rand.Reader
	r = bytes.NewReader(randseed)
	//} else {
		//r = mrand.New(mrand.NewSource(randseed))
	//}

	// Generate a key pair for this host. We will use it
	// to obtain a valid host ID.
	//TYPE É 0 = RSA, 1, 2, 3 //VER OS OUTROS TIPOS
	priv, _, err := crypto.GenerateKeyPairWithReader(cryptographicType, cryptographicBits, r)

	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/%s/tcp/%d", listenIP, listenPort)),
		libp2p.Identity(priv),
	}

	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := basicHost.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)
	log.Printf("I am %s\n", fullAddr)
	log.Printf("Now run \"go run demochain.go -p %d -d %s\" on a different terminal\n", listenPort+1, fullAddr)

	return basicHost, nil
}
