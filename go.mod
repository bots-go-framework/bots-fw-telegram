module github.com/bots-go-framework/bots-fw-telegram

go 1.24.3

//replace github.com/bots-go-framework/bots-fw => ../bots-fw

//replace github.com/bots-go-framework/bots-api-telegram => ../bots-api-telegram

//replace github.com/bots-go-framework/bots-fw-store => ../bots-fw-store

require (
	github.com/bots-go-framework/bots-api-telegram v0.14.6
	github.com/bots-go-framework/bots-fw v0.71.25
	github.com/bots-go-framework/bots-fw-store v0.10.0
	github.com/bots-go-framework/bots-fw-telegram-models v0.3.30
	github.com/bots-go-framework/bots-go-core v0.2.3
	github.com/dal-go/dalgo v0.40.2
	github.com/strongo/i18n v0.8.6
	github.com/strongo/logus v0.4.0
	go.uber.org/mock v0.6.0
)

require (
	github.com/RoaringBitmap/roaring v1.9.4 // indirect
	github.com/alexsergivan/transliterator v1.0.1 // indirect
	github.com/bits-and-blooms/bitset v1.24.4 // indirect
	github.com/mschoch/smat v0.2.0 // indirect
	github.com/strongo/analytics v0.2.2 // indirect
	github.com/strongo/random v0.0.1 // indirect
	github.com/strongo/slice v0.3.3 // indirect
	github.com/strongo/strongoapp v0.31.11 // indirect
	github.com/strongo/validation v0.0.7 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
)
