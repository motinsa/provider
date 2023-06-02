package esxi

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: esxi.Provider})
}
