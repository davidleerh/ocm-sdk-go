/*
Copyright (c) 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package v1 // github.com/openshift-online/ocm-sdk-go/authorizations/v1

// ResourceReviewRequest represents the values of the 'resource_review_request' type.
//
// Request to perform a resource access review.
type ResourceReviewRequest struct {
	accountUsername *string
	action          *string
	resourceType    *string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ResourceReviewRequest) Empty() bool {
	return o == nil || (o.accountUsername == nil &&
		o.action == nil &&
		o.resourceType == nil &&
		true)
}

// AccountUsername returns the value of the 'account_username' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the account that is trying to access the resource.
func (o *ResourceReviewRequest) AccountUsername() string {
	if o != nil && o.accountUsername != nil {
		return *o.accountUsername
	}
	return ""
}

// GetAccountUsername returns the value of the 'account_username' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the account that is trying to access the resource.
func (o *ResourceReviewRequest) GetAccountUsername() (value string, ok bool) {
	ok = o != nil && o.accountUsername != nil
	if ok {
		value = *o.accountUsername
	}
	return
}

// Action returns the value of the 'action' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Action that will be performed on the resource.
func (o *ResourceReviewRequest) Action() string {
	if o != nil && o.action != nil {
		return *o.action
	}
	return ""
}

// GetAction returns the value of the 'action' attribute and
// a flag indicating if the attribute has a value.
//
// Action that will be performed on the resource.
func (o *ResourceReviewRequest) GetAction() (value string, ok bool) {
	ok = o != nil && o.action != nil
	if ok {
		value = *o.action
	}
	return
}

// ResourceType returns the value of the 'resource_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Type of resource.
func (o *ResourceReviewRequest) ResourceType() string {
	if o != nil && o.resourceType != nil {
		return *o.resourceType
	}
	return ""
}

// GetResourceType returns the value of the 'resource_type' attribute and
// a flag indicating if the attribute has a value.
//
// Type of resource.
func (o *ResourceReviewRequest) GetResourceType() (value string, ok bool) {
	ok = o != nil && o.resourceType != nil
	if ok {
		value = *o.resourceType
	}
	return
}

// ResourceReviewRequestListKind is the name of the type used to represent list of objects of
// type 'resource_review_request'.
const ResourceReviewRequestListKind = "ResourceReviewRequestList"

// ResourceReviewRequestListLinkKind is the name of the type used to represent links to list
// of objects of type 'resource_review_request'.
const ResourceReviewRequestListLinkKind = "ResourceReviewRequestListLink"

// ResourceReviewRequestNilKind is the name of the type used to nil lists of objects of
// type 'resource_review_request'.
const ResourceReviewRequestListNilKind = "ResourceReviewRequestListNil"

// ResourceReviewRequestList is a list of values of the 'resource_review_request' type.
type ResourceReviewRequestList struct {
	href  *string
	link  bool
	items []*ResourceReviewRequest
}

// Len returns the length of the list.
func (l *ResourceReviewRequestList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *ResourceReviewRequestList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ResourceReviewRequestList) Get(i int) *ResourceReviewRequest {
	if l == nil || i < 0 || i >= len(l.items) {
		return nil
	}
	return l.items[i]
}

// Slice returns an slice containing the items of the list. The returned slice is a
// copy of the one used internally, so it can be modified without affecting the
// internal representation.
//
// If you don't need to modify the returned slice consider using the Each or Range
// functions, as they don't need to allocate a new slice.
func (l *ResourceReviewRequestList) Slice() []*ResourceReviewRequest {
	var slice []*ResourceReviewRequest
	if l == nil {
		slice = make([]*ResourceReviewRequest, 0)
	} else {
		slice = make([]*ResourceReviewRequest, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ResourceReviewRequestList) Each(f func(item *ResourceReviewRequest) bool) {
	if l == nil {
		return
	}
	for _, item := range l.items {
		if !f(item) {
			break
		}
	}
}

// Range runs the given function for each index and item of the list, in order. If
// the function returns false the iteration stops, otherwise it continues till all
// the elements of the list have been processed.
func (l *ResourceReviewRequestList) Range(f func(index int, item *ResourceReviewRequest) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
