package authUsecase

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/Jaxongir1006/Chat-X-v2/internal/domain"
	apperr "github.com/Jaxongir1006/Chat-X-v2/internal/errors"
	"github.com/redis/go-redis/v9"
)


func (a *AuthUsecase) Register(ctx context.Context, req RegisterRequest) error {
	if req.ConfirmPass != req.Password {
		return apperr.New(apperr.CodeConflict, http.StatusConflict, "passwords do not match")
	}

	// password validation
	if len(req.Password) < 8 {
		return apperr.New(apperr.CodeConflict, http.StatusConflict, "password must be at least 8 characters long")
	}

	// check if user with that email exists
	if _, err := a.authStore.GetByEmail(ctx, req.Email); err == nil {
		return apperr.New(apperr.CodeConflict, http.StatusConflict, "user with that email already exists")
	}

	if req.Phone != "" {
		// check if user with that phone number exists
		if _, err := a.authStore.GetByPhone(ctx, req.Phone); err == nil {
			return apperr.New(apperr.CodeConflict, http.StatusConflict, "user with that phone number already exists")
		}
	}
	
	if req.Username != "" {
		// check if user with that username exists
		if _, err := a.authStore.GetByUsername(ctx, req.Username); err == nil {
			return apperr.New(apperr.CodeConflict, http.StatusConflict, "user with that username already exists")
		}
	}

	hashed, err := a.hasher.Hash(req.Password)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	// create user
	user := &domain.User{
		Username: req.Username,
		Email: req.Email,
		Phone: req.Phone,
		Password: hashed,
		Role: "user",
		Verified: false,
	}
	if err := a.authStore.InsertUser(ctx, user); err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	go func ()  {
		// save the email code to redis with a 5 minute ttl
		bgCtx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
		defer cancel()

		code := generateRandomCode()
		hashedCode, err := a.hasher.Hash(fmt.Sprintf("%d", code))
		if err != nil {
			log.Printf("failed to hash email code for user %s: %v", req.Email, err)
			return
		}
		err = a.redis.SaveEmailCode(bgCtx, req.Email, hashedCode, time.Minute * 5)
		if err != nil {
			log.Printf("failed to save email code for user %s: %v", req.Email, err)
			return
		}
		fmt.Println(code)
		// send email to user with the email code.
	}()

	return nil
}

func (a *AuthUsecase) VerifyUser(ctx context.Context, email string, code string) error {
	codeHash, err := a.redis.GetEmailCodeHash(ctx, email)
	if err != nil {
		if err == redis.Nil {
			return apperr.New(apperr.CodeConflict, http.StatusConflict, "email code is invalid")
		}
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	if err := a.hasher.CheckPasswordHash(code, codeHash); err != nil {
		return apperr.New(apperr.CodeConflict, http.StatusConflict, "email code is invalid")
	}

	err = a.authStore.VerifyUser(ctx, email)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	err = a.redis.DeleteEmailCode(ctx, email)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	// create session
	// accessToken := 
	

	// a.session.Create(ctx, )

	return nil
}

func generateRandomCode() int {
	min := 100000
	max := 999999

	// rand.IntN returns a non-negative random number in [0, n)
	// so rand.IntN(max-min+1) generates a number in [0, 900000]
	// Adding 'min' shifts the range to [100000, 999999]
	randomNumber := rand.IntN(max-min+1) + min
	fmt.Println(randomNumber)

	return randomNumber
}