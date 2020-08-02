//     Copyright 2020 Nexus Operator and/or its authors
//
//     This file is part of Nexus Operator.
//
//     Nexus Operator is free software: you can redistribute it and/or modify
//     it under the terms of the GNU General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.
//
//     Nexus Operator is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU General Public License for more details.
//
//     You should have received a copy of the GNU General Public License
//     along with Nexus Operator.  If not, see <https://www.gnu.org/licenses/>.

package server

import "github.com/m88i/aicura/nexus"

var communityMavenProxies = map[string]nexus.MavenProxyRepository{
	"apache":  createDefaultMavenProxy("apache", "https://repo.maven.apache.org/maven2/"),
	"red-hat": createDefaultMavenProxy("red-hat", "https://maven.repository.redhat.com/ga/"),
	"jboss":   createDefaultMavenProxy("jboss", "https://repository.jboss.org/"),
}

const (
	mavenCentralRepoId = "maven-central"
)

// RepositoryOperations describes the public operations in the repository domain for the Nexus instance
type RepositoryOperations interface {
	EnsureCommunityMavenProxies() error
}

type repositoryOperation struct {
	server
}

func repositoryOperations(server *server) RepositoryOperations {
	return &repositoryOperation{server: *server}
}

func (r *repositoryOperation) EnsureCommunityMavenProxies() error {
	if err := r.createCommunityReposIfNotExists(); err != nil {
		return err
	}
	mavenCentral, err := r.nexuscli.MavenGroupRepositoryService.GetRepoByName(mavenCentralRepoId)
	if err != nil {
		return err
	}
	if mavenCentral == nil {
		log.Warnf("Maven Central repository group not found in the server instance, won't add community repos to the group")
		return nil
	}
	var newMembers []string
	for _, member := range mavenCentral.Group.MemberNames {
		if _, ok := communityMavenProxies[member]; !ok {
			newMembers = append(newMembers, member)
		}
	}
	mavenCentral.Group.MemberNames = append(mavenCentral.Group.MemberNames, newMembers...)
	if err := r.nexuscli.MavenGroupRepositoryService.Update(*mavenCentral); err != nil {
		return err
	}
	return nil
}

func (r *repositoryOperation) createCommunityReposIfNotExists() error {
	var repos []nexus.MavenProxyRepository
	for _, repo := range communityMavenProxies {
		fetchRepo, err := r.nexuscli.MavenProxyRepositoryService.GetRepoByName(repo.Name)
		if err != nil {
			return err
		}
		if fetchRepo == nil {
			repos = append(repos, repo)
		}
	}
	if len(repos) > 0 {
		if err := r.nexuscli.MavenProxyRepositoryService.Add(repos...); err != nil {
			return err
		}
	}
	return nil
}

func createDefaultMavenProxy(name, url string) nexus.MavenProxyRepository {
	return nexus.MavenProxyRepository{
		Repository: nexus.Repository{
			URL:    nexus.NewString(url),
			Format: nexus.NewRepositoryFormat(nexus.RepositoryFormatMaven2),
			Name:   name,
			Type:   nexus.NewRepositoryType(nexus.RepositoryTypeProxy),
		},
		Storage: nexus.Storage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: false,
		},
		Maven: nexus.Maven{
			VersionPolicy: nexus.VersionPolicyRelease,
			LayoutPolicy:  nexus.LayoutPolicyPermissive,
		},
	}
}
