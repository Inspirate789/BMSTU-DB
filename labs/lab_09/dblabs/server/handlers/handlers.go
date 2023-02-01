package handlers

import (
	"dblabs/database"
)

type handler interface {
	Title() string
	Execute(*database.Database, []string) (string, error)
}

// Handlers map[labNumber]handlers
var Handlers = map[int][]handler{
	0: {
		handler00{},
		handler01{},
		handler02{},
	},
	6: {
		handler60{},
		handler61{},
		handler62{},
		handler63{},
		handler64{},
		handler65{},
		handler66{},
		handler67{},
		handler68{},
		handler69{},
	},
	7: {
		handler70{},
		handler71{},
		handler72{},
		handler73{},
		handler74{},
		handler75{},
		handler76{},
		handler77{},
		handler78{},
		handler79{},
		handler710{},
		handler711{},
		handler712{},
		handler713{},
		handler714{},
	},
	9: {
		handler90{},
		handler91{},
		handler92{},
		handler93{},
	},
}
