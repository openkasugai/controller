* 表示用アプリのdockerイメージをビルド

```
ubuntu@wb-m5:~/demo/rcv_video_tool$ ./build_image.sh
```

* 表示用アプリのコンテナを起動

```
ubuntu@wb-m5:~/demo/rcv_video_tool$ docker run --net=host --restart=always -v /home/ubuntu/work/DATA/rcv_video_tool:/opt/rcv_video_tool --name rcv_video_tool -itd rcv_video_tool:latest bash
```


* コンテナ内に移動

```
ubuntu@wb-m5:~/demo/rcv_video_tool$ docker exec -it rcv_video_tool bash
```

* 昼用表示アプリのプロセスを起動

```
# 8ストリーム受信ポート + モニター出力送信ポート
root@wb-m5:/opt/src# ./daytime_reciever.sh 2001 2002 2003 2004 2005 2006 2007 2008 3001
```


* 夜用表示アプリのプロセスを起動

```
# 8ストリーム受信ポート + モニター出力送信ポート
root@wb-m5:/opt/src# ./night_reciever.sh 2001 2002 2003 2004 2005 2006 2007 2008 3001
```

* 起動中の表示アプリのプロセスを停止

```
root@wb-m5:/opt/src# ./stop_reciever.sh
```

* ログの表示（コンテナ外で実行）

```
ubuntu@wb-m5:~/demo/rcv_video_tool$ docker logs -f rcv_video_tool
```

* ログファイル保存場所（コンテナ外で実行）
```
ubuntu@wb-m5:~/demo/rcv_video_tool$ docker inspect rcv_video_tool | grep 'LogPath'
```
---
Copyright 2022 NTT Corporation, FUJITSU LIMITED
