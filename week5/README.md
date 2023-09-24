作业：测试前端接口 LoginJWT 或者 LoginSMS
注意：

你要测试的是 LoginJWT 或者 LoginSMS 本身，所以你要摆脱对 Service 的依赖，使用 mock 生成的代码。
你要覆盖所有的分支。
因为在 LoginSMS 里面，返回的响应是 Result 序列化后的 JSON 串，所以你要考虑怎么校验响应数据。
LoginJWT 和 LoginSMS 任意挑一个。