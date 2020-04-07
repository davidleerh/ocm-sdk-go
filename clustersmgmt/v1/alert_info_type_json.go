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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalAlertInfo writes a value of the 'alert_info' type to the given writer.
func MarshalAlertInfo(object *AlertInfo, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeAlertInfo(object, stream)
	stream.Flush()
	return stream.Error
}

// writeAlertInfo writes a value of the 'alert_info' type to the given stream.
func writeAlertInfo(object *AlertInfo, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	if object.name != nil {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(*object.name)
		count++
	}
	if object.severity != nil {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("severity")
		stream.WriteString(string(*object.severity))
		count++
	}
	stream.WriteObjectEnd()
}

// UnmarshalAlertInfo reads a value of the 'alert_info' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAlertInfo(source interface{}) (object *AlertInfo, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readAlertInfo(iterator)
	err = iterator.Error
	return
}

// readAlertInfo reads a value of the 'alert_info' type from the given iterator.
func readAlertInfo(iterator *jsoniter.Iterator) *AlertInfo {
	object := &AlertInfo{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "name":
			value := iterator.ReadString()
			object.name = &value
		case "severity":
			text := iterator.ReadString()
			value := AlertSeverity(text)
			object.severity = &value
		default:
			iterator.ReadAny()
		}
	}
	return object
}
