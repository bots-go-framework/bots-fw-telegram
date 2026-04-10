module github.com/bots-go-framework/bots-fw-telegram

go 1.24.3

//replace github.com/bots-go-framework/bots-fw => ../bots-fw

//replace github.com/bots-go-framework/bots-api-telegram => ../bots-api-telegram

//replace github.com/bots-go-framework/bots-fw-store => ../bots-fw-store

require (
	github.com/bots-go-framework/bots-api-telegram v0.14.8
	github.com/bots-go-framework/bots-fw v0.71.42
	github.com/bots-go-framework/bots-fw-store v0.10.1
	github.com/bots-go-framework/bots-fw-telegram-models v0.3.44
	github.com/bots-go-framework/bots-go-core v0.2.4
	github.com/dal-go/dalgo v0.41.10
	github.com/strongo/i18n v0.8.9
	github.com/strongo/logus v0.4.1
	go.uber.org/mock v0.6.0
)

require (
	github.com/RoaringBitmap/roaring v1.9.4 // indirect
	github.com/RoaringBitmap/roaring/v2 v2.16.0 // indirect
	github.com/alexsergivan/transliterator v1.0.1 // indirect
	github.com/bits-and-blooms/bitset v1.24.4 // indirect
	github.com/mschoch/smat v0.2.0 // indirect
	github.com/strongo/analytics v0.2.4 // indirect
	github.com/strongo/random v0.0.1 // indirect
	github.com/strongo/slice v0.3.4 // indirect
	github.com/strongo/strongoapp v0.31.22 // indirect
	github.com/strongo/validation v0.0.8 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
)
