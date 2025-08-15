package onvifcam

import (
	"encoding/xml"
	"errors"
	"fmt"
	"time"
)

var (
	ErrAlarm error = errors.New("error unmarshaling alarm")
	ErrBeat  error = errors.New("error unmarshaling beat")
)

type AlarmStatusInfo struct {
	MotionAlarm     bool `xml:"motionAlarm"`
	PerimeterAlarm  bool `xml:"perimeterAlarm"`
	TripwireAlarm   bool `xml:"tripwireAlarm"`
	OscAlarm        bool `xml:"oscAlarm"`
	SceneChange     bool `xml:"sceneChange"`
	ClarityAbnormal bool `xml:"clarityAbnormal"`
	ColorAbnormal   bool `xml:"colorAbnormal"`
}

type DeviceInfo struct {
	Name       string `xml:"deviceName"`
	No         string `xml:"deviceNo."`
	Sn         string `xml:"sn"`
	IPAddress  string `xml:"ipAddress"`
	MacAddress string `xml:"macAddress"`
}

type Beat struct {
	Datatime   string     `xml:"dataTime"`
	DeviceInfo DeviceInfo `xml:"deviceInfo"`
}

type Alarm struct {
	AlarmStatus AlarmStatusInfo `xml:"alarmStatusInfo"`
	Datatime    string          `xml:"dataTime"`
	DeviceInfo  DeviceInfo      `xml:"deviceInfo"`
}

func ParseDatetime(d string) (time.Time, error) {
	return time.Parse(time.DateTime, d)
}

func IsBeat(data string) bool {
	return len(data) > 4 && data[0:4] != "POST"
}

func BeatUnmarshal(data string) (*Beat, error) {
	var b Beat

	err := xml.Unmarshal([]byte(data), &b)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrBeat, err)
	}

	return &b, nil
}

func IsAlarm(data string) bool {
	return len(data) > 15 && data[0:15] == "POST /SendAlarm"
}

func AlarmUnmarshal(data string) (*Alarm, error) {
	var a Alarm

	err := xml.Unmarshal([]byte(data), &a)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrAlarm, err)
	}

	return &a, nil
}
