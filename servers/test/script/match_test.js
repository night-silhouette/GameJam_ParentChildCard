import http from 'k6/http';
import { check, sleep } from 'k6';
import ws from 'k6/ws';

export let options = {
    vus: 9,          // Virtual Users
    duration: '500s',  // 持续压测时间
};




// 2. 核心逻辑：每个“虚拟人”都会重复执行这个函数
export default function () {
    const BASE_URL = 'http://127.0.0.1:10086'; // 你的 Go 后端地址

    // 生成随机账号，防止 psql 因为唯一约束报错
    const uniqueId = `${__VU}-${__ITER}`;
    const payload = JSON.stringify({
        name: `k6_user_${uniqueId}`,
        password: '123456',
    });

    const params = { headers: { 'Content-Type': 'application/json' } };

    let loginRes = http.post(`${BASE_URL}/v1/token/`, payload, params);

    // 检查是否登录成功并拿到了 Token
    const isOk = check(loginRes, {
        '登陆 code is 0': (r) => r.json().code == 0,
        'has token': (r) => r.json().data.token !== undefined,
    });

    const wsParams = {
    headers: {
        Origin: 'http://127.0.0.1', // 试试加 origin
    },
    };

    const token = loginRes.json().data.token;
    const WS_URL = `ws://127.0.0.1:10086/v1/ws/?token=${token}`;
 

    const res = ws.connect(WS_URL, wsParams, function (socket) {


        socket.on('message', function (message) {
            console.log(`Received message: ${message}`);
        });

        if (!message.battle_id) {
            handleBattleAction(socket);
        }
    });



    check(res, {
        'WebSocket 协议升级成功': (r) => r && r.status === 101,
    });


    sleep(500);
}

function handleBattleAction(socket) {
        socket.send(JSON.stringify({
            action_code: 1,
            action_name:"",
            action_data:{},
        }));
    }