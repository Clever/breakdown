module github.com/Clever/breakdown

go 1.16

require (
	github.com/Clever/breakdown/gen-go/client v0.0.0-00010101000000-000000000000 // indirect
	github.com/Clever/breakdown/gen-go/models v0.0.0-00010101000000-000000000000
	github.com/Clever/discovery-go v1.8.1 // indirect
	github.com/Clever/go-process-metrics v0.4.0
	github.com/Clever/kayvee-go/v7 v7.6.0
	github.com/Clever/launch-gen v0.0.0-20230222233441-17c275320509
	github.com/Clever/wag v4.1.0+incompatible
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5 // indirect
	github.com/cespare/reflex v0.3.1
	github.com/donovanhide/eventsource v0.0.0-20171031113327-3ed64d21fb0b // indirect
	github.com/get-woke/woke v0.19.0
	github.com/go-errors/errors v1.1.1
	github.com/go-openapi/runtime v0.19.21 // indirect
	github.com/go-openapi/strfmt v0.21.2
	github.com/go-openapi/swag v0.21.1
	github.com/golang/mock v1.6.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonschema v1.2.1-0.20191114132342-001aa27b4d11 // indirect
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.34.0
	go.opentelemetry.io/otel v1.9.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.9.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.9.0
	go.opentelemetry.io/otel/sdk v1.9.0
	go.opentelemetry.io/otel/trace v1.9.0
	golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
)

replace github.com/Clever/breakdown/gen-go/models => ./gen-go/models

replace github.com/Clever/breakdown/gen-go/client => ./gen-go/client
