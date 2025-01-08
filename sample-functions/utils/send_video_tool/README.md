* 映像配信ツールのdockerイメージをビルド

```
ubuntu@wb-m5:~/demo/send_video_tool$ ./build_image.sh
```

* 映像配信ツールのコンテナを起動

```
ubuntu@wb-m5:~/demo/send_video_tool$ docker run --net=host -v /home/ubuntu/work/DATA/video:/opt/video --restart=always --name send_video_tool -itd send_video_tool:latest bash
```

* コンテナに乗り込んで映像配信アプリのプロセスを起動

```
ubuntu@wb-m5:~/demo/send_video_tool$ docker exec -it send_video_tool bash
root@wb-m5:/opt/src# ./start_gst_main_sender.sh 1
```

* 映像配信アプリのプロセスを停止

```
root@wb-m5:/opt/src# ./stop_gst_sender.sh
```


* 参考：映像配信リスト設定

```
root@wb-m5:/opt/src# vim start_gst_main_sender.sh

```

```
#!/bin/bash
#./start_gst_sender.sh <ファイルパス> <宛先IP> <宛先ポート> <セッション数> <遅延間隔>

sleep_time=$1

./start_gst_sender.sh /opt/video/sample_720p.mp4 127.0.0.1 988 2 ${sleep_time:-3} &

#sleep ${sleep_time:-3}
#./start_gst_sender.sh /opt/video/sample_1080p.mp4 127.0.0.1 989 2 ${sleep_time:-3} &

```


* ログの表示
```
ubuntu@wb-m5:~/demo/send_video_tool$ docker logs -f send_video_tool
```

* ログファイル保存場所
```
ubuntu@wb-m5:~/demo/send_video_tool$ docker inspect send_video_tool | grep 'LogPath'
```

---
Copyright 2022 NTT Corporation, FUJITSU LIMITED
