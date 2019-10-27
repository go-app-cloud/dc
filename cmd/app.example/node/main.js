let express = require('express');
let app = express();
let jwt = require('jsonwebtoken');

let SecretKey = "uvTM$ei6pxxzCJC..!W]wK+Y*,^ZwGJgOyWDVk]%!eV$gQXToUpqIwTfZWoN30r(";
let url = "mongodb://localhost:27017/demo0";

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
}

let option = {
    uri: "ws://localhost:8000/source.cgi",
    appid: "87116b7089d543a68eb519ac231d0c2b",
    secretkey: SecretKey,
    onRead: function (data) {
        let info = JSON.parse(data);
        tokens = [];
        if (info.data != null && typeof info.data.length != "undefined") {
            for (let i = 0; i < info.data.length; i++) {
                tokens.push(info.data[i].secret_key)
            }
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
let MongoClient = require('mongodb').MongoClient;
let bodyParser = require('body-parser');
// parse application/x-www-form-urlencoded
app.use(bodyParser.urlencoded({ extended: false }))
// parse application/json
app.use(bodyParser.json());
/**
 * HTTP 业务逻辑
 */
app.use(function (req, res) {
    res.setHeader('Content-Type', 'text/plain')
    res.write('you posted:\n')
    res.end(JSON.stringify(req.body, null, 2))
})
app.post('/a', function (req, res) {
    console.log(req.body);
    // MongoClient.connect(url, {useNewUrlParser: true}, function (err, db) {
    //     if (err) throw err;
    //     db.collection('abc').insertOne({
    //         "name": name,
    //         "age": age,
    //         "score": {
    //             "shuxue": shuxuechengji,
    //             "yuwen": yuwenchengji
    //         }
    //     }, function (err, result) {
    //         if (err) {
    //             res.send('写入数据失败！');
    //         }
    //         db.close();
    //         res.send('数据写入成功！');
    //     });
    // });
});

app.listen(3000, function () {
    console.log('http://localhost:3000');
});