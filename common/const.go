/*
@Time : 2023/4/20 17:47
@Author : Hhx06
@File : const
@Description:
@Software: GoLand
*/

package common

/*
  服务器性能
 1、执行成功
 2、执行失败，未知失败
 3、 文件上传master失败,包含超时
 4、文件上传slave失败,包含超时
 5、master执行命令失败
 6、slave执行命令失败
 7、运行中
 8、 master连接失败
 9、 slave连接失败
 10、 获取报告失败，
 11、 手动中断，
 12、 清理老进程失败
 13、 超过失败阈值
 14、运行至不可执行时间段
 15、 启动中
*/
const (
	GameSuc                 = 1 + iota //  执行成功
	GameFail                           //  报告执行失败, 未知失败
	GameUploadMaster                   // 文件上传master失败,包含超时
	GameUploadSlave                    // 文件上传slave失败,包含超时
	GameRunMaster                      // master执行命令失败
	GameRunSlave                       // slave执行命令失败
	GameRunning                        //  运行中
	GameConnMaster                     // master连接失败
	GameConnSlave                      // slave连接失败
	GameGetReport                      // 获取报告失败
	GameInterrupt                      //  手动中断
	GameClearProcess                   // 清理老进程失败
	GameOverFailedThreshold            // 超过失败阈值
	GameRunInForbiddenTime             // 运行至不可执行时间段
	GameStaring                        // 启动中

	GameReport = 1 // 游戏压测报告
	PlatReport = 2 // 平台压测报告
)
