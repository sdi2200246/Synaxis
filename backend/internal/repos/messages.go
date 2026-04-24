package repos

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/VauntDev/tqla"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sdi2200246/synaxis/internal/entities"
	"github.com/sdi2200246/synaxis/internal/error"
)

type MessagesRepo struct {
	db *pgxpool.Pool 
}

func NewMessagesRepo(db *pgxpool.Pool ) *MessagesRepo {
	return &MessagesRepo{db: db}
}

func (r *MessagesRepo) CreateConversation(ctx context.Context,conv entities.Conversation,organizer uuid.UUID,attendee uuid.UUID,)  error {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return apperr.ErrInternal
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
        INSERT INTO conversation (
            id, booking_id, created_at
        ) VALUES ($1, $2, $3)
    `,
		conv.ID,
		conv.BookingID,
		conv.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apperr.ErrConflict
		}

		slog.Error("ConversationRepo.Create failed",
			"error", err,
			"conversation_id", conv.ID,
			"booking_id", conv.BookingID,
		)
		return apperr.ErrInternal
	}

	_, err = tx.Exec(ctx, `
        INSERT INTO conversation_participant (conversation_id, user_id, role)
        VALUES 
            ($1, $2, 'organizer'),
            ($1, $3, 'attendee')
    `,
		conv.ID,
		organizer,
		attendee,
	)
	if err != nil {
		slog.Error("ConversationRepo.Create participants failed",
			"error", err,
			"conversation_id", conv.ID,
			"organizer_id", organizer,
			"attendee_id", attendee,
		)
		return apperr.ErrInternal
	}

	return tx.Commit(ctx)
}

func (r *MessagesRepo) GetConversationByBookingID(ctx context.Context, bookingID uuid.UUID) (entities.Conversation, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, booking_id, created_at
		FROM conversation
		WHERE booking_id = $1
	`, bookingID)

	var c entities.Conversation
	err := row.Scan(&c.ID, &c.BookingID, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.Conversation{}, apperr.ErrNotFound
		}
		return entities.Conversation{}, apperr.ErrInternal
	}
	return c, nil
}


func (r *MessagesRepo) Create(ctx context.Context, msg entities.Message) error {
    _, err := r.db.Exec(ctx, `
        INSERT INTO "message" (
            id, conversation_id, sender_id, content, is_read, status, sent_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7
        )`,
        msg.ID,
        msg.ConversationID,
        msg.SenderID,
        msg.Content,
        msg.IsRead,
        msg.Status,
        msg.SentAt,
    )

    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == "23505" {
            return apperr.ErrConflict
        }
        slog.Error("MessageRepo.Create failed", "error", err, "id", msg.ID)
        return apperr.ErrInternal
    }
    return nil
}


func (r *MessagesRepo) UpdateMessage(ctx context.Context, id uuid.UUID, mu entities.MessageUpdate) error {
    t, err := tqla.New(tqla.WithPlaceHolder(tqla.Dollar))
    if err != nil {
        return apperr.ErrInternal
    }

    query, args, err := t.Compile(`
        UPDATE "message" SET
            {{ if .Update.Status }} status = {{ .Update.Status }}, {{ end }}
            {{ if .Update.Content }} content = {{ .Update.Content }}, {{ end }}
            updated_at = now()
        WHERE id = {{ .ID }}
    `, struct {
        Update entities.MessageUpdate
        ID     uuid.UUID
    }{mu, id})

    if err != nil {
        slog.Error("UpdateMessage template compilation failed", "error", err)
        return apperr.ErrInternal
    }

    result, err := r.db.Exec(ctx, query, args...)
    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == "23505" {
            return apperr.ErrConflict
        }
        slog.Error("UpdateMessage execution failed", "error", err, "message_id", id)
        return apperr.ErrInternal
    }

    if result.RowsAffected() == 0 {
        return apperr.ErrNotFound
    }

    return nil
}


func (r *MessagesRepo) GetByConversationID(ctx context.Context, conversationID uuid.UUID) ([]entities.Message, error) {
    rows, err := r.db.Query(ctx, `
        SELECT 
            id, conversation_id, sender_id, content, is_read, status, sent_at, updated_at
        FROM "message"
        WHERE conversation_id = $1
        ORDER BY sent_at ASC
    `, conversationID)
    
    if err != nil {
        slog.Error("GetByConversationID query failed", "error", err, "conversation_id", conversationID)
        return nil, apperr.ErrInternal
    }
    defer rows.Close()

    var messages []entities.Message
    for rows.Next() {
        var msg entities.Message
        err := rows.Scan(
            &msg.ID,
            &msg.ConversationID,
            &msg.SenderID,
            &msg.Content,
            &msg.IsRead,
            &msg.Status,
            &msg.SentAt,
            &msg.UpdatedAt,
        )
        if err != nil {
            slog.Error("GetByConversationID scan failed", "error", err)
            return nil, apperr.ErrInternal
        }
        messages = append(messages, msg)
    }

    if err = rows.Err(); err != nil {
        return nil, apperr.ErrInternal
    }

    return messages, nil
}

func (r *MessagesRepo) GetUserConversations(ctx context.Context, userID uuid.UUID) ([]entities.Conversation, error) {
	rows, err := r.db.Query(ctx, `
		SELECT 
			c.id,
			c.booking_id,
			c.created_at
		FROM conversation c
		JOIN conversation_participant cp 
			ON cp.conversation_id = c.id

		LEFT JOIN LATERAL (
			SELECT m.sent_at
			FROM message m
			WHERE m.conversation_id = c.id
			ORDER BY m.sent_at DESC
			LIMIT 1
		) last_msg ON true

		WHERE cp.user_id = $1
		ORDER BY last_msg.sent_at DESC NULLS LAST
	`, userID)

	if err != nil {
		slog.Error("MessagesRepo.GetUserConversations query failed",
			"error", err,
			"user_id", userID,
		)
		return nil, apperr.ErrInternal
	}
	defer rows.Close()

	var conversations []entities.Conversation

	for rows.Next() {
		var c entities.Conversation
		err := rows.Scan(
			&c.ID,
			&c.BookingID,
			&c.CreatedAt,
		)
		if err != nil {
			slog.Error("MessagesRepo.GetUserConversations scan failed",
				"error", err,
			)
			return nil, apperr.ErrInternal
		}
		conversations = append(conversations, c)
	}

	if err = rows.Err(); err != nil {
		return nil, apperr.ErrInternal
	}

	return conversations, nil
}

func (r *MessagesRepo) GetParticipantsByConversationID(ctx context.Context,conversationID uuid.UUID,) ([]entities.ConvParticipant, error) {

	rows, err := r.db.Query(ctx, `
		SELECT conversation_id, user_id, role
		FROM conversation_participant
		WHERE conversation_id = $1
	`, conversationID)
	if err != nil {
		return nil, apperr.ErrInternal
	}
	defer rows.Close()

	var result []entities.ConvParticipant

	for rows.Next() {
		var p entities.ConvParticipant
		err := rows.Scan(
			&p.ConversationID,
			&p.UserId,
			&p.Role,
		)
		if err != nil {
			return nil, apperr.ErrInternal
		}
		result = append(result, p)
	}

	return result, nil
}

func (r *MessagesRepo) GetConversationByID(ctx context.Context, id uuid.UUID) (entities.Conversation, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, booking_id, created_at
		FROM conversation
		WHERE id = $1
	`, id)

	var c entities.Conversation
	err := row.Scan(&c.ID, &c.BookingID, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.Conversation{}, apperr.ErrNotFound
		}
		return entities.Conversation{}, apperr.ErrInternal
	}

	return c, nil
}
func (r *MessagesRepo) GetUnreadMessagesCountByUser(ctx context.Context,userID uuid.UUID,) (map[uuid.UUID]int, error) {

	rows, err := r.db.Query(ctx, `
		SELECT
			m.conversation_id,
			COUNT(*) AS unread_count
		FROM message m
		JOIN conversation_participant cp
			ON cp.conversation_id = m.conversation_id
		WHERE cp.user_id = $1
			AND m.is_read = false
			AND m.sender_id != $1
		GROUP BY m.conversation_id
	`, userID)
	if err != nil {
		slog.Error("MessagesRepo.GetUnreadMessagesCountByUser query failed",
			"error", err,
			"user_id", userID,
		)
		return nil, apperr.ErrInternal
	}
	defer rows.Close()

	result := make(map[uuid.UUID]int)

	for rows.Next() {
		var convID uuid.UUID
		var count int

		err := rows.Scan(&convID, &count)
		if err != nil {
			slog.Error("MessagesRepo.GetUnreadMessagesCountByUser scan failed",
				"error", err,
				"user_id", userID,
			)
			return nil, apperr.ErrInternal
		}

		result[convID] = count
	}

	if err := rows.Err(); err != nil {
		return nil, apperr.ErrInternal
	}

	return result, nil
}

func (r *MessagesRepo) GetMessagesByConversationID(ctx context.Context,conversationID uuid.UUID,) ([]entities.Message, error) {

	rows, err := r.db.Query(ctx, `
		SELECT 
			id, conversation_id, sender_id, content, is_read, status, sent_at, updated_at
		FROM message
		WHERE conversation_id = $1
		ORDER BY sent_at ASC
	`, conversationID)
	if err != nil {
		slog.Error("MessagesRepo.GetByConversationID query failed",
			"error", err,
			"conversation_id", conversationID,
		)
		return nil, apperr.ErrInternal
	}
	defer rows.Close()

	var messages []entities.Message

	for rows.Next() {
		var m entities.Message
		err := rows.Scan(
			&m.ID,
			&m.ConversationID,
			&m.SenderID,
			&m.Content,
			&m.IsRead,
			&m.Status,
			&m.SentAt,
			&m.UpdatedAt,
		)
		if err != nil {
			slog.Error("MessagesRepo.GetByConversationID scan failed",
				"error", err,
			)
			return nil, apperr.ErrInternal
		}
		messages = append(messages, m)
	}

	if err := rows.Err(); err != nil {
		return nil, apperr.ErrInternal
	}

	return messages, nil
}

func (r *MessagesRepo) GetParticipantsByConversationIDs(ctx context.Context,conversationIDs []uuid.UUID,) (map[uuid.UUID][]entities.ConvParticipant, error) {

	rows, err := r.db.Query(ctx, `
		SELECT conversation_id, user_id, role
		FROM conversation_participant
		WHERE conversation_id = ANY($1)
	`, conversationIDs)
	if err != nil {
		slog.Error("GetParticipantsByConversationIDs query failed",
			"error", err,
		)
		return nil, apperr.ErrInternal
	}
	defer rows.Close()

	result := make(map[uuid.UUID][]entities.ConvParticipant)

	for rows.Next() {
		var p entities.ConvParticipant
		var convID uuid.UUID

		err := rows.Scan(
			&convID,
			&p.UserId,
			&p.Role,
		)
		if err != nil {
			slog.Error("GetParticipantsByConversationIDs scan failed",
				"error", err,
			)
			return nil, apperr.ErrInternal
		}

		result[convID] = append(result[convID], p)
	}

	if err := rows.Err(); err != nil {
		return nil, apperr.ErrInternal
	}

	return result, nil
}


func (r *MessagesRepo) MarkAsReadUpToMessage(ctx context.Context,conversationID uuid.UUID,userID uuid.UUID,lastMessageTime time.Time,) error {
	_, err := r.db.Exec(ctx, `
		UPDATE message
		SET is_read = true,
		    updated_at = now()
		WHERE conversation_id = $1
		  AND sender_id != $2
		  AND is_read = false
		  AND sent_at <= $3
	`, conversationID, userID, lastMessageTime)

	if err != nil {
		slog.Error("MessagesRepo.MarkAsReadUpToMessage failed",
			"error", err,
		)
		return apperr.ErrInternal
	}

	return nil
}

func (r *MessagesRepo) GetMessageByID(ctx context.Context, id uuid.UUID) (entities.Message, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, conversation_id, sender_id, content, is_read, status, sent_at, updated_at
		FROM message
		WHERE id = $1
	`, id)

	var m entities.Message
	err := row.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.IsRead, &m.Status, &m.SentAt, &m.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.Message{}, apperr.ErrNotFound
		}
		return entities.Message{}, apperr.ErrInternal
	}

	return m, nil
}