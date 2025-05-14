import { v1 } from "@authzed/authzed-node";

const client = v1.NewClient(
	"averysecretpresharedkey",
	"localhost:50051",
	v1.ClientSecurity.INSECURE_LOCALHOST_ALLOWED,
);
const { promises: promiseClient } = client; // access client.promises

// Create the relationship between the resource and the user.
const firstPost = v1.ObjectReference.create({
	objectType: "post",
	objectId: "1",
});

// Create the user reference.
const emilia = v1.ObjectReference.create({
	objectType: "user",
	objectId: "emilia",
});

// Create the subject reference using the user reference
const subject = v1.SubjectReference.create({
	object: emilia,
});

const checkPermissionRequest = v1.CheckPermissionRequest.create({
	resource: firstPost,
	permission: "read",
	subject,
});

// client.checkPermission(checkPermissionRequest, (err, response) => {
// 	console.log(response);
// 	console.log(err);
// });

const result = await promiseClient.checkPermission(checkPermissionRequest);
console.log(result);
