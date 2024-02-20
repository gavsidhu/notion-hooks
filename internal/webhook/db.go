package webhook

import (
	"context"
	"log"
	"time"

	"github.com/gavsidhu/notion-hooks/internal/logging"
	"github.com/gavsidhu/notion-hooks/internal/models"
	"github.com/gavsidhu/notion-hooks/internal/notion"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rabbitmq/amqp091-go"
)

func GetWebhooksForProcessing(ctx context.Context, db *pgxpool.Pool, channel *amqp091.Channel) error {
	query := `
    WITH UpdatedWebhooks AS (
        UPDATE webhooks
        SET status = 'processing'
        WHERE last_polled + make_interval(mins => polling_interval) < NOW()
        AND is_active = true AND status = 'idle'
        RETURNING id
    )
    SELECT id FROM UpdatedWebhooks;`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var webhookId string
		err = rows.Scan(&webhookId)
		if err != nil {
			return err
		}

		err = channel.PublishWithContext(ctx, "", "processingQueue", false, false, amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(webhookId),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func GetPageIDsSnapshot(ctx context.Context, db *pgxpool.Pool, webhookId string) ([]string, error) {
	query := `SELECT * FROM notion_database_page_ids WHERE webhook_id = $1;`

	var webhookPageIDs models.NotionDatabasePageIDRow
	err := db.QueryRow(ctx, query, webhookId).Scan(&webhookPageIDs.ID, &webhookPageIDs.WebhookID, &webhookPageIDs.UserID, &webhookPageIDs.PageIDs, &webhookPageIDs.CreatedAt, &webhookPageIDs.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return webhookPageIDs.PageIDs, nil
}

func GetDatabaseDetailsSnapshot(ctx context.Context, db *pgxpool.Pool, webhookId string) (*notion.DatabaseQueryResponse, error) {
	query := `SELECT * FROM notion_database_details WHERE webhook_id = $1;`
	var webhookDatabaseDetails models.NotionDatabaseDetailRow
	err := db.QueryRow(ctx, query, webhookId).Scan(&webhookDatabaseDetails.ID, &webhookDatabaseDetails.WebhooksID, &webhookDatabaseDetails.UserID, &webhookDatabaseDetails.DatabasePageDetails, &webhookDatabaseDetails.CreatedAt, &webhookDatabaseDetails.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &webhookDatabaseDetails.DatabasePageDetails, nil
}

func SavePageIDsSnapshot(ctx context.Context, db *pgxpool.Pool, webhookId string, userId string, pageIDs []string) error {
	query := `INSERT INTO notion_database_page_ids (webhook_id,user_id, page_ids) VALUES ($1, $2, $3);`
	_, err := db.Exec(ctx, query, webhookId, userId, pageIDs)
	if err != nil {
		return err
	}

	return nil
}

func SaveDatabaseDetailsSnapshot(ctx context.Context, db *pgxpool.Pool, webhookId string, userId string, databaseDetails *notion.DatabaseQueryResponse) error {
	query := `INSERT INTO notion_database_details (webhook_id, user_id, database_page_details) VALUES ($1, $2, $3);`
	_, err := db.Exec(ctx, query, webhookId, userId, databaseDetails)
	if err != nil {
		return err
	}

	return nil
}

func UpdateSavedPageIDsSnapshot(ctx context.Context, db *pgxpool.Pool, webhookId string, userId string, pageIDs []string) error {
	query := `UPDATE notion_database_page_ids SET page_ids = $1 WHERE webhook_id = $2;`
	_, err := db.Exec(ctx, query, pageIDs, webhookId)
	if err != nil {
		return err
	}

	return nil
}

func GetURLForWebhook(ctx context.Context, db *pgxpool.Pool, webhookId string) (string, error) {
	query := `SELECT url FROM webhooks WHERE id = $1;`

	var url string
	err := db.QueryRow(ctx, query, webhookId).Scan(&url)
	if err != nil {
		return "", err
	}

	return url, nil
}

func UpdateSavedDatabaseDetailsSnapshot(ctx context.Context, db *pgxpool.Pool, webhookId string, userId string, databaseDetails *notion.DatabaseQueryResponse) error {
	query := `UPDATE notion_database_details SET database_page_details = $1 WHERE webhook_id = $2;`
	_, err := db.Exec(ctx, query, databaseDetails, webhookId)
	if err != nil {
		return err
	}

	return nil
}

func StartPollingDatabase(ctx context.Context, db *pgxpool.Pool, channel *amqp091.Channel) error {
	logging.Logger.Info("Starting polling database every 30 seconds.")

	ticker := time.NewTicker(30 * time.Second)

	for range ticker.C {
		rows, err := db.Query(ctx, `SELECT id FROM webhooks WHERE last_polled + make_interval(mins => polling_interval) < NOW() AND is_active = true;`)
		if err != nil {
			log.Println(err)
			return err
		}

		defer rows.Close()

		for rows.Next() {
			var webhookId string
			err = rows.Scan(&webhookId)
			if err != nil {
				log.Println(err)
				return err
			}

			err = channel.PublishWithContext(ctx, "", "proccessingQueue", false, false, amqp091.Publishing{
				ContentType: "application/json",
				Body:        []byte(webhookId),
			})
			if err != nil {
				return err
			}

		}
	}
	return nil
}

func GetWebhook(ctx context.Context, db *pgxpool.Pool, webhookId string) (models.Webhook, error) {
	query := `SELECT * FROM webhooks WHERE id = $1;`

	var webhook models.Webhook
	err := db.QueryRow(ctx, query, webhookId).Scan(&webhook.ID, &webhook.Name, &webhook.Description, &webhook.UserID, &webhook.URL, &webhook.Secret, &webhook.Events, &webhook.IsActive, &webhook.PollingInterval, &webhook.LastPolled, &webhook.Status, &webhook.NotionObjectID, &webhook.NotionObjectType, &webhook.CreatedAt, &webhook.UpdatedAt)
	if err != nil {
		return models.Webhook{}, err
	}

	return webhook, nil
}

func UpdateWebhookLastPolled(ctx context.Context, db *pgxpool.Pool, webhookId string) error {
	query := `UPDATE webhooks SET last_polled = NOW() WHERE id = $1;`
	_, err := db.Exec(ctx, query, webhookId)
	if err != nil {
		return err
	}

	return nil
}

func UpdateWebhookStatus(ctx context.Context, db *pgxpool.Pool, webhookId string, status string) error {
	query := `UPDATE webhooks SET status = $1 WHERE id = $2;`
	_, err := db.Exec(ctx, query, status, webhookId)
	if err != nil {
		return err
	}

	return nil
}

func GetNotionAccessToken(ctx context.Context, db *pgxpool.Pool, userId string) (string, error) {
	query := `SELECT access_token FROM notion_integrations WHERE user_id = $1;`

	var accessToken string
	err := db.QueryRow(ctx, query, userId).Scan(&accessToken)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
