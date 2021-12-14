package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {

	mux := chi.NewRouter()
	mux.Use(SessionLoad)

	mux.Get("/", app.Home)
	//mux.Get("/ws", app.WsEndPoint)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.Auth)
		//	mux.Get("/virtual-terminal", app.VirtualTerminal)
		//	mux.Get("/all-sales", app.AllSales)
		//	mux.Get("/all-subscriptions", app.AllSubscriptions)
		//	mux.Get("/sales/{id}", app.ShowSale)
		//	mux.Get("/subscriptions/{id}", app.ShowSubscription)
		mux.Get("/all-customers", app.AllCustomers)
		mux.Get("/all-customers/0", app.OneCustomer)

		mux.Get("/all-customers/accounts", app.AllAccounts)
		mux.Get("/all-customers/{id}/accounts", app.AllAccountsWitCustomerID)
		mux.Get("/all-customers/{id}/accounts/{account_id}", app.OneAccount)
		mux.Get("/all-customers/accounts/0", app.OneAccountBar)

		mux.Get("/all-customers/{id}/accounts/{account_id}/transactions", app.AllTransactions)

		mux.Get("/transfer", app.Transfer)
		mux.Get("/deposit", app.Deposit)

		mux.Get("/all-users", app.AllUsers)
	})

	// create a new user
	mux.Get("/all-users/{id}", app.OneUser)
	mux.Get("/forgot-password", app.ForgotPassword)

	// auth routes
	mux.Get("/login", app.LoginPage)
	mux.Post("/login", app.PostLoginPage)
	mux.Get("/logout", app.Logout)
	//mux.Get("/forgot-password", app.ForgotPassword)
	//mux.Get("/reset-password", app.ShowResetPassword)
	//
	//fileServer := http.FileServer(http.Dir("./static"))
	//mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
