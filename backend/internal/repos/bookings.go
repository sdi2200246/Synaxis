package repos

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sdi2200246/synaxis/internal/entities"
	apperr "github.com/sdi2200246/synaxis/internal/error"
)

type BookingsRepo struct {
	db *pgxpool.Pool
}

func NewBookingsRepo(db *pgxpool.Pool)*BookingsRepo{
	return &BookingsRepo{db}
}


func (r *BookingsRepo) Create(ctx context.Context, booking entities.Booking) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return apperr.ErrInternal
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`UPDATE tickettype SET available = available - $1 WHERE id = $2`,
		booking.NumberOfTickets, booking.TicketTypeID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23514" {
			return apperr.ErrConflict
		}
		return apperr.ErrInternal
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO booking (id, user_id, ticket_type_id, number_of_tickets, total_cost, status, booked_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		booking.ID, booking.UserID, booking.TicketTypeID,
		booking.NumberOfTickets, booking.TotalCost, booking.Status, booking.BookedAt,
	)
	if err != nil {
		return apperr.ErrInternal
	}

	return tx.Commit(ctx)
}

func (r *BookingsRepo) GetByTicketTypeID(ctx context.Context, ticketTypeID uuid.UUID) ([]entities.Booking, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, ticket_type_id, number_of_tickets, total_cost, status, booked_at
		 FROM booking WHERE ticket_type_id = $1
		 ORDER BY booked_at DESC`,
		ticketTypeID,
	)
	if err != nil {
		return nil, apperr.ErrInternal
	}
	defer rows.Close()

	var bookings []entities.Booking
	for rows.Next() {
		var b entities.Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.TicketTypeID, &b.NumberOfTickets, &b.TotalCost, &b.Status, &b.BookedAt); err != nil {
			return nil, apperr.ErrInternal
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *BookingsRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]entities.UserBooking, error) {
	rows, err := r.db.Query(ctx,
		`SELECT b.id, b.ticket_type_id, tt.name, b.number_of_tickets, b.total_cost, b.status, b.booked_at,
		        e.id, e.title, e.start_datetime,
		        v.name, v.city, v.latitude, v.longitude
		 FROM booking b
		 JOIN tickettype tt ON b.ticket_type_id = tt.id
		 JOIN event e ON tt.event_id = e.id
		 JOIN venue v ON e.venue_id = v.id
		 WHERE b.user_id = $1
		 ORDER BY b.booked_at DESC`,
		userID,
	)
	if err != nil {
		return nil, apperr.ErrInternal
	}
	defer rows.Close()

	var bookings []entities.UserBooking
	for rows.Next() {
		var b entities.UserBooking
		if err := rows.Scan(
			&b.ID, &b.TicketTypeID, &b.TicketName, &b.NumberOfTickets, &b.TotalCost, &b.Status, &b.BookedAt,
			&b.EventID, &b.EventTitle, &b.EventStart,
			&b.VenueName, &b.VenueCity, &b.VenueLatitude, &b.VenueLongitude,
		); err != nil {
			return nil, apperr.ErrInternal
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *BookingsRepo) GetByEventID(ctx context.Context, eventID uuid.UUID) ([]entities.EventBooking, error) {
	rows, err := r.db.Query(ctx,
		`SELECT b.id, tt.name, b.number_of_tickets, b.total_cost, b.booked_at,
		        u.first_name || ' ' || u.last_name, u.email, u.phone
		 FROM booking b
		 JOIN tickettype tt ON b.ticket_type_id = tt.id
		 JOIN "user" u ON b.user_id = u.id
		 WHERE tt.event_id = $1
		 ORDER BY b.booked_at DESC`,
		eventID,
	)
	if err != nil {
		return nil, apperr.ErrInternal
	}
	defer rows.Close()

	var bookings []entities.EventBooking
	for rows.Next() {
		var b entities.EventBooking
		if err := rows.Scan(
			&b.ID, &b.TicketName, &b.NumberOfTickets, &b.TotalCost, &b.BookedAt,
			&b.AttendeeName, &b.AttendeeEmail, &b.AttendeePhone,
		); err != nil {
			return nil, apperr.ErrInternal
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *BookingsRepo) GetForExport(ctx context.Context, eventID uuid.UUID) ([]entities.ExportBooking, error) {
	rows, err := r.db.Query(ctx,
		`SELECT b.id, b.ticket_type_id, b.user_id, b.number_of_tickets, b.total_cost, b.status, b.booked_at
		 FROM booking b
		 JOIN tickettype tt ON b.ticket_type_id = tt.id
		 WHERE tt.event_id = $1
		 ORDER BY b.booked_at ASC`,
		eventID,
	)
	if err != nil {
		return nil, apperr.ErrInternal
	}
	defer rows.Close()

	var bookings []entities.ExportBooking
	for rows.Next() {
		var b entities.ExportBooking
		if err := rows.Scan(&b.ID, &b.TicketTypeID, &b.AttendeeID, &b.NumberOfTickets, &b.TotalCost, &b.Status, &b.BookedAt); err != nil {
			return nil, apperr.ErrInternal
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}