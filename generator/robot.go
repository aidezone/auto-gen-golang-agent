package generator

import (
	"context"
	"time"
	
	"aidezone/auto-gen-golang-agent/defines"
	"aidezone/auto-gen-golang-agent/logger"
	"aidezone/auto-gen-golang-agent/ai-platforms"

	"github.com/google/uuid"
)

type RobotTalk struct {
	Robot *Robot
	Say *string
	TransId uuid.UUID
}

type Robot struct {
	aiPlatform platforms.AIPlatform
	chat *defines.Chat
	name string
	actor string
	chatMessage chan *RobotTalk
	answerMap map[uuid.UUID]chan *RobotTalk
}

func NewRobot(name, actor string, aiPlatform defines.PlatformName) *Robot {
	return &Robot{
		chat: defines.NewChat(),
		aiPlatform: platforms.NewAIPlatform(aiPlatform),
		name: name,
		actor: actor,
		chatMessage: make(chan *RobotTalk, 10),
		answerMap: make(map[uuid.UUID]chan *RobotTalk),
	}
}

func (s *Robot) Do(talk *RobotTalk) *RobotTalk {
	s.answerMap[talk.TransId] = make(chan *RobotTalk, 1)
	s.chatMessage <- talk

	for {
		resp := func () *RobotTalk {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			select {
			case respTalk := <- s.answerMap[talk.TransId]:
				delete(s.answerMap, talk.TransId)
				s.chat.AppendAnswer(respTalk.Say)
				return respTalk
			case <-ctx.Done():
				return nil
			}
		}()
		if resp != nil {
			return resp
		}
	}
}


func (s *Robot) Run() {
	for {
		talk := <- s.chatMessage

		logger.Infof("Robot[%v][%v]: [%v][%v] said to me '%v'", s.actor, s.name, talk.Robot.actor, talk.Robot.name, *talk.Say)

		chatRequest := s.chat.AppendAsk(talk.Say)
		answer, err := s.aiPlatform.Call(chatRequest)
		if err != nil {
			logger.Errorf("Robot[%v][%v]: call platform error [%v]", s.actor, s.name, err)
		}
		logger.Infof("Robot[%v][%v]: I answer to [%v][%v] '%v'", s.actor, s.name, talk.Robot.actor, talk.Robot.name, *answer)
		s.answerMap[talk.TransId] <- &RobotTalk{
			Robot: s,
			Say: answer,
			TransId: uuid.New(),
		}
	}
}

