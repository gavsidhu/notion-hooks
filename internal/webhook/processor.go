package webhook

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gavsidhu/notion-hooks/internal/logging"
	"github.com/gavsidhu/notion-hooks/internal/models"
	"github.com/gavsidhu/notion-hooks/internal/notion"
	"github.com/gavsidhu/notion-hooks/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

func ProccessWebhook(msg amqp091.Delivery, ch *amqp091.Channel, pool *pgxpool.Pool) {
	logging.Logger.WithFields(logrus.Fields{
		"webhook_id": string(msg.Body),
	}).Info("Received message from processing queue")

	webhook, err := GetWebhook(context.Background(), pool, string(msg.Body))
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error":        err,
			"message_body": string(msg.Body),
		}).Error("Error getting webhook from database")
		return
	}

	accesstoken, err := GetNotionAccessToken(context.Background(), pool, webhook.UserID)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error":      err,
			"webhook_id": webhook.ID,
			"user_id":    webhook.UserID,
		}).Error("Error getting Notion access token from database")
		return
	}

	notionClient := notion.NewNotionClient(accesstoken)

	if webhook.NotionObjectType == "database" {
		err := handleDatabaseEvents(context.Background(), pool, ch, notionClient, webhook.ID, webhook.UserID, webhook.NotionObjectID, webhook.Events)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"error":      err,
				"webhook_id": webhook.ID,
				"user_id":    webhook.UserID,
			}).Error("Error handling database events")
			return
		}
	} else {
		// TODO: Add support for handling page events
		logging.Logger.WithFields(logrus.Fields{
			"webhook_id": webhook.ID,
			"user_id":    webhook.UserID,
		}).Info("Notion object type not supported")
		return
	}

	if err := msg.Ack(false); err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error":      err,
			"webhook_id": webhook.ID,
		}).Error("Error acknowledging message")
	}

	logging.Logger.WithFields(logrus.Fields{
		"webhook_id": webhook.ID,
	}).Info("Successfully processed webhook")
}

func handleDatabaseEvents(ctx context.Context, pool *pgxpool.Pool, ch *amqp091.Channel, notionClient *notion.NotionClient, webhookId string, userId string, notionObjectID string, events []string) error {
	var eventsToSend []models.EventsToSend

	logging.Logger.WithFields(logrus.Fields{
		"webhookId":      webhookId,
		"userId":         userId,
		"notionObjectID": notionObjectID,
		"events":         events,
	}).Info("Starting to handle database events")

	newPages, err := notionClient.GetAllDatabasePages(ctx, notionObjectID)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error":          err,
			"notionObjectID": notionObjectID,
		}).Error("Error getting all pages from notion database")
		return err
	}

	if utils.StringInSlice("page.added", events) || utils.StringInSlice("page.deleted", events) {
		oldPageIDs, err := GetPageIDsSnapshot(ctx, pool, webhookId)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"error":     err,
				"webhookId": webhookId,
			}).Error("Error getting page ids from database")
			return err
		}

		newPageIDs, err := notionClient.GetAllDatabasePageIDs(ctx, newPages)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"error":     err,
				"webhookId": webhookId,
			}).Error("Error getting page ids from new pages")
			return err
		}

		added, deleted := findAddedOrDeletedPages(newPageIDs, oldPageIDs)

		oldPage, err := GetDatabaseDetailsSnapshot(ctx, pool, webhookId)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"error":     err,
				"webhookId": webhookId,
			}).Error("Error getting database details from database")
			return err
		}

		if utils.StringInSlice("page.added", events) && !utils.StringInSlice("page.deleted", events) {
			for _, id := range added {
				eventsToSend = append(eventsToSend, models.EventsToSend{
					Type:      "page.added",
					UserID:    userId,
					WebhookID: webhookId,
					Data: models.EventData{
						ObjectID:   id,
						ObjectType: "page",
						CreatedAt:  time.Now().Unix(),
					},
				})
			}
		} else if utils.StringInSlice("page.deleted", events) && !utils.StringInSlice("page.added", events) {
			for _, id := range deleted {
				eventsToSend = append(eventsToSend, models.EventsToSend{
					Type:      "page.deleted",
					UserID:    userId,
					WebhookID: webhookId,
					Data: models.EventData{
						ObjectID:   id,
						ObjectType: "page",
						CreatedAt:  time.Now().Unix(),
					},
				})
			}

		} else {
			for _, id := range added {
				eventsToSend = append(eventsToSend, models.EventsToSend{
					Type:      "page.added",
					UserID:    userId,
					WebhookID: webhookId,
					Data: models.EventData{
						ObjectID:   id,
						ObjectType: "page",
						CreatedAt:  time.Now().Unix(),
					},
				})
			}
			for _, id := range deleted {
				eventsToSend = append(eventsToSend, models.EventsToSend{
					Type:      "page.deleted",
					UserID:    userId,
					WebhookID: webhookId,
					Data: models.EventData{
						ObjectID:   id,
						ObjectType: "page",
						CreatedAt:  time.Now().Unix(),
					},
				})
			}
		}

		if utils.StringInSlice("page.updated", events) {
			updatedPages := compareSnapshotsByLastEditedTime(oldPage, newPages)
			for _, id := range updatedPages {
				eventsToSend = append(eventsToSend, models.EventsToSend{
					Type:      "page.updated",
					UserID:    userId,
					WebhookID: webhookId,
					Data: models.EventData{
						ObjectID:   id,
						ObjectType: "page",
						CreatedAt:  time.Now().Unix(),
					},
				})
			}
		}
	}

	newPageIDs, err := notionClient.GetAllDatabasePageIDs(ctx, newPages)

	go func() {
		err = UpdateSavedPageIDsSnapshot(ctx, pool, webhookId, userId, newPageIDs)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"error":     err,
				"webhookId": webhookId,
			}).Error("Error updating page ids to database")
			return
		}
	}()

	go func() {
		err = UpdateSavedDatabaseDetailsSnapshot(ctx, pool, webhookId, userId, newPages)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"error":     err,
				"webhookId": webhookId,
			}).Error("Error updating database details to database")
			return
		}
	}()

	go func() {
		err := UpdateWebhookLastPolled(ctx, pool, webhookId)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"error":     err,
				"webhookId": webhookId,
			}).Error("Error updating webhook last polled")
			return
		}
	}()

	for _, event := range eventsToSend {
		jsonEvent, err := json.Marshal(event)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"error":     err,
				"webhookId": webhookId,
			}).Error("Error marshalling event to json")
			return err
		}
		err = ch.PublishWithContext(context.Background(), "", "eventsQueue", false, false, amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(jsonEvent),
		})
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"error":     err,
				"webhookId": webhookId,
			}).Error("Error publishing event to events queue")
		}
	}

	return nil
}

func HandleInitialPolling(msg amqp091.Delivery, ch *amqp091.Channel, pool *pgxpool.Pool) {
	logging.Logger.WithFields(logrus.Fields{
		"message_body": string(msg.Body),
	}).Info("Received message from initial polling queue")

	var pollMsg models.InitialPollMessage
	err := json.Unmarshal(msg.Body, &pollMsg)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error":        err,
			"message_body": string(msg.Body),
		}).Error("Error unmarshalling message")
		return
	}

	accesstoken, err := GetNotionAccessToken(context.Background(), pool, pollMsg.UserID)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error":      err,
			"webhook_id": pollMsg.WebhookID,
			"user_id":    pollMsg.UserID,
		}).Error("Error getting Notion access token from database")
		return
	}

	notionClient := notion.NewNotionClient(accesstoken)

	pages, err := notionClient.GetAllDatabasePages(context.Background(), pollMsg.NotionObjectID)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error":      err,
			"webhook_id": pollMsg.WebhookID,
			"user_id":    pollMsg.UserID,
		}).Error("Error getting all pages from notion database")
		return
	}

	pageIDs, err := notionClient.GetAllDatabasePageIDs(context.Background(), pages)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error":      err,
			"webhook_id": pollMsg.WebhookID,
			"user_id":    pollMsg.UserID,
		}).Error("Error getting all page ids from notion database")
		return
	}

	go func() {
		err = SavePageIDsSnapshot(context.Background(), pool, pollMsg.WebhookID, pollMsg.UserID, pageIDs)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"error":      err,
				"webhook_id": pollMsg.WebhookID,
				"user_id":    pollMsg.UserID,
			}).Error("Error saving page ids to database")
			return
		}
	}()

	go func() {
		err = SaveDatabaseDetailsSnapshot(context.Background(), pool, pollMsg.WebhookID, pollMsg.UserID, pages)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"error":      err,
				"webhook_id": pollMsg.WebhookID,
				"user_id":    pollMsg.UserID,
			}).Error("Error saving database details to database")
			return
		}
	}()

	go func() {
		err := UpdateWebhookLastPolled(context.Background(), pool, pollMsg.WebhookID)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"error":      err,
				"webhook_id": pollMsg.WebhookID,
				"user_id":    pollMsg.UserID,
			}).Error("Error updating webhook last polled")
			return
		}
	}()

	go func() {
		err := UpdateWebhookStatus(context.Background(), pool, pollMsg.WebhookID, "idle")
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"error":      err,
				"webhook_id": pollMsg.WebhookID,
				"user_id":    pollMsg.UserID,
			}).Error("Error updating webhook status")
			return
		}
	}()

	if err := msg.Ack(false); err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error":      err,
			"webhook_id": pollMsg.WebhookID,
		}).Error("Error acknowledging message")
	}

}

func SendEventsToUser(msg amqp091.Delivery, ch *amqp091.Channel, pool *pgxpool.Pool) {
	logging.Logger.WithFields(logrus.Fields{
		"message_body": string(msg.Body),
	}).Info("Received message from events queue for sending events to user")

	var eventMsg models.EventsToSend
	err := json.Unmarshal(msg.Body, &eventMsg)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error":        err,
			"message_body": string(msg.Body),
		}).Error("Error unmarshalling message")
		return
	}

	var event models.Event
	event = models.Event{
		ID:        uuid.New().String(),
		Type:      eventMsg.Type,
		WebhookID: eventMsg.WebhookID,
		Data:      eventMsg.Data,
		CreatedAt: time.Now().Unix(),
	}

	url, err := GetURLForWebhook(context.Background(), pool, eventMsg.WebhookID)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error":      err,
			"webhook_id": eventMsg.WebhookID,
		}).Error("Error getting URL for webhook")
		return
	}

	err = SendEventToUser(url, event)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error":      err,
			"webhook_id": eventMsg.WebhookID,
		}).Error("Error sending event to user")
		return
	}

	if err := msg.Ack(false); err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error":      err,
			"webhook_id": eventMsg.WebhookID,
		}).Error("Error acknowledging message")
	}
}

func compareSnapshotsByLastEditedTime(previous, current *notion.DatabaseQueryResponse) []string {
	var updatedPages []string

	prevPagesByLastEditedTime := make(map[string]string)
	for _, page := range previous.Results {

		prevPagesByLastEditedTime[page.ID] = page.LastEditedTime

	}

	for _, currentPage := range current.Results {
		prevLastEditedTime, exists := prevPagesByLastEditedTime[currentPage.ID]
		if !exists {
			updatedPages = append(updatedPages, currentPage.ID)
			continue
		}

		if currentPage.LastEditedTime != prevLastEditedTime {
			updatedPages = append(updatedPages, currentPage.ID)
		}
	}

	return updatedPages
}

func findAddedOrDeletedPages(newPageIDs []string, oldPageIDs []string) ([]string, []string) {
	var added []string
	var deleted []string

	newSet := make(map[string]bool)
	oldSet := make(map[string]bool)

	for _, id := range newPageIDs {
		newSet[id] = true
	}
	for _, id := range oldPageIDs {
		oldSet[id] = true
	}

	for id := range newSet {
		if !oldSet[id] {
			added = append(added, id)
		}
	}

	for id := range oldSet {
		if !newSet[id] {
			deleted = append(deleted, id)
		}
	}

	return added, deleted
}
