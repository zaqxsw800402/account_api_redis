package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"red/internal/encryption"
	"red/internal/urlsigner"
)

// Home displays the home page
func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "home", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

// LoginPage displays the login page
func (app *application) LoginPage(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "login", &templateData{}); err != nil {
		app.errorLog.Print(err)
	}
}

func (app *application) PostLoginPage(w http.ResponseWriter, r *http.Request) {
	app.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	id, err := app.DB.Authenticate(email, password)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), "userID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) Logout(w http.ResponseWriter, r *http.Request) {
	app.Session.Destroy(r.Context())
	app.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "forgot-password", &templateData{}); err != nil {
		app.errorLog.Print(err)
	}
}

func (app *application) ShowResetPassword(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	theURL := r.RequestURI
	testURL := fmt.Sprintf("%s%s", app.config.frontend, theURL)

	signer := urlsigner.Signer{
		Secret: []byte(app.config.secretKey),
	}

	valid := signer.VerifyToken(testURL)

	if !valid {
		app.errorLog.Println("Invalid url - tampering detected")
		return
	}

	// make sure not expired
	expired := signer.Expired(testURL, 60)
	if expired {
		app.errorLog.Println("Link expired")
		return
	}

	encryptor := encryption.Encryption{Key: []byte(app.config.secretKey)}
	encryptEmail, err := encryptor.Encrypt(email)
	if err != nil {
		app.errorLog.Println("Encryption failed")
		return
	}

	data := make(map[string]interface{})
	data["email"] = encryptEmail

	if err := app.renderTemplate(w, r, "reset-password", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Print(err)
	}
}

func (app *application) AllUsers(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "all-users", &templateData{}); err != nil {
		app.errorLog.Print(err)
	}
}

func (app *application) OneUser(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "one-user", &templateData{}); err != nil {
		app.errorLog.Print(err)
	}
}

func (app *application) UserProfile(w http.ResponseWriter, r *http.Request) {
	userID := app.Session.GetInt(r.Context(), "userID")
	profile, err := app.mongodb.ReadProfile(userID)
	if err != nil {
		log.Println(err.Error())
	}

	stringMap := make(map[string]string)
	stringMap["firstname"] = profile.FirstName
	stringMap["lastname"] = profile.LastName
	stringMap["email"] = profile.Email

	if err = app.renderTemplate(w, r, "user-profile", &templateData{StringMap: stringMap}); err != nil {
		app.errorLog.Print(err)
	}
}

func (app *application) AllCustomers(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "all-customers", &templateData{}); err != nil {
		app.errorLog.Print(err)
	}
}

func (app *application) OneCustomer(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "one-customer", &templateData{}); err != nil {
		app.errorLog.Print(err)
	}
}

func (app *application) AllAccountsWitCustomerID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	stringMap := make(map[string]string)
	stringMap["id"] = id
	if err := app.renderTemplate(w, r, "all-accounts-customerID", &templateData{StringMap: stringMap}); err != nil {
		app.errorLog.Print(err)
	}
}

func (app *application) AllAccounts(w http.ResponseWriter, r *http.Request) {

	if err := app.renderTemplate(w, r, "all-accounts", &templateData{}); err != nil {
		app.errorLog.Print(err)
	}
}

func (app *application) OneAccountBar(w http.ResponseWriter, r *http.Request) {
	//id := chi.URLParam(r, "id")
	//accountId := chi.URLParam(r, "account_id")
	//stringMap := make(map[string]string)
	//stringMap["id"] = id
	//stringMap["account_id"] = accountId
	if err := app.renderTemplate(w, r, "one-account-bar", &templateData{}); err != nil {
		app.errorLog.Print(err)
	}
}

func (app *application) OneAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	//accountId := chi.URLParam(r, "account_id")
	stringMap := make(map[string]string)
	stringMap["id"] = id
	//stringMap["account_id"] = accountId
	if err := app.renderTemplate(w, r, "one-account", &templateData{StringMap: stringMap}); err != nil {
		app.errorLog.Print(err)
	}
}

func (app *application) AllTransactions(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	accountId := chi.URLParam(r, "account_id")
	stringMap := make(map[string]string)
	stringMap["id"] = id
	stringMap["account_id"] = accountId
	if err := app.renderTemplate(w, r, "all-transactions", &templateData{StringMap: stringMap}); err != nil {
		app.errorLog.Print(err)
	}
}

func (app *application) Deposit(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "deposit", &templateData{}); err != nil {
		app.errorLog.Print(err)
	}
}

func (app *application) Withdrawal(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "withdrawal", &templateData{}); err != nil {
		app.errorLog.Print(err)
	}
}

func (app *application) Transfer(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "transfer", &templateData{}); err != nil {
		app.errorLog.Print(err)
	}
}
