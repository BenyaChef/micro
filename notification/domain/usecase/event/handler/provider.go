package eventusecasehandler

import commonautheventmodel "github.com/BenyaChef/micro/common/auth/adapter/event/model"

type Provider struct {
	handlers map[commonautheventmodel.UserEventAction]Handler
}

func NewProvider(handlers map[commonautheventmodel.UserEventAction]Handler) (*Provider, error) {
	if len(handlers) == 0 {
		return nil, ErrHandlersAreRequired
	}

	return &Provider{handlers: handlers}, nil
}

func (p *Provider) ByAction(action commonautheventmodel.UserEventAction) (Handler, bool) {
	handler, found := p.handlers[action]

	return handler, found
}

func (p *Provider) SupportedActions() []commonautheventmodel.UserEventAction {
	actions := make([]commonautheventmodel.UserEventAction, 0, len(p.handlers))
	for action := range p.handlers {
		actions = append(actions, action)
	}

	return actions
}
