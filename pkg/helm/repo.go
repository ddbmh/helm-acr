// Modifications copyright (C) 2019 Alibaba Group Holding Limited / Yuning Xie (xyn1016@gmail.com)

package helm

import (
	"fmt"
	helm_env "k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/helmpath"
	"k8s.io/helm/pkg/repo"
	urllib "net/url"
	"os"
	"strings"
)

type (
	// Repo represents a collection of parameters for chart repository
	Repo struct {
		*repo.Entry
	}
)

// GetRepoByURL returns repository by url
func GetRepoByURL(url string) (*Repo, error) {
	r, err := repoFile()
	if err != nil {
		return nil, err
	}
	entry, exists := findRepoEntryByURL(url, r)
	if !exists {
		return nil, fmt.Errorf("no repo url %q found", url)
	}
	return &Repo{entry}, nil
}

// GetRepoByName returns repository by name
func GetRepoByName(name string) (*Repo, error) {
	r, err := repoFile()
	if err != nil {
		return nil, err
	}
	entry, exists := findRepoEntry(name, r)
	if !exists {
		return nil, fmt.Errorf("no repo named %q found", name)
	}
	return &Repo{entry}, nil
}

// TempRepoFromURL builds a temporary Repo from a given URL
func TempRepoFromURL(url string) (*Repo, error) {
	u, err := urllib.Parse(url)
	if err != nil {
		return nil, err
	}
	entry := &repo.Entry{}
	if u.User != nil {
		// remove the username/password section from URL
		pass, _ := u.User.Password()
		entry.URL = strings.Split(url, "://")[0] + "://" + strings.Split(url, fmt.Sprintf("%s@", pass))[1]
		entry.Username = u.User.Username()
		entry.Password = pass
	} else {
		entry.URL = url
	}
	return &Repo{entry}, nil
}

func repoFile() (*repo.RepoFile, error) {
	home := helmHome()
	return repo.LoadRepositoriesFile(home.RepositoryFile())
}

func helmHome() helmpath.Home {
	var helmHomePath string
	if v, ok := os.LookupEnv("HELM_HOME"); ok {
		helmHomePath = v
	} else {
		helmHomePath = helm_env.DefaultHelmHome
	}
	return helmpath.Home(helmHomePath)
}

func findRepoEntry(name string, r *repo.RepoFile) (*repo.Entry, bool) {
	var entry *repo.Entry
	exists := false
	for _, re := range r.Repositories {
		if re.Name == name {
			entry = re
			exists = true
			break
		}
	}
	return entry, exists
}

func findRepoEntryByURL(url string, r *repo.RepoFile) (*repo.Entry, bool) {
	var entry *repo.Entry
	exists := false
	for _, re := range r.Repositories {
		if strings.TrimSuffix(re.URL, "/") == strings.TrimSuffix(url, "/") {
			entry = re
			exists = true
			break
		}
	}
	return entry, exists
}
