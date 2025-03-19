package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/Jalenarms1/caters-go/internal/db"
	"github.com/Jalenarms1/caters-go/internal/types"
	"github.com/Jalenarms1/caters-go/internal/utils"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HandleLoginPage(w http.ResponseWriter, r *http.Request) *types.Error {

	t, err := template.ParseGlob("templates/login/*.html")
	if err != nil {
		return &types.Error{
			Err:        err,
			ReturnCode: http.StatusInternalServerError,
		}
	}

	t.ExecuteTemplate(w, "layout.html", nil)

	return nil
}

func HandleSignupPage(w http.ResponseWriter, r *http.Request) *types.Error {

	t, err := template.ParseGlob("templates/signup/*.html")
	if err != nil {
		return &types.Error{
			Err:        err,
			ReturnCode: http.StatusInternalServerError,
		}
	}

	var isError bool
	ctxErr := r.Context().Value("IsError")
	isError = ctxErr != nil

	fmt.Println(isError)

	t.ExecuteTemplate(w, "layout.html", &types.FormError{IsError: isError})

	return nil
}

func HandleSignupV2(w http.ResponseWriter, r *http.Request) *types.Error {

	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmedPassword := r.FormValue("confirmPassword")

	fmt.Println(email)
	fmt.Println(password)
	fmt.Println(confirmedPassword)

	if !strings.EqualFold(password, confirmedPassword) {

		// t, err := template.ParseGlob("templates/signup/*.html")
		// if err != nil {
		// 	return &types.Error{
		// 		Err:        err,
		// 		ReturnCode: http.StatusInternalServerError,
		// 	}
		// }

		http.Redirect(w, r.WithContext(context.WithValue(context.Background(), "IsError", true)), "/pages/signup", http.StatusPermanentRedirect)

		// t.ExecuteTemplate(w, "layout.html", &types.FormError{IsError: true})

		return nil
	}

	http.Redirect(w, r, "/pages/signup", http.StatusPermanentRedirect)

	return nil
}

func HandleNewAccount(w http.ResponseWriter, r *http.Request) *types.Error {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return &types.Error{
			Err:        err,
			ReturnCode: http.StatusBadRequest,
		}
	}
	defer r.Body.Close()

	fmt.Println(string(body))

	var acctInfo db.Account
	err = json.Unmarshal(body, &acctInfo)
	if err != nil {
		return &types.Error{
			Err:        err,
			ReturnCode: http.StatusBadRequest,
		}
	}

	if acctInfo.Email == "" || acctInfo.Password == "" {
		return &types.Error{
			Err:        errors.New("provide both an email and password"),
			ReturnCode: http.StatusNotAcceptable,
		}
	}
	fmt.Println("getting existing accts")
	existingAcct, _ := db.GetAccountByEmail(acctInfo.Email)
	if existingAcct != nil {
		return &types.Error{
			Err:        errors.New("account with this email already exists"),
			ReturnCode: http.StatusBadRequest,
		}
	}

	fmt.Println(acctInfo.Email, acctInfo.Password)

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(acctInfo.Password), 10)
	if err != nil {
		return &types.Error{
			Err:        err,
			ReturnCode: http.StatusInternalServerError,
		}
	}

	acctInfo.Password = string(hashedPass)

	fmt.Println(acctInfo)
	fmt.Println(string(hashedPass))

	uid, err := uuid.NewV4()
	if err != nil {
		return &types.Error{
			Err:        err,
			ReturnCode: http.StatusInternalServerError,
		}
	}

	acctInfo.Id = uid.String()

	token, err := utils.GenerateJWT(uid.String())
	if err != nil {
		return &types.Error{
			Err:        err,
			ReturnCode: http.StatusInternalServerError,
		}
	}

	isProd := os.Getenv("IS_DEV") != "true"

	cookie := &http.Cookie{
		Name:     "foodgo-auth",
		Value:    token,
		Path:     "/",
		Domain:   "caters-go.pages.dev",
		Expires:  time.Now().Add(3600 * time.Hour),
		Secure:   isProd,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(w, cookie)

	err = acctInfo.Insert()
	if err != nil {
		fmt.Println(err)
		return &types.Error{
			Err:        err,
			ReturnCode: http.StatusInternalServerError,
		}
	}
	fmt.Println(token)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})

	return nil
}

func HandleLoginV2(w http.ResponseWriter, r *http.Request) *types.Error {
	email := r.FormValue("email")
	password := r.FormValue("password")

	fmt.Println(email, password)

	http.Redirect(w, r, "/pages/login", http.StatusPermanentRedirect)

	return nil
}

func HandleLogin(w http.ResponseWriter, r *http.Request) *types.Error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return &types.Error{
			Err:        err,
			ReturnCode: http.StatusInternalServerError,
		}
	}
	defer r.Body.Close()

	var acctInfo *db.Account

	err = json.Unmarshal(body, &acctInfo)
	if err != nil {
		return &types.Error{
			Err:        err,
			ReturnCode: http.StatusInternalServerError,
		}
	}

	existingUser, _ := db.GetUserWPasswordByEmail(acctInfo.Email)
	if existingUser == nil {
		return &types.Error{
			Err:        errors.New("existing user not found with the email provided"),
			ReturnCode: http.StatusNotFound,
		}

	}
	fmt.Println(acctInfo.Password)
	isPasswordMatch := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(acctInfo.Password))

	if isPasswordMatch != nil {
		return &types.Error{
			Err:        errors.New("invalid credentials"),
			ReturnCode: http.StatusBadRequest,
		}

	}

	token, err := utils.GenerateJWT(existingUser.Id)
	if err != nil {
		return &types.Error{
			Err:        err,
			ReturnCode: http.StatusInternalServerError,
		}
	}

	fmt.Println(token)
	isDev := os.Getenv("IS_DEV") != "true"

	cookie := &http.Cookie{
		Name:     "foodgo-auth",
		Value:    token,
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().Add(3600 * time.Hour),
		Secure:   isDev,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(w, cookie)

	_ = json.NewEncoder(w).Encode((map[string]string{"token": token}))

	return nil
}

func HandleGetMe(w http.ResponseWriter, r *http.Request) *types.Error {
	fmt.Println("getme")
	uid := r.Context().Value(types.AuthKey)
	if uid == nil {
		return &types.Error{
			Err:        errors.New("no authentication"),
			ReturnCode: http.StatusBadRequest,
		}

	}

	fmt.Println(uid)
	acct, err := db.GetAccountById(uid.(string))
	if err != nil {
		return &types.Error{
			Err:        err,
			ReturnCode: http.StatusInternalServerError,
		}
	}

	fmt.Println(acct)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(acct)

	return nil
}

func HandleLogout(w http.ResponseWriter, r *http.Request) error {

	http.SetCookie(w, &http.Cookie{
		Name:     string(types.AuthKey),
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		Secure:   os.Getenv("IS_DEV") != "true",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	w.WriteHeader(http.StatusOK)

	return nil
}
