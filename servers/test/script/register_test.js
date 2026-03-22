import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    vus: 10,          // Virtual Users
    duration: '10s',  // 持续压测时间
};




// 2. 核心逻辑：每个“虚拟人”都会重复执行这个函数
export default function () {
    const BASE_URL = 'http://120.26.145.68:10086'; // 你的 Go 后端地址

    // 生成随机账号，防止 psql 因为唯一约束报错
    const uniqueId = `${__VU}-${__ITER}`;
    const payload = JSON.stringify({
        name: `k6_user_${uniqueId}`,
        password: '123456',
    });

    const params = { headers: { 'Content-Type': 'application/json' } };

    // --- 步骤 A: 模拟注册 ---
    let regRes = http.post(`${BASE_URL}/v1/user/`, payload, params);
    check(regRes, {
        '注册code is 0': (r) => r.json().code == 0,

    },);

    // --- 步骤 B: 模拟登录并获取 Token ---
    let loginRes = http.post(`${BASE_URL}/v1/token/`, payload, params);

    // 检查是否登录成功并拿到了 Token
    const isOk = check(loginRes, {
        '登陆 code is 0': (r) => r.json().code == 0,
        'has token': (r) => r.json().data.token !== undefined,
    });

    if (regRes.json().code!=0){
        console.log("注册失败");
        console.log(reges.json().code);
    }
    if (loginRes.json().code!=0){
        console.log("登陆失败")
        console.log(loginRes.json().code);
    }

    // --- 步骤 C: 模拟真实用户停顿 ---
    sleep(1);
}