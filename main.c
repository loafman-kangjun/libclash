#include "libclash.h"   /* 由 go build 生成的头文件 */
#include <stdio.h>

static void my_log_cb(int lvl, const char* msg) {
	printf("[C LOG][%d] %s\n", lvl, msg);
}

int main(void)
{
	Mihomo_SetLogCallback((void*)my_log_cb);

	char* err = Mihomo_Run("", "", false);   /* ① 调用 */
	if (err != NULL) {                       /* ② 判断 */
		printf("Run error: %s\n", err);
		free_go_str(err);                    /* ③ 别忘了释放 */
		return 1;
	}

	/* ……业务代码…… */

	Mihomo_Stop();
	return 0;
}