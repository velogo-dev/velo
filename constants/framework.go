package constants

// Framework represents a UI framework
type Framework string

const (
	React   Framework = "react"
	Vue     Framework = "vue"
	Svelte  Framework = "svelte"
	Angular Framework = "angular"
	Solid   Framework = "solid"
)

// SubFramework represents a specific implementation of a UI framework
type SubFramework struct {
	Parent Framework
	Name   string
}

// React-based frameworks
var (
	NextJS         = SubFramework{Parent: React, Name: "next"}
	Remix          = SubFramework{Parent: React, Name: "remix"}
	CreateReactApp = SubFramework{Parent: React, Name: "create-react-app"}
	Gatsby         = SubFramework{Parent: React, Name: "gatsby"}
	ReactVite      = SubFramework{Parent: React, Name: "vite"}
)

// Vue-based frameworks
var (
	Nuxt    = SubFramework{Parent: Vue, Name: "nuxt"}
	Quasar  = SubFramework{Parent: Vue, Name: "quasar"}
	VueVite = SubFramework{Parent: Vue, Name: "vite"}
)

// Svelte-based frameworks
var (
	SvelteKit  = SubFramework{Parent: Svelte, Name: "sveltekit"}
	SvelteVite = SubFramework{Parent: Svelte, Name: "vite"}
)

// Angular-based frameworks
var (
	AngularUniversal = SubFramework{Parent: Angular, Name: "universal"}
	Nest             = SubFramework{Parent: Angular, Name: "nest"}
)

// Solid-based frameworks
var (
	SolidStart = SubFramework{Parent: Solid, Name: "solid-start"}
	SolidVite  = SubFramework{Parent: Solid, Name: "vite"}
)

// IsReactBased checks if a SubFramework is React-based
func (sf SubFramework) IsReactBased() bool {
	return sf.Parent == React
}

// IsVueBased checks if a SubFramework is Vue-based
func (sf SubFramework) IsVueBased() bool {
	return sf.Parent == Vue
}

// IsSvelteBased checks if a SubFramework is Svelte-based
func (sf SubFramework) IsSvelteBased() bool {
	return sf.Parent == Svelte
}

// GetParentFramework returns the parent framework name
func (sf SubFramework) GetParentFramework() string {
	return string(sf.Parent)
}

// GetFullName returns the full name of the framework (parent+subframework)
func (sf SubFramework) GetFullName() string {
	return string(sf.Parent) + "/" + sf.Name
}
