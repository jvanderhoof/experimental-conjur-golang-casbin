package main

import (
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type credential struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

var secrets = []credential{
	{ID: "demo:variable:foo", Value: "foobar-1"},
	{ID: "demo:variable:bar", Value: "foobar-2"},
	{ID: "demo:variable:baz", Value: "foobar-3"},
}

func main() {
	router := gin.Default()

	db, err := gorm.Open(sqlite.Open("/data/conjur.db"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to open sqlite database: %v", err))
	}

	// Initialize  casbin adapter
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize casbin adapter: %v", err))
	}

	// Load model configuration file and policy store adapter
	enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		panic(fmt.Sprintf("failed to create casbin enforcer: %v", err))
	}

	//add policy
	// if hasPolicy := enforcer.HasPolicy("demo:user:admin", "demo:variable:foo", "execute"); !hasPolicy {
	// 	enforcer.AddPolicy("demo:user:admin", "demo:variable:foo", "execute")
	// }
	if hasPolicy := enforcer.HasPolicy("demo:user:admin", "demo:variable:baz", "execute"); !hasPolicy {
		enforcer.AddPolicy("demo:user:admin", "demo:variable:baz", "execute")
	}
	if hasPolicy := enforcer.HasPolicy("demo:group:viewer", "demo:variable:foo", "execute"); !hasPolicy {
		enforcer.AddPolicy("demo:user:admin", "demo:variable:foo", "execute")
	}
	// if hasPolicy := enforcer.HasPolicy("demo:user:admin", "demo:variable:foo", "execute"); !hasPolicy {
	// 	enforcer.AddPolicy("demo:user:admin", "demo:variable:foo", "execute")
	// }
	if hasPolicy, err := enforcer.HasRoleForUser("demo:user:admin", "demo:group:viewer"); !hasPolicy {
		if err != nil {
			panic(fmt.Sprintf("failed to verify role is already present: %v", err))
		}

		enforcer.AddRoleForUser("demo:user:admin", "demo:group:viewer")
	}

	router.GET("/secrets/:account/variable/:id", getSecret(enforcer))

	router.Run("localhost:8088")
}

func getSecret(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		account := c.Param("account")
		secretId := fmt.Sprintf("%s:variable:%s", account, id)

		// Load policy from Database
		err := enforcer.LoadPolicy()
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"msg": "Failed to load policy from DB"})
			return
		}

		// Casbin enforces policy
		ok, err := enforcer.Enforce("demo:user:admin", secretId, "execute")

		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"msg": "Error occurred when authorizing user"})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"msg": "User is not allowed to access this secret"})
			return
		}

		for _, secret := range secrets {
			if secret.ID == secretId {
				// c.IndentedJSON(http.StatusOK, secret)
				c.String(http.StatusOK, secret.Value)
				return
			}
		}
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "secret not found"})
	}
}
