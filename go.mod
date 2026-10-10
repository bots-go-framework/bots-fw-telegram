module github.com/bots-go-framework/bots-fw-telegram

go 1.26.0

toolchain go1.27.2

//replace github.com/bots-go-framework/bots-api-telegram => ../bots-api-telegram

require (
	github.com/bots-go-framework/bots-api-telegram v0.15.24
	github.com/bots-go-framework/bots-fw v0.77.9
	github.com/bots-go-framework/bots-fw-store v0.14.2
	github.com/bots-go-framework/bots-fw-telegram-models v0.3.91
	github.com/bots-go-framework/bots-go-core v0.3.6
	github.com/strongo/i18n v0.8.23
	github.com/strongo/logus v0.4.7
	go.uber.org/mock v0.6.0
)

require (
	github.com/alexsergivan/transliterator v1.0.1 // indirect
	github.com/strongo/analytics v0.2.5 // indirect
	github.com/strongo/random v0.0.3 // indirect
	github.com/strongo/slice v0.3.12 // indirect
	github.com/strongo/strongoapp v0.31.67 // indirect
	github.com/strongo/validation v0.0.15 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
)
