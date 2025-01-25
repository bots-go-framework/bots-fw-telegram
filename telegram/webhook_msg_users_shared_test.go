package telegram

import (
	"testing"

	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
)

func Test_tgWebhookUsersSharedMessage_GetSharedUsers(t *testing.T) {
	tests := []struct {
		name                string
		webhookMessage      tgWebhookUsersSharedMessage
		wantSharedUsers     []botinput.SharedUserMessageItem
		expectPanic         bool
		expectedSharedUsers []botinput.SharedUserMessageItem
	}{
		{
			name: "should_pass",
			webhookMessage: tgWebhookUsersSharedMessage{
				tgWebhookMessage: tgWebhookMessage{
					message: &tgbotapi.Message{
						UsersShared: &tgbotapi.UsersShared{
							Users: []tgbotapi.SharedUser{
								{
									UserID:    123,
									FirstName: "Jack",
									LastName:  "Smith",
									Photo: []tgbotapi.PhotoSize{
										{
											FileID:       "file123",
											FileUniqueID: "f123",
											Width:        640,
											Height:       460,
											FileSize:     1025,
										},
									},
								},
							},
						},
					},
				},
			},
			expectedSharedUsers: []botinput.SharedUserMessageItem{
				tgSharedUser{
					tgbotapi.SharedUser{
						UserID:    123,
						FirstName: "Jack",
						LastName:  "Smith",
						Photo: []tgbotapi.PhotoSize{
							{
								FileID:       "file123",
								FileUniqueID: "f123",
								Width:        640,
								Height:       460,
								FileSize:     1025,
							},
						},
					},
				},
			},
		},
		{
			name: "should_panic",
			webhookMessage: tgWebhookUsersSharedMessage{
				tgWebhookMessage: tgWebhookMessage{
					message: &tgbotapi.Message{
						UsersShared: nil,
					},
				},
			},
			expectPanic: true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("The code did not panic as expected")
					}
				}()
			}
			gotSharedUsers := tt.webhookMessage.GetSharedUsers()
			if gotSharedUsers == nil {
				t.Fatalf("GetSharedUsers() returned nil")
			}
			if len(gotSharedUsers) != len(tt.expectedSharedUsers) {
				t.Fatalf("GetSharedUsers() returned wrong number of shared users: got %d want %d", len(gotSharedUsers), len(tt.expectedSharedUsers))
			}
			for i, actual := range gotSharedUsers {
				if expected := tt.expectedSharedUsers[i].GetBotUserID(); actual.GetBotUserID() != expected {
					t.Errorf("actual bot user ID different from expected: %s != %s", actual.GetBotUserID(), expected)
				}
				if expected := tt.expectedSharedUsers[i].GetUsername(); actual.GetUsername() != expected {
					t.Errorf("actual bot username different from expected: %s != %s", actual.GetUsername(), expected)
				}
				if expected := tt.expectedSharedUsers[i].GetFirstName(); actual.GetFirstName() != expected {
					t.Errorf("actual first name different from expected: %s != %s", actual.GetFirstName(), expected)
				}
				if expected := tt.expectedSharedUsers[i].GetLastName(); actual.GetLastName() != expected {
					t.Errorf("actual last name different from expected: %s != %s", actual.GetFirstName(), expected)
				}
			}
		})
	}
}
