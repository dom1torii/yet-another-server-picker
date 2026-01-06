package presets

type Preset struct {
	Name string
	Pops map[string]struct{}
}

var Presets = map[string]Preset{
	"eu": {
		Name: "Europe",
		Pops: map[string]struct{}{
			"ams":  {},
			"fra":  {},
			"fsn":  {},
			"hel":  {},
			"lhr":  {},
			"mad":  {},
			"par":  {},
			"sto":  {},
			"vie":  {},
			"waw":  {},
			"sto2": {},
			"ams4": {},
		},
	},
	"eu-west": {
		Name: "Western Europe",
		Pops: map[string]struct{}{
			"ams":  {},
			"fra":  {},
			"fsn":  {},
			"lhr":  {},
			"mad":  {},
			"par":  {},
			"vie":  {},
			"ams4": {},
		},
	},
	"de": {
		Name: "Germany",
		Pops: map[string]struct{}{
			"fra": {},
			"fsn": {},
		},
	},
	"nl": {
		Name: "Netherlands",
		Pops: map[string]struct{}{
			"ams":  {},
			"ams4": {},
		},
	},
	"eu-east": {
		Name: "Eastern Europe",
		Pops: map[string]struct{}{
			"hel":  {},
			"sto":  {},
			"waw":  {},
			"sto2": {},
		},
	},
	"sw": {
		Name: "Sweden",
		Pops: map[string]struct{}{
			"sto":  {},
			"sto2": {},
		},
	},
	"na": {
		Name: "North America",
		Pops: map[string]struct{}{
			"atl": {},
			"dfw": {},
			"eat": {},
			"iad": {},
			"lax": {},
			"ord": {},
			"sea": {},
		},
	},
	"sa": {
		Name: "South America",
		Pops: map[string]struct{}{
			"eze": {},
			"gru": {},
			"lim": {},
			"scl": {},
		},
	},
	"me": {
		Name: "Middle East",
		Pops: map[string]struct{}{
			"dxb": {},
		},
	},
	"cn-pw": {
		Name: "China (PW)",
		Pops: map[string]struct{}{
			"ctu":  {},
			"pek":  {},
			"pvg":  {},
			"pwg":  {},
			"pwj":  {},
			"pwu":  {},
			"pww":  {},
			"pwz":  {},
			"sha":  {},
			"shb":  {},
			"tgd":  {},
			"ctum": {},
			"pekm": {},
			"pvgm": {},
			"tgdm": {},
			"ctut": {},
			"pekt": {},
			"pvgt": {},
			"tgdt": {},
			"ctuu": {},
			"peku": {},
			"pvgu": {},
			"tgdu": {},
		},
	},
	"as": {
		Name: "Asia",
		Pops: map[string]struct{}{
			"hkg":  {},
			"seo":  {},
			"sgp":  {},
			"tyo":  {},
			"bom2": {},
			"maa2": {},
			"hkg4": {},
		},
	},
	"in": {
		Name: "India",
		Pops: map[string]struct{}{
			"bom2": {},
			"maa2": {},
		},
	},
	"hk": {
		Name: "Hong Kong",
		Pops: map[string]struct{}{
			"hkg":  {},
			"hkg4": {},
		},
	},
	"af": {
		Name: "Africa",
		Pops: map[string]struct{}{
			"jnb": {},
		},
	},
	"au": {
		Name: "Australia",
		Pops: map[string]struct{}{
			"syd": {},
		},
	},
}
