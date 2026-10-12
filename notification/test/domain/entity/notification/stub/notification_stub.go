package notificationentitystub

import (
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	idprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/id"
	notificationentity "github.com/BenyaChef/micro/notification/domain/entity/notification"
)

const RawPayload = `{"event_id":"x","action":"REGISTERED","data":{"user_id":"y","email":"a@b.c"}}`

func GetEventID() notificationentity.EventID {
	eventID, err := notificationentity.EventIDFrom(idprimitive.NewEntityID().String())
	if err != nil {
		panic(err)
	}

	return eventID
}

func GetPayload() notificationentity.Payload {
	payload, err := notificationentity.PayloadFrom([]byte(RawPayload))
	if err != nil {
		panic(err)
	}

	return payload
}

func GetNotificationBuilder() *notificationentity.Builder {
	return notificationentity.NewBuilder().
		ID(notificationentity.NewNotificationID()).
		EventID(GetEventID()).
		UserID(commonuserentity.NewUserID()).
		Kind(notificationentity.KindWelcome).
		Payload(GetPayload())
}

func GetValidNotification() *notificationentity.Notification {
	notification, err := GetNotificationBuilder().Build()
	if err != nil {
		panic(err)
	}

	return notification
}
