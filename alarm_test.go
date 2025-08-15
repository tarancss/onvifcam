package onvifcam

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBeat_Unmarshal(t *testing.T) {
	var beat string = `<?xml version="1.0" encoding="UTF-8"?>
	<config version="1.7" xmlns="http://www.ipc.com/ver10">
	<dataTime><![CDATA[2025-06-24 14:02:17]]></dataTime>
	<deviceInfo>
	<deviceName><![CDATA[cam2]]></deviceName>
	<deviceNo.><![CDATA[1]]></deviceNo.>
	<sn><![CDATA[I2483098GAA3]]></sn>
	<ipAddress><![CDATA[192.168.1.64]]></ipAddress>
	<macAddress><![CDATA[58:5b:69:3f:24:83]]></macAddress>
	</deviceInfo>
	</config>`

	require.True(t, IsBeat(beat))
	require.False(t, IsAlarm(beat))

	b, err := BeatUnmarshal(beat)
	require.NoError(t, err)

	require.Equal(t, "2025-06-24 14:02:17", b.Datatime)

	_, err = ParseDatetime(b.Datatime)
	require.NoError(t, err)

	testDeviceInfo(t, &b.DeviceInfo, "cam2", "1", "I2483098GAA3", "192.168.1.64", "58:5b:69:3f:24:83")
}

func TestAlarm_Unmarshal(t *testing.T) {
	var alarm string = `POST /SendAlarmStatus HTTP/1.1
Host: 192.168.2.79

<?xml version="1.0" encoding="UTF-8"?>
<config version="1.7" xmlns="http://www.ipc.com/ver10">
<alarmStatusInfo>
<motionAlarm type="boolean" id="1">true</motionAlarm>
<perimeterAlarm type="boolean" id="1">false</perimeterAlarm>
<tripwireAlarm type="boolean" id="1">false</tripwireAlarm>
<oscAlarm type="boolean" id="1">false</oscAlarm>
<sceneChange type="boolean" id="1">false</sceneChange>
<clarityAbnormal type="boolean" id="1">false</clarityAbnormal>
<colorAbnormal type="boolean" id="1">false</colorAbnormal>
</alarmStatusInfo>
<dataTime><![CDATA[2025-06-24 14:01:22]]></dataTime>
<deviceInfo>
<deviceName><![CDATA[cam2]]></deviceName>
<deviceNo.><![CDATA[1]]></deviceNo.>
<sn><![CDATA[I2483098GAA3]]></sn>
<ipAddress><![CDATA[192.168.1.64]]></ipAddress>
<macAddress><![CDATA[58:5b:69:3f:24:83]]></macAddress>
</deviceInfo>
</config>`

	require.False(t, IsBeat(alarm))
	require.True(t, IsAlarm(alarm))

	a, err := AlarmUnmarshal(alarm)
	require.NoError(t, err)

	require.Equal(t, true, a.AlarmStatus.MotionAlarm)
	require.Equal(t, false, a.AlarmStatus.PerimeterAlarm)
	require.Equal(t, false, a.AlarmStatus.TripwireAlarm)
	require.Equal(t, false, a.AlarmStatus.OscAlarm)
	require.Equal(t, false, a.AlarmStatus.SceneChange)
	require.Equal(t, false, a.AlarmStatus.ClarityAbnormal)
	require.Equal(t, false, a.AlarmStatus.ColorAbnormal)

	testDeviceInfo(t, &a.DeviceInfo, "cam2", "1", "I2483098GAA3", "192.168.1.64", "58:5b:69:3f:24:83")
}

func testDeviceInfo(t *testing.T, di *DeviceInfo, cam, no, sn, ipAddress, macAddress string) {
	require.Equal(t, cam, di.Name)
	require.Equal(t, no, di.No)
	require.Equal(t, sn, di.Sn)
	require.Equal(t, ipAddress, di.IPAddress)
	require.Equal(t, macAddress, di.MacAddress)
}
