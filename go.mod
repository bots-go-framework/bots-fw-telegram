module github.com/bots-go-framework/bots-fw-telegram

go 1.26.0

toolchain go1.27.4

//replace github.com/bots-go-framework/bots-api-telegram => ../bots-api-telegram

require (
	github.com/bots-go-framework/bots-api-telegram v0.15.27
	github.com/bots-go-framework/bots-fw v0.77.9
	github.com/bots-go-framework/bots-fw-store v0.14.2
	github.com/bots-go-framework/bots-fw-telegram-models v0.3.98
	github.com/bots-go-framework/bots-go-core v0.3.8
	github.com/strongo/i18n v0.8.25
	github.com/strongo/logus v0.4.9
	go.uber.org/mock v0.6.0
)

require (
	github.com/alexsergivan/transliterator v1.0.1 // indirect
	github.com/strongo/analytics v0.2.5 // indirect
	github.com/strongo/random v0.0.3 // indirect
	github.com/strongo/slice v0.3.15 // indirect
	github.com/strongo/strongoapp v0.31.73 // indirect
	github.com/strongo/validation v0.0.17 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
)
