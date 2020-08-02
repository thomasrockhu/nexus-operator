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
	"testing"

	"github.com/m88i/aicura/nexus"
	"github.com/m88i/nexus-operator/pkg/apis/apps/v1alpha1"
	"github.com/m88i/nexus-operator/pkg/framework"
	"github.com/m88i/nexus-operator/pkg/test"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_userOperation_EnsureOperatorUser(t *testing.T) {
	nexusInstance := &v1alpha1.Nexus{ObjectMeta: v1.ObjectMeta{Name: "nexus3", Namespace: t.Name()}}
	client := test.NewFakeClientBuilder(
		nexusInstance,
		&corev1.Secret{ObjectMeta: v1.ObjectMeta{Name: "nexus3", Namespace: t.Name()}}).
		Build()
	server := &server{
		nexus:     nexusInstance,
		k8sclient: client,
		nexuscli:  nexus.NewFakeClient(),
	}

	err := userOperations(server).EnsureOperatorUser()
	assert.NoError(t, err)
	user, err := server.nexuscli.UserService.GetUserByID(operatorUsername)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, operatorUsername, user.UserID)

	secret := &corev1.Secret{ObjectMeta: v1.ObjectMeta{Name: "nexus3", Namespace: t.Name()}}
	err = client.Get(context.TODO(), framework.Key(secret), secret)
	assert.NoError(t, err)
	assert.Equal(t, operatorUsername, secret.StringData[secretKeyUsername])
	assert.NotEmpty(t, secret.StringData[secretKeyPassword])
}
