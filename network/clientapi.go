package network

import (
	"log"
	"time"
)

type ClientOpts struct {
	//default ClientDefaultTimeout
	Timeout time.Duration

	Addr string

	Room, Role string

	//Specific self disconnect (server lost). may be needed later , but in general use Pause
	OnReconnect  func()
	OnDisconnect func()

	//Any reason's pause of game process (disconnect self, disconnect other, loading new state self or other
	//Specific reason may be getted by PauseReason()
	OnPause   func()
	OnUnpause func()

	OnCommonSend func() []byte
	OnCommonRecv func(data []byte, readOwnPart bool)

	OnStateChanged func(wanted string)

	//async, must close result chan then done
	OnGetStateData func([]byte)

	OnCommand func(command string)
}

type PauseReason struct {
	PingLost   bool
	IsFull     bool
	IsCoherent bool
	CurState   string
	WantState  string
}

func (c *Client) recalcPauseReason() {
	c.prmu.Lock()
	c.pr = PauseReason{
		PingLost:   c.pingLost,
		IsFull:     c.isFull,
		IsCoherent: c.isCoherent,
		CurState:   c.curState,
		WantState:  c.wantState,
	}
	c.prmu.Unlock()
}

func (c *Client) PauseReason() PauseReason {
	c.prmu.RLock()
	defer c.prmu.RUnlock()

	return c.pr
}

func (c *Client) RequestNewState(wanted string) {
	if c.wantState != c.curState {
		log.Println("client is already changing state")
	}
	//_, err := c.doReq(POST, statePattern, []byte(wanted))

	c.sendCommand(COMMAND_REQUESTSTATE, wanted)
}

func (c *Client) SendCommand(command string) {
	c.sendCommand(COMMAND_CLIENT, command)
}

func (c *Client) Start() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.started {
		return
	}

	c.started = true
	go clientPing(c)
}
