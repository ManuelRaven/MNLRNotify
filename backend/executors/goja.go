package executors

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

type GojaExecutor struct{}

func init() {
	RegisterExecutor("goja", &GojaExecutor{})
}

type KVStoreJSObject struct {
	pb      *pocketbase.PocketBase
	ownerid string
	Name    string
}

func (o *KVStoreJSObject) Get(inputkey string) string {
	key := strings.ToLower(inputkey)
	record, err := o.pb.FindFirstRecordByFilter(
		"KVStore",
		"id = {:key} && owner = {:ownerid}",
		dbx.Params{"key": key, "ownerid": o.ownerid},
	)

	if err != nil || record == nil {
		fmt.Printf("Error getting value for key: %s, error: %v\n", key, err)
		return ""
	}

	value := record.GetString("value")
	fmt.Printf("Getting value for key: %s, value: %s\n", key, value)
	return value
}

func (o *KVStoreJSObject) Set(inputkey string, value string) {
	// Try to find the existing record
	key := strings.ToLower(inputkey)
	record, err := o.pb.FindFirstRecordByFilter(
		"KVStore",
		"id = {:key} && owner = {:ownerid}",
		dbx.Params{"key": key, "ownerid": o.ownerid},
	)

	if err != nil || record == nil {
		// Record doesn't exist, create a new one
		collection, err := o.pb.FindCollectionByNameOrId("KVStore")
		if err != nil {
			fmt.Printf("Error finding collection: %v\n", err)
			return
		}

		record := core.NewRecord(collection)
		record.Set("id", key)
		record.Set("owner", o.ownerid)
		record.Set("value", value)

		if err := o.pb.Save(record); err != nil {
			fmt.Printf("Error creating record for key: %s, error: %v\n", key, err)
			return
		}
		fmt.Printf("Created new value for key: %s, value: %s\n", key, value)
	} else {
		// Record exists, update it
		record.Set("value", value)
		if err := o.pb.Save(record); err != nil {
			fmt.Printf("Error updating record for key: %s, error: %v\n", key, err)
			return
		}
		fmt.Printf("Updated value for key: %s, value: %s\n", key, value)
	}
}

func (o *KVStoreJSObject) Delete(inputkey string) {
	key := strings.ToLower(inputkey)
	record, err := o.pb.FindRecordById("KVStore", key)
	if err != nil || record == nil {
		fmt.Printf("Error finding record for key: %s, error: %v\n", key, err)
		return
	}

	err = o.pb.Delete(record)
	if err != nil {
		return
	}

}

func (e *GojaExecutor) Execute(mep *MessageExecutorParameters) string {
	message := mep.Message
	script := mep.Script
	pb := mep.PB
	ownerid := mep.OwnerID
	vm := goja.New()

	// Create a new KVStoreJSObject
	kvStore := &KVStoreJSObject{
		pb:      pb,
		ownerid: ownerid,
		Name:    "kvStore",
	}
	// Set the KVStoreJSObject as a global variable
	if err := vm.Set("kvStore", kvStore); err != nil {
		return "Error: Failed to set kvStore variable"
	}
	// Set the message as a global variable

	// Set message as a global variable
	err := vm.Set("message", message)
	if err != nil {
		return "Error: Failed to set message variable"
	}

	// Execute the script
	val, err := vm.RunString(script)
	if err != nil {
		return "Error: " + err.Error()
	}

	// Convert result to string
	result := val.String()
	if result == "undefined" {
		return message
	}

	return result
}
