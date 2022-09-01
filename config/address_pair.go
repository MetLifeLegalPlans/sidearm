package config

func (a *AddressPair) SetDefaults() {
	var (
		hasBind, hasConnect bool = a.Bind != "", a.Connect != ""
	)

	if hasBind && hasConnect {
		return
	}

	if hasBind && !hasConnect {
		a.Connect = a.Bind
	}

	if hasConnect && !hasBind {
		a.Bind = a.Connect
	}
}
