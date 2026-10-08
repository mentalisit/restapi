package restapi

import (
	"fmt"
	"time"

	"github.com/mentalisit/conf/logger"
	"github.com/mentalisit/restapi/bridge2"
	"github.com/mentalisit/restapi/compendium"
	"github.com/mentalisit/restapi/models"
	"github.com/mentalisit/restapi/rs_bot2"
)

type Recover struct {
	log               *logger.Logger
	bridge2Message    []models.ToBridgeMessage
	compendiumMessage []models.IncomingMessage
	rsBotV2Message    []models.InMessageV2
	bridge2           *bridge2.Client
	rs                *rs_bot2.Client
	compendiumNew     *compendium.Client
}

func NewRecover(log *logger.Logger) *Recover {
	r := &Recover{
		log:           log,
		bridge2:       bridge2.NewClient(log),
		rs:            rs_bot2.NewClient(log),
		compendiumNew: compendium.NewClient(log),
	}
	go r.trySend()
	return r
}

func (r *Recover) SendBridge2AppRecover(m models.ToBridgeMessage) {
	fmt.Printf("%s SendBridge2App Text:%s Sender:%s Tip:%s ChatId:%s\n",
		time.Now().Format(time.DateTime), m.Text, m.Sender, m.Tip, m.ChatId)

	err := r.bridge2.SendToBridge(m)
	if err != nil {
		r.log.InfoStruct("SendBridge2App err "+err.Error(), m)
		r.bridge2Message = append(r.bridge2Message, m)
	}
}

func (r *Recover) SendCompendiumAppRecover(m models.IncomingMessage) {
	fmt.Printf("%s SendCompendiumApp :%+v\n", time.Now().Format(time.DateTime), m)
	err := r.compendiumNew.SendToCompendium(m)
	if err != nil {
		r.log.InfoStruct("SendCompendiumApp err "+err.Error(), m)
		r.compendiumMessage = append(r.compendiumMessage, m)
	}
}

func (r *Recover) SendRsBotV2AppRecover(m models.InMessageV2) {
	err := r.rs.SendToRs2(m)
	if err != nil {
		r.log.InfoStruct("SendRsBotV2App err "+err.Error(), m)
		r.rsBotV2Message = append(r.rsBotV2Message, m)
	}
}

func (r *Recover) trySend() {
	for {
		// Проверка и отправка сообщений в rsBot2 (единственный клиент)
		if len(r.rsBotV2Message) > 0 {
			for i := 0; i < len(r.rsBotV2Message); i++ {
				message := r.rsBotV2Message[i]
				err := r.rs.SendToRs2(message)
				if err == nil {
					r.rsBotV2Message = append(r.rsBotV2Message[:i], r.rsBotV2Message[i+1:]...)
					i--
				}
				time.Sleep(1 * time.Second)
			}
		}

		if len(r.compendiumMessage) > 0 {
			for i := 0; i < len(r.compendiumMessage); i++ {
				message := r.compendiumMessage[i]
				err := r.compendiumNew.SendToCompendium(message)
				if err == nil {
					r.compendiumMessage = append(r.compendiumMessage[:i], r.compendiumMessage[i+1:]...)
					i--
				}
				time.Sleep(1 * time.Second)
			}
		}

		if len(r.bridge2Message) > 0 {
			for i := 0; i < len(r.bridge2Message); i++ {
				message := r.bridge2Message[i]
				err := r.bridge2.SendToBridge(message)
				if err == nil {
					r.bridge2Message = append(r.bridge2Message[:i], r.bridge2Message[i+1:]...)
					i--
				}
				time.Sleep(1 * time.Second)
			}
		}

		time.Sleep(10 * time.Second)
	}
}
func (r *Recover) Close() {
	_ = r.bridge2.Close()
	_ = r.rs.Close()
	_ = r.compendiumNew.Close()
}
