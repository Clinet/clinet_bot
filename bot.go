package bot

import (
	"fmt"

	"github.com/Clinet/clinet_cmds"
	"github.com/Clinet/clinet_config"
	"github.com/Clinet/clinet_convos"
	"github.com/Clinet/clinet_features"
	"github.com/Clinet/clinet_services"
)

//Bot holds an instance of the bot framework
//TODO: Add a logger interface
type Bot struct {
	cfg config.Config
	registeredFeatures []string
	registeredServices []services.Service
}

//NewBot returns a new bot instance (limit one per process)
//TODO: Add argument for a logger interface
func NewBot(cfg *config.Config) *Bot {
	features.SetFeatures(cfg.Features)

	bot := &Bot{
		cfg: *cfg,
		registeredFeatures: make([]string, 0),
		registeredServices: make([]services.Service, 0),
	}
	return bot
}

func (b *Bot) Shutdown() {
	for i := 0; i < len(b.registeredServices); i++ {
		b.registeredServices[i].Shutdown()
	}
}

//RegisterFeature adds an enabled feature's commands, runs its init, and logs into its chat and convo services, as available and halting on error
func (b *Bot) RegisterFeature(feature features.Feature) error {
	if features.IsEnabled(feature.Name) {
		if b.isRegisteredFeature(feature.Name) {
			return fmt.Errorf("bot: cannot register feature %s twice", feature.Name)
		}
		if feature.Cmds != nil && len(feature.Cmds) > 0 {
			cmds.Commands = append(cmds.Commands, feature.Cmds...)
		}
		if feature.Init != nil {
			if err := feature.Init(); err != nil {
				return fmt.Errorf("bot: init call for feature %s failed: %v", feature.Name, err)
			}
		}
		if feature.ServiceChat != nil {
			if err := feature.ServiceChat.Login(); err != nil {
				return fmt.Errorf("bot: registering chat service for feature %s failed: %v", feature.Name, err)
			}
			b.registeredServices = append(b.registeredServices, feature.ServiceChat)
		}
		if feature.ServiceConvo != nil {
			if err := feature.ServiceConvo.Login(); err != nil {
				return fmt.Errorf("bot: registering conversation service for feature %s failed: %v", feature.Name, err)
			}
			convos.ConvoServices = append(convos.ConvoServices, feature.ServiceConvo)
		}
		b.registerFeature(feature.Name)
		features.UpdateFeature(feature)
	}
	return nil
}

func (b *Bot) registerFeature(feature string) {
	b.registeredFeatures = append(b.registeredFeatures, feature)
}
func (b *Bot) isRegisteredFeature(feature string) bool {
	for _, registeredFeature := range b.registeredFeatures {
		if registeredFeature == feature {
			return true
		}
	}
	return false
}