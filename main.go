package main

import (
    "aidezone/auto-gen-golang-agent/playbook"
    "aidezone/auto-gen-golang-agent/logger"
)

func main() {

    // 初始化日志
    logger.InitLogger("", true, false)

    playbook.KTVGenerate()

}

