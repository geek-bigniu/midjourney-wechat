package replay

import "strings"

const (
	TaskMainFinishMsg     = 0
	TaskSubVFinishMsg     = 1
	TaskSubUFinishMsg     = 2
	TaskTransImgErrMsg    = 3
	TaskParamsErrMsg      = 4
	TaskBannedErrMsg      = 5
	TaskLinkErrMsg        = 6
	TaskErrMsg            = 7
	TaskErrMsg1           = 8
	TaskMainCreateMsg     = 9
	TaskSubCreateMsg      = 10
	TaskSendErrMsg        = 11
	TaskNewUserErrMsg     = 12
	TaskSubParamsErrMsg   = 13
	TaskMainCommandErrMsg = 14
	TaskSubCommandErrMsg  = 15
)

type Info struct {
	TaskId    string
	NewTaskId string
	Prompt    string
	Action    string
	NickName  string
	Url       string
	Msg       string
}

func (info *Info) GenrateMessage(typeName int) string {
	switch typeName {
	case TaskMainFinishMsg:
		info.Msg = "@" + info.NickName + "\n" +
			"🎨 绘画成功!\n" +
			"📨 消息ID：\n" +
			info.TaskId + "\n" +
			"🪄 变换：\n" +
			"[ U1 ] [ U2 ] [ U3 ] [ U4 ] \n" +
			"[ V1 ] [ V2 ] [ V3 ] [ V4 ] \n" +
			"✏️ 可使用 [/up-任务ID-操作] 进行变换\n" +
			"/up " + info.TaskId + " U1"
		break
	case TaskSubVFinishMsg:
		info.Msg = "@" + info.NickName + "\n" +
			"🎨 绘画成功!\n" +
			"📨 消息ID：\n" +
			info.TaskId + "\n" +
			"🪄 变换：\n" +
			"[ U1 ] [ U2 ] [ U3 ] [ U4 ] \n" +
			"[ V1 ] [ V2 ] [ V3 ] [ V4 ] \n" +
			"✏️ 可使用 [/up-任务ID-操作] 进行变换\n" +
			"/up " + info.TaskId + " U1"
		break
	case TaskSubUFinishMsg:
		info.Msg = "@" + info.NickName + "\n" +
			"🎨 绘画成功!\n" +
			"📨 消息ID：\n" + info.TaskId
		break
	case TaskTransImgErrMsg:
		info.Msg = "✅任务已完成\n" +
			"ℹ️图片转码失败\n" +
			"🌟任务ID:\n" +
			info.TaskId + "\n" +
			"🧷任务返回图片地址:\n" +
			info.Url
		break
	case TaskParamsErrMsg:
		info.Msg = "@" + info.NickName + "\n" +
			"❌任务被拒绝\n" +
			"⭕️参数错误，请检查\n" +
			"⚠️删除任务:\n" + info.TaskId
		break
	case TaskBannedErrMsg:
		info.Msg = "@" + info.NickName + "\n" +
			"❌任务被拒绝\n" +
			"⭕️可能包含违禁词，请检查\n" +
			"⚠️删除任务:\n" + info.TaskId
		break
	case TaskLinkErrMsg:
		info.Msg = "@" + info.NickName + "\n" +
			"❌任务被拒绝\n" +
			"⭕️图片链接地址错误\n" +
			"请提供能直接访问的图片链接地址\n" +
			"⚠️删除任务:\n" + info.TaskId
		break
	case TaskErrMsg:
		info.Msg = "@" + info.NickName + "\n" +
			"❌任务处理失败\n" +
			"⭕️任务被拒绝或处理超时\n" +
			"请尝试重新发送指令进行生成\n" +
			"⚠️删除任务:\n" + info.TaskId
		break
	case TaskErrMsg1:
		info.Msg = "@" + info.NickName + "\n" +
			"❌任务处理失败\n" +
			"⭕️队列人数过多,请稍后再试\n" +
			"⚠️删除任务:\n" + info.TaskId
		break
	case TaskMainCreateMsg:
		info.Msg = "@" + info.NickName + "\n" +
			"✅你发送的任务已提交\n" +
			"✨Prompt: " + info.Prompt + "\n" +
			"🌟任务ID:\n" +
			info.TaskId + "\n" +
			"🚀正在快速处理中,请稍后!"
		break
	case TaskSubCreateMsg:
		info.Msg = "@" + info.NickName + "\n" +
			"✅你发送的任务已提交\n" +
			"✨变换ID:\n" +
			info.TaskId + "\n" +
			"🌟任务ID:\n" +
			info.NewTaskId + "\n" +
			"💫变换类型: " + strings.ToUpper(info.Action) + "\n" +
			"🚀正在快速处理中,请稍后!"
		break
	case TaskSendErrMsg:
		info.Msg = "@" + info.NickName + "\n" +
			"❌任务创建失败，请联系管理员或稍后再试"
		break
	case TaskNewUserErrMsg:
		info.Msg = "@" + info.NickName + "\n" +
			"❌这位新朋友，请先冒泡后再发送指令哦"
		break
	case TaskMainCommandErrMsg:
		info.Msg = "@" + info.NickName + "\n" +
			"❌指令错误，请输入/imagine+空格+内容"
		break
	case TaskSubCommandErrMsg:
		info.Msg = "@" + info.NickName + "\n" +
			"❌命令格式错误，示例:/up 任务id u1"
		break
	case TaskSubParamsErrMsg:
		info.Msg = "@" + info.NickName + "\n" +
			"❌参数错误\n" +
			"✨可选参数:\n" +
			"[ U1 ] [ U2 ] [ U3 ] [ U4 ] \n" +
			"[ V1 ] [ V2 ] [ V3 ] [ V4 ] \n" +
			"✏️ 可使用 [/up-任务ID-操作] 进行变换\n" +
			"/up [任务id] U1"
		break
	}
	return info.Msg
}
