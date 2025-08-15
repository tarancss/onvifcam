package onvifcam

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/use-go/onvif"
	"github.com/use-go/onvif/event"
	"github.com/use-go/onvif/media"
	sevent "github.com/use-go/onvif/sdk/event"
	smedia "github.com/use-go/onvif/sdk/media"
	"github.com/use-go/onvif/xsd"
	xonvif "github.com/use-go/onvif/xsd/onvif"
)

const (
	TopicMotionAlarm = "tns1:VideoSource/MotionAlarm"
	// profile examples
	mainProfileOnvif18 = "Profile_1"
	mainProfileOnvif23 = "profile1"

	V18 = "18" // ONVIF v18
	V23 = "23" // ONVIF v23
)

var (
	ErrFailedNew             = errors.New("failed to set new device")
	ErrNoURIFrame            = errors.New("failed to get URI for snapshot")
	ErrNoURIStream           = errors.New("failed to get URI for stream")
	ErrSubscribe             = errors.New("failed to subscribe")
	ErrUnmarshalEventMessage = errors.New("failed to unmarshal event message")
	ErrVersion               = errors.New("method does not exist in ONVIF version")
)

var (
	EventChanged     = "Changed"
	EventInitialized = "Initialized"
)

type Config struct {
	Addr     string
	Username string
	Password string
	Profile  string
	Version  string
}

type Onvifcam struct {
	cfg         *Config
	d           *onvif.Device
	httpClient  *http.Client
	snapshotURI *url.URL
}

// New returns a new bare ONVIF device using basic authentication.
// httpClient is used also by the ONVIF device implementation. It is set to a default client if not provided.
func New(cfg *Config, httpClient *http.Client) *Onvifcam {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return &Onvifcam{
		cfg:        cfg,
		d:          nil,
		httpClient: httpClient,
	}
}

// Init connects to the device using basic authentication and sets it up.
// A context is currently unused (can be set to nil) but passed for future improvement of go-onvif module.
func (c *Onvifcam) Init(_ context.Context) error {
	d, err := onvif.NewDevice(onvif.DeviceParams{Xaddr: c.cfg.Addr, Username: c.cfg.Username, Password: c.cfg.Password,
		HttpClient: c.httpClient})
	if err != nil {
		return fmt.Errorf("%s: %w", ErrFailedNew, err)
	}

	c.d = d

	return nil
}

// GetSnapshot returns an image frame (jpeg) from the camera.
func (c *Onvifcam) GetSnapshot(ctx context.Context) ([]byte, error) {
	if c.snapshotURI == nil {
		// only first time
		req := media.GetSnapshotUri{
			XMLName:      "",
			ProfileToken: xonvif.ReferenceToken(c.cfg.Profile),
		}
		r, err := smedia.Call_GetSnapshotUri(ctx, c.d, req)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", ErrNoURIFrame, err)
		}

		var uri = string(r.MediaUri.Uri)

		if uri == "" {
			return nil, ErrNoURIFrame
		}

		urlURI, err := url.Parse(uri)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", ErrNoURIFrame, err)
		}

		c.snapshotURI = urlURI
	}

	httpReq := &http.Request{
		Method: http.MethodGet,
		URL:    c.snapshotURI,
		Header: http.Header{},
	}

	httpReq.SetBasicAuth(c.cfg.Username, c.cfg.Password)

	respHTTP, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrNoURIFrame, err)
	}

	defer respHTTP.Body.Close()
	frame, err := io.ReadAll(respHTTP.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrNoURIFrame, err)
	}

	return frame, nil
}

// GetStreamURI returns an rstp URI.
// For a rtsp client see https://pkg.go.dev/github.com/aler9/gortsplib#section-readme.
func (c *Onvifcam) GetStreamURI(ctx context.Context) (string, error) {
	if c.cfg.Version != V18 {
		return "", ErrVersion
	}

	reqStream := media.GetStreamUri{
		XMLName: "",
		StreamSetup: xonvif.StreamSetup{
			Stream: "RTP", // "RTP_unicast",
			Transport: xonvif.Transport{
				Protocol: "TCP",
			},
		},
		ProfileToken: xonvif.ReferenceToken(c.cfg.Profile),
	}

	rStream, err := smedia.Call_GetStreamUri(ctx, c.d, reqStream)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrNoURIStream, err)
	}

	uri := string(rStream.MediaUri.Uri)
	if uri == "" {
		return "", ErrNoURIStream
	}

	return uri, nil
}

// Subscribe returns a subscription reference address.
// Events are sent via http POST / request to the given http server listening in addr.
func (c *Onvifcam) Subscribe(ctx context.Context, addr, topic, dateTimeOrDuration string) (string, error) {
	if c.cfg.Version != V18 {
		return "", ErrVersion
	}

	req := event.Subscribe{
		ConsumerReference: event.EndpointReferenceType{
			Address:             event.AttributedURIType(addr),
			ReferenceParameters: event.ReferenceParametersType{},
			Metadata:            event.MetadataType{},
		},
		Filter: event.FilterType{
			TopicExpression: event.TopicExpressionType{
				Dialect:    "http://docs.oasis-open.org/wsn/t-1/TopicExpression/Concrete",
				TopicKinds: xsd.String(topic),
			},
		},
		SubscriptionPolicy: event.SubscriptionPolicy{
			ChangedOnly: true, // we don't need Initialized or Deleted events
		},
		InitialTerminationTime: event.AbsoluteOrRelativeTimeType(dateTimeOrDuration),
	}

	res, err := sevent.Call_Subscribe(ctx, c.d, req)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrSubscribe, err)
	}

	if res.SubscriptionReference.Address == "" {
		return "", ErrSubscribe
	}

	return string(res.SubscriptionReference.Address), nil
}

type EventMessage struct {
	Header HeaderXML
	Body   BodyXML
}

type HeaderXML struct {
	To     string
	Action string
}

type BodyXML struct {
	Notify NotifyXML
}

type NotifyXML struct {
	NotificationMessage event.NotificationMessage
}

// UnmarshalEventMessage is normally used by the handler of the http server listening for events.
func UnmarshalEventMessage(data []byte, r *EventMessage) error {
	err := xml.Unmarshal(data, r)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrUnmarshalEventMessage, err)
	}

	return nil
}
