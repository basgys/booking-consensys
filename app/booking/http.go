package booking

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"time"

	"github.com/deixis/errors"
	"github.com/deixis/errors/httperrors"
	"github.com/deixis/pkg/httputil"
	"github.com/deixis/pkg/utc"
	"github.com/deixis/spine/net/http"
)

func (s *Service) HandleHTTP(srv *http.Server) {
	h := httpHandler{
		svc: s,
	}

	srv.HandleFunc("/booking/rooms", http.GET, h.listRooms)
	srv.HandleFunc("/booking/rooms/availabilities", http.GET, h.roomsAvailabilities)
	srv.HandleFunc("/booking/rooms/{rid}/availabilities", http.GET, h.roomAvailabilities)
	srv.HandleFunc("/booking/rooms/{rid}/reservations", http.GET, h.listRoomReservations)
	srv.HandleFunc("/booking/rooms/{rid}/reservations", http.POST, h.reserveRoom)
	srv.HandleFunc("/booking/rooms/{rid}/reservations/{id}", http.DELETE, h.cancelRoomReservation)
}

type httpHandler struct {
	svc *Service
}

func (h *httpHandler) listRooms(
	ctx context.Context, w http.ResponseWriter, req *http.Request,
) {
	rooms, err := h.svc.ListRooms(ctx)
	if err != nil {
		httperrors.Marshal(req.HTTP, w, err)
		return
	}
	w.JSON(http.StatusOK, struct {
		Rooms []*Room `json:"rooms"`
	}{
		Rooms: rooms,
	})
}

func (h *httpHandler) roomsAvailabilities(
	ctx context.Context, w http.ResponseWriter, req *http.Request,
) {
	query := req.HTTP.URL.Query()
	params := struct {
		From     utc.UTC `qs:"from"`
		To       utc.UTC `qs:"to"`
		Duration int64   `qs:"duration"`
	}{}
	if err := httputil.ParseQuery(query, &params); err != nil {
		httperrors.Marshal(req.HTTP, w, err)
		return
	}
	d := time.Duration(params.Duration) * time.Minute

	slots, err := h.svc.RoomsAvailabilities(ctx, params.From, params.To, d)
	if err != nil {
		httperrors.Marshal(req.HTTP, w, err)
		return
	}
	w.JSON(http.StatusOK, struct {
		Slots []*TimeInterval `json:"slots"`
	}{
		Slots: slots,
	})
}

func (h *httpHandler) roomAvailabilities(
	ctx context.Context, w http.ResponseWriter, req *http.Request,
) {
	query := req.HTTP.URL.Query()
	params := struct {
		From     utc.UTC `qs:"from"`
		To       utc.UTC `qs:"to"`
		Duration int64   `qs:"duration"`
	}{}
	if err := httputil.ParseQuery(query, &params); err != nil {
		httperrors.Marshal(req.HTTP, w, err)
		return
	}
	d := time.Duration(params.Duration) * time.Minute

	availabilities, err := h.svc.RoomAvailabilities(ctx, req.Params["rid"], params.From, params.To, d)
	if err != nil {
		httperrors.Marshal(req.HTTP, w, err)
		return
	}
	w.JSON(http.StatusOK, struct {
		Availabilities []*TimeInterval `json:"availabilities"`
	}{
		Availabilities: availabilities,
	})
}

func (h *httpHandler) listRoomReservations(
	ctx context.Context, w http.ResponseWriter, req *http.Request,
) {
	reservations, err := h.svc.ListRoomReservations(ctx, req.Params["rid"])
	if err != nil {
		httperrors.Marshal(req.HTTP, w, err)
		return
	}
	w.JSON(http.StatusOK, struct {
		Reservations []*Reservation `json:"reservations"`
	}{
		Reservations: reservations,
	})
}

type httpReserveRoomRequest struct {
	From  utc.UTC `qs:"from"`
	Hours int64   `qs:"hours"`
}

func (h *httpHandler) reserveRoom(
	ctx context.Context, w http.ResponseWriter, req *http.Request,
) {
	defer req.HTTP.Body.Close()
	r := httpReserveRoomRequest{}
	if err := unmarshalJSON(req.HTTP.Body, &r); err != nil {
		httperrors.Marshal(req.HTTP, w, err)
		return
	}

	res, err := h.svc.ReserveRoom(ctx, req.Params["rid"], r.From, r.Hours)
	if err != nil {
		httperrors.Marshal(req.HTTP, w, err)
		return
	}
	w.JSON(http.StatusCreated, res)
}

func (h *httpHandler) cancelRoomReservation(
	ctx context.Context, w http.ResponseWriter, req *http.Request,
) {
	err := h.svc.CancelRoomReservation(ctx, req.Params["rid"], req.Params["id"])
	if err != nil {
		httperrors.Marshal(req.HTTP, w, err)
		return
	}
	w.Head(http.StatusNoContent)
}

func unmarshalJSON(r io.Reader, v interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, v); err != nil {
		return errors.WithBad(err)
	}
	return nil
}
