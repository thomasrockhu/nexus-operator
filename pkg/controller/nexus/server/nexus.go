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

	"github.com/google/uuid"
	nexusapi "github.com/m88i/aicura/nexus"
	"github.com/m88i/nexus-operator/pkg/apis/apps/v1alpha1"
	"github.com/m88i/nexus-operator/pkg/logger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type server struct {
	nexus     *v1alpha1.Nexus
	k8sclient client.Client
	nexuscli  *nexusapi.Client
	// TODO: add an option to change the user/password of a given client instance in AICURA project
	nexusadmincli *nexusapi.Client
}

const (
	defaultUsername   = "admin"
	defaultPassword   = "admin123"
	operatorUsername  = "nexus-operator"
	operatorEmail     = "nexus-operator@googlegroups.com"
	operatorStatus    = "active"
	operatorName      = "Nexus"
	operatorLastName  = "Operator"
	defaultSource     = "default"
	adminRole         = "nx-admin"
	secretKeyPassword = "password"
	secretKeyUsername = "username"
)

var log = logger.GetLogger("server_operations")

func UpdateServer(nexus *v1alpha1.Nexus, client client.Client) error {
	log.Debug("Initializing server operations in instance %s", nexus.Name)
	s := server{nexus: nexus, k8sclient: client}
	if !nexus.Spec.GenerateRandomAdminPassword && s.isServerReady() {
		internalEndpoint, err := s.getNexusEndpoint()
		if err != nil {
			return err
		}
		if len(internalEndpoint) == 0 {
			log.Warnf("Impossible to resolve endpoint for Nexus instance %s", nexus.Name)
			return nil
		}
		// we need to create another client since we would use the default credentials
		s.nexusadmincli = nexusapi.NewClient(internalEndpoint).WithCredentials(defaultUsername, defaultPassword).Build()
		if err := s.createOperatorUserWithAdminDefaultCredentials(); err != nil {
			return err
		}
		s.nexuscli = nexusapi.
			NewClient(internalEndpoint).
			WithCredentials(defaultUsername, defaultPassword).
			Build()
		if err := s.createDefaultRepos(); err != nil {
			return err
		}
	}
	return nil
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

func (s *server) createOperatorUserWithAdminDefaultCredentials() error {
	credentials, err := s.getOperatorCredentials()
	if err != nil {
		return err
	}
	if credentials == nil {
		log.Infof("Operator credentials not found for %s", s.nexus.Name)
		user, err := s.nexusadmincli.UserService.GetUserByID(operatorUsername)
		if err != nil {
			return err
		}
		if user == nil {
			log.Infof("No operator user for instance %s, creating...", s.nexus.Name)
			secret, err := s.getOrCreateOperatorCredentials()
			if err != nil {
				return err
			}
			newUser := nexusapi.User{
				UserID:    secret.StringData[secretKeyUsername],
				Email:     operatorEmail,
				FirstName: operatorName,
				LastName:  operatorLastName,
				Status:    operatorStatus,
				Source:    defaultSource,
				Roles:     []string{adminRole},
				Password:  secret.StringData[secretKeyPassword],
			}
			if err := s.nexusadmincli.UserService.Add(newUser); err != nil {
				return err
			}
			log.Infof("User %s created with random password", secret.StringData[secretKeyUsername])
		}
		// TODO: maybe updating the existing user's password with our credentials?
		// TODO: AICURA does not have an interface to update an user
	}
	return nil
}

func (s *server) getOperatorCredentials() (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	if err := s.k8sclient.Get(context.TODO(), types.NamespacedName{Name: s.nexus.Name, Namespace: s.nexus.Namespace}, secret); errors.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return secret, nil
}

func (s *server) getOrCreateOperatorCredentials() (*corev1.Secret, error) {
	secret, err := s.getOperatorCredentials()
	if secret == nil {
		password, err := generateRandomPassword()
		if err != nil {
			return nil, err
		}
		secret = &corev1.Secret{
			ObjectMeta: v1.ObjectMeta{
				Name:      s.nexus.Name,
				Namespace: s.nexus.Namespace,
			},
			StringData: map[string]string{
				secretKeyUsername: operatorUsername,
				secretKeyPassword: password,
			},
			Type: corev1.SecretTypeOpaque,
		}
	} else if err != nil {
		return nil, err
	}
	return secret, nil
}

func (s *server) createDefaultRepos() error {
	//TODO: do :)
	//TODO: add the possibility to update a given maven repository group in AICURA
	return nil
}

func generateRandomPassword() (string, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", nil
	}
	return uid.String(), nil
}
