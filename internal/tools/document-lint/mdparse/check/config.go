package check

import "strings"

type set map[string]struct{}

func (s set) Exists(key string) bool {
	_, ok := s[key]
	return ok
}

var config = struct {
	dryRun       bool
	resource     set
	rp           set
	skipResource set
	skipRP       set
}{}

func newSet(val string) set {
	if val == "" {
		return nil
	}
	vals := strings.Split(val, ",")
	res := make(set, len(vals))
	for _, name := range vals {
		res[name] = struct{}{}
	}
	return res
}

func SetConfig(resource, rp, skipResource, skipRP string, dryRun bool) {
	config.resource = newSet(resource)
	config.rp = newSet(rp)
	config.skipResource = newSet(skipResource)
	config.skipRP = newSet(skipRP)
	config.dryRun = dryRun
}

func SkipResource(name string) bool {
	if config.skipResource.Exists(name) {
		return true
	}
	// if specifies a list of resources and given name not in the list
	if len(config.resource) > 0 && !config.resource.Exists(name) {
		return true
	}
	return isSkipResource(name)
}

func SkipRP(name string) bool {
	if config.skipRP.Exists(name) {
		return true
	}
	if len(config.rp) > 0 && !config.rp.Exists(name) {
		return true
	}
	return false
}
