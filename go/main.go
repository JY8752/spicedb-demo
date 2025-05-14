package main

import (
	"context"
	"log"

	"github.com/JY8752/spicedb-go-demo/spicedb"
)

func main() {
	client, err := spicedb.NewSpiceDBClient("localhost", 50051, "averysecretpresharedkey")
	if err != nil {
		log.Fatalf("failed to create spiceDB client: %s", err)
	}

	// req := spicedb.NewCheckPermissionRequest(
	// 	spicedb.PostObjectType,
	// 	"1",
	// 	spicedb.UserObjectType,
	// 	"emilia",
	// 	spicedb.ReadPermission,
	// )

	// allowed, err := client.CheckPermission(context.Background(), req)
	// if err != nil {
	// 	log.Fatalf("failed to check permission: %s", err)
	// }

	// if allowed {
	// 	log.Println("allowed!")
	// } else {
	// 	log.Println("not allowed!")
	// }

	// req := spicedb.NewLookupResourcesRequest(
	// 	spicedb.PostObjectType,
	// 	spicedb.ReadPermission,
	// 	spicedb.UserObjectType,
	// 	"emilia",
	// )

	// resources, err := client.LookupResources(context.Background(), req)
	// if err != nil {
	// 	log.Fatalf("failed to lookup resources: %s", err)
	// }

	// log.Printf("resources: %v", resources)

	req := spicedb.NewLookupSubjectsRequest(
		spicedb.PostObjectType,
		"1",
		spicedb.ReadPermission,
		spicedb.UserObjectType,
	)

	subjects, err := client.LookupSubjects(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to lookup subjects: %s", err)
	}

	for subject := range subjects {
		log.Printf("subject: %s\n", subject)
	}
}
