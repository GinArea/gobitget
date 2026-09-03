package bitgetv3

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/msw-x/moon/ulog"
)

// Offline tests of the WebSocket event classification, driven by a local mock server
// (no network access, no credentials). They pin down how a server-side error event is
// treated, and in particular notice 30033 - the upgrade warning Bitget broadcasts about
// a minute before it resets the connection: it is neither an alert-worthy failure nor a
// login rejection, only the server asking for a reconnect it is about to force anyway.
//
// The assertions deliberately observe the wire and the callbacks instead of Running()
// and Ready(): those read fields written by the connection goroutine, and reading them
// from the test goroutine would be a data race rather than a check.

const mockWaitTimeout = 3 * time.Second

const mockUpgradeMsg = "Service upgrade in progress. Connection reset imminent. Please reconnect."

const mockLoginOkFrame = `{"event":"login","code":0,"msg":"","connId":"mock"}`

func mockErrorFrame(code int, msg string) string {
	return fmt.Sprintf(`{"event":"error","code":%d,"msg":%q}`, code, msg)
}

// mockTickerPush - v3 ticker push shaped for the args wsArgs("ticker", ...) builds
func mockTickerPush(instType, symbol string) string {
	return fmt.Sprintf(`{"action":"snapshot","arg":{"instType":%q,"topic":"ticker","symbol":%q},`+
		`"data":[{"symbol":%q,"bid1Price":"100","ask1Price":"101"}],"ts":1700000000000}`,
		instType, symbol, symbol)
}

// mockCandlePushV2 - legacy v2 candle push shaped for the args wsArgsV2 builds
func mockCandlePushV2(instType, channel, instId string) string {
	return fmt.Sprintf(`{"action":"snapshot","arg":{"instType":%q,"channel":%q,"instId":%q},`+
		`"data":[["1700000000000","1","2","0.5","1.5","10","20","20"]],"ts":1700000000000}`,
		instType, channel, instId)
}

// mockSession - one accepted connection: the script drives it while a reader goroutine
// collects the frames the client sends and answers the keepalive
type mockSession struct {
	conn  *websocket.Conn
	index int
	in    chan string
	done  <-chan struct{}
	mutex sync.Mutex
}

func (o *mockSession) send(frame string) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.conn.WriteMessage(websocket.TextMessage, []byte(frame))
}

func (o *mockSession) read() {
	for {
		_, data, err := o.conn.ReadMessage()
		if err != nil {
			close(o.in)
			return
		}
		s := string(data)
		if s == "ping" {
			o.send("pong")
			continue
		}
		select {
		case o.in <- s:
		default:
		}
	}
}

// wait - block until the client sends a frame containing substr
func (o *mockSession) wait(substr string, timeout time.Duration) bool {
	deadline := time.After(timeout)
	for {
		select {
		case s, ok := <-o.in:
			if !ok {
				return false
			}
			if strings.Contains(s, substr) {
				return true
			}
		case <-deadline:
			return false
		}
	}
}

// hold - keep the connection open until the test tears down; a script that returns
// closes the socket instead, which is how a reset is simulated
func (o *mockSession) hold() {
	select {
	case <-o.done:
	case <-time.After(30 * time.Second):
	}
}

// mockServer - local WS endpoint running the script for every accepted connection
type mockServer struct {
	srv   *httptest.Server
	conns int32
	// connected - the index of every accepted connection, in order
	connected chan int
	done      chan struct{}
}

// newMockServer - the script runs in the connection's own goroutine, so it must never
// touch *testing.T: it reports through channels the test goroutine reads
func newMockServer(t *testing.T, script func(*mockSession)) *mockServer {
	o := &mockServer{
		connected: make(chan int, 8),
		done:      make(chan struct{}),
	}
	var upgrader websocket.Upgrader
	o.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		s := &mockSession{
			conn:  conn,
			index: int(atomic.AddInt32(&o.conns, 1)),
			in:    make(chan string, 32),
			done:  o.done,
		}
		select {
		case o.connected <- s.index:
		default:
		}
		go s.read()
		script(s)
	}))
	// cleanups run last-registered-first: the client is closed first (it registers its
	// own cleanup later), then the held scripts are released, and only then is the
	// server shut down - httptest waits for the handlers still in flight
	t.Cleanup(o.srv.Close)
	t.Cleanup(func() { close(o.done) })
	return o
}

func (o *mockServer) url() string {
	return "ws" + strings.TrimPrefix(o.srv.URL, "http")
}

// waitConnected - block until at least n connections have been accepted
func (o *mockServer) waitConnected(n int, timeout time.Duration) bool {
	deadline := time.After(timeout)
	for {
		select {
		case i := <-o.connected:
			if i >= n {
				return true
			}
		case <-deadline:
			return false
		}
	}
}

// mockLogSpy - collects what the client logs; the production alert hook forwards warning
// and above, so the level of a line decides whether it becomes a Telegram alert
type mockLogSpy struct {
	mutex    sync.Mutex
	messages []ulog.Message
}

func newMockLogSpy(t *testing.T) *mockLogSpy {
	o := new(mockLogSpy)
	ulog.SetHook(func(m ulog.Message) {
		o.mutex.Lock()
		defer o.mutex.Unlock()
		o.messages = append(o.messages, m)
	})
	t.Cleanup(func() { ulog.SetHook(nil) })
	return o
}

// waitLevel - level of the first logged line containing substr; the line is written from
// the connection goroutine, hence the polling
func (o *mockLogSpy) waitLevel(substr string, timeout time.Duration) (ulog.Level, bool) {
	deadline := time.Now().Add(timeout)
	for {
		o.mutex.Lock()
		for _, m := range o.messages {
			if strings.Contains(m.Text, substr) {
				level := m.Level
				o.mutex.Unlock()
				return level, true
			}
		}
		o.mutex.Unlock()
		if time.Now().After(deadline) {
			return 0, false
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func mockPublicClient(t *testing.T, m *mockServer) *WsPublic {
	c := NewWsPublic()
	c.WithLog(ulog.New("mock-public"))
	c.Transport().Base = m.url()
	c.Transport().Path = "ws"
	t.Cleanup(c.Close)
	return c
}

func mockPublicV2Client(t *testing.T, m *mockServer) *WsPublicV2 {
	c := NewWsPublicV2()
	c.WithLog(ulog.New("mock-public-v2"))
	c.Transport().Base = m.url()
	c.Transport().Path = "ws"
	t.Cleanup(c.Close)
	return c
}

func mockPrivateClient(t *testing.T, m *mockServer) *WsPrivate {
	c := NewWsPrivate("mock-key", "mock-secret", "mock-passphrase")
	c.WithLog(ulog.New("mock-private"))
	c.Transport().Base = m.url()
	c.Transport().Path = "ws"
	t.Cleanup(c.Close)
	return c
}

// TestWsServiceUpgradeDuringLogin - the notice can land in the private handshake window,
// and that window is exactly where it used to be mistaken for a credential rejection:
// the rejection path stops the reconnect loop for good, so a single mistimed notice left
// the private socket dead until the process restarted. The login must still complete and
// the registered subscriptions must reach the server
func TestWsServiceUpgradeDuringLogin(t *testing.T) {
	tests := []struct {
		name string
		// beforeLogin - push the notice on connect, before the client's login request,
		// instead of while that request is in flight
		beforeLogin bool
	}{
		{
			name:        "notice before login request",
			beforeLogin: true,
		},
		{
			name:        "notice after login request",
			beforeLogin: false,
		},
	}
	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			subscribed := make(chan struct{}, 1)
			m := newMockServer(t, func(s *mockSession) {
				if v.beforeLogin {
					s.send(mockErrorFrame(30033, mockUpgradeMsg))
					if !s.wait("login", mockWaitTimeout) {
						return
					}
				} else {
					if !s.wait("login", mockWaitTimeout) {
						return
					}
					s.send(mockErrorFrame(30033, mockUpgradeMsg))
				}
				s.send(mockLoginOkFrame)
				if s.wait("position", mockWaitTimeout) {
					subscribed <- struct{}{}
				}
				s.hold()
			})
			c := mockPrivateClient(t, m)
			loginFailed := make(chan struct{}, 1)
			c.WithOnLoginFailed(func() {
				loginFailed <- struct{}{}
			})
			c.Position().Subscribe(func(Topic[[]WsPosition]) {})
			c.Run()
			select {
			case <-subscribed:
			case <-loginFailed:
				t.Fatal("the service upgrade notice was taken for a login rejection")
			case <-time.After(mockWaitTimeout * 2):
				t.Fatal("no subscription reached the server after the login")
			}
		})
	}
}

// TestWsLoginWindowRejection - regression guard for the branch the fix has to leave
// intact: a genuine rejection still stops the reconnect loop instead of retrying the
// same bad credentials forever
func TestWsLoginWindowRejection(t *testing.T) {
	tests := []struct {
		name  string
		frame string
	}{
		{
			// Bitget reports the rejection as an error event, not as a login ack
			name:  "error event",
			frame: mockErrorFrame(30015, "Invalid sign"),
		},
		{
			name:  "login ack",
			frame: `{"event":"login","code":30015,"msg":"Invalid sign"}`,
		},
	}
	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			m := newMockServer(t, func(s *mockSession) {
				if !s.wait("login", mockWaitTimeout) {
					return
				}
				s.send(v.frame)
				s.hold()
			})
			c := mockPrivateClient(t, m)
			loginFailed := make(chan struct{}, 1)
			c.WithOnLoginFailed(func() {
				loginFailed <- struct{}{}
			})
			c.Position().Subscribe(func(Topic[[]WsPosition]) {})
			c.Run()
			select {
			case <-loginFailed:
			case <-time.After(mockWaitTimeout):
				t.Fatal("a rejected login did not report a failure")
			}
			// a re-dial would show up as a second accepted connection
			if m.waitConnected(2, 2*time.Second) {
				t.Fatal("the client re-dialled after a rejected login")
			}
		})
	}
}

// TestWsServiceUpgradePublicReconnect - on the public socket the notice changes nothing:
// the subscription keeps working until the server resets the connection, and the
// reconnect replays it. The level matters as much as the behaviour - anything from
// warning up is forwarded to the alert channel, once per socket, and the public stack
// runs one socket per symbol
func TestWsServiceUpgradePublicReconnect(t *testing.T) {
	const symbol = "BTCUSDT"
	instType := strings.ToLower(string(UsdtFutures))
	resubscribed := make(chan struct{}, 1)
	m := newMockServer(t, func(s *mockSession) {
		if !s.wait("ticker", mockWaitTimeout) {
			return
		}
		if s.index == 1 {
			s.send(mockErrorFrame(30033, mockUpgradeMsg))
			// a push after the notice proves the subscription survived it
			s.send(mockTickerPush(instType, symbol))
			time.Sleep(300 * time.Millisecond)
			// returning closes the socket: the reset the notice announced
			return
		}
		resubscribed <- struct{}{}
		s.hold()
	})
	spy := newMockLogSpy(t)
	c := mockPublicClient(t, m)
	pushes := make(chan struct{}, 4)
	errors := make(chan int64, 4)
	c.WithOnError(func(r WsResponse) {
		errors <- r.Code.Value()
	})
	c.Ticker(UsdtFutures, symbol).Subscribe(func(Topic[[]WsTicker]) {
		pushes <- struct{}{}
	})
	c.Run()
	select {
	case <-pushes:
	case <-time.After(mockWaitTimeout):
		t.Fatal("the ticker push after the notice was not routed")
	}
	select {
	case <-resubscribed:
	case <-time.After(mockWaitTimeout * 2):
		t.Fatal("the client did not resubscribe after the reset")
	}
	level, ok := spy.waitLevel("30033", time.Second)
	if !ok {
		t.Fatal("the notice was not logged at all")
	}
	if level >= ulog.LevelWarning {
		t.Fatalf("the notice is logged at %v: warning and above reaches the alert channel", level)
	}
	select {
	case code := <-errors:
		t.Fatalf("the notice was reported through onError: %d", code)
	default:
	}
}

// TestWsServiceUpgradePublicV2 - the legacy candle socket talks to the same server and
// gets the same notice, so it must treat it the same way
func TestWsServiceUpgradePublicV2(t *testing.T) {
	const symbol = "BTCUSDT"
	instType := string(Spot)
	channel := "candle" + string(Interval1m)
	m := newMockServer(t, func(s *mockSession) {
		if !s.wait(channel, mockWaitTimeout) {
			return
		}
		s.send(mockErrorFrame(30033, mockUpgradeMsg))
		s.send(mockCandlePushV2(instType, channel, symbol))
		s.hold()
	})
	spy := newMockLogSpy(t)
	c := mockPublicV2Client(t, m)
	pushes := make(chan struct{}, 4)
	errors := make(chan int64, 4)
	c.WithOnError(func(r WsResponseV2) {
		errors <- r.Code.Value()
	})
	c.Candle(Spot, symbol, Interval1m).Subscribe(func(TopicV2[[]WsCandleV2]) {
		pushes <- struct{}{}
	})
	c.Run()
	select {
	case <-pushes:
	case <-time.After(mockWaitTimeout):
		t.Fatal("the candle push after the notice was not routed")
	}
	level, ok := spy.waitLevel("30033", time.Second)
	if !ok {
		t.Fatal("the notice was not logged at all")
	}
	if level >= ulog.LevelWarning {
		t.Fatalf("the notice is logged at %v: warning and above reaches the alert channel", level)
	}
	select {
	case code := <-errors:
		t.Fatalf("the notice was reported through onError: %d", code)
	default:
	}
}

// TestWsErrorEventRouting - regression guard: every other error code keeps its old
// treatment, an error-level line plus the onError callback
func TestWsErrorEventRouting(t *testing.T) {
	t.Run("v3", func(t *testing.T) {
		m := newMockServer(t, func(s *mockSession) {
			if !s.wait("ticker", mockWaitTimeout) {
				return
			}
			// 30001: the subscribed channel/symbol does not exist
			s.send(mockErrorFrame(30001, "BTCUSDT doesn't exist"))
			s.hold()
		})
		spy := newMockLogSpy(t)
		c := mockPublicClient(t, m)
		errors := make(chan int64, 4)
		c.WithOnError(func(r WsResponse) {
			errors <- r.Code.Value()
		})
		c.Ticker(UsdtFutures, "BTCUSDT").Subscribe(func(Topic[[]WsTicker]) {})
		c.Run()
		select {
		case code := <-errors:
			if code != 30001 {
				t.Fatalf("onError got code %d, want 30001", code)
			}
		case <-time.After(mockWaitTimeout):
			t.Fatal("onError was not called")
		}
		level, ok := spy.waitLevel("30001", time.Second)
		if !ok {
			t.Fatal("the error was not logged")
		}
		if level != ulog.LevelError {
			t.Fatalf("the error is logged at %v, want error", level)
		}
	})
	t.Run("v2", func(t *testing.T) {
		channel := "candle" + string(Interval1m)
		m := newMockServer(t, func(s *mockSession) {
			if !s.wait(channel, mockWaitTimeout) {
				return
			}
			// 30016: invalid subscription parameter
			s.send(mockErrorFrame(30016, "Param error"))
			s.hold()
		})
		spy := newMockLogSpy(t)
		c := mockPublicV2Client(t, m)
		errors := make(chan int64, 4)
		c.WithOnError(func(r WsResponseV2) {
			errors <- r.Code.Value()
		})
		c.Candle(Spot, "BTCUSDT", Interval1m).Subscribe(func(TopicV2[[]WsCandleV2]) {})
		c.Run()
		select {
		case code := <-errors:
			if code != 30016 {
				t.Fatalf("onError got code %d, want 30016", code)
			}
		case <-time.After(mockWaitTimeout):
			t.Fatal("onError was not called")
		}
		level, ok := spy.waitLevel("30016", time.Second)
		if !ok {
			t.Fatal("the error was not logged")
		}
		if level != ulog.LevelError {
			t.Fatalf("the error is logged at %v, want error", level)
		}
	})
}
