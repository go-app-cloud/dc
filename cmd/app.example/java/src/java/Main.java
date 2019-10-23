import com.alibaba.fastjson.JSON;
import okhttp3.*;

import java.io.IOException;
import java.util.concurrent.TimeUnit;


class StanderRequest {
    private int Code;
    private Object Data;

    public StanderRequest(int code, Object data) {
        Code = code;
        Data = data;
    }

    public int getCode() {
        return Code;
    }

    public void setCode(int code) {
        Code = code;
    }

    public Object getData() {
        return Data;
    }

    public void setData(Object data) {
        Data = data;
    }
}

class Data1 {
    private String name;
    private String email;

    public Data1(String name, String email) {
        this.name = name;
        this.email = email;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
    }
}

class Data2 {
    private String sex;
    private String wchart;

    public Data2(String sex, String wchart) {
        this.sex = sex;
        this.wchart = wchart;
    }

    public String getSex() {
        return sex;
    }

    public void setSex(String sex) {
        this.sex = sex;
    }

    public String getWchart() {
        return wchart;
    }

    public void setWchart(String wchart) {
        this.wchart = wchart;
    }
}

public class Main {
    private OkHttpClient client = null;
    private String RouteURL = "http://114.jghb.top:8000/source/login.cgi";
    private String RefreshURL = "http://114.jghb.top:8000/source/refresh.cgi";
    private String id = "b3a69ff2740e41319132c51271fab38a";
    private String secret = "hoTql.yUHw^H)OdM0rR1S%p%*$8^rE3(Es3xK*I%V@gyS]SR*[]qnx7X$@flkaiJ";

    private Router getRouter() {
        Router r = null;
        client = new OkHttpClient.Builder()
                .connectTimeout(10, TimeUnit.SECONDS)
                .readTimeout(10, TimeUnit.SECONDS)
                .build();
        RequestBody requestBodyPost = new FormBody.Builder()
                .add("id", id)
                .add("secret", secret)
                .build();
        Request requestPost = new Request.Builder()
                .url(RouteURL)
                .post(requestBodyPost)
                .build();
        Call call = client.newCall(requestPost);
        try {
            Response response = call.execute();
            final String string = response.body().string();
            System.out.println(string);
            r = JSON.parseObject(string, Router.class);
            System.out.println(r.getData().getToken());
        } catch (IOException e) {
            e.printStackTrace();
        }
        return r;
    }

    private boolean sendData(String uri, String token, Object object) {
        client = new OkHttpClient.Builder()
                .connectTimeout(10, TimeUnit.SECONDS)
                .readTimeout(10, TimeUnit.SECONDS)
                .build();
        RequestBody requestBody = RequestBody.create(MediaType.parse("application/json; charset=utf-8"), JSON.toJSONString(object));
        Request requestPost = new Request.Builder()
                .url(uri)
                .addHeader("authorization", token)
                .post(requestBody)
                .build();
        Call call = client.newCall(requestPost);
        try {
            Response response = call.execute();
            final String string = response.body().string();
            System.out.println(string);
        } catch (IOException e) {
            e.printStackTrace();
        }
        return true;
    }


    public static void main(String[] agv) {
        Main o = new Main();

        // 第一步：获取Router
        Router router = o.getRouter();
        if (router == null || router.getData() == null) {
            System.out.println("获取Router失败！");
            System.exit(-1);
        }
        // 第二步：发送数据到应用矩阵单元
        String uri = router.getData().getUri();
        String token = router.getData().getToken();
        System.out.println(uri);
        System.out.println(token);
        boolean b = o.sendData(uri, token, new StanderRequest(0, new Data1("admin", "pass")));


        System.out.println("=====END=====");
    }
}
