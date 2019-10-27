let express = require('express');
let app = express();
let jwt = require('jsonwebtoken');

let SecretKey = "OPD47*bJgaf0V^evBj1ym[l[F[X*vzcYip-ilF=CbugyhW36iw^=@8Cev!yJ,9Kx";
let url = "mongodb://localhost:27017/demo0";

process.on('uncaughtException', function (err) {
    console.log(err.stack);
    console.log('NOT exit...');
});

/**
 * 认证过滤器
 */
let tokens = [];
app.use(function (req, res, next) {
    let isAuth = false;
    let authorization = req.headers["authorization"];
    try {
        let decoded = jwt.verify(authorization, SecretKey);
        isAuth = true
    } catch (err) {
        isAuth = false;
    }
    if (isAuth) {
        next();
        return
    }
    for (let i = 0; i < tokens.length; i++) {
        try {
            let decoded = jwt.verify(authorization, tokens[i]);
            isAuth = true;
            break;
        } catch (err) {
            isAuth = false;
        }
    }
    if (!isAuth) {
        res.send({code: -1, msg: "无权访问"});
        return
    }
    next();
});

/**
 * 客户端TCP连接集群服务
 * @type {WebSocket|WebSocket}
 */
const WebSocket = require('ws');

function BuildClient(option) {
    const ws = new WebSocket(option.uri);
    ws.on('open', function open() {
        ws.send(JSON.stringify({
            appid: option.appid,
            secretkey: option.secretkey
        }));
    });
    ws.on("close", function close() {
        option.onClose();
    });
    ws.on('message', function incoming(data) {
        option.onRead(data);
    });
    ws.on('error', function (err) {

    });
}

let option = {
    uri: "ws://localhost:8000/source.cgi",
    appid: "a879b30784e940678a18e15e0bc4fe1b",
    secretkey: SecretKey,
    onRead: function (data) {
        let info = JSON.parse(data);
        switch (info.code) {
            case 100:
                console.log("=======注册服务成功=======");
                break;
            case 1001:
                tokens = [];
                if (info.data != null && typeof info.data.length != "undefined") {
                    for (let i = 0; i < info.data.length; i++) {
                        tokens.push(info.data[i].secret_key)
                    }
                }
                break;
        }
    },
    onClose: function () {
        console.log("device connection is close.");
        setTimeout(function () {
            BuildClient(option);
        }, 5000)
    }
};

BuildClient(option);


// 数据库管理
let mongoose = require('mongoose');
mongoose.connect('mongodb://localhost/test', {});
mongoose.Promise = global.Promise;

let First = mongoose.model('first', {id: String, name: String});

app.use(express.json());
/**
 * HTTP 业务逻辑
 */
app.post('/receive.cgi', function (req, res) {
    console.log(req.body);
    var kitty = new First(req.body);
    kitty.save(function (err) {
        if (err) {
            console.log(err);
        } else {
            console.log('meow');
        }
    });
    res.end();
});

app.listen(9066, function () {
    console.log('http://localhost:9066');
});