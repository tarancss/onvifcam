module github.com/tarancss/onvifcam

go 1.18

require (
	github.com/gorilla/mux v1.8.0
	github.com/stretchr/testify v1.4.0
	github.com/use-go/onvif v0.0.0-20230921065222-217f1231c56f
)

require (
	github.com/beevik/etree v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/elgs/gostrgen v0.0.0-20161222160715-9d61ae07eeae // indirect
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/juju/errors v0.0.0-20220331221717-b38fca44723b // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rs/zerolog v1.26.1 // indirect
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d // indirect
	golang.org/x/sys v0.0.0-20210809222454-d867a43fc93e // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

// replace github.com/use-go/onvif v0.0.0-20230921065222-217f1231c56f => github.com/tarancss/onvif v0.0.0-20231123182348-22058f5239ad
replace github.com/use-go/onvif v0.0.0-20230921065222-217f1231c56f => ../onvif
