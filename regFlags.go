package gramework

import "flag"

var flagsToRegister = []Flag{
	{
		Name:        "bind",
		Description: "address to listen",
		Default:     ":80",
	},
}

// AddFlag adds a Flag to flag queue that will be
// parsed if flags wasn't parsed yet
func (app *App) AddFlag(f Flag) {
	if app.flagsQueue == nil {
		app.flagsQueue = make([]Flag, 0)
	}
	app.flagsQueue = append(app.flagsQueue, f)
}

// RegFlags registers current flag queue in flag parser
func (app *App) RegFlags() {
	if app.Flags.values == nil {
		app.Flags.values = make(map[string]Flag)
	}
	app.flagsRegistered = true
	for _, v := range app.flagsQueue {
		app.Flags.values[v.Name] = Flag{
			Name:        v.Name,
			Description: v.Description,
			Default:     v.Default,
			Value:       flag.String(v.Name, v.Default, v.Description),
		}
	}
}
