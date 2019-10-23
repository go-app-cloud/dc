import com.alibaba.fastjson.JSON;
import okhttp3.*;

import java.io.IOException;
import java.util.concurrent.TimeUnit;


class MyData {
    private String name;
    private String id;
    private String email;

    public MyData(String name, String id, String email) {
        this.name = name;
        this.id = id;
        this.email = email;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
    }
}

public class Main {
    private OkHttpClient client = null;
    private String RouteURL = "http://localhost:8000/source/login.cgi";
    private String RefreshURL = "http://localhost:8000/source/login.cgi";
    private String id = "b6c837870b0a4fcca737e9459867500d";
    private String secret = ".F!QAoCB+RfjDov%84u#A95tf1RFWyROLl2SJHBFLIR%WD$U)]HnMu4]DoTJSGDM";

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
        boolean b = o.sendData(uri, token, new MyData("a", "b", "c"));


        System.out.println("=====END=====");
    }
}
