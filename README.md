# onvifcam
Client for ONVIF compatible IP cameras.

[![GoDoc](https://godoc.org/github.com/tarancss/onvifcam?status.svg)](https://godoc.org/github.com/tarancss/onvifcam)
[![Go Report Card](https://goreportcard.com/badge/github.com/tarancss/onvifcam)](https://goreportcard.com/report/github.com/tarancss/onvifcam)

Contributions are welcome.

## Usage
This package provides simplified client to interact with an ONVIF compatible IP camera.

Once the device is created with `New`, it must be initialized with `Init`. Then, `GetSnapshot` and `GetStreamURI` can be
called and used straight away.

The `Subscribe` method creates an event subscription to the device for motion detection. The camera sends events via an
HTTP POST / request. The package provides the function `UnmarshalEventMessage` to obtain the event's data. This function
can be used by the HTTP server's handler function.

## Tests
The tests require an ONVIF compatible IP camera to run.

