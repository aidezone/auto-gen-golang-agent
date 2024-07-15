package playbook

import (
	"fmt"
    "io/ioutil"
    "aidezone/auto-gen-golang-agent/defines"
    "aidezone/auto-gen-golang-agent/logger"
    "aidezone/auto-gen-golang-agent/generator"

    "github.com/google/uuid"
)

var prompt2 = `
你的角色： %v
请帮我完成工作： %v
具体细节如下：
%v
请用json格式回答：[{"feature": "xxxx","description": "xxxx"}]
`

var KTVDeepseek = Playbook{
	Robots: []struct{
		Actor string
		Number int
	}{
		struct{
			Actor string
			Number int
		}{
			Actor: "产品经理",
			Number: 1,
		},
		struct{
			Actor string
			Number int
		}{
			Actor: "Golang资深工程师",
			Number: 1,
		},
	},
}

func KTVDeepseekGenerate(){
	// 创建机器人
    robots := make(map[string]*generator.Robot)
    for _, robot := range KTVDeepseek.Robots {
        for i := 0; i < robot.Number; i++ {
            name := fmt.Sprintf("%v_%v", robot.Actor, i)
            robots[name] = generator.NewRobot(name, robot.Actor, defines.Deepseek)
            go robots[name].Run()

            // 与每个机器人说一句话，先定义好其角色
            firstSay := fmt.Sprintf("你好，我们一起来完成一个软件开发工作，你来扮演“%v”", robot.Actor)
            sayToRobotDeepseek(robots[name], firstSay)
        }
    }

    // 读取文件内容
    sayToRobotDeepseek(robots["产品经理_0"], `我们来做第一件事情
        稍后我会给出具体的需求文档，请你先根据我的需求进行功能细化。
        注意：所有回答的内容不要做省略。`)
    data, err := ioutil.ReadFile("PRD.txt")
    if err != nil {
        logger.Errorf("read prd file err: %v", err)
    }
    prd := sayToRobotDeepseek(robots["产品经理_0"], fmt.Sprintf("以下是我想到的所有需求,输出json格式每个模块名称作为key,模块中包含多个功能，每个功能下面写清楚名称、详细介绍、逻辑关系：\n\n %v", string(data)))

    sayToRobotDeepseek(robots["Golang资深工程师_0"], *prd.Say)
    sayToRobotDeepseek(robots["Golang资深工程师_0"], "以上是产品经理提供的需求，请你设计出当前项目的目录结构以及所有的类定义，框架使用gin、gorm，数据库使用mysql,返回格式使用json，输出filepath 和 简略的文件内容")

    //TODO 换用一个更强大的大语言模型，获取其json返回结果，并根据返回内容生成文件。
}

func sayToRobotDeepseek(robot *generator.Robot, say string) *generator.RobotTalk {
    firstRobot := generator.NewRobot("人类", "人类", defines.Deepseek)
    firstTalk := &generator.RobotTalk{
        Robot: firstRobot,
        Say: &say,
        TransId: uuid.New(),
    }
    return robot.Do(firstTalk)
}