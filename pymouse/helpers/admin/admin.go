package admin

import (
	"fmt"
	"pymouse/pymouse/helpers/utils"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
)

type AdminPerm uint64

const (
	PermDeleteMessages AdminPerm = 1 << iota
	PermRestrictMembers
	PermPromoteMembers
	PermChangeInfo
	PermInviteUsers
	PermPinMessages
)

func sendAdminError(
	bot *telego.Bot,
	ctx *telegohandler.Context,
	update telego.Update,
	chatID int64,
	isCallback bool,
	text string,
) {
	if isCallback {
		bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            utils.StripHTML(text),
			ShowAlert:       true,
		})
	} else {
		bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID:    telegoutil.ID(chatID),
			Text:      text,
			ParseMode: "HTML",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
	}
}

func CheckAdmin(
	ctx *telegohandler.Context,
	update telego.Update,
	bot *telego.Bot,
	required AdminPerm,
	l func(string) string,
	acceptInPrivate ...bool,
) bool {
	allowPrivate := true
	if len(acceptInPrivate) > 0 {
		allowPrivate = acceptInPrivate[0]
	}

	var (
		chatID     int64
		userID     int64
		isCallback bool
	)

	if update.Message != nil {
		chatID = update.Message.Chat.ID
		userID = update.Message.From.ID

		if update.Message.Chat.Type == telego.ChatTypePrivate {
			return allowPrivate
		}
	}

	if update.CallbackQuery != nil {
		chatID = update.CallbackQuery.Message.GetChat().ID
		userID = update.CallbackQuery.From.ID

		if update.CallbackQuery.Message.GetChat().Type == telego.ChatTypePrivate {
			return allowPrivate
		}
		isCallback = true
	}

	if chatID == 0 || userID == 0 {
		return false
	}

	member, err := bot.GetChatMember(ctx, &telego.GetChatMemberParams{
		ChatID: telegoutil.ID(chatID),
		UserID: userID,
	})
	if err != nil {
		return false
	}

	if member.MemberStatus() == telego.MemberStatusCreator {
		return true
	}

	admin, ok := member.(*telego.ChatMemberAdministrator)
	if !ok {
		sendAdminError(bot, ctx, update, chatID, isCallback,
			l("admin.checkers.not-admin"),
		)
		return false
	}

	if required == 0 {
		return true
	}

	var missing []string

	if required&PermDeleteMessages != 0 && !admin.CanDeleteMessages {
		missing = append(missing, "Delete Messages")
	}
	if required&PermRestrictMembers != 0 && !admin.CanRestrictMembers {
		missing = append(missing, "Restrict Members")
	}
	if required&PermPromoteMembers != 0 && !admin.CanPromoteMembers {
		missing = append(missing, "Promote Members")
	}
	if required&PermChangeInfo != 0 && !admin.CanChangeInfo {
		missing = append(missing, "Change Chat Info")
	}
	if required&PermInviteUsers != 0 && !admin.CanInviteUsers {
		missing = append(missing, "Invite Users")
	}
	if required&PermPinMessages != 0 && !admin.CanPinMessages {
		missing = append(missing, "Pin Messages")
	}

	if len(missing) > 0 {
		perms := ""
		text := l("admin.checkers.missing-permissions")
		if len(missing) == 1 {
			perms = missing[0]
		} else {
			for _, p := range missing {
				perms += p + ", "
			}
		}

		text = fmt.Sprintf(text, perms)
		sendAdminError(bot, ctx, update, chatID, isCallback, text)
		return false
	}

	return true
}
