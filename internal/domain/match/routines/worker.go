package match_routines

import (
	"time"

	matchdata "github.com/andreamper220/snakeai/internal/domain/match/data"
	"github.com/andreamper220/snakeai/pkg/logger"
)

var PlayerJobsChannel = make(chan *matchdata.Player, 100)

func MatchWorker() {
	for p := range PlayerJobsChannel {
		if p.InProcess {
			continue
		}
		if p.InParty {
			if len(p.Party.Players) == p.Party.Size {
				PartiesChannel <- p.Party
				matchdata.CurrentParties.RemoveParty(p.Party)
				logger.Log.Infof("formed party: %v", p.Party)
			} else {
				isPartyExisted := false
				parties := matchdata.CurrentParties.GetParties()
				for _, pa := range parties {
					if p.Party == pa {
						isPartyExisted = true
						break
					}
				}

				if !isPartyExisted {
					matchdata.CurrentParties.AddParty(p.Party)
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
			matchdata.CurrentParties.RemoveParty(p.Party)
			logger.Log.Infof("formed party: %v", pa)
		} else {
			addDelayToReenqueuePlayer(p, PlayerJobsChannel)
		}
	}
}
func addDelayToReenqueuePlayer(p *matchdata.Player, playerJobsChannel chan *matchdata.Player) {
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
