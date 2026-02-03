package heartbeat

import (
	"net/http"
	"strconv"
)

var beat = false

func Beat(w http.ResponseWriter, r *http.Request) {
	if beat {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusGone)
	}
}

func ModifyBeat(w http.ResponseWriter, r *http.Request) {
	modifyBeat, err := strconv.ParseBool(r.URL.Query().Get("beat"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		beat = modifyBeat
		w.WriteHeader(http.StatusOK)
	}
}
