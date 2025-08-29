package scm

// ClientFactory interface for creating SCM clients
type ClientFactory interface {
	CreateClient(providerType, name, url, token string, insecure bool) (Client, error)
}

// MultiClientManager manages multiple SCM clients
type MultiClientManager struct {
	clients []Client
}

func NewMultiClientManager() *MultiClientManager {
	return &MultiClientManager{
		clients: make([]Client, 0),
	}
}

func (m *MultiClientManager) AddClient(client Client) {
	m.clients = append(m.clients, client)
}

func (m *MultiClientManager) GetClients() []Client {
	return m.clients
}

func (m *MultiClientManager) ListAllRepositories() ([]*Repository, error) {
	var allRepos []*Repository

	for _, client := range m.clients {
		repos, err := client.ListAllRepositories()
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
	}

	return allRepos, nil
}

func (m *MultiClientManager) ListRepositoriesInGroup(groupPath string) ([]*Repository, error) {
	var allRepos []*Repository

	for _, client := range m.clients {
		repos, err := client.ListRepositoriesInGroup(groupPath)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
	}

	return allRepos, nil
}
