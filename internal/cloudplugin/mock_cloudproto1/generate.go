// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:generate go run github.com/golang/mock/mockgen -destination mock.go github.com/hashicorp/mnptu/internal/cloudplugin/cloudproto1 CommandServiceClient,CommandService_ExecuteClient

package mock_cloudproto1
