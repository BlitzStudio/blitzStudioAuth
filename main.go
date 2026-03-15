package main

import (
	"database/sql"
	"time"

	"github.com/BlitzStudio/blitzStudioAuth/out/repository"
	"github.com/BlitzStudio/blitzStudioAuth/types"
	"github.com/BlitzStudio/blitzStudioAuth/utils"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/subosito/gotenv"
)

type GlobalValidator struct {
	validator *validator.Validate
}

func (v *GlobalValidator) Validate(out any) error {
	return v.validator.Struct(out)
}

func main() {
	log := utils.GetLogger()
	err := gotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("mysql", "root:pass@/test?parseTime=true")
	// ctx := context.Background()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// i dont know what those 2 options do
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	repo := repository.New(db)
	validate := validator.New()
	app := fiber.New(fiber.Config{
		StructValidator: &GlobalValidator{validator: validate},
	})

	app.Use(func(c fiber.Ctx) error {
		c.Locals("repo", repo)
		return c.Next()
	})

	app.Post("/auth/signup", func(c fiber.Ctx) error {
		repo, ok := c.Locals("repo").(*repository.Queries)
		if !ok {
			log.Fatal("Couldnt access the db")
		}

		// preia valorile din post body
		userData := new(types.User)
		if err := c.Bind().JSON(userData); err != nil {
			log.Error(err)
			return c.Status(400).SendString("Incomplete body request")
		}

		userPasswordHash, err := utils.GenerateHash(userData.Password)
		if err != nil {
			log.Error(err.Error())
			return c.SendStatus(500)
		}

		createdUser, err := repo.CreateUser(c.Context(), repository.CreateUserParams{
			Email:    userData.Email,
			Name:     userData.Name,
			Password: userPasswordHash,
		})

		if err != nil {
			log.Error(err)
			return c.Status(500).SendString("This user already exists")
		}

		userId, err := createdUser.LastInsertId()
		if err != nil {
			log.Fatal(err.Error())
		}

		timeNow := time.Now()
		tokenFamily := uuid.New().String()
		refreshTokenId := utils.GenerateUlid()
		tokens := types.AuthTokens{
			AccessToken:  utils.GenerateAccessToken(userId, utils.GenerateUlid(), timeNow),
			RefreshToken: utils.GenerateRefreshToken(userId, refreshTokenId, tokenFamily, timeNow),
		}
		if err != nil {
			log.Error(err)
			return c.SendStatus(500)
		}

		err = repo.SaveRefreshToken(c.Context(), repository.SaveRefreshTokenParams{
			TokenId:     refreshTokenId,
			UserId:      sql.NullInt32{Int32: int32(userId), Valid: true},
			TokenFamily: tokenFamily,
			ExpiresAt:   sql.NullTime{Time: timeNow.Add(7 * 24 * time.Hour), Valid: true},
		})
		if err != nil {
			log.Error(err)
			return c.SendStatus(500)
		}

		log.Info("Created user: " + userData.Email + "with access token: " + tokens.AccessToken + " and refresh token: " + tokens.RefreshToken)
		return c.JSON(tokens)
	})

	app.Post("/auth/signin", func(c fiber.Ctx) error {
		repo, ok := c.Locals("repo").(*repository.Queries)
		if !ok {
			log.Fatal("Couldnt access the db")
		}
		// preia valorile din post body
		userData := new(types.User)
		if err := c.Bind().JSON(userData); err != nil {
			log.Error(err)
			return c.Status(400).SendString("Incomplete body request")
		}

		user, err := repo.FindUserByEmail(c.Context(), userData.Email)
		if err != nil {
			log.Error(err)
			return c.Status(401).SendString("Invalid email or password")
		}

		match, err := utils.CompareHash(userData.Password, user.Password)
		if err != nil {
			log.Error(err)
			return c.SendStatus(500)
		}

		if !match {
			return c.Status(401).SendString("Invalid email or password")
		}

		timeNow := time.Now()
		tokenFamily := uuid.New().String()
		refreshTokenId := utils.GenerateUlid()
		tokens := types.AuthTokens{
			AccessToken:  utils.GenerateAccessToken(int64(user.ID), utils.GenerateUlid(), timeNow),
			RefreshToken: utils.GenerateRefreshToken(int64(user.ID), refreshTokenId, tokenFamily, timeNow),
		}

		repo.SaveRefreshToken(c.Context(), repository.SaveRefreshTokenParams{
			TokenId: refreshTokenId,
			UserId:  sql.NullInt32{Int32: int32(user.ID), Valid: true},
			// DeviceId:    userData.DeviceId,
			TokenFamily: tokenFamily,
			// TokenHash:   sql.NullString{String: refreshTokenHash, Valid: true},
			ExpiresAt: sql.NullTime{Time: timeNow.Add(7 * 24 * time.Hour), Valid: true},
		})

		return c.JSON(tokens)

	})
	app.Post("/auth/refresh-token", func(c fiber.Ctx) error {
		refreshToken := new(types.RefreshToken)
		if err := c.Bind().JSON(refreshToken); err != nil {
			log.Error(err)
			return c.Status(400).SendString("Incomplete body request")
		}
		token, err := utils.ValidateRefreshToken(refreshToken.Value)
		if err != nil {
			log.Error(err.Error())
			return c.Status(401).SendString("invalid token")
		}

		savedToken, err := repo.FindTokenById(c.Context(), token.ID)
		if err != nil {
			log.Error(err)
			return c.Status(401).SendString("invalid token")
		}

		if savedToken.Isrevoked.Bool && savedToken.Expiresat.Time.Unix()-time.Now().Unix() > 0 {
			repo.RevokeTokenFamily(c.Context(), savedToken.Tokenfamily)
		}

		if (time.Now().Unix()-savedToken.UpdatedAt.Time.Unix() < 10 && savedToken.Isrevoked.Bool == true) || (savedToken.Isrevoked.Bool == false && savedToken.Expiresat.Time.Unix()-time.Now().Unix() > 0) {
			timeNow := time.Now()
			refreshTokenUlid := utils.GenerateUlid()
			tokens := types.AuthTokens{
				AccessToken:  utils.GenerateAccessToken(int64(savedToken.Userid.Int32), utils.GenerateUlid(), timeNow),
				RefreshToken: utils.GenerateRefreshToken(int64(savedToken.Userid.Int32), refreshTokenUlid, savedToken.Tokenfamily, timeNow),
			}
			if err != nil {
				log.Error(err)
				return c.SendStatus(500)
			}
			err = repo.SaveRefreshToken(c.Context(), repository.SaveRefreshTokenParams{
				TokenId:     refreshTokenUlid,
				UserId:      sql.NullInt32{Int32: savedToken.Userid.Int32, Valid: true},
				TokenFamily: savedToken.Tokenfamily,
				ExpiresAt:   sql.NullTime{Time: timeNow.Add(7 * 24 * time.Hour), Valid: true},
			})
			if err != nil {
				log.Error(err)
				return c.SendStatus(500)
			}
			err = repo.RevokeRefreshTokenById(c.Context(), savedToken.ID)
			if err != nil {
				log.Error(err)
				return c.SendStatus(500)
			}
			return c.JSON(tokens)
		} else if time.Now().Unix()-savedToken.Expiresat.Time.Unix() > 10 && savedToken.Isrevoked.Bool == true {
			return c.Status(400).SendString("invalid token")
		}

		return c.Status(401).SendString("invalid token")
	})

	app.Get("*", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	log.Fatal(app.Listen(":3000"))
}
