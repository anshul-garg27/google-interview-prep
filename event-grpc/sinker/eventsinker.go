package sinker

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/encoding/protojson"
	"init.bulbul.tv/bulbul-backend/event-grpc/clickhouse"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/model"
	"init.bulbul.tv/bulbul-backend/event-grpc/proto/go/bulbulgrpc"
	"init.bulbul.tv/bulbul-backend/event-grpc/transformer"
	"log"
	"reflect"
	"strings"
	"time"
)

func SinkEventToClickhouse(delivery amqp.Delivery) bool {

	grpcEvent := &bulbulgrpc.Events{}
	if err := protojson.Unmarshal(delivery.Body, grpcEvent); err == nil {
		// iterate through all events and sink each to event

		if grpcEvent.Events == nil || len(grpcEvent.Events) == 0 {
			log.Println("No event in events: ", grpcEvent)
			return false
		}

		log.Printf("Received in q: %v", grpcEvent)

		events := []model.Event{}
		for _, e := range grpcEvent.Events {
			event, err := TransformToEventModel(grpcEvent, e)
			if err != nil {
				return false
			}

			events = append(events, event)
		}

		if len(events) > 0 {
			if result := clickhouse.Clickhouse(config.New(), nil).Create(events); result != nil && result.Error == nil {
				log.Printf("Published successfully: %v", err)
				return true
			} else {
				log.Printf("Published failed: %v", result.Error)
			}
		}
	} else {
		log.Println("Some error processing event %v", err)
	}

	return true
}

func SinkErrorEventToClickhouse(delivery amqp.Delivery) bool {

	grpcEvent := &bulbulgrpc.Events{}
	if err := protojson.Unmarshal(delivery.Body, grpcEvent); err == nil {
		// iterate through all events and sink each to event

		if grpcEvent.Events == nil || len(grpcEvent.Events) == 0 {
			log.Println("No event in events: ", grpcEvent)
			return false
		}

		errorEvents := []model.ErrorEvent{}
		for _, e := range grpcEvent.Events {

			errorEvent, err := TransformToErrorEventModel(grpcEvent, e)
			if err != nil {
				return false
			}

			errorEvents = append(errorEvents, errorEvent)
		}

		if len(errorEvents) > 0 {
			if result := clickhouse.Clickhouse(config.New(), nil).Create(errorEvents); result != nil && result.Error == nil {
				log.Printf("Published successfully: %v", err)
				return true
			} else {
				log.Printf("Published failed: %v", result.Error)
			}
		}

	} else {
		log.Println("Some error processing error events for clickhouse %v", err)
	}

	return false
}
func TransformToErrorEventModel(grpcEvent *bulbulgrpc.Events, e *bulbulgrpc.Event) (model.ErrorEvent, error) {
	if grpcEvent.Header == nil || e.Timestamp == "" {
		log.Println("Incomplete errorEvent: ", grpcEvent)
		return model.ErrorEvent{}, errors.New("Incomplete errorEvent")
	}

	eventName := fmt.Sprintf("%v", reflect.TypeOf(e.GetEventOf()))
	eventName = strings.ReplaceAll(eventName, "*bulbulgrpc.Event_", "")

	errorEvent := model.ErrorEvent{
		EventId:              e.Id,
		EventName:            eventName,
		EventTimestamp:       time.Unix(transformer.InterfaceToInt64(e.Timestamp)/1000, 0),
		InsertTimestamp:      time.Now(),
		SessionId:            grpcEvent.Header.SessionId,
		BbDeviceId:           grpcEvent.Header.BbDeviceId,
		UserId:               grpcEvent.Header.UserId,
		DeviceId:             grpcEvent.Header.DeviceId,
		AndroidAdvertisingId: grpcEvent.Header.AndroidAdvertisingId,
		ClientId:             grpcEvent.Header.ClientId,
		Channel:              grpcEvent.Header.Channel,
		Os:                   grpcEvent.Header.Os,
		ClientType:           grpcEvent.Header.ClientType,
		AppLanguage:          grpcEvent.Header.AppLanguage,
		AppVersion:           grpcEvent.Header.AppVersion,
		AppVersionCode:       grpcEvent.Header.AppVersionCode,
		CurrentURL:           grpcEvent.Header.CurrentURL,
		MerchantId:           grpcEvent.Header.MerchantId,
		UtmReferrer:          grpcEvent.Header.UtmReferrer,
		UtmPlatform:          grpcEvent.Header.UtmPlatform,
		UtmSource:            grpcEvent.Header.UtmSource,
		UtmMedium:            grpcEvent.Header.UtmMedium,
		UtmCampaign:          grpcEvent.Header.UtmCampaign,
		PpId:                 grpcEvent.Header.PpId,
	}

	if e.CurrentURL != "" {
		errorEvent.CurrentURL = e.CurrentURL
	}

	data := map[string]interface{}{}
	if e.GetLaunchEvent() != nil {
		js, err := json.Marshal(e.GetLaunchEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetLaunchEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetPageOpenedEvent() != nil {
		js, err := json.Marshal(e.GetPageOpenedEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetPageOpenedEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetPageLoadedEvent() != nil {
		js, err := json.Marshal(e.GetPageLoadedEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetPageLoadedEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetPageLoadFailedEvent() != nil {
		js, err := json.Marshal(e.GetPageLoadFailedEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetPageLoadFailedEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetWidgetViewEvent() != nil {
		js, err := json.Marshal(e.GetWidgetViewEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetWidgetViewEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetWidgetCtaClickedEvent() != nil {
		js, err := json.Marshal(e.GetWidgetCtaClickedEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetWidgetCtaClickedEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetWidgetElementViewEvent() != nil {
		js, err := json.Marshal(e.GetWidgetElementViewEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetWidgetElementViewEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetWidgetElementClickedEvent() != nil {
		js, err := json.Marshal(e.GetWidgetElementClickedEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetWidgetElementClickedEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetEnterPWAProductPageEvent() != nil {
		js, err := json.Marshal(e.GetEnterPWAProductPageEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetEnterPWAProductPageEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetStreamEnterEvent() != nil {
		js, err := json.Marshal(e.GetStreamEnterEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetStreamEnterEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetChooseProductEvent() != nil {
		js, err := json.Marshal(e.GetChooseProductEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetChooseProductEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetPhoneVerificationInitiateEvent() != nil {
		js, err := json.Marshal(e.GetPhoneVerificationInitiateEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetPhoneVerificationInitiateEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetPhoneNumberEntered() != nil {
		js, err := json.Marshal(e.GetPhoneNumberEntered())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetPhoneNumberEntered().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetOtpVerifiedEvent() != nil {
		js, err := json.Marshal(e.GetOtpVerifiedEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetOtpVerifiedEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetAddToCartEvent() != nil {
		js, err := json.Marshal(e.GetAddToCartEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetAddToCartEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetGoToPaymentsEvent() != nil {
		js, err := json.Marshal(e.GetGoToPaymentsEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetGoToPaymentsEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetInitiatePurchaseEvent() != nil {
		js, err := json.Marshal(e.GetInitiatePurchaseEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetInitiatePurchaseEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetCompletePurchaseEvent() != nil {
		js, err := json.Marshal(e.GetCompletePurchaseEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetCompletePurchaseEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetSocialStreamEnterEvent() != nil {
		js, err := json.Marshal(e.GetSocialStreamEnterEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetSocialStreamEnterEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetSocialStreamActionEvent() != nil {
		js, err := json.Marshal(e.GetSocialStreamActionEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetSocialStreamActionEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetEnterProductEvent() != nil {
		js, err := json.Marshal(e.GetEnterProductEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetEnterProductEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetLaunchReferEvent() != nil {
		js, err := json.Marshal(e.GetLaunchReferEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetLaunchReferEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetPurchaseFailedEvent() != nil {
		js, err := json.Marshal(e.GetPurchaseFailedEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetPurchaseFailedEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetHighlightScreenEnterEvent() != nil {
		js, err := json.Marshal(e.GetHighlightScreenEnterEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetHighlightScreenEnterEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetHighlightScreenLoadedEvent() != nil {
		js, err := json.Marshal(e.GetHighlightScreenLoadedEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetHighlightScreenLoadedEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetHighlightScreenActionEvent() != nil {
		js, err := json.Marshal(e.GetHighlightScreenActionEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetHighlightScreenActionEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetHighlightScreenErrorEvent() != nil {
		js, err := json.Marshal(e.GetHighlightScreenErrorEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetHighlightScreenErrorEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetInitiateSearchEvent() != nil {
		js, err := json.Marshal(e.GetInitiateSearchEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetInitiateSearchEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetSearchRequestComplete() != nil {
		js, err := json.Marshal(e.GetSearchRequestComplete())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetSearchRequestComplete().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetCompleteSearchEvent() != nil {
		js, err := json.Marshal(e.GetCompleteSearchEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetCompleteSearchEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetInitiateReview() != nil {
		js, err := json.Marshal(e.GetInitiateReview())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetInitiateReview().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetEnterReviewScreen() != nil {
		js, err := json.Marshal(e.GetEnterReviewScreen())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetEnterReviewScreen().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetInitiateReviewMediaUpload() != nil {
		js, err := json.Marshal(e.GetInitiateReviewMediaUpload())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetInitiateReviewMediaUpload().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetSubmitReview() != nil {
		js, err := json.Marshal(e.GetSubmitReview())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetSubmitReview().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetEditReview() != nil {
		js, err := json.Marshal(e.GetEditReview())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetEditReview().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetViewAllReviews() != nil {
		js, err := json.Marshal(e.GetViewAllReviews())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetViewAllReviews().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetExpandReviewImage() != nil {
		js, err := json.Marshal(e.GetExpandReviewImage())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetExpandReviewImage().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetNudgeForReview() != nil {
		js, err := json.Marshal(e.GetNudgeForReview())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetNudgeForReview().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetSessionIdChangeEvent() != nil {
		js, err := json.Marshal(e.GetSessionIdChangeEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetSessionIdChangeEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetCheckDeliveryInfoEvent() != nil {
		js, err := json.Marshal(e.GetCheckDeliveryInfoEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetCheckDeliveryInfoEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetPdpTagClickEvent() != nil {
		js, err := json.Marshal(e.GetPdpTagClickEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetPdpTagClickEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetActionBarBagCLickEvent() != nil {
		js, err := json.Marshal(e.GetActionBarBagCLickEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetActionBarBagCLickEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetPdpActionEvent() != nil {
		js, err := json.Marshal(e.GetPdpActionEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetPdpActionEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetProfileActionEvent() != nil {
		js, err := json.Marshal(e.GetProfileActionEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetProfileActionEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetBadgeClickEvent() != nil {
		js, err := json.Marshal(e.GetBadgeClickEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetBadgeClickEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetHashTagClickEvent() != nil {
		js, err := json.Marshal(e.GetHashTagClickEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetHashTagClickEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetWhatsAppBuddyChatOpen() != nil {
		js, err := json.Marshal(e.GetWhatsAppBuddyChatOpen())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetWhatsAppBuddyChatOpen().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetHeaderSearchBarClick() != nil {
		js, err := json.Marshal(e.GetHeaderSearchBarClick())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetHeaderSearchBarClick().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetHeaderCartIconClick() != nil {
		js, err := json.Marshal(e.GetHeaderCartIconClick())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetHeaderCartIconClick().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetAppStartEvent() != nil {
		js, err := json.Marshal(e.GetAppStartEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetAppStartEvent().String()
			errorEvent.EventParams = data
		}
	}
	if e.GetHelpActionEvent() != nil {
		js, err := json.Marshal(e.GetHelpActionEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetHelpActionEvent().String()
			errorEvent.EventParams = data
		}
	}

	if e.GetReturnExchangeFlowActionEvent() != nil {
		js, err := json.Marshal(e.GetReturnExchangeFlowActionEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetReturnExchangeFlowActionEvent().String()
			errorEvent.EventParams = data
		}
	}

	if e.GetTestEvent() != nil {
		js, err := json.Marshal(e.GetTestEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetTestEvent().String()
			errorEvent.EventParams = data
		}
	}

	if e.GetNotificationActionEvent() != nil {
		js, err := json.Marshal(e.GetNotificationActionEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetNotificationActionEvent().String()
			errorEvent.EventParams = data
		}
	}

	if e.GetAffiliateActionEvent() != nil {
		js, err := json.Marshal(e.GetAffiliateActionEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetAffiliateActionEvent().String()
			errorEvent.EventParams = data
		}
	}

	if e.GetHostActionEvent() != nil {
		js, err := json.Marshal(e.GetHostActionEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				errorEvent.EventParams = data
			}
		}
		if err != nil {
			data["errorEventParams"] = e.GetHostActionEvent().String()
			errorEvent.EventParams = data
		}
	}

	return errorEvent, nil
}

func TransformToEventModel(grpcEvent *bulbulgrpc.Events, e *bulbulgrpc.Event) (model.Event, error) {
	if grpcEvent.Header == nil || e.Timestamp == "" || e.Id == "" || grpcEvent.Header.DeviceId == "" {
		log.Println("Incomplete event: ", grpcEvent)
		return model.Event{}, errors.New("Incomplete event")
	}

	eventName := fmt.Sprintf("%v", reflect.TypeOf(e.GetEventOf()))
	eventName = strings.ReplaceAll(eventName, "*bulbulgrpc.Event_", "")

	event := model.Event{
		EventId:              e.Id,
		EventName:            eventName,
		EventTimestamp:       time.Unix(transformer.InterfaceToInt64(e.Timestamp)/1000, 0),
		InsertTimestamp:      time.Now(),
		SessionId:            grpcEvent.Header.SessionId,
		BbDeviceId:           grpcEvent.Header.BbDeviceId,
		UserId:               grpcEvent.Header.UserId,
		DeviceId:             grpcEvent.Header.DeviceId,
		AndroidAdvertisingId: grpcEvent.Header.AndroidAdvertisingId,
		ClientId:             grpcEvent.Header.ClientId,
		Channel:              grpcEvent.Header.Channel,
		Os:                   grpcEvent.Header.Os,
		ClientType:           grpcEvent.Header.ClientType,
		AppLanguage:          grpcEvent.Header.AppLanguage,
		AppVersion:           grpcEvent.Header.AppVersion,
		AppVersionCode:       grpcEvent.Header.AppVersionCode,
		CurrentURL:           grpcEvent.Header.CurrentURL,
		MerchantId:           grpcEvent.Header.MerchantId,
		UtmReferrer:          grpcEvent.Header.UtmReferrer,
		UtmPlatform:          grpcEvent.Header.UtmPlatform,
		UtmSource:            grpcEvent.Header.UtmSource,
		UtmMedium:            grpcEvent.Header.UtmMedium,
		UtmCampaign:          grpcEvent.Header.UtmCampaign,
		PpId:                 grpcEvent.Header.PpId,
	}

	if e.CurrentURL != "" {
		event.CurrentURL = e.CurrentURL
	}

	data := map[string]interface{}{}
	if e.GetLaunchEvent() != nil {
		js, err := json.Marshal(e.GetLaunchEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}

	}
	if e.GetPageOpenedEvent() != nil {
		js, err := json.Marshal(e.GetPageOpenedEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetPageLoadedEvent() != nil {
		js, err := json.Marshal(e.GetPageLoadedEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetPageLoadFailedEvent() != nil {
		js, err := json.Marshal(e.GetPageLoadFailedEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetWidgetViewEvent() != nil {
		js, err := json.Marshal(e.GetWidgetViewEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetWidgetCtaClickedEvent() != nil {
		js, err := json.Marshal(e.GetWidgetCtaClickedEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetWidgetElementViewEvent() != nil {
		js, err := json.Marshal(e.GetWidgetElementViewEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetWidgetElementClickedEvent() != nil {
		js, err := json.Marshal(e.GetWidgetElementClickedEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetEnterPWAProductPageEvent() != nil {
		js, err := json.Marshal(e.GetEnterPWAProductPageEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetStreamEnterEvent() != nil {
		js, err := json.Marshal(e.GetStreamEnterEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetChooseProductEvent() != nil {
		js, err := json.Marshal(e.GetChooseProductEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetPhoneVerificationInitiateEvent() != nil {
		js, err := json.Marshal(e.GetPhoneVerificationInitiateEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetPhoneNumberEntered() != nil {
		js, err := json.Marshal(e.GetPhoneNumberEntered())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetOtpVerifiedEvent() != nil {
		js, err := json.Marshal(e.GetOtpVerifiedEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetAddToCartEvent() != nil {
		js, err := json.Marshal(e.GetAddToCartEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetGoToPaymentsEvent() != nil {
		js, err := json.Marshal(e.GetGoToPaymentsEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetInitiatePurchaseEvent() != nil {
		js, err := json.Marshal(e.GetInitiatePurchaseEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetCompletePurchaseEvent() != nil {
		js, err := json.Marshal(e.GetCompletePurchaseEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetSocialStreamEnterEvent() != nil {
		js, err := json.Marshal(e.GetSocialStreamEnterEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetSocialStreamActionEvent() != nil {
		js, err := json.Marshal(e.GetSocialStreamActionEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetEnterProductEvent() != nil {
		js, err := json.Marshal(e.GetEnterProductEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetLaunchReferEvent() != nil {
		js, err := json.Marshal(e.GetLaunchReferEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetPurchaseFailedEvent() != nil {
		js, err := json.Marshal(e.GetPurchaseFailedEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetHighlightScreenEnterEvent() != nil {
		js, err := json.Marshal(e.GetHighlightScreenEnterEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetHighlightScreenLoadedEvent() != nil {
		js, err := json.Marshal(e.GetHighlightScreenLoadedEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetHighlightScreenActionEvent() != nil {
		js, err := json.Marshal(e.GetHighlightScreenActionEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetHighlightScreenErrorEvent() != nil {
		js, err := json.Marshal(e.GetHighlightScreenErrorEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetInitiateSearchEvent() != nil {
		js, err := json.Marshal(e.GetInitiateSearchEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetSearchRequestComplete() != nil {
		js, err := json.Marshal(e.GetSearchRequestComplete())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetCompleteSearchEvent() != nil {
		js, err := json.Marshal(e.GetCompleteSearchEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetInitiateReview() != nil {
		js, err := json.Marshal(e.GetInitiateReview())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetEnterReviewScreen() != nil {
		js, err := json.Marshal(e.GetEnterReviewScreen())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetInitiateReviewMediaUpload() != nil {
		js, err := json.Marshal(e.GetInitiateReviewMediaUpload())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetSubmitReview() != nil {
		js, err := json.Marshal(e.GetSubmitReview())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetEditReview() != nil {
		js, err := json.Marshal(e.GetEditReview())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetViewAllReviews() != nil {
		js, err := json.Marshal(e.GetViewAllReviews())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetExpandReviewImage() != nil {
		js, err := json.Marshal(e.GetExpandReviewImage())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetNudgeForReview() != nil {
		js, err := json.Marshal(e.GetNudgeForReview())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetSessionIdChangeEvent() != nil {
		js, err := json.Marshal(e.GetSessionIdChangeEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetCheckDeliveryInfoEvent() != nil {
		js, err := json.Marshal(e.GetCheckDeliveryInfoEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetPdpTagClickEvent() != nil {
		js, err := json.Marshal(e.GetPdpTagClickEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetActionBarBagCLickEvent() != nil {
		js, err := json.Marshal(e.GetActionBarBagCLickEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetPdpActionEvent() != nil {
		js, err := json.Marshal(e.GetPdpActionEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetProfileActionEvent() != nil {
		js, err := json.Marshal(e.GetProfileActionEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetBadgeClickEvent() != nil {
		js, err := json.Marshal(e.GetBadgeClickEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetHashTagClickEvent() != nil {
		js, err := json.Marshal(e.GetHashTagClickEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetWhatsAppBuddyChatOpen() != nil {
		js, err := json.Marshal(e.GetWhatsAppBuddyChatOpen())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetHeaderSearchBarClick() != nil {
		js, err := json.Marshal(e.GetHeaderSearchBarClick())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetHeaderCartIconClick() != nil {
		js, err := json.Marshal(e.GetHeaderCartIconClick())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetAppStartEvent() != nil {
		js, err := json.Marshal(e.GetAppStartEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetHelpActionEvent() != nil {
		js, err := json.Marshal(e.GetHelpActionEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetReturnExchangeFlowActionEvent() != nil {
		js, err := json.Marshal(e.GetReturnExchangeFlowActionEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetTestEvent() != nil {
		js, err := json.Marshal(e.GetTestEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetNotificationActionEvent() != nil {
		js, err := json.Marshal(e.GetNotificationActionEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetAffiliateActionEvent() != nil {
		js, err := json.Marshal(e.GetAffiliateActionEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	if e.GetHostActionEvent() != nil {
		js, err := json.Marshal(e.GetHostActionEvent())
		if err == nil {
			if err = json.Unmarshal(js, &data); err == nil {
				event.EventParams = data
			}
		}
		if err != nil {
			return event, err
		}
	}
	return event, nil
}
