package webhook

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gavsidhu/notion-hooks/internal/logging"
	"github.com/gavsidhu/notion-hooks/internal/models"
	"github.com/sirupsen/logrus"
)

func SendEventToUser(url string, event models.Event) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(eventBytes))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		// Log error reading response body
		logging.Logger.WithFields(logrus.Fields{
			"error": err,
			"url":   url,
		}).Error("Failed to read response body from user event")
		return err
	}

	if response.StatusCode != 200 {
		logging.Logger.WithFields(logrus.Fields{
			"event":      event,
			"response":   string(bodyBytes),
			"url":        url,
			"status":     response.StatusCode,
			"webhook_id": event.WebhookID,
		}).Warn("Failed to send event to user's endpoint.")

		// TODO: Implement retry logic and add to log table
		return nil
	}

	logging.Logger.WithFields(logrus.Fields{
		"event":    event,
		"url":      url,
		"status":   response.StatusCode,
		"response": string(bodyBytes),
	}).Info("Successfully sent event to user's endpoint.")

	// TODO add to log table

	return nil
}
