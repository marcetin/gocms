package routes

import (
	"strconv"

	"github.com/dgraph-io/badger/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/marcetin/gocms/db"
	"github.com/marcetin/gocms/models"
	"github.com/vmihailenco/msgpack/v5"
)

func AdminDashboard(c *fiber.Ctx) error {
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

	return c.Render("admin/index", fiber.Map{"Posts": posts})
}

// Renders the page for creating a new post
func CreatePostPage(c *fiber.Ctx) error {
	return c.Render("admin/create", fiber.Map{})
}

// Handles creating a new post
func CreatePost(c *fiber.Ctx) error {
	title := c.FormValue("title")
	content := c.FormValue("content")

	// Generate a new ID (you can replace this with your custom logic)
	id := strconv.FormatInt(int64(c.Context().Time().Unix()), 10)

	post := models.Post{
		ID:      id,
		Title:   title,
		Content: content,
	}

	// Save the post to the database
	err := db.GetDB().Update(func(txn *badger.Txn) error {
		data, err := msgpack.Marshal(post)
		if err != nil {
			return err
		}
		return txn.Set([]byte(id), data)
	})
	if err != nil {
		return c.Status(500).SendString("Failed to save post")
	}

	return c.Redirect("/admin")
}

// Renders the page for editing a post
func EditPostPage(c *fiber.Ctx) error {
	id := c.Params("id")

	var post models.Post
	err := db.GetDB().View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return msgpack.Unmarshal(val, &post)
		})
	})
	if err != nil {
		return c.Status(404).SendString("Post not found")
	}

	return c.Render("admin/edit", fiber.Map{"Post": post})
}

// Handles editing an existing post
func EditPost(c *fiber.Ctx) error {
	id := c.Params("id")
	title := c.FormValue("title")
	content := c.FormValue("content")

	post := models.Post{
		ID:      id,
		Title:   title,
		Content: content,
	}

	// Update the post in the database
	err := db.GetDB().Update(func(txn *badger.Txn) error {
		data, err := msgpack.Marshal(post)
		if err != nil {
			return err
		}
		return txn.Set([]byte(id), data)
	})
	if err != nil {
		return c.Status(500).SendString("Failed to update post")
	}

	return c.Redirect("/admin")
}

// Handles deleting a post
func DeletePost(c *fiber.Ctx) error {
	id := c.Params("id")

	// Delete the post from the database
	err := db.GetDB().Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(id))
	})
	if err != nil {
		return c.Status(500).SendString("Failed to delete post")
	}

	return c.Redirect("/admin")
}
