package yin

import (
	"net/http"

	newrelic "github.com/newrelic/go-agent"
)

type NewRelic struct {
	App newrelic.Application
}

func InitNewRelic(app newrelic.Application) *NewRelic {
	return &NewRelic{app}
}

func (nr *NewRelic) CustomEvent(event string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if nr.App == nil {
				next.ServeHTTP(w, r)
				return
			}

			txn := newrelic.FromContext(r.Context())
			if txn != nil {
				txn.Ignore()
			}

			customTxn := nr.App.StartTransaction(event, w, r)
			r = newrelic.RequestWithTransactionContext(r, customTxn)
			defer customTxn.End()

			next.ServeHTTP(customTxn, r)
		})
	}
}

func (nr *NewRelic) EventFromURLPath(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if nr.App == nil {
			next.ServeHTTP(w, r)
			return
		}
		txn := nr.App.StartTransaction(r.URL.Path, w, r)
		defer txn.End()
		r = newrelic.RequestWithTransactionContext(r, txn)
		next.ServeHTTP(txn, r)
	})
}

func (nr *NewRelic) Ignore(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if nr.App == nil {
			next.ServeHTTP(w, r)
			return
		}
		var txn newrelic.Transaction
		txn = newrelic.FromContext(r.Context())
		if txn != nil {
			txn.Ignore()
		}

		next.ServeHTTP(txn, r)
	})
}
