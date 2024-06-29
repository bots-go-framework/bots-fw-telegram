module github.com/bots-go-framework/bots-fw-telegram

go 1.22.3

toolchain go1.22.4

//replace github.com/bots-go-framework/bots-fw-store => ../bots-fw-store

//replace github.com/bots-go-framework/bots-fw => ../bots-fw

require (
	github.com/bots-go-framework/bots-api-telegram v0.4.4
	github.com/bots-go-framework/bots-fw v0.25.2
	github.com/bots-go-framework/bots-fw-store v0.4.0
	github.com/bots-go-framework/bots-fw-telegram-models v0.1.6
	github.com/dal-go/dalgo v0.12.1
	github.com/pquerna/ffjson v0.0.0-20190930134022-aa0246cd15f7
	github.com/strongo/i18n v0.0.4
	github.com/strongo/logus v0.0.0-20240628225821-04cf45b5968f
)

require (
	github.com/alexsergivan/transliterator v1.0.1 // indirect
	github.com/bots-go-framework/bots-go-core v0.0.2 // indirect
	github.com/strongo/gamp v0.0.1 // indirect
	github.com/strongo/random v0.0.1 // indirect
	github.com/strongo/strongoapp v0.18.3 // indirect
	github.com/strongo/validation v0.0.6 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
)
