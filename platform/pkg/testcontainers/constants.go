package testcontainers

// MongoDB constants
const (
	// MongoDB container constants
	MongoContainerName = "mongo"
	MongoPort          = "27017"

	// MongoDB environment variables
	MongoImageNameKey = "MONGO_IMAGE_NAME"
	MongoHostKey      = "MONGO_HOST"
	MongoPortKey      = "MONGO_PORT"
	MongoDatabaseKey  = "MONGO_DATABASE"
	MongoRootPrefix   = "MONGO_INITDB_ROOT_"
	MongoUsernameKey  = MongoRootPrefix + "USERNAME"
	MongoPasswordKey  = MongoRootPrefix + "PASSWORD"
	MongoAuthDBKey    = "MONGO_AUTH_DB"
)
