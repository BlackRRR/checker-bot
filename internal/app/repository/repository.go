package repository

import (
	"context"
	"github.com/BlackRRR/checker-bot/internal/app/model"
	"github.com/bots-empire/base-bot/msgs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"strings"
)

type Repository struct {
	Pool      *pgxpool.Pool
	globalBot *model.GlobalBot
	msgs      *msgs.Service
	ctx       context.Context
}

func NewRepository(pool *pgxpool.Pool, msgs *msgs.Service, globalBot *model.GlobalBot) *Repository {
	return &Repository{pool, globalBot, msgs, context.Background()}
}

func (r *Repository) CheckingTheUser(message *tgbotapi.Message) (*model.User, error) {
	rows, err := r.Pool.Query(r.ctx, `
SELECT id FROM checker.users 
	WHERE id = $1;`,
		message.From.ID)
	if err != nil {
		return nil, errors.Wrap(err, "get user")
	}

	users, err := readUsers(rows)
	if err != nil {
		return nil, errors.Wrap(err, "read user")
	}

	switch len(users) {
	case 0:
		user := createSimpleUser(message)
		if err := r.addNewUser(user); err != nil {
			return nil, errors.Wrap(err, "add new user")
		}
		return user, nil
	case 1:
		return users[0], nil
	default:
		return nil, model.ErrFoundTwoUsers
	}
}

func (r *Repository) addNewUser(u *model.User) error {
	_, err := r.Pool.Exec(r.ctx, `INSERT INTO checker.users VALUES ($1);`, u.ID)
	if err != nil {
		return errors.Wrap(err, "insert new user")
	}

	_ = r.msgs.SendSimpleMsg(u.ID, r.globalBot.LangText(r.globalBot.BotLang, "start_text"))

	return nil
}

func (r *Repository) SaveIncomeInfo(info *model.IncomeInfo, userName string) error {
	_, err := r.Pool.Exec(r.ctx, `INSERT INTO checker.income_info VALUES ($1,$2,$3,$4,$5,$6)`,
		info.UserID,
		userName,
		info.BotLink,
		info.BotName,
		info.IncomeSource,
		info.TypeBot)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil
		}
		return errors.Wrap(err, "insert income info")
	}

	return nil
}

func (r *Repository) GetURL() (string, error) {
	var url string
	err := r.Pool.QueryRow(r.ctx, `SELECT url FROM checker.url`).Scan(&url)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return "", nil
		}
		return "", err
	}

	return url, nil
}

func (r *Repository) GetText() (string, error) {
	var text string
	err := r.Pool.QueryRow(r.ctx, `SELECT url_text FROM checker.url`).Scan(&text)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return "", nil
		}
		return "", err
	}

	return text, nil
}

func (r *Repository) SetText(text string) error {
	_, err := r.Pool.Exec(r.ctx, `INSERT INTO checker.url (url_text) VALUES ($1)`, text)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) SetURL(url string) error {
	_, err := r.Pool.Exec(r.ctx, `INSERT INTO checker.url (url) VALUES ($1)`, url)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateURL(url string) error {
	_, err := r.Pool.Exec(r.ctx, `UPDATE checker.url SET url = $1`, url)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateText(urlText string) error {
	_, err := r.Pool.Exec(r.ctx, `UPDATE checker.url SET url_text = $1`, urlText)
	if err != nil {
		return err
	}

	return nil
}

func createSimpleUser(message *tgbotapi.Message) *model.User {
	if message.From.UserName != "" {
		return &model.User{
			ID: message.From.ID,
		}
	}

	return &model.User{
		ID: message.From.ID,
	}
}

func readUsers(rows pgx.Rows) ([]*model.User, error) {
	defer rows.Close()
	var users []*model.User

	for rows.Next() {
		user := &model.User{}

		if err := rows.Scan(
			&user.ID,
		); err != nil {
			return nil, errors.Wrap(err, model.ErrScanSqlRow.Error())
		}

		users = append(users, user)
	}

	return users, nil
}
