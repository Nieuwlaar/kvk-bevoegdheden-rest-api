module github.com/kvk-innovatie/kvk-bevoegdheden-rest-api

go 1.18

replace github.com/kvk-innovatie/kvk-bevoegdheden => ../kvk-bevoegdheden

require (
	github.com/go-chi/chi/v5 v5.0.7
	github.com/kvk-innovatie/kvk-bevoegdheden v0.0.0-00010101000000-000000000000
	github.com/unrolled/render v1.4.1
)

require (
	github.com/beevik/etree v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/google/uuid v1.3.0 // indirect
	golang.org/x/sys v0.0.0-20210525143221-35b2ab0089ea // indirect
)
