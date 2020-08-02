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

type UserOperations interface {
	EnsureOperatorUser() error
}

type userOperation struct {
	server
}

func userOperations(server *server) UserOperations {
	return &userOperation{server: *server}
}

func (u *userOperation) EnsureOperatorUser() error {
	if u.nexus.Spec.ServerOperations.DisableOperatorUserCreation {
		return nil
	}
	user, err := u.createOperatorUserIfNotExists()
	if err != nil {
		return err
	}
	_, pass, err := u.getOperatorUserCredentials()
	if err != nil {
		return err
	}
	u.nexuscli.SetCredentials(user.UserID, pass)
	return nil
}

func (u *userOperation) createOperatorUserIfNotExists() (*nexus.User, error) {
	// TODO: open an issue to handle access to a custom admin credentials to be used by the operator
	u.nexuscli.SetCredentials(defaultAdminUsername, defaultAdminPassword)
	user, err := u.nexuscli.UserService.GetUserByID(operatorUsername)
	if err != nil {
		return nil, err
	}
	if user != nil {
		u.status.OperatorUserCreated = true
		return user, nil
	}
	user, err = u.createOperatorUserInstance()
	if err != nil {
		return nil, err
	}
	if err := u.nexuscli.UserService.Add(*user); err != nil {
		return nil, err
	}
	if err := u.storeOperatorUserCredentials(user); err != nil {
		//  TODO: in case of an error here, we should remove the user from the Nexus database. Edge case: an user could manually add the credentials later to the secret with a manually created user for us.
		return nil, err
	}
	u.status.OperatorUserCreated = true
	return user, nil
}

func (u *userOperation) storeOperatorUserCredentials(user *nexus.User) error {
	secret := &corev1.Secret{}
	if err := framework.Fetch(u.k8sclient, framework.Key(u.nexus), secret); err != nil {
		return err
	}
	if secret.StringData == nil {
		secret.StringData = make(map[string]string)
	}
	secret.StringData[secretKeyPassword] = user.Password
	secret.StringData[secretKeyUsername] = user.UserID
	if err := u.k8sclient.Update(context.TODO(), secret); err != nil {
		return err
	}
	return nil
}

func (u *userOperation) getOperatorUserCredentials() (user, password string, err error) {
	secret := &corev1.Secret{}
	if err := framework.Fetch(u.k8sclient, framework.Key(u.nexus), secret); err != nil {
		return "", "", err
	}
	return secret.StringData[secretKeyUsername], secret.StringData[secretKeyPassword], nil
}

func (u *userOperation) createOperatorUserInstance() (*nexus.User, error) {
	password, err := u.generateRandomPassword()
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

func (u *userOperation) generateRandomPassword() (string, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", nil
	}
	return uid.String(), nil
}
