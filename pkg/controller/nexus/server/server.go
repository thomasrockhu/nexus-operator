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

import (
	"github.com/m88i/nexus-operator/pkg/apis/apps/v1alpha1"
	"github.com/m88i/nexus-operator/pkg/logger"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type server struct {
	nexus  *v1alpha1.Nexus
	client client.Client
}

var log = logger.GetLogger("server_operations")

func UpdateServer(nexus *v1alpha1.Nexus, client client.Client) error {
	log.Debug("Initializing server operations in instance %s", nexus.Name)
	s := server{
		nexus:  nexus,
		client: client,
	}
	if !nexus.Spec.GenerateRandomAdminPassword && s.isUp() {
		if err := s.createOperatorUser(); err != nil {
			return err
		}
		if err := s.createDefaultRepos(); err != nil {
			return err
		}
	}
	return nil
}

// isUp checks if the given Nexus instance is ready to receive requests
func (s *server) isUp() bool {
	return s.nexus.Status.DeploymentStatus.AvailableReplicas > 0
}

func (s *server) createOperatorUser() error {
	return nil
}

func (s *server) createDefaultRepos() error {
	return nil
}
