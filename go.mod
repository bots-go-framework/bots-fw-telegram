module github.com/bots-go-framework/bots-fw-telegram

go 1.24.3

//replace github.com/bots-go-framework/bots-fw => ../bots-fw

//replace github.com/bots-go-framework/bots-api-telegram => ../bots-api-telegram

//replace github.com/bots-go-framework/bots-fw-store => ../bots-fw-store

require (
	github.com/bots-go-framework/bots-api-telegram v0.13.0
	github.com/bots-go-framework/bots-fw v0.63.0
	github.com/bots-go-framework/bots-fw-store v0.10.0
	github.com/bots-go-framework/bots-fw-telegram-models v0.3.25
	github.com/dal-go/dalgo v0.23.0
	github.com/strongo/i18n v0.8.2
	github.com/strongo/logus v0.2.1
	go.uber.org/mock v0.5.2
)

require (
	github.com/alexsergivan/transliterator v1.0.1 // indirect
	github.com/bots-go-framework/bots-go-core v0.0.3 // indirect
	github.com/strongo/analytics v0.0.11 // indirect
	github.com/strongo/random v0.0.1 // indirect
	github.com/strongo/slice v0.3.1 // indirect
	github.com/strongo/strongoapp v0.31.6 // indirect
	github.com/strongo/validation v0.0.7 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
)
