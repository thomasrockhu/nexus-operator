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
	"context"
	"fmt"

	nexusapi "github.com/m88i/aicura/nexus"
	"github.com/m88i/nexus-operator/pkg/apis/apps/v1alpha1"
	"github.com/m88i/nexus-operator/pkg/logger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type server struct {
	nexus     *v1alpha1.Nexus
	k8sclient client.Client
	nexuscli  *nexusapi.Client
	status    *OperationStatus
}

type OperationStatus struct {
	ServerReady                  bool
	OperatorUserCreated          bool
	CommunityRepositoriesCreated bool
	MavenCentralUpdated          bool
}

var log = logger.GetLogger("server_operations")

func handleServerOperations(nexus *v1alpha1.Nexus, client client.Client, nexusApiBuilder func(url, user, pass string) *nexusapi.Client) (OperationStatus, error) {
	log.Debug("Initializing server operations in instance %s", nexus.Name)
	s := server{nexus: nexus, k8sclient: client, status: &OperationStatus{}}
	if !nexus.Spec.GenerateRandomAdminPassword && s.isServerReady() {
		s.status.ServerReady = true
		internalEndpoint, err := s.getNexusEndpoint()
		if err != nil {
			return *s.status, err
		}
		if len(internalEndpoint) == 0 {
			log.Warnf("Impossible to resolve endpoint for Nexus instance %s", nexus.Name)
			return *s.status, nil
		}
		s.nexuscli = nexusApiBuilder(internalEndpoint, defaultAdminUsername, defaultAdminPassword)

		if err := userOperations(&s).EnsureOperatorUser(); err != nil {
			return *s.status, err
		}
		if err := repositoryOperations(&s).EnsureCommunityMavenProxies(); err != nil {
			return *s.status, err
		}
	}
	return *s.status, nil
}

// HandleServerOperations makes all required operations in the Nexus server side, such as creating the operator user
func HandleServerOperations(nexus *v1alpha1.Nexus, client client.Client) (OperationStatus, error) {
	return handleServerOperations(nexus, client, func(url, user, pass string) *nexusapi.Client {
		return nexusapi.NewClient(url).WithCredentials(user, pass).Build()
	})
}

func (s *server) getNexusEndpoint() (string, error) {
	svc := &corev1.Service{}
	if err := s.k8sclient.Get(context.TODO(), types.NamespacedName{Name: s.nexus.Name, Namespace: s.nexus.Namespace}, svc); errors.IsNotFound(err) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s:%s", svc.Name, svc.Spec.Ports[0].TargetPort.String()), nil
}

// isServerReady checks if the given Nexus instance is ready to receive requests
func (s *server) isServerReady() bool {
	return s.nexus.Status.DeploymentStatus.AvailableReplicas > 0
}
