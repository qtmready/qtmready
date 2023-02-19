package github

func (g *Repo) PreCreate() error         { return nil }
func (g *Repo) PreUpdate() error         { return nil }
func (g *Installation) PreCreate() error { return nil }
func (g *Installation) PreUpdate() error { return nil }
