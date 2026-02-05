## Milestone
### Day 1-2: 
ZeroTier的接口验证和Web框架初步搭建，在本地配置`zerotier-one`，验证Go调用`zerotier-cli`或者其API验证网络状态、加入或离开网络的可行性。确定Web技术栈，初始化项目骨架，并且验证ZeroTier的Go调用逻辑。
### Day 3-8: 
分为后端逻辑开发和CLI/GUI的搭建。
#### Day 3-5:
在`pkg/`下完成所有业务代码，确保可测试性。包括：ZeroTier核心功能，Worker Pool的子网扫描，AES加密以及Samba挂载等。
#### Day 6:
封装Cobra命令并且调用`pkg/`内的函数。实现Survey的交互式引导与JSON输出支持，完成第一个可用的二进制文件`rnctl-cli`。
#### Day 7-9：
Wails调试。初始化Wails项目，编写App绑定逻辑，将`pkg/`功能暴露给前端测试。在Wails给定的模板基础上设计UI，联调GUI事件流给足用户界面反馈。
### Day 9-10:
压力测试和兼容性测试，在不同网络环境（Docker内部署）进行测试。完善README和技术文档，录制Demo视频。
