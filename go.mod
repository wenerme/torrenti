module github.com/wenerme/torrenti

go 1.18

replace github.com/gocolly/colly/v2 => ./modules/colly

require (
	github.com/adrg/xdg v0.4.0
	github.com/blugelabs/bluge v0.1.9
	github.com/bodgit/sevenzip v1.1.1
	github.com/caarlos0/env/v6 v6.9.1
	github.com/dgraph-io/badger/v3 v3.2103.2
	github.com/dustin/go-humanize v1.0.0
	github.com/glebarez/go-sqlite v1.15.1
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/httplog v0.2.4
	github.com/gocolly/colly/v2 v2.0.0-00010101000000-000000000000
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.0
	github.com/jackc/pgx/v4 v4.15.0
	github.com/longbridgeapp/opencc v0.1.7
	github.com/mitchellh/mapstructure v1.1.2
	github.com/multiformats/go-multihash v0.1.0
	github.com/nwaples/rardecode/v2 v2.0.0-beta.2
	github.com/oklog/run v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.12.1
	github.com/rs/zerolog v1.26.1
	github.com/samber/lo v1.11.0
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/gjson v1.14.0
	github.com/urfave/cli/v2 v2.4.0
	github.com/xgfone/bt v0.4.1
	go.uber.org/fx v1.17.1
	go.uber.org/multierr v1.5.0
	golang.org/x/exp v0.0.0-20220325121720-054d8573a5d8
	golang.org/x/text v0.3.7
	google.golang.org/genproto v0.0.0-20220317150908-0efb43f6373e
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	gorm.io/datatypes v1.0.6
	gorm.io/driver/postgres v1.3.1
	gorm.io/driver/sqlite v1.3.1
	gorm.io/gorm v1.23.3
)

require (
	github.com/PuerkitoBio/goquery v1.8.0 // indirect
	github.com/RoaringBitmap/roaring v0.9.1 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/antchfx/htmlquery v1.2.4 // indirect
	github.com/antchfx/xmlquery v1.3.10 // indirect
	github.com/antchfx/xpath v1.2.0 // indirect
	github.com/axiomhq/hyperloglog v0.0.0-20191112132149-a4c4c47bc57f // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bits-and-blooms/bitset v1.2.0 // indirect
	github.com/blevesearch/go-porterstemmer v1.0.3 // indirect
	github.com/blevesearch/mmap-go v1.0.2 // indirect
	github.com/blevesearch/segment v0.9.0 // indirect
	github.com/blevesearch/snowballstem v0.9.0 // indirect
	github.com/blevesearch/vellum v1.0.5 // indirect
	github.com/blugelabs/bluge_segment_api v0.2.0 // indirect
	github.com/blugelabs/ice v0.2.0 // indirect
	github.com/bodgit/plumbing v1.1.0 // indirect
	github.com/bodgit/windows v1.0.0 // indirect
	github.com/caio/go-tdigest v3.1.0+incompatible // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/connesc/cipherio v0.2.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgraph-io/ristretto v0.1.0 // indirect
	github.com/dgryski/go-metro v0.0.0-20180109044635-280f6062b5bc // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/flatbuffers v1.12.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.11.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.2.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.10.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/kennygrant/sanitize v1.2.4 // indirect
	github.com/klauspost/compress v1.15.1 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/liuzl/cedar-go v0.0.0-20170805034717-80a9c64b256d // indirect
	github.com/liuzl/da v0.0.0-20180704015230-14771aad5b1d // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mattn/go-sqlite3 v1.14.12 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/mschoch/smat v0.2.0 // indirect
	github.com/multiformats/go-varint v0.0.6 // indirect
	github.com/nlnwa/whatwg-url v0.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20200410134404-eec4a21b6bb0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/saintfish/chardet v0.0.0-20120816061221-3af4cd4741ca // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/temoto/robotstxt v1.1.2 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/willf/bitset v1.1.10 // indirect
	go.opencensus.io v0.22.5 // indirect
	go.uber.org/atomic v1.6.0 // indirect
	go.uber.org/dig v1.14.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	go4.org v0.0.0-20200411211856-f5505b9728dd // indirect
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gorm.io/driver/mysql v1.3.2 // indirect
	lukechampine.com/blake3 v1.1.6 // indirect
	modernc.org/libc v1.14.12 // indirect
	modernc.org/mathutil v1.4.1 // indirect
	modernc.org/memory v1.0.7 // indirect
	modernc.org/sqlite v1.15.2 // indirect
)
