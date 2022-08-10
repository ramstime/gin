package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

//var db = make(map[string]string)
type Author struct {
	Name  string   `json:"name"`
	Age   int      `json:"age"`
	Books []string `json:"books"`
}
type handler struct {
	Client *redis.Client
}

func (h handler) setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	// Get user value
	r.GET("user/:name", h.getDB)

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	authorized.POST("admin", h.insertDB)

	return r
}

/* example curl for /user/name with basicauth header
curl -X GET   http://localhost:8080/user/foo
curl -X GET   http://localhost:8080/user/fo*
{"user":"foo","value":"{\"name\":\"rams\",\"age\":35,\"books\":null}"}
*/
func (h handler) getDB(c *gin.Context) {
	author := c.Params.ByName("name")
	fmt.Println("author: ", author)
	if strings.Contains(author, "*") {
		value, err := h.Client.Keys(c, author).Result()
		if err != nil {
			c.JSON(http.StatusNoContent, gin.H{"user": author, "status": "no value"})
			fmt.Println(err)
			return
		}
		var authorDetails []Author
		for _, user := range value {
			value, err := h.Client.Get(c, user).Result()
			if err != nil {
				c.JSON(http.StatusNoContent, gin.H{"user": user, "status": "no value"})
				fmt.Println(err)
				return
			}
			fmt.Printf("user: %v value: %v \n", user, value)
			var authorData Author
			if err := json.Unmarshal([]byte(value), &authorData); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"user": user, "status": "no value"})
				fmt.Println(err)
				return
			}
			authorDetails = append(authorDetails, authorData)
		}
		c.JSON(http.StatusOK, gin.H{"user": author, "value": authorDetails})
		return
	}
	value, err := h.Client.Get(c, author).Result()
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{"user": author, "status": "no value"})
		fmt.Println(err)
		return
	}
	fmt.Printf("user: %v value: %v", author, value)
	//value, ok := db[user]
	c.JSON(http.StatusOK, gin.H{"user": author, "value": value})

}

/* example curl for /admin with basicauth header
   Zm9vOmJhcg== is base64("foo:bar")

	curl -X POST \
  	http://localhost:8080/admin \
  	-H 'authorization: Basic Zm9vOmJhcg==' \
  	-H 'content-type: application/json' \
  	-d '{"value":"bar"}'
*/
func (h handler) insertDB(c *gin.Context) {

	user := c.MustGet(gin.AuthUserKey).(string)

	author := Author{}
	fmt.Println("user:", user)
	if c.Bind(&author) == nil {

		fmt.Println("author:", author)
		jsonVal, err := json.Marshal(author)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"user": user, "status": "no value"})
			return
		}
		err = h.Client.Set(c, author.Name, jsonVal, 0).Err()
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"user": user, "status": "no value"})
			return
		}
		val, err := h.Client.Get(c, author.Name).Result()
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"user": user, "status": "no value"})
			return
		}
		fmt.Println(val)

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"user": user, "status": "no value"})
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	handler := &handler{Client: client}
	r := handler.setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")

}
