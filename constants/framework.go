package constants

// Library represents a UI library
type Library string

const (
	React   Library = "react"
	Vue     Library = "vue"
	Svelte  Library = "svelte"
	Angular Library = "angular"
	Solid   Library = "solid"
)

// Framework represents a specific implementation of a UI framework
type Framework struct {
	Parent Library
	Name   string
}

// React-based frameworks
var (
	NextJS         = Framework{Parent: React, Name: "next"}
	Remix          = Framework{Parent: React, Name: "remix"}
	CreateReactApp = Framework{Parent: React, Name: "create-react-app"}
	Gatsby         = Framework{Parent: React, Name: "gatsby"}
	ReactVite      = Framework{Parent: React, Name: "vite"}
)

// Vue-based frameworks
var (
	Nuxt    = Framework{Parent: Vue, Name: "nuxt"}
	Quasar  = Framework{Parent: Vue, Name: "quasar"}
	VueVite = Framework{Parent: Vue, Name: "vite"}
)

// Svelte-based frameworks
var (
	SvelteKit  = Framework{Parent: Svelte, Name: "sveltekit"}
	SvelteVite = Framework{Parent: Svelte, Name: "vite"}
)

// Angular-based frameworks
var (
	AngularUniversal = Framework{Parent: Angular, Name: "universal"}
	Nest             = Framework{Parent: Angular, Name: "nest"}
)

// Solid-based frameworks
var (
	SolidStart = Framework{Parent: Solid, Name: "solid-start"}
	SolidVite  = Framework{Parent: Solid, Name: "vite"}
)

// AvailableLibraries defines the supported frontend libraries
var AvailableLibraries = []Library{
	React,
	Vue,
	Svelte,
	Angular,
	Solid,
}

// LibraryFrameworks maps libraries to their popular frameworks
var LibraryFrameworks = map[Library][]Framework{
	React: {
		CreateReactApp,
		NextJS,
		Remix,
		ReactVite,
	},
	Vue: {
		Nuxt,
		Quasar,
		VueVite,
	},
	Svelte: {
		SvelteKit,
		SvelteVite,
	},
	Angular: {
		AngularUniversal,
		Nest,
	},
	Solid: {
		SolidStart,
		SolidVite,
	},
}
