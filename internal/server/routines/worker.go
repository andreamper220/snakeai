package routines

import (
	"snake_ai/internal/shared/match/data"
	"time"

	"snake_ai/internal/logger"
)

var PlayerJobsChannel = make(chan *data.Player, 100)

func MatchWorker(
	parties *[]*data.Party,
) {
	for p := range PlayerJobsChannel {
		if p.InProcess {
			continue
		}
		if p.InParty {
			if len(p.Party.Players) == p.Party.Size {
				PartiesChannel <- p.Party
				removeParty(parties, p.Party)
				logger.Log.Infof("formed party: %v", p.Party)
			} else {
				isPartyExisted := false
				for _, pa := range *parties {
					if p.Party == pa {
						isPartyExisted = true
						break
					}
				}

				if !isPartyExisted {
					addParty(parties, p.Party)
				}
			}
			continue
		}

		pa, err := p.FindParty(*parties)
		if err != nil {
			logger.Log.Infof("player with ID %s already searching for party", p.Id)
			continue
		}
		if pa != nil && p.Party != pa {
			pa.AddPlayer(p)
		}

		if pa != nil && len(pa.Players) == pa.Size {
			PartiesChannel <- p.Party
			removeParty(parties, p.Party)
			logger.Log.Infof("formed party: %v", pa)
		} else {
			addDelayToReenqueuePlayer(p, PlayerJobsChannel)
		}
	}
}
func addParty(parties *[]*data.Party, pa *data.Party) {
	*parties = append(*parties, pa)
}
func removeParty(parties *[]*data.Party, pa *data.Party) {
	result := make([]*data.Party, 0)
	for _, par := range *parties {
		if par != pa {
			result = append(result, par)
		}
	}
	*parties = result
}
func addDelayToReenqueuePlayer(p *data.Player, playerJobsChannel chan *data.Player) {
	timer := time.NewTimer(3 * time.Second)
	go func() {
		<-timer.C
		p.Lock()
		if !p.InParty {
			p.Party.RemovePlayer(p)
			p.SkillDelta = p.SkillDelta * 2
			p.Party = nil
			p.InParty = false
			p.InProcess = false
			p.Unlock()
			playerJobsChannel <- p
		} else {
			p.Unlock()
		}
	}()
}
