package onvifcam

import (
	"context"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func getConfig() *Config {
	return &Config{
		Addr:     os.Getenv("ONVIF_ADDR"),
		Username: os.Getenv("ONVIF_USERNAME"),
		Password: os.Getenv("ONVIF_PASSWORD"),
		Profile:  os.Getenv("ONVIF_PROFILE"),
		Version:  os.Getenv("ONVIF_VERSION"),
	}
}

func TestNew(t *testing.T) {
	cfg := getConfig()

	cam := New(cfg, &http.Client{})
	require.Nil(t, cam.d)
	require.Equal(t, cfg.Username, cam.cfg.Username)
	require.NotNil(t, cam.httpClient)
}

func TestInit(t *testing.T) {
	cfg := getConfig()
	cam := New(cfg, &http.Client{})

	err := cam.Init(context.TODO())
	require.NoError(t, err)
	require.NotNil(t, cam.d)
}

func writeFrame(f []byte) error {
	filename := "/tmp/frame.jpeg"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fs.FileMode(0600))
	if err != nil {
		return err
	}

	n, err := file.Write(f)
	if err != nil {
		return err
	}

	if n != len(f) {
		return errors.New("cannot save")
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func TestGetFrame(t *testing.T) {
	cfg := getConfig()
	cam := New(cfg, &http.Client{})

	err := cam.Init(context.TODO())
	require.NoError(t, err)

	frame, err := cam.GetSnapshot(context.Background())
	require.NoError(t, err)
	require.Greater(t, len(frame), 0)
	writeFrame(frame)
}

func TestGetStreamURI(t *testing.T) {
	cfg := getConfig()
	cam := New(cfg, &http.Client{})

	err := cam.Init(context.TODO())
	require.NoError(t, err)

	uri, err := cam.GetStreamURI(context.Background())
	if cfg.Version != V18 {
		require.True(t, errors.Is(err, ErrVersion))

		return
	}

	require.NoError(t, err)
	require.NotEmpty(t, uri)
}

func TestSubscribe(t *testing.T) {
	cfg := getConfig()
	cam := New(cfg, &http.Client{})

	err := cam.Init(context.TODO())
	require.NoError(t, err)

	until := time.Now().Add(10 * time.Minute).Format(time.RFC3339) // same than "PT10M" but with client's time

	ep, err := cam.Subscribe(context.Background(), "http://192.168.1.79:3030", TopicMotionAlarm, until)
	if cfg.Version != V18 {
		require.True(t, errors.Is(err, ErrVersion))

		return
	}

	require.NoError(t, err)
	require.Contains(t, ep, "onvif/Events/SubManager")

	ep, err = cam.Subscribe(context.Background(), "http://192.168.1.79:3030", TopicMotionAlarm, "PT10M")
	require.NoError(t, err)
	require.Contains(t, ep, "onvif/Events/SubManager")
}

func TestUnmarshalEventMessage(t *testing.T) {
	// data contains event message XML
	data := `<?xml version="1.0" encoding="UTF-8"?>
	<env:Envelope xmlns:env="http://www.w3.org/2003/05/soap-envelope"
		xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
		xmlns:xs="http://www.w3.org/2001/XMLSchema"
		xmlns:tt="http://www.onvif.org/ver10/schema"
		xmlns:tds="http://www.onvif.org/ver10/device/wsdl"
		xmlns:trt="http://www.onvif.org/ver10/media/wsdl"
		xmlns:timg="http://www.onvif.org/ver20/imaging/wsdl"
		xmlns:tev="http://www.onvif.org/ver10/events/wsdl"
		xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl"
		xmlns:tan="http://www.onvif.org/ver20/analytics/wsdl"
		xmlns:tst="http://www.onvif.org/ver10/storage/wsdl"
		xmlns:ter="http://www.onvif.org/ver10/error"
		xmlns:dn="http://www.onvif.org/ver10/network/wsdl"
		xmlns:tns1="http://www.onvif.org/ver10/topics"
		xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl"
		xmlns:wsdl="http://schemas.xmlsoap.org/wsdl"
		xmlns:wsoap12="http://schemas.xmlsoap.org/wsdl/soap12"
		xmlns:http="http://schemas.xmlsoap.org/wsdl/http"
		xmlns:d="http://schemas.xmlsoap.org/ws/2005/04/discovery"
		xmlns:wsadis="http://schemas.xmlsoap.org/ws/2004/08/addressing"
		xmlns:wsnt="http://docs.oasis-open.org/wsn/b-2"
		xmlns:wsa="http://www.w3.org/2005/08/addressing"
		xmlns:wstop="http://docs.oasis-open.org/wsn/t-1"
		xmlns:wsrf-bf="http://docs.oasis-open.org/wsrf/bf-2"
		xmlns:wsntw="http://docs.oasis-open.org/wsn/bw-2"
		xmlns:wsrf-rw="http://docs.oasis-open.org/wsrf/rw-2"
		xmlns:wsaw="http://www.w3.org/2006/05/addressing/wsdl"
		xmlns:wsrf-r="http://docs.oasis-open.org/wsrf/r-2"
		xmlns:trc="http://www.onvif.org/ver10/recording/wsdl"
		xmlns:tse="http://www.onvif.org/ver10/search/wsdl"
		xmlns:trp="http://www.onvif.org/ver10/replay/wsdl"
		xmlns:tnsn="http://www.eventextension.com/2011/event/topics"
		xmlns:extwsd="http://www.onvifext.com/onvif/ext/ver10/wsdl"
		xmlns:extxsd="http://www.onvifext.com/onvif/ext/ver10/schema"
		xmlns:tas="http://www.onvif.org/ver10/advancedsecurity/wsdl"
		xmlns:tr2="http://www.onvif.org/ver20/media/wsdl"
		xmlns:axt="http://www.onvif.org/ver20/analytics">
		<env:Header>
			<wsa:To env:mustUnderstand="true">http://192.168.1.79:3030/</wsa:To>
			<wsa:Action>http://docs.oasis-open.org/wsn/bw-2/NotificationConsumer/Notify</wsa:Action>
		</env:Header>
		<env:Body>
			<wsnt:Notify>
				<wsnt:NotificationMessage>
					<wsnt:Topic Dialect="http://www.onvif.org/ver10/tev/topicExpression/ConcreteSet">tns1:VideoSource/MotionAlarm</wsnt:Topic>
					<wsnt:Message>
						<tt:Message UtcTime="2022-09-09T20:49:30Z" PropertyOperation="Changed">
							<tt:Source>
								<tt:SimpleItem Name="Source" Value="VideoSource_1"/>
							</tt:Source>
							<tt:Data>
								<tt:SimpleItem Name="State" Value="false"/>
							</tt:Data>
						</tt:Message>
					</wsnt:Message>
				</wsnt:NotificationMessage>
			</wsnt:Notify>
		</env:Body>
	</env:Envelope>`

	var msg EventMessage

	err := UnmarshalEventMessage([]byte(data), &msg)
	require.NoError(t, err)
	require.Equal(t, "http://192.168.1.79:3030/", msg.Header.To)
	require.Equal(t, "http://docs.oasis-open.org/wsn/bw-2/NotificationConsumer/Notify", msg.Header.Action)
	require.Equal(t, "http://www.onvif.org/ver10/tev/topicExpression/ConcreteSet", string(msg.Body.Notify.NotificationMessage.Topic.Dialect))
	require.Equal(t, "tns1:VideoSource/MotionAlarm", string(msg.Body.Notify.NotificationMessage.Topic.TopicKinds))
	require.Equal(t, "2022-09-09T20:49:30Z", string(msg.Body.Notify.NotificationMessage.Message.Message.UtcTime))
	require.Equal(t, EventChanged, string(msg.Body.Notify.NotificationMessage.Message.Message.PropertyOperation))
	require.Equal(t, "Source", msg.Body.Notify.NotificationMessage.Message.Message.Source.Name)
	require.Equal(t, "VideoSource_1", string(msg.Body.Notify.NotificationMessage.Message.Message.Source.Value))
	require.Equal(t, "State", msg.Body.Notify.NotificationMessage.Message.Message.Data.Name)
	require.Equal(t, "false", string(msg.Body.Notify.NotificationMessage.Message.Message.Data.Value))
}
