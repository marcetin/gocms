package routes

import (
	"strconv"

	"github.com/dgraph-io/badger/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/marcetin/gocms/db"
	"github.com/marcetin/gocms/models"
	"github.com/vmihailenco/msgpack/v5"
)

func FrontendIndex(c *fiber.Ctx) error {
	var posts []models.Post
	err := db.GetDB().View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			var post models.Post
			err := item.Value(func(val []byte) error {
				return msgpack.Unmarshal(val, &post)
			})
			if err != nil {
				return err
			}
			posts = append(posts, post)
		}
		return nil
	})
	if err != nil {
		return c.Status(500).SendString("Failed to load posts")
	}

	return c.Render("frontend/index", fiber.Map{"Posts": posts})
}

func FrontendPost(c *fiber.Ctx) error {
	idStr := c.Params("id")        // Retrieve the "id" parameter from the URL as a string
	id, err := strconv.Atoi(idStr) // Convert the string to an integer
	if err != nil {
		return c.Status(400).SendString("Invalid post ID") // Handle invalid ID format
	}

	var post models.Post
	err = db.GetDB().View(func(txn *badger.Txn) error {
		// Fetch the post using the ID from the database
		item, err := txn.Get([]byte(strconv.Itoa(id)))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return msgpack.Unmarshal(val, &post)
		})
	})
	if err != nil {
		return c.Status(404).SendString("Post not found") // Handle post not found
	}

	// Use the post to render the frontend post page
	return c.Render("frontend/post", fiber.Map{
		"Post": post,
	})
}
