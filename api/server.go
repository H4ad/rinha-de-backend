package main

import (
	"context"
	"strconv"

	"api/models"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/joho/godotenv"
	"github.com/meilisearch/meilisearch-go"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func main() {
	envErr := godotenv.Load()

	if envErr != nil {
		log.Warn("Error loading .env file")
	}

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	app.Use(pprof.New())

	db := GetDatabase()
	meilisearchRef := GetMeilisearch()

	peopleIndex := meilisearchRef.Index("people")

	validate := validator.New()

	app.Post("pessoas", func(c *fiber.Ctx) error {
		var payload CreatePersonPayload

		err := c.BodyParser(&payload)

		if err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "Failed to parse the body"})
		}

		err = validate.Struct(payload)

		if err != nil {
			return c.Status(422).JSON(fiber.Map{"message": "One or more fields are invalid, more details: " + err.Error()})
		}

		person := &models.Pessoa{
			Nome:       payload.Name,
			Apelido:    payload.Nickname,
			Nascimento: payload.Birthdate,
			Stack:      payload.Stack,
		}

		err = person.Insert(context.Background(), db, boil.Infer())

		if err != nil {
			return c.Status(422).JSON(fiber.Map{"message": "Error during the creation of the user: " + err.Error()})
		}

		// _, errMeili := peopleIndex.AddDocuments(person)
		//
		// if errMeili != nil {
		// 	return c.Status(500).JSON(fiber.Map{"message": "Error during the creation of the user in the search engine: " + errMeili.Error()})
		// }

		c.Location("/pessoas/" + person.ID)

		return c.Status(201).JSON(person)
	})

	app.Get("pessoas", func(c *fiber.Ctx) error {
		searchText := c.Query("t", "")

		if searchText == "" {
			return c.Status(400).JSON(fiber.Map{"message": "The query parameter 't' is required"})
		}

		searchRes, errMeili := peopleIndex.Search(searchText,
			&meilisearch.SearchRequest{
				Limit: 50,
			})

		if errMeili != nil {
			return c.Status(500).JSON(fiber.Map{"message": "Error during the search in the search engine: " + errMeili.Error()})
		}

		return c.Status(200).JSON(searchRes.Hits)
	})

	app.Get("pessoas/:id", func(c *fiber.Ctx) error {
		c.Response().Header.Add("Cache-Control", "max-age=60")

		found, err := models.FindPessoa(context.Background(), db, c.Params("id"))

		if err != nil {
			return c.Status(404).JSON(fiber.Map{"message": "Person not found"})
		}

		return c.Status(200).JSON(found)
	})

	app.Get("contagem-pessoas", func(c *fiber.Ctx) error {
		count, err := models.Pessoas().Count(context.Background(), db)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "Some error happened during the count. More details: " + err.Error()})
		}

		return c.Status(200).SendString(strconv.FormatInt(count, 10))
	})

	err := app.Listen(":3000")

	if err != nil {
		log.Fatal("Error during the startup of the server: " + err.Error())
	}
}
