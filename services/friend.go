package services

import (
	"database/sql"
	"fmt"
	"log"
	r "math/rand"
	"strings"
	"time"

	"github.com/ViniNepo/secretfriend/domain"
)

type FriendService interface {
	Create(friend domain.Friend) (int, error)
	Reminder() error
	Shuffle() error
	Validate(request domain.ValidateRequest) error
}

type friendService struct {
	emailService EmailService
	db           *sql.DB
}

func NewFriendService(emailService EmailService, db *sql.DB) FriendService {
	return &friendService{
		emailService: emailService,
		db:           db,
	}
}

func (f *friendService) Create(friend domain.Friend) (int, error) {
	query := `
        INSERT INTO friends (name, email, description, requirement, select_friend, validate_code, is_valid) 
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id`

	var id int
	err := f.db.QueryRow(
		query,
		friend.Name,
		friend.Email,
		friend.Description,
		friend.Requirement,
		friend.SelectFriend,
		friend.ValidateCode,
		friend.IsValid,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to execute statement: %w", err)
	}

	data := map[string]string{
		"name": friend.Name,
		"code": friend.ValidateCode,
	}
	content := GenerateEmail(validationCodeEmailTemplate, data)

	// Send the email
	go f.emailService.SendEmail(friend.Email, "Amigo chocolate - Código de validação", content)

	log.Println("email sent successfully!")

	return id, nil
}

func (f *friendService) Reminder() error {
	query := `SELECT id, name, email, description, requirement, select_friend, validate_code, is_valid FROM friends`

	rows, err := f.db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query friends: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var friend domain.Friend

		err := rows.Scan(
			&friend.Id,
			&friend.Name,
			&friend.Email,
			&friend.Description,
			&friend.Requirement,
			&friend.SelectFriend,
			&friend.ValidateCode,
			&friend.IsValid,
		)
		if err != nil {
			return fmt.Errorf("failed to scan friend: %w", err)
		}

		var selectedFriend domain.Friend
		if friend.SelectFriend != nil {
			err = f.db.QueryRow(
				"SELECT id, name, email, description, requirement, select_friend, validate_code, is_valid FROM friends WHERE id = $1",
				*friend.SelectFriend,
			).Scan(
				&selectedFriend.Id,
				&selectedFriend.Name,
				&selectedFriend.Email,
				&selectedFriend.Description,
				&selectedFriend.Requirement,
				&selectedFriend.SelectFriend,
				&selectedFriend.ValidateCode,
				&selectedFriend.IsValid,
			)
			if err != nil && err != sql.ErrNoRows {
				return fmt.Errorf("failed to fetch selected friend: %w", err)
			}
		}

		data := map[string]string{
			"name": friend.Name,
		}
		content := GenerateEmail(reminderEmail, data)

		// Send the email
		go f.emailService.SendEmail(friend.Email, "Amigo chocolate - Lembrete do evento", content)

		log.Println("email sent successfully!")
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("error iterating friends: %w", err)
	}

	return nil
}

func (f *friendService) Shuffle() error {
	query := `SELECT id, name, email, description, requirement, select_friend, validate_code, is_valid FROM friends`

	rows, err := f.db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query friends: %w", err)
	}
	defer rows.Close()

	var friends []domain.Friend
	for rows.Next() {
		var friend domain.Friend
		err := rows.Scan(
			&friend.Id,
			&friend.Name,
			&friend.Email,
			&friend.Description,
			&friend.Requirement,
			&friend.SelectFriend,
			&friend.ValidateCode,
			&friend.IsValid,
		)
		if err != nil {
			return fmt.Errorf("failed to scan friend: %w", err)
		}
		friends = append(friends, friend)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating friends: %w", err)
	}

	if len(friends) < 2 {
		return fmt.Errorf("not enough friends to perform shuffle")
	}

	r.Seed(time.Now().UnixNano())
	r.Shuffle(len(friends), func(i, j int) {
		friends[i], friends[j] = friends[j], friends[i]
	})

	for i, friend := range friends {
		secretFriend := friends[(i+1)%len(friends)]

		updateQuery := `UPDATE friends SET select_friend = $1 WHERE id = $2`
		_, err := f.db.Exec(updateQuery, secretFriend.Id, friend.Id)
		if err != nil {
			return fmt.Errorf("failed to update select_friend for %s: %w", friend.Name, err)
		}

		data := map[string]string{
			"name":        friend.Name,
			"select_name": secretFriend.Name,
			"description": secretFriend.Description,
			"requirement": secretFriend.Requirement,
		}
		content := GenerateEmail(friendDrawnEmailTemplate, data)

		// Send the email
		go f.emailService.SendEmail(friend.Email, "Amigo chocolate - Amigo sorteado", content)

		log.Println("email sent successfully!")
	}

	return nil
}

func (f *friendService) Validate(request domain.ValidateRequest) error {
	var friend domain.Friend
	query := "SELECT id, name, email, validate_code FROM friends WHERE id = $1"
	err := f.db.QueryRow(query, request.FriendID).Scan(&friend.Id, &friend.Name, &friend.Email, &friend.ValidateCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("friend with ID %d not found", request.FriendID)
		}
		return fmt.Errorf("failed to query friend: %w", err)
	}

	if friend.ValidateCode != request.Code {
		return fmt.Errorf("invalid code for friend ID %d", request.FriendID)
	}

	updateQuery := "UPDATE friends SET is_valid = true WHERE id = $1"
	_, err = f.db.Exec(updateQuery, request.FriendID)
	if err != nil {
		return fmt.Errorf("failed to update friend validation status: %w", err)
	}

	return nil
}

func GenerateEmail(template string, data map[string]string) string {
	for key, value := range data {
		template = strings.ReplaceAll(template, "{{"+key+"}}", value)
	}
	return template
}

var reminderEmail string = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f7f7f7; color: #333; }
    .container { max-width: 600px; margin: 20px auto; padding: 20px; background-color: #ffffff; border-radius: 8px; box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); }
    .header { text-align: center; padding: 10px 0; background-color: #4CAF50; color: #ffffff; border-radius: 8px 8px 0 0; }
    .header h1 { margin: 0; font-size: 24px; }
    .content { padding: 20px; }
    .footer { text-align: center; margin-top: 20px; font-size: 14px; color: #666; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>Lembrete</h1>
    </div>
    <div class="content">
      <p>Olá, <strong>{{name}}</strong>,</p>
      <p>Este é um lembrete para você participar do nosso evento de amigo secreto!</p>
      <p>Data do sorteio: <strong>14 de dezembro</strong></p>
    </div>
    <div class="footer">
      <p>Estamos ansiosos pela sua participação!</p>
    </div>
  </div>
</body>
</html>
`
var validationCodeEmailTemplate string = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f7f7f7; color: #333; }
    .container { max-width: 600px; margin: 20px auto; padding: 20px; background-color: #ffffff; border-radius: 8px; box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); }
    .header { text-align: center; padding: 10px 0; background-color: #4CAF50; color: #ffffff; border-radius: 8px 8px 0 0; }
    .header h1 { margin: 0; font-size: 24px; }
    .content { padding: 20px; }
    .footer { text-align: center; margin-top: 20px; font-size: 14px; color: #666; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>Código de Validação</h1>
    </div>
    <div class="content">
      <p>Olá, <strong>{{name}}</strong>,</p>
      <p>Seu código de validação é:</p>
      <h2>{{code}}</h2>
    </div>
    <div class="footer">
      <p>Use este código para validar sua participação no sistema.</p>
    </div>
  </div>
</body>
</html>
`
var friendDrawnEmailTemplate string = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f7f7f7; color: #333; }
    .container { max-width: 600px; margin: 20px auto; padding: 20px; background-color: #ffffff; border-radius: 8px; box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); }
    .header { text-align: center; padding: 10px 0; background-color: #4CAF50; color: #ffffff; border-radius: 8px 8px 0 0; }
    .header h1 { margin: 0; font-size: 24px; }
    .content { padding: 20px; }
    .footer { text-align: center; margin-top: 20px; font-size: 14px; color: #666; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>Amigo Sorteado</h1>
    </div>
    <div class="content">
      <p>Olá, <strong>{{name}}</strong>,</p>
      <p>O amigo sorteado para você é:</p>
      <h2>{{select_name}} - {{description}}</h2>
      <p>Preferência do chocolate:</p>
	  <p>{{requirement}}</p>
    </div>
    <div class="footer">
      <p>Divirta-se escolhendo o presente!</p>
    </div>
  </div>
</body>
</html>
`
