package cdn

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator 0.14.0.0
// Changes may cause incorrect behavior and will be lost if the code is
// regenerated.

import (
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// NameAvailabilityClient is the use these APIs to manage Azure CDN resources
// through the Azure Resource Manager. You must make sure that requests made
// to these resources are secure. For more information, see <a
// href="https://msdn.microsoft.com/en-us/library/azure/dn790557.aspx">Authenticating
// Azure Resource Manager requests.</a>
type NameAvailabilityClient struct {
	ManagementClient
}

// NewNameAvailabilityClient creates an instance of the NameAvailabilityClient
// client.
func NewNameAvailabilityClient(subscriptionID string) NameAvailabilityClient {
	return NewNameAvailabilityClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewNameAvailabilityClientWithBaseURI creates an instance of the
// NameAvailabilityClient client.
func NewNameAvailabilityClientWithBaseURI(baseURI string, subscriptionID string) NameAvailabilityClient {
	return NameAvailabilityClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CheckNameAvailability sends the check name availability request.
//
// checkNameAvailabilityInput is input to check
func (client NameAvailabilityClient) CheckNameAvailability(checkNameAvailabilityInput CheckNameAvailabilityInput) (result CheckNameAvailabilityOutput, err error) {
	req, err := client.CheckNameAvailabilityPreparer(checkNameAvailabilityInput)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "cdn/NameAvailabilityClient", "CheckNameAvailability", nil, "Failure preparing request")
	}

	resp, err := client.CheckNameAvailabilitySender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "cdn/NameAvailabilityClient", "CheckNameAvailability", resp, "Failure sending request")
	}

	result, err = client.CheckNameAvailabilityResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cdn/NameAvailabilityClient", "CheckNameAvailability", resp, "Failure responding to request")
	}

	return
}

// CheckNameAvailabilityPreparer prepares the CheckNameAvailability request.
func (client NameAvailabilityClient) CheckNameAvailabilityPreparer(checkNameAvailabilityInput CheckNameAvailabilityInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPath("/providers/Microsoft.Cdn/checkNameAvailability"),
		autorest.WithJSON(checkNameAvailabilityInput),
		autorest.WithQueryParameters(queryParameters))
}

// CheckNameAvailabilitySender sends the CheckNameAvailability request. The method will close the
// http.Response Body if it receives an error.
func (client NameAvailabilityClient) CheckNameAvailabilitySender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// CheckNameAvailabilityResponder handles the response to the CheckNameAvailability request. The method always
// closes the http.Response Body.
func (client NameAvailabilityClient) CheckNameAvailabilityResponder(resp *http.Response) (result CheckNameAvailabilityOutput, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
