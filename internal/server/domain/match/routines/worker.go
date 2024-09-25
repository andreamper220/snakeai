package match_routines

import (
	"github.com/andreamper220/snakeai/internal/server/domain/match/data"
	"time"

	"github.com/andreamper220/snakeai/pkg/logger"
)

var PlayerJobsChannel = make(chan *data.Player, 100)

func MatchWorker() {
	// TODO observer pattern ?

	for p := range PlayerJobsChannel {
		if p.InProcess {
			continue
		}
		if p.InParty {
			if len(p.Party.Players) == p.Party.Size {
				PartiesChannel <- p.Party
				data.CurrentParties.RemoveParty(p.Party)
				logger.Log.Infof("formed party: %v", p.Party)
			} else {
				isPartyExisted := false
				parties := data.CurrentParties.GetParties()
				for _, pa := range parties {
					if p.Party == pa {
						isPartyExisted = true
						break
					}
				}

				if !isPartyExisted {
					data.CurrentParties.AddParty(p.Party)
				}
			}
			continue
		}

		pa, err := p.FindParty()
		if err != nil {
			logger.Log.Infof("player with ID %s already searching for party", p.Id)
			continue
		}
		if pa != nil && p.Party != pa {
			pa.AddPlayer(p)
		}

		if pa != nil && len(pa.Players) == pa.Size {
			PartiesChannel <- p.Party
			data.CurrentParties.RemoveParty(p.Party)
			logger.Log.Infof("formed party: %v", pa)
		} else {
			addDelayToReenqueuePlayer(p, PlayerJobsChannel)
		}
	}
}
func addDelayToReenqueuePlayer(p *data.Player, playerJobsChannel chan *data.Player) {
	timer := time.NewTimer(3 * time.Second)
	go func() {
		<-timer.C
		p.Lock()
		if !p.InParty {
			if p.Party != nil {
				p.Party.RemovePlayer(p)
				p.Party = nil
			}
			p.SkillDelta = p.SkillDelta * 2
			p.InParty = false
			p.InProcess = false
			p.Unlock()
			playerJobsChannel <- p
		} else {
			p.Unlock()
		}
	}()
}
