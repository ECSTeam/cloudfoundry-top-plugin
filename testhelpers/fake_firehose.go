package testhelpers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"

	"encoding/binary"
	"strings"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/nu7hatch/gouuid"
)

type FakeFirehose struct {
	server *httptest.Server
	lock   sync.Mutex

	AppMode bool
	AppName string

	validToken string

	lastAuthorization string
	requested         bool

	events         []events.Envelope
	closeMessage   []byte
	stayAlive      bool
	subscriptionID string
	wg             sync.WaitGroup
}

func NewFakeFirehose(validToken string) *FakeFirehose {
	return &FakeFirehose{
		validToken:   validToken,
		closeMessage: websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	}
}

func NewFakeFirehoseInAppMode(validToken, appName string) *FakeFirehose {
	return &FakeFirehose{
		AppMode:      true,
		AppName:      appName,
		validToken:   validToken,
		closeMessage: websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	}
}

func (f *FakeFirehose) Start() {
	f.server = httptest.NewUnstartedServer(f)
	f.server.Start()
}

func (f *FakeFirehose) Close() {
	f.server.Close()
}

func (f *FakeFirehose) URL() string {
	return fmt.Sprintf("ws://%s", f.server.Listener.Addr().String())
}

func (f *FakeFirehose) LastAuthorization() string {
	f.lock.Lock()
	defer f.lock.Unlock()
	return f.lastAuthorization
}

func (f *FakeFirehose) Requested() bool {
	f.lock.Lock()
	defer f.lock.Unlock()
	return f.requested
}

func (f *FakeFirehose) SendEvent(eventType events.Envelope_EventType, message string) {
	envelope := events.Envelope{
		Origin:     proto.String("origin"),
		Timestamp:  proto.Int64(1000000000),
		Deployment: proto.String("deployment-name"),
		Job:        proto.String("doppler"),
	}

	switch eventType {
	case events.Envelope_LogMessage:
		envelope.EventType = events.Envelope_LogMessage.Enum()
		envelope.LogMessage = &events.LogMessage{
			Message:     []byte(message),
			MessageType: events.LogMessage_OUT.Enum(),
			Timestamp:   proto.Int64(1000000000),
		}
	case events.Envelope_ValueMetric:
		envelope.EventType = events.Envelope_ValueMetric.Enum()
		envelope.ValueMetric = &events.ValueMetric{
			Name:  proto.String(message),
			Value: proto.Float64(42),
			Unit:  proto.String("unit"),
		}
	case events.Envelope_CounterEvent:
		envelope.EventType = events.Envelope_CounterEvent.Enum()
		envelope.CounterEvent = &events.CounterEvent{
			Name:  proto.String(message),
			Delta: proto.Uint64(42),
		}
	case events.Envelope_ContainerMetric:
		envelope.EventType = events.Envelope_ContainerMetric.Enum()
		envelope.ContainerMetric = &events.ContainerMetric{
			ApplicationId: proto.String(message),
			InstanceIndex: proto.Int32(1),
			CpuPercentage: proto.Float64(1),
			MemoryBytes:   proto.Uint64(1),
			DiskBytes:     proto.Uint64(1),
		}
	case events.Envelope_Error:
		envelope.EventType = events.Envelope_Error.Enum()
		envelope.Error = &events.Error{
			Source:  proto.String("source"),
			Code:    proto.Int32(404),
			Message: proto.String(message),
		}
	case events.Envelope_HttpStart:
		envelope.EventType = events.Envelope_HttpStart.Enum()
		uuid, _ := uuid.NewV4()
		envelope.HttpStart = &events.HttpStart{
			Timestamp:     proto.Int64(12),
			RequestId:     NewUUID(uuid),
			PeerType:      events.PeerType_Client.Enum(),
			Method:        events.Method_GET.Enum(),
			Uri:           proto.String("some uri"),
			RemoteAddress: proto.String("some address"),
			UserAgent:     proto.String(message),
		}
	case events.Envelope_HttpStop:
		envelope.EventType = events.Envelope_HttpStop.Enum()
		uuid, _ := uuid.NewV4()
		envelope.HttpStop = &events.HttpStop{
			Timestamp:     proto.Int64(12),
			Uri:           proto.String("http://stop.example.com"),
			RequestId:     NewUUID(uuid),
			PeerType:      events.PeerType_Client.Enum(),
			StatusCode:    proto.Int32(404),
			ContentLength: proto.Int64(98475189),
		}
	case events.Envelope_HttpStartStop:
		envelope.EventType = events.Envelope_HttpStartStop.Enum()
		uuid, _ := uuid.NewV4()
		envelope.HttpStartStop = &events.HttpStartStop{
			StartTimestamp: proto.Int64(1234),
			StopTimestamp:  proto.Int64(5555),
			RequestId:      NewUUID(uuid),
			PeerType:       events.PeerType_Server.Enum(),
			Method:         events.Method_GET.Enum(),
			Uri:            proto.String("http://startstop.example.com"),
			RemoteAddress:  proto.String("http://startstop.example.com"),
			UserAgent:      proto.String("test"),
			StatusCode:     proto.Int32(1234),
			ContentLength:  proto.Int64(5678),
			ApplicationId:  NewUUID(uuid),
		}

	}

	f.addEvent(envelope)
}

func (f *FakeFirehose) addEvent(event events.Envelope) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.events = append(f.events, event)
}

func (f *FakeFirehose) SetCloseMessage(message []byte) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.closeMessage = make([]byte, len(message))
	copy(f.closeMessage, message)
}

func (f *FakeFirehose) KeepConnectionAlive() {
	f.wg.Add(1)
}

func (f *FakeFirehose) CloseAliveConnection() {
	f.wg.Done()
}

func (f *FakeFirehose) SubscriptionID() string {
	return f.subscriptionID
}

func (f *FakeFirehose) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	f.lock.Lock()
	defer f.lock.Unlock()

	if f.AppMode && r.URL.String() != fmt.Sprintf("/apps/%s/stream", f.AppName) {
		log.Printf("App not found: %s", f.AppName)
		rw.WriteHeader(404)
		r.Body.Close()
		return
	}

	f.lastAuthorization = r.Header.Get("Authorization")
	f.requested = true
	f.subscriptionID = strings.Split(r.URL.String(), "/")[2]
	if f.lastAuthorization != f.validToken {
		log.Printf("Bad token passed to firehose: %s", f.lastAuthorization)
		rw.WriteHeader(403)
		r.Body.Close()
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(*http.Request) bool { return true },
	}

	ws, _ := upgrader.Upgrade(rw, r, nil)

	defer ws.Close()
	defer ws.WriteControl(websocket.CloseMessage, f.closeMessage, time.Time{})

	for _, envelope := range f.events {
		buffer, _ := proto.Marshal(&envelope)
		err := ws.WriteMessage(websocket.BinaryMessage, buffer)
		if err != nil {
			panic(err)
		}
	}
	f.wg.Wait()
}

func NewUUID(id *uuid.UUID) *events.UUID {
	return &events.UUID{Low: proto.Uint64(binary.LittleEndian.Uint64(id[:8])), High: proto.Uint64(binary.LittleEndian.Uint64(id[8:]))}
}
