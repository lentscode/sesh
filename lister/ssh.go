package lister

import "github.com/joshmedeski/sesh/v2/model"

func listSSH(l *RealLister) (model.SeshSessions, error) {
	sessions := model.SeshSessions{
		Directory:    make(model.SeshSessionMap),
		OrderedIndex: make([]string, 0, len(l.config.SSHConfigs)),
	}
	for _, config := range l.config.SSHConfigs {
		key := "ssh:" + config.Name
		sessions.OrderedIndex = append(sessions.OrderedIndex, key)
		sessions.Directory[key] = config.Session()
	}
	return sessions, nil
}
