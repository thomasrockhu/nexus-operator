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

	"github.com/google/uuid"
	"github.com/m88i/aicura/nexus"
	"github.com/m88i/nexus-operator/pkg/framework"
	corev1 "k8s.io/api/core/v1"
)

const (
	defaultAdminUsername = "admin"
	defaultAdminPassword = "admin123"
	operatorUsername     = "nexus-operator"
	operatorEmail        = "nexus-operator@googlegroups.com"
	operatorStatus       = "active"
	operatorName         = "Nexus"
	operatorLastName     = "Operator"
	defaultSource        = "default"
	adminRole            = "nx-admin"

	secretKeyPassword = "server-user-password"
	secretKeyUsername = "server-user-username"
)

func (s *server) ensureOperatorUser() (*nexus.User, error) {
	// TODO: open an issue to handle access to a custom admin credentials to be used by the operator
	s.nexuscli.SetCredentials(defaultAdminUsername, defaultAdminPassword)
	user, err := s.nexuscli.UserService.GetUserByID(operatorUsername)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}
	user, err = createOperatorUserInstance()
	if err != nil {
		return nil, err
	}
	if err := s.nexuscli.UserService.Add(*user); err != nil {
		return nil, err
	}
	if err := s.storeOperatorUserCredentials(user); err != nil {
		//  TODO: in case of an error here, we should remove the user from the Nexus database. Edge case, but an user could manually add the credentials later to the secret with a manually created user for us.
		return nil, err
	}
	return user, nil
}

func (s *server) storeOperatorUserCredentials(user *nexus.User) error {
	secret := &corev1.Secret{}
	if err := framework.Fetch(s.k8sclient, framework.Key(s.nexus), secret); err != nil {
		return err
	}
	if secret.StringData == nil {
		secret.StringData = make(map[string]string)
	}
	secret.StringData[secretKeyPassword] = user.Password
	secret.StringData[secretKeyUsername] = user.UserID
	if err := s.k8sclient.Update(context.TODO(), secret); err != nil {
		return err
	}
	return nil
}

func (s *server) getOperatorUserCredentials() (user, password string, err error) {
	secret := &corev1.Secret{}
	if err := framework.Fetch(s.k8sclient, framework.Key(s.nexus), secret); err != nil {
		return "", "", err
	}
	return secret.StringData[secretKeyUsername], secret.StringData[secretKeyPassword], nil
}

func createOperatorUserInstance() (*nexus.User, error) {
	password, err := generateRandomPassword()
	if err != nil {
		return nil, err
	}
	return &nexus.User{
		Email:     operatorEmail,
		Roles:     []string{adminRole},
		FirstName: operatorName,
		LastName:  operatorLastName,
		Password:  password,
		Source:    defaultSource,
		Status:    operatorStatus,
		UserID:    operatorUsername,
	}, nil
}

func generateRandomPassword() (string, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", nil
	}
	return uid.String(), nil
}
