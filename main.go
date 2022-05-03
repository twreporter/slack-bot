package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	// Load env
	godotenv.Load(".env")
	token := os.Getenv("SLACK_AUTH_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")

	// Set debug to true for dev
	client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))
	// Enable Socket Mode
	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(true),
		// Custom logger
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx context.Context, client *slack.Client, socketClient *socketmode.Client) {
		for {
			select {
			case <-ctx.Done():
				log.Println("Shutting down socketmode listener")
				return
			case event := <-socketClient.Events:
				switch event.Type {
				// EventAPI
				case socketmode.EventTypeEventsAPI:
					eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						log.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
						continue
					}
					socketClient.Ack(*event.Request)
					// log.Println(eventsAPIEvent)
					err := handleEventMessage(eventsAPIEvent, client)
					if err != nil {
						//TODO: actual err handeling
						log.Fatal(err)
					}

				// Slash Commands
				case socketmode.EventTypeSlashCommand:
					command, ok := event.Data.(slack.SlashCommand)
					if !ok {
						log.Printf("Could not type cast the message to a SlashCommand: %v\n", command)
						continue
					}
					payload, err := handleSlashCommand(command, client)
					if err != nil {
						log.Fatal(err)
					}
					socketClient.Ack(*event.Request, payload)

				// Interactive
				case socketmode.EventTypeInteractive:
					interaction, ok := event.Data.(slack.InteractionCallback)
					if !ok {
						log.Printf("Could not type cast the message to a Interaction callback: %v\n", interaction)
						continue
					}

					err := handleInteractionEvent(interaction, client)
					if err != nil {
						log.Fatal(err)
					}
					socketClient.Ack(*event.Request)
				}
			}

		}
	}(ctx, client, socketClient)

	socketClient.Run()
}

func handleEventMessage(event slackevents.EventsAPIEvent, client *slack.Client) error {
	switch event.Type {
	case slackevents.CallbackEvent:
		innerEvent := event.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			err := handleAppMentionEvent(ev, client)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("unsupported event type")
	}
	return nil
}

func handleAppMentionEvent(event *slackevents.AppMentionEvent, client *slack.Client) error {
	user, err := client.GetUserInfo(event.User)
	if err != nil {
		return err
	}

	attachment := slack.Attachment{}
	attachment.Fields = []slack.AttachmentField{
		{
			Title: "Date",
			Value: time.Now().String(),
		}, {
			Title: "Initializer",
			Value: user.Name,
		},
	}

	text := strings.ToLower(event.Text)
	if strings.Contains(text, "hello") {
		attachment.Text = fmt.Sprintf("Hello %s", user.Name)
		attachment.Pretext = "Greetings"
		attachment.Color = "#4af030"
	} else {
		attachment.Text = fmt.Sprintf("How can I help you %s?", user.Name)
		attachment.Pretext = "`/sync [id]` for netiCRM sync-tool"
		attachment.Color = "#3d3d3d"
	}

	//TODO

	_, _, err = client.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}

	return nil
}

func handleSlashCommand(command slack.SlashCommand, client *slack.Client) (interface{}, error) {
	switch command.Command {
	case "/sync":
		return nil, handleSyncCommand(command, client)
	case "/test":
		return handleCheckboxCommand(command, client)
	}

	return nil, nil
}

func handleSyncCommand(command slack.SlashCommand, client *slack.Client) error {
	attachment := slack.Attachment{}
	attachment.Fields = []slack.AttachmentField{
		{
			Title: "ID",
			Value: command.Text,
		}, {
			Title: "Date",
			Value: time.Now().String(),
		}, {
			Title: "Initializer",
			Value: command.UserName,
		},
	}

	attachment.Text = fmt.Sprintf("Syncing to netiCRM...")
	attachment.Color = "#AA0000"

	_, _, err := client.PostMessage(command.ChannelID, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}

	return nil
}

func handleCheckboxCommand(command slack.SlashCommand, client *slack.Client) (interface{}, error) {
	attachment := slack.Attachment{}

	checkbox := slack.NewCheckboxGroupsBlockElement("answer",
		slack.NewOptionBlockObject("yes", &slack.TextBlockObject{Text: "Yes", Type: slack.MarkdownType}, &slack.TextBlockObject{Text: "Will execute the sync-tool", Type: slack.MarkdownType}),
		slack.NewOptionBlockObject("no", &slack.TextBlockObject{Text: "No", Type: slack.MarkdownType}, &slack.TextBlockObject{Text: "Cancel command", Type: slack.MarkdownType}),
	)

	accessory := slack.NewAccessory(checkbox)
	attachment.Blocks = slack.Blocks{
		BlockSet: []slack.Block{
			slack.NewSectionBlock(
				&slack.TextBlockObject{
					Type: slack.MarkdownType,
					Text: "ID correct?",
				},
				nil,
				accessory,
			),
		},
	}

	attachment.Text = "Confirmation"
	attachment.Color = "#4af030"
	return attachment, nil
}

func handleInteractionEvent(interaction slack.InteractionCallback, client *slack.Client) error {
	log.Printf("The action called is: %s\n", interaction.ActionID)
	log.Printf("The response was of type: %s\n", interaction.Type)
	switch interaction.Type {
	case slack.InteractionTypeBlockActions:
		for _, action := range interaction.ActionCallback.BlockActions {
			log.Printf("%+v", action)
			log.Println("Selected option: ", action.SelectedOptions)

		}

	default:

	}

	return nil
}
